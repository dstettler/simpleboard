package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"simpleboard/internal/auth"
	"simpleboard/internal/chess"
	"simpleboard/internal/repository"
	"simpleboard/internal/timer"
	"simpleboard/pkg/db"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

// Game endpoint
func Game(c *gin.Context) {
	var input struct {
		Action             string `json:"action"`
		GameID             int    `json:"game_id"`
		PlayerID           uint   `json:"player_id"`
		GuestID            string `json:"guest_id"`
		OtherID            uint   `json:"other_id"`
		Move               string `json:"move"`
		StartingSide       string `json:"starting_side"`
		TimeControlSeconds int    `json:"time_control_seconds"`
	}

	// bad request; could not bind context into input struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	action := input.Action
	if action == "create" {

		// get claims to validate
		claims := auth.GetClaims(c)
		if claims == nil || (claims.UserID == nil && claims.GuestID == nil) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "auth required"})
			return
		}

		var wpid uint = 0
		var bpid uint = 0
		var wgid string = ""
		var bgid string = ""

		// determine starting color for the "creating" player
		if input.PlayerID != 0 {
			if input.StartingSide == "w" {
				wpid = input.PlayerID
				bpid = input.OtherID
			} else if input.StartingSide == "b" {
				bpid = input.PlayerID
				wpid = input.OtherID
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "starting side not specified"})
				return
			}
		} else if input.GuestID != "" {
			if input.StartingSide == "w" {
				wgid = input.GuestID
			} else if input.StartingSide == "b" {
				bgid = input.GuestID
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "starting side not specified"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user or guest id not specified"})
			return
		}

		game := chess.ReadChessGame(chess.StartFEN, nil, nil)
		game.NextMoves = game.LegalMoves()
		game.PrevMoves = []chess.Move{}

		nextMovesStr := []string{}
		for _, m := range game.NextMoves {
			nextMovesStr = append(nextMovesStr, m.WriteMoveStr())
		}

		prevMovesStr := []string{}
		for _, m := range game.PrevMoves {
			prevMovesStr = append(prevMovesStr, m.WriteMoveStr())
		}

		nextMovesData, _ := json.Marshal(nextMovesStr)
		nextMovesJSON := datatypes.JSON(nextMovesData)

		prevMovesData, _ := json.Marshal(prevMovesStr)
		prevMovesJSON := datatypes.JSON(prevMovesData)

		entry := repository.Game{
			WhitePlayerID: wpid,
			BlackPlayerID: bpid,
			WhiteGuestID:  wgid,
			BlackGuestID:  bgid,
			State:         game.FEN(),
			Status:        chess.InProgress.String(),
			Side:          game.Side,
			NextMoves:     nextMovesJSON,
			PrevMoves:     prevMovesJSON,
		}

		// stamp initial clock state; white's time starts ticking from now
		timer.InitGameTime(&entry, input.TimeControlSeconds, time.Now())

		// create entry
		if err := db.DB.Create(&entry).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// game successfully added
		c.JSON(http.StatusCreated, gin.H{
			"message": "game created",
			"state":   gameStatePayload(&entry, entry.WhiteRemainingMs, entry.BlackRemainingMs),
		})
	} else if action == "state" {

		// get claims to validate
		claims := auth.GetClaims(c)
		if claims == nil || (claims.UserID == nil && claims.GuestID == nil) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "auth required"})
			return
		}

		// get game
		var entry repository.Game
		if err := db.DB.Where("id = ?", input.GameID).First(&entry).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid game id"})
			return
		}

		// check that the polling player is authenticated player
		_, success := colorFromID(&entry, claims)
		if !success {
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid player / guest id for state"})
			return
		}

		// compute authoritative clocks; if active side flag-fell while polling,
		// persist the timeout result so a refresh sees the final state
		whiteMs, blackMs, timedOut, loser := timer.LiveRemaining(&entry, time.Now())
		if timedOut {
			entry.Status = timer.FlagFallStatus(loser).String()
			if loser == "w" {
				entry.WhiteRemainingMs = 0
			} else {
				entry.BlackRemainingMs = 0
			}
			db.DB.Save(&entry)
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "state",
			"state":   gameStatePayload(&entry, whiteMs, blackMs),
		})
	} else if action == "move" {

		// get claims to validate
		claims := auth.GetClaims(c)
		if claims == nil || (claims.UserID == nil && claims.GuestID == nil) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "auth required"})
			return
		}

		// get game
		var entry repository.Game
		if err := db.DB.Where("id = ?", input.GameID).First(&entry).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid game id"})
			return
		}

		var moves []string
		err := json.Unmarshal(entry.NextMoves, &moves)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var prevmoves []string
		err = json.Unmarshal(entry.PrevMoves, &prevmoves)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// check that the moving player is authenticated player
		color, success := colorFromID(&entry, claims)
		if !success {
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid player / guest id for move"})
			return
		}

		if color != entry.Side {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid side to move"})
			return
		}

		// deduct the moving side's elapsed time first; if their flag fell,
		// they lose on time and the move is never applied
		if timedOut, loser := timer.ApplyMove(&entry, time.Now()); timedOut {
			entry.Status = timer.FlagFallStatus(loser).String()
			db.DB.Save(&entry)
			c.JSON(http.StatusOK, gin.H{
				"message": "flag fall",
				"state":   gameStatePayload(&entry, entry.WhiteRemainingMs, entry.BlackRemainingMs),
			})
			return
		}

		game := chess.ReadChessGame(entry.State, moves, prevmoves)

		moveErr := game.Move(input.Move)
		if moveErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": moveErr.Error()})
			return
		}

		nextMovesStr := []string{}
		for _, m := range game.NextMoves {
			nextMovesStr = append(nextMovesStr, m.WriteMoveStr())
		}

		prevMovesStr := []string{}
		for _, m := range game.PrevMoves {
			prevMovesStr = append(prevMovesStr, m.WriteMoveStr())
		}

		nextMovesData, _ := json.Marshal(nextMovesStr)
		nextMovesJSON := datatypes.JSON(nextMovesData)

		prevMovesData, _ := json.Marshal(prevMovesStr)
		prevMovesJSON := datatypes.JSON(prevMovesData)

		// update in place so timer fields (already mutated by ApplyMove) persist
		entry.State = game.FEN()
		entry.Status = game.Status.String()
		entry.Side = game.Side
		entry.NextMoves = nextMovesJSON
		entry.PrevMoves = prevMovesJSON

		db.DB.Save(&entry)

		c.JSON(http.StatusOK, gin.H{
			"message": "move applied",
			"state":   gameStatePayload(&entry, entry.WhiteRemainingMs, entry.BlackRemainingMs),
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid action"})
		return
	}
}

