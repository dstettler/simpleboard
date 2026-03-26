package handler

import (
	"encoding/json"
	//"fmt"
	"net/http"

	//"simpleboard/internal/domain"
	"simpleboard/internal/chess"
	"simpleboard/internal/repository"
	"simpleboard/pkg/db"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

// Game endpoint
func Game(c *gin.Context) {
	var input struct {
		Action string `json:"action"`
		GameID int `json:"game_id"`
		PlayerID int `json:"player_id"`
		Move string `json:"move"`
	}

	// bad request; could not bind context into input struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	action := input.Action
	if action == "create" {
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
			WhitePlayerID: uint(input.PlayerID),
			BlackPlayerID: uint(input.PlayerID),
			State: game.FEN(),
			Status: chess.InProgress.String(),
			Side: game.Side,
			NextMoves: nextMovesJSON,
			PrevMoves: prevMovesJSON,
		}

		// create entry
		if err := db.DB.Create(&entry).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// game successfully added
		c.JSON(http.StatusCreated, gin.H{
			"message": "game created",
			"user": gin.H{
				"game_id":  entry.ID,
				"white_player_id": entry.WhitePlayerID,
				"black_player_id": entry.BlackPlayerID,
				"state": entry.State,
				"status": entry.Status,
				"side": entry.Side,
				"next_moves": entry.NextMoves,
				"prev_moves": entry.PrevMoves,
				"created_at": entry.CreatedAt,
				"updated_at": entry.UpdatedAt,
			},
		})
	} else if action == "state" {

		// get game
		var entry repository.Game
		if err := db.DB.Where("id = ?", input.GameID).First(&entry).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid game id"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "game created",
			"user": gin.H{
				"game_id":  entry.ID,
				"white_player_id": entry.WhitePlayerID,
				"black_player_id": entry.BlackPlayerID,
				"state": entry.State,
				"status": entry.Status,
				"side": entry.Side,
				"next_moves": entry.NextMoves,
				"prev_moves": entry.PrevMoves,
				"created_at": entry.CreatedAt,
				"updated_at": entry.UpdatedAt,
			},
		})
	} else if action == "move" {
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

		game := chess.ReadChessGame(entry.State, moves, prevmoves)

		game.Move(input.Move)

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
			ID: uint(input.GameID),
			WhitePlayerID: uint(input.PlayerID),
			BlackPlayerID: uint(input.PlayerID),
			State: game.FEN(),
			Status: game.Status.String(),
			Side: game.Side,
			NextMoves: nextMovesJSON,
			PrevMoves: prevMovesJSON,
		}

		db.DB.Save(&updatedEntry)

		c.JSON(http.StatusOK, gin.H{
			"message": "move applied",
			"user": gin.H{
				"game_id":  entry.ID,
				"white_player_id": entry.WhitePlayerID,
				"black_player_id": entry.BlackPlayerID,
				"state": entry.State,
				"status": entry.Status,
				"side": entry.Side,
				"next_moves": entry.NextMoves,
				"prev_moves": entry.PrevMoves,
				"created_at": entry.CreatedAt,
				"updated_at": entry.UpdatedAt,
			},
		})
	} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid action"})
			return
	}
}
