package chess

import (
	"testing"
)

// KingCoords
// -----------

// Test both kings on the starting board
func TestKingCoordsStartPosition(t *testing.T) {
	game := ReadChessGame(StartFEN, []string{}, []string{})

	// white king at row 7, col 4 (e1)
	wr, wf := game.KingCoords(true)
	if wr != 7 || wf != 4 {
		t.Errorf("White king expected at (7,4), got (%d,%d)", wr, wf)
	}

	// black king at row 0, col 4 (e8)
	br, bf := game.KingCoords(false)
	if br != 0 || bf != 4 {
		t.Errorf("Black king expected at (0,4), got (%d,%d)", br, bf)
	}
}

// Test when king is missing
func TestKingCoordsMissingKing(t *testing.T) {
	// empty board FEN - no kings
	game := ReadChessGame("8/8/8/8/8/8/8/8 w - - 0 1", []string{}, []string{})

	r, f := game.KingCoords(true)
	if r != -1 || f != -1 {
		t.Errorf("Missing white king should return (-1,-1), got (%d,%d)", r, f)
	}
}

// MakeMove
// ----------

// Test that a piece is moved and board is updated
func TestMakeMoveUpdatesBoard(t *testing.T) {
	game := ReadChessGame(StartFEN, []string{}, []string{})

	// move pawn from e2 to e4
	m := Move{SR: 6, SF: 4, TR: 4, TF: 4, Capture: false}
	game.MakeMove(m)

	// source square becomes empty
	if game.Board.grid[6][4] != EMPTY {
		t.Errorf("Source square (6,4) should be empty after move, got %q", game.Board.grid[6][4])
	}

	// target square has pawn
	if game.Board.grid[4][4] != "P" {
		t.Errorf("Target square (4,4) should be 'P', got %q", game.Board.grid[4][4])
	}

	// side switches to black
	if game.Side != "b" {
		t.Errorf("Side should be 'b' after white moves, got %q", game.Side)
	}
}

// Test that halfmove clock resets on pawn move
func TestMakeMoveHalfmoveClock(t *testing.T) {
	game := ReadChessGame(StartFEN, []string{}, []string{})

	// move a knight - clock should increment
	knightMove := Move{SR: 7, SF: 1, TR: 5, TF: 2, Capture: false}
	game.MakeMove(knightMove)
	if game.HalfmoveClock != 1 {
		t.Errorf("Halfmove clock should be 1 after knight move, got %d", game.HalfmoveClock)
	}

	// move a pawn - clock should reset to 1
	pawnMove := Move{SR: 1, SF: 4, TR: 3, TF: 4, Capture: false}
	game.MakeMove(pawnMove)
	if game.HalfmoveClock != 1 {
		t.Errorf("Halfmove clock should reset to 1 after pawn move, got %d", game.HalfmoveClock)
	}
}

// Test that fullmove counter increments after black moves
func TestMakeMoveFullmoveCounter(t *testing.T) {
	game := ReadChessGame(StartFEN, []string{}, []string{})

	// white moves - fullmove stay at 1
	whiteMove := Move{SR: 6, SF: 4, TR: 4, TF: 4, Capture: false}
	game.MakeMove(whiteMove)
	if game.FullmoveCount != 1 {
		t.Errorf("Fullmove should still be 1 after white moves, got %d", game.FullmoveCount)
	}

	// black moves - fullmove increment to 2
	blackMove := Move{SR: 1, SF: 4, TR: 3, TF: 4, Capture: false}
	game.MakeMove(blackMove)
	if game.FullmoveCount != 2 {
		t.Errorf("Fullmove should be 2 after black moves, got %d", game.FullmoveCount)
	}
}

// Test that capture replaces the target piece
func TestMakeMoveCapture(t *testing.T) {
	// board with white knight on c3 and black pawn on d5
	game := ReadChessGame("8/8/8/3p4/8/2N5/8/8 w - - 0 1", []string{}, []string{})

	// knight captures pawn
	m := Move{SR: 5, SF: 2, TR: 3, TF: 3, Capture: true}
	game.MakeMove(m)

	if game.Board.grid[3][3] != "N" {
		t.Errorf("Target square should be 'N' after capture, got %q", game.Board.grid[3][3])
	}
	if game.Board.grid[5][2] != EMPTY {
		t.Errorf("Source square should be empty after capture, got %q", game.Board.grid[5][2])
	}
}

// Test that side toggles back to white after both sides move
func TestMakeMoveSideToggle(t *testing.T) {
	game := ReadChessGame(StartFEN, []string{}, []string{})

	game.MakeMove(Move{SR: 6, SF: 4, TR: 4, TF: 4})
	if game.Side != "b" {
		t.Errorf("Expected side 'b' after white move, got %q", game.Side)
	}

	game.MakeMove(Move{SR: 1, SF: 4, TR: 3, TF: 4})
	if game.Side != "w" {
		t.Errorf("Expected side 'w' after black move, got %q", game.Side)
	}
}

// PositionMoves
// -------------

// Test that 20 pseudo-legal moves are generatedfrom start
func TestPositionMovesStartCount(t *testing.T) {
	game := ReadChessGame(StartFEN, []string{}, []string{})

	moves := game.PositionMoves()

	// 16 pawn moves (8 single + 8 double) + 4 knight moves
	if len(moves) != 20 {
		t.Errorf("Expected 20 position moves from start, got %d", len(moves))
	}
}

// LegalMoves
// -------------

// Test that LegalMoves from start also returns 20
func TestLegalMovesStartCount(t *testing.T) {
	game := ReadChessGame(StartFEN, []string{}, []string{})
	moves := game.LegalMoves()

	if len(moves) != 20 {
		t.Errorf("Expected 20 legal moves from start, got %d", len(moves))
	}
}
