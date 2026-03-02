package chess

import (
	"fmt"
)

// Side const for 'W' and 'B'
type Side byte

const (
	W Side = iota
	B
)

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
type board struct {
	board [8][8]byte // 8x8
}

// Go type for storing the game in memory. Produced
// using a fully parsed FEN string. Enables move
// validation, move application, scoring, etc.
type ChessGame struct {
	Board         board  // board state representation
	Side          Side   // side to move ''
	Castle        string // castling ability string
	EPTS          string // en passant target square
	HalfmoveClock int    // halfmove clock; 50 move draw rule
	FullmoveCount int    // fullmove counter
	Status        Status // status indicator
}

// string conversions
func (s Side) String() string {
	switch s {
	case W:
		return "w"
	case B:
		return "b"
	}
	return "Invalid"
}

func (s Status) String() string {
	switch s {
	case NotStarted:
		return "NotStarted"
	case InProgress:
		return "InProgress"
	case Draw:
		return "Draw"
	case WinWhite:
		return "WinWhite"
	case WinBlack:
		return "WinBlack"
	}
	return "Invalid"
}
