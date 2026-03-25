package chess

import (
	"testing"
)

// boardFromRanks builds a Board from 8 strings, one per rank
// Each string should have 8 characters (pieces or '.' for empty).
func boardFromRanks(ranks [8]string) Board {
	var b Board
	for r, rank := range ranks {
		for f, ch := range rank {
			b.grid[r][f] = string(ch)
		}
	}
	return b
}

// boardsEqual checks if two boards have the same piece layout
func boardsEqual(a, b Board) bool {
	for r := 0; r < 8; r++ {
		for f := 0; f < 8; f++ {
			if a.grid[r][f] != b.grid[r][f] {
				return false
			}
		}
	}
	return true
}

// Status.String() - test
// -----------------------
func TestStatusString(t *testing.T) {
	cases := []struct {
		s    Status
		want string
	}{
		{NotStarted, "NotStarted"},
		{InProgress, "InProgress"},
		{Draw, "Draw"},
		{WinWhite, "WinWhite"},
		{WinBlack, "WinBlack"},
		{Status(255), "Invalid"}, // out-of-range value
	}

	for _, tc := range cases {
		got := tc.s.String()
		if got != tc.want {
			t.Errorf("Status(%d).String() = %q, want %q", tc.s, got, tc.want)
		}
	}
}

// ReadChessGame — starting position tests
// ---------------------------------------
func TestReadStartingPosition(t *testing.T) {
	game := ReadChessGame(StartFEN, nil, nil)

	// The starting FEN should always give NotStarted status
	if game.Status != NotStarted {
		t.Errorf("expected NotStarted status for StartFEN, got %s", game.Status)
	}

	if game.Side != "w" {
		t.Errorf("side = %q, want \"w\"", game.Side)
	}
	if game.Castle != "KQkq" {
		t.Errorf("castle = %q, want \"KQkq\"", game.Castle)
	}
	if game.EPTS != "-" {
		t.Errorf("en passant = %q, want \"-\"", game.EPTS)
	}
	if game.HalfmoveClock != 0 {
		t.Errorf("halfmove clock = %d, want 0", game.HalfmoveClock)
	}
	if game.FullmoveCount != 1 {
		t.Errorf("fullmove count = %d, want 1", game.FullmoveCount)
	}

	// Verify the actual board layout
	expected := boardFromRanks([8]string{
		"rnbqkbnr",
		"pppppppp",
		"........",
		"........",
		"........",
		"........",
		"PPPPPPPP",
		"RNBQKBNR",
	})
	if !boardsEqual(game.Board, expected) {
		t.Error("starting board layout does not match expected")
		t.Log("got:")
		game.Board.Print()
	}
}

// ReadChessGame — mid-game position tests
// ---------------------------------------
func TestReadAfterE4(t *testing.T) {
	// Position after 1. e4
	fen := "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"
	game := ReadChessGame(fen, nil, nil)

	if game.Status != InProgress {
		t.Errorf("expected InProgress, got %s", game.Status)
	}
	if game.Side != "b" {
		t.Errorf("side = %q, want \"b\"", game.Side)
	}
	if game.EPTS != "e3" {
		t.Errorf("en passant = %q, want \"e3\"", game.EPTS)
	}

	// Pawn should be on e4 (rank 4 = index 4, file e = index 4)
	if game.Board.grid[4][4] != "P" {
		t.Errorf("expected white pawn on e4, got %q", game.Board.grid[4][4])
	}
	// e2 should now be empty
	if game.Board.grid[6][4] != "." {
		t.Errorf("expected e2 to be empty after 1.e4, got %q", game.Board.grid[6][4])
	}
}

func TestReadCastledPosition(t *testing.T) {
	// White has castled kingside, black hasn't yet
	fen := "r1bqkb1r/pppp1ppp/2n2n2/4p3/2B1P3/5N2/PPPP1PPP/RNBQ1RK1 b kq - 5 4"
	game := ReadChessGame(fen, nil, nil)

	// White already castled, so only black has castling rights
	if game.Castle != "kq" {
		t.Errorf("castle = %q, want \"kq\"", game.Castle)
	}

	// White king on g1, rook on f1
	if game.Board.grid[7][6] != "K" {
		t.Errorf("expected white king on g1, got %q", game.Board.grid[7][6])
	}
	if game.Board.grid[7][5] != "R" {
		t.Errorf("expected white rook on f1, got %q", game.Board.grid[7][5])
	}
}

func TestReadNoCastlingRights(t *testing.T) {
	// Late middlegame — no castling rights for either side
	fen := "r4rk1/pp2ppbp/2np1np1/q7/3PP1b1/2N1BN2/PPP1BPPP/R2Q1RK1 b - - 7 10"
	game := ReadChessGame(fen, nil, nil)

	if game.Castle != "-" {
		t.Errorf("castle = %q, want \"-\" (no castling rights)", game.Castle)
	}
	if game.HalfmoveClock != 7 {
		t.Errorf("halfmove = %d, want 7", game.HalfmoveClock)
	}
	if game.FullmoveCount != 10 {
		t.Errorf("fullmove = %d, want 10", game.FullmoveCount)
	}
}

