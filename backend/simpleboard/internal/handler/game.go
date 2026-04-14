package handler

import (
	"encoding/json"
	//"fmt"
	"net/http"

	"simpleboard/internal/auth"
	"simpleboard/internal/chess"
	"simpleboard/internal/repository"
	"simpleboard/pkg/db"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

// Game endpoint
func Game(c *gin.Context) {
	var input struct {
		Action       string `json:"action"`
		GameID       int    `json:"game_id"`
		PlayerID     uint   `json:"player_id"`
		GuestID      string `json:"guest_id"`
		OtherID      uint   `json:"other_id"`
		Move         string `json:"move"`
		StartingSide string `json:"starting_side"`
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

		// create entry
		if err := db.DB.Create(&entry).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// game successfully added
		c.JSON(http.StatusCreated, gin.H{
			"message": "game created",
			"state": gin.H{
				"game_id":         entry.ID,
				"white_player_id": entry.WhitePlayerID,
				"black_player_id": entry.BlackPlayerID,
				"white_guest_id":  entry.WhiteGuestID,
				"black_guest_id":  entry.BlackGuestID,
				"state":           entry.State,
				"status":          entry.Status,
				"side":            entry.Side,
				"next_moves":      entry.NextMoves,
				"prev_moves":      entry.PrevMoves,
				"created_at":      entry.CreatedAt,
				"updated_at":      entry.UpdatedAt,
			},
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

		c.JSON(http.StatusOK, gin.H{
			"message": "state",
			"state": gin.H{
				"game_id":         entry.ID,
				"white_player_id": entry.WhitePlayerID,
				"black_player_id": entry.BlackPlayerID,
				"white_guest_id":  entry.WhiteGuestID,
				"black_guest_id":  entry.BlackGuestID,
				"state":           entry.State,
				"status":          entry.Status,
				"side":            entry.Side,
				"next_moves":      entry.NextMoves,
				"prev_moves":      entry.PrevMoves,
				"created_at":      entry.CreatedAt,
				"updated_at":      entry.UpdatedAt,
			},
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

		updatedEntry := repository.Game{
			ID:            uint(input.GameID),
			WhitePlayerID: uint(entry.WhitePlayerID),
			BlackPlayerID: uint(entry.BlackPlayerID),
			WhiteGuestID:  entry.WhiteGuestID,
			BlackGuestID:  entry.BlackGuestID,
			State:         game.FEN(),
			Status:        game.Status.String(),
			Side:          game.Side,
			NextMoves:     nextMovesJSON,
			PrevMoves:     prevMovesJSON,
		}

		db.DB.Save(&updatedEntry)

		if err := db.DB.Where("id = ?", input.GameID).First(&entry).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid game id"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "move applied",
			"state": gin.H{
				"game_id":         entry.ID,
				"white_player_id": entry.WhitePlayerID,
				"black_player_id": entry.BlackPlayerID,
				"white_guest_id":  entry.WhiteGuestID,
				"black_guest_id":  entry.BlackGuestID,
				"state":           entry.State,
				"status":          entry.Status,
				"side":            entry.Side,
				"next_moves":      entry.NextMoves,
				"prev_moves":      entry.PrevMoves,
				"created_at":      entry.CreatedAt,
				"updated_at":      entry.UpdatedAt,
			},
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid action"})
		return
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