// builds the JSON state object returned to clients; whiteMs/blackMs are the
// authoritative live clocks (computed by the caller, may differ from stored
// values when the active side's time is still ticking)
func gameStatePayload(entry *repository.Game, whiteMs, blackMs int64) gin.H {
	return gin.H{
		"game_id":              entry.ID,
		"white_player_id":      entry.WhitePlayerID,
		"black_player_id":      entry.BlackPlayerID,
		"white_guest_id":       entry.WhiteGuestID,
		"black_guest_id":       entry.BlackGuestID,
		"state":                entry.State,
		"status":               entry.Status,
		"side":                 entry.Side,
		"next_moves":           entry.NextMoves,
		"prev_moves":           entry.PrevMoves,
		"time_control_seconds": entry.TimeControlSeconds,
		"white_remaining_ms":   whiteMs,
		"black_remaining_ms":   blackMs,
		"last_move_at":         entry.LastMoveAt,
		"server_time":          time.Now().UTC(),
		"created_at":           entry.CreatedAt,
		"updated_at":           entry.UpdatedAt,
	}
}

// returns the color of a player / guest id from a Game struct
func colorFromID(game *repository.Game, claims *auth.Claims) (string, bool) {

	// registered users
	if claims.UserID != nil {
		if game.WhitePlayerID != 0 && game.WhitePlayerID == *claims.UserID {
			return "w", true
		}
		if game.BlackPlayerID != 0 && game.BlackPlayerID == *claims.UserID {
			return "b", true
		}
	}

	// guests
	if claims.GuestID != nil {
		if game.WhiteGuestID != "" && game.WhiteGuestID == *claims.GuestID {
			return "w", true
		}
		if game.BlackGuestID != "" && game.BlackGuestID == *claims.GuestID {
			return "b", true
		}
	}

	return "", false
}