func TestReadEndgamePosition(t *testing.T) {
	// King + rook vs king endgame
	fen := "8/8/8/4k3/8/8/2K5/7R w - - 0 50"
	game := ReadChessGame(fen, nil, nil)

	if game.Side != "w" {
		t.Errorf("side = %q, want \"w\"", game.Side)
	}
	if game.FullmoveCount != 50 {
		t.Errorf("fullmove = %d, want 50", game.FullmoveCount)
	}

	// Most squares should be empty
	emptyCount := 0
	for r := 0; r < 8; r++ {
		for f := 0; f < 8; f++ {
			if game.Board.grid[r][f] == "." {
				emptyCount++
			}
		}
	}
	// 64 squares minus 3 pieces = 61 empty
	if emptyCount != 61 {
		t.Errorf("expected 61 empty squares, got %d", emptyCount)
	}

	// Verify piece placements
	if game.Board.grid[3][4] != "k" {
		t.Errorf("expected black king on e5, got %q", game.Board.grid[3][4])
	}
	if game.Board.grid[6][2] != "K" {
		t.Errorf("expected white king on c2, got %q", game.Board.grid[6][2])
	}
	if game.Board.grid[7][7] != "R" {
		t.Errorf("expected white rook on h1, got %q", game.Board.grid[7][7])
	}
}

// ReadChessGame — board grid edge cases - all empty, full rank, etc
// ---------------------------------------------------------------
func TestBoardAllEmpty(t *testing.T) {
	// Hypothetical position: kings only (minimum legal-ish board)
	fen := "4k3/8/8/8/8/8/8/4K3 w - - 0 1"
	game := ReadChessGame(fen, nil, nil)

	// Two kings and 62 empty squares
	if game.Board.grid[0][4] != "k" {
		t.Errorf("expected black king on e8, got %q", game.Board.grid[0][4])
	}
	if game.Board.grid[7][4] != "K" {
		t.Errorf("expected white king on e1, got %q", game.Board.grid[7][4])
	}

	pieceCount := 0
	for r := 0; r < 8; r++ {
		for f := 0; f < 8; f++ {
			if game.Board.grid[r][f] != "." {
				pieceCount++
			}
		}
	}
	if pieceCount != 2 {
		t.Errorf("expected exactly 2 pieces on board, counted %d", pieceCount)
	}
}

func TestBoardFullRank(t *testing.T) {
	// Rank 1 and 8 are fully packed with pieces (no digit shorthand)
	game := ReadChessGame(StartFEN, nil, nil)

	// Back ranks should have no empty squares
	for f := 0; f < 8; f++ {
		if game.Board.grid[0][f] == "." {
			t.Errorf("rank 8 file %d should not be empty in start position", f)
		}
		if game.Board.grid[7][f] == "." {
			t.Errorf("rank 1 file %d should not be empty in start position", f)
		}
	}

	// Middle ranks (3-6) should be all empty
	for r := 2; r < 6; r++ {
		for f := 0; f < 8; f++ {
			if game.Board.grid[r][f] != "." {
				t.Errorf("rank %d file %d should be empty in start position, got %q",
					8-r, f, game.Board.grid[r][f])
			}
		}
	}
}

// Table-driven test: known position tests
// ------------------------------------------
func TestReadChessGameTable(t *testing.T) {
	tests := []struct {
		name     string
		fen      string
		side     string
		castle   string
		epts     string
		halfmove int
		fullmove int
		status   Status
	}{
		{
			name:     "start position",
			fen:      StartFEN,
			side:     "w",
			castle:   "KQkq",
			epts:     "-",
			halfmove: 0,
			fullmove: 1,
			status:   NotStarted,
		},
		{
			name:     "after 1.e4",
			fen:      "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
			side:     "b",
			castle:   "KQkq",
			epts:     "e3",
			halfmove: 0,
			fullmove: 1,
			status:   InProgress,
		},
		{
			name:     "king and pawn endgame",
			fen:      "8/5k2/8/5P2/5K2/8/8/8 w - - 3 45",
			side:     "w",
			castle:   "-",
			epts:     "-",
			halfmove: 3,
			fullmove: 45,
			status:   InProgress,
		},
		{
			name:     "white kingside castling only",
			fen:      "r3kb1r/pppq1ppp/5n2/3p4/3P4/5N2/PPP2PPP/R1BQ1RK1 b Kkq - 0 9",
			side:     "b",
			castle:   "Kkq",
			epts:     "-",
			halfmove: 0,
			fullmove: 9,
			status:   InProgress,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			game := ReadChessGame(tc.fen, nil, nil)

			if game.Side != tc.side {
				t.Errorf("side = %q, want %q", game.Side, tc.side)
			}
			if game.Castle != tc.castle {
				t.Errorf("castle = %q, want %q", game.Castle, tc.castle)
			}
			if game.EPTS != tc.epts {
				t.Errorf("epts = %q, want %q", game.EPTS, tc.epts)
			}
			if game.HalfmoveClock != tc.halfmove {
				t.Errorf("halfmove = %d, want %d", game.HalfmoveClock, tc.halfmove)
			}
			if game.FullmoveCount != tc.fullmove {
				t.Errorf("fullmove = %d, want %d", game.FullmoveCount, tc.fullmove)
			}
			if game.Status != tc.status {
				t.Errorf("status = %s, want %s", game.Status, tc.status)
			}
		})
	}
}
