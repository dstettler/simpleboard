package chess

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

// Initial game state FEN
const StartFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// Game status: different from game state, but is the
// official status of the game. This informs the actual
// status of the game, win / loss and active or unstarted
type Status byte

const (
	NotStarted Status = iota
	InProgress
	Draw
	WinWhite
	WinBlack
)

// Type representing game board in a `[rank][file]` index schema
type Board struct {
	grid [8][8]string // 8x8
}

// Board print
func (b Board) Print() {
	for _, r := range b.grid {
		fmt.Println(strings.Join(r[:], " "))
	}
}

// Go type for storing the game in memory. Produced
// using a fully parsed FEN string. Enables move
// validation, move application, scoring, etc.
type ChessGame struct {
	Board         Board  // board state representation
	Side          string // side to move "w" or "b"
	Castle        string // castling ability string
	EPTS          string // en passant target square
	HalfmoveClock int    // halfmove clock; 50 move draw rule
	FullmoveCount int    // fullmove counter
	Status        Status // status indicator
}

// Print the entire game details
func (game ChessGame) Print() {
	game.Board.Print()
	fmt.Println("Side:", game.Side)
	fmt.Println("Castle:", game.Castle)
	fmt.Println("EPTS:", game.EPTS)
	fmt.Println("Halfmove:", game.HalfmoveClock)
	fmt.Println("Fullmove:", game.FullmoveCount)
	fmt.Println("Status:", game.Status)
}

func (s Status) String() string {
	switch s {
	case NotStarted:
		return "NotStarted" // used if FEN string matches the initial game state
	case InProgress:
		return "InProgress" // used otherwise
	case Draw:
		return "Draw" // "Draw", "WinWhite", and "WinBlack are included Status values for returning a finalized game state
	case WinWhite:
		return "WinWhite"
	case WinBlack:
		return "WinBlack"
	}
	return "Invalid"
}

// FEN string to ChessGame function
// Takes an FEN string and parses it,
// returning the populated ChessGame
func ReadChessGame(fenStr string) ChessGame {

	// parse simple fields
	var b Board
	var status Status
	fields := strings.Split(fenStr, " ")
	boardStr := fields[0]
	side := fields[1]
	castle := fields[2]
	epts := fields[3]
	halfclock, err := strconv.Atoi(fields[4])
	if err != nil {
		log.Fatalf("Halfmove clock '%s' is invalid.", fields[4])
	}
	fullmoveCount, err := strconv.Atoi(fields[5])
	if err != nil {
		log.Fatalf("Fullmove count '%s' is invalid.", fields[5])
	}

	// only two statuses possible during the reading stage
	if fenStr == StartFEN {
		status = NotStarted
	} else {
		status = InProgress
	}

	// parse the board grid
	ranks := strings.Split(boardStr, "/")
	if len(ranks) != 8 {
		log.Fatalf("Invalid FEN board string: '%s'", boardStr)
	}

	// for each rank items, deal with the shorthand int notation for
	// empty space, read pieces
	for rankindex, rank := range ranks {
		fileindex := 0
		for i := 0; i < len(rank); i++ {
			token := rank[i] // get the next token in the rank
			parsedNum, err := strconv.Atoi(string(token))

			// is not an integer value
			if err != nil {
				b.grid[rankindex][fileindex] = string(token)
				fileindex++
			} else {
				// populate the empty spaces with the placeholder '.'
				for j := 0; j < parsedNum; j++ {
					b.grid[rankindex][fileindex] = "."
					fileindex++
				}
			}
		}
	}

	// return the completed result struct
	return ChessGame{
		Board:         b,
		Side:          side,
		Castle:        castle,
		EPTS:          epts,
		HalfmoveClock: halfclock,
		FullmoveCount: fullmoveCount,
		Status:        status,
	}
}
