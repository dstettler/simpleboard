package chess

import (
	"strings"
	"testing"
)

// --- helpers ---

// boardFromRanks builds a Board from 8 rank strings (rank 8 first).
// Each string must be 8 chars wide; use '.' for empty squares.
func boardFromRanks(ranks [8]string) Board {
	var b Board
	for r, rank := range ranks {
		for f, ch := range rank {
			b.grid[r][f] = string(ch)
		}
	}
	return b
}

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

// containsMove checks whether a move string appears in a move slice.
func containsMove(moves []Move, moveStr string) bool {
	m := ParseMoveStr(moveStr)
	for _, lm := range moves {
		if m.IsEqual(lm) {
			return true
		}
	}
	return false
}

// --- Status ---

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
		{Status(255), "Invalid"},
	}
	for _, tc := range cases {
		got := tc.s.String()
		if got != tc.want {
			t.Errorf("Status(%d).String() = %q, want %q", tc.s, got, tc.want)
		}
	}
}

// --- ParseCoords / ParseAlg ---

func TestParseCoords(t *testing.T) {
	cases := []struct {
		sq    string
		wantr int
		wantf int
	}{
		{"a1", 7, 0},
		{"h8", 0, 7},
		{"e4", 4, 4},
		{"d2", 6, 3},
		{"g6", 2, 6},
	}
	for _, tc := range cases {
		r, f := ParseCoords(tc.sq)
		if r != tc.wantr || f != tc.wantf {
			t.Errorf("ParseCoords(%q) = (%d,%d), want (%d,%d)", tc.sq, r, f, tc.wantr, tc.wantf)
		}
	}
}

func TestParseAlg(t *testing.T) {
	cases := []struct {
		r, f int
		want string
	}{
		{7, 0, "a1"},
		{0, 7, "h8"},
		{4, 4, "e4"},
		{6, 3, "d2"},
		{2, 6, "g6"},
	}
	for _, tc := range cases {
		got := ParseAlg(tc.r, tc.f)
		if got != tc.want {
			t.Errorf("ParseAlg(%d,%d) = %q, want %q", tc.r, tc.f, got, tc.want)
		}
	}
}

func TestParseAlgCoordsRoundTrip(t *testing.T) {
	squares := []string{"a1", "h8", "e4", "d7", "c3", "f5", "b2", "g6"}
	for _, sq := range squares {
		r, f := ParseCoords(sq)
		got := ParseAlg(r, f)
		if got != sq {
			t.Errorf("round-trip %q → (%d,%d) → %q", sq, r, f, got)
		}
	}
}

// --- Board.FEN ---

func TestBoardFENStartPosition(t *testing.T) {
	game := ReadChessGame(StartFEN, nil, nil)
	want := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR"
	if got := game.Board.FEN(); got != want {
		t.Errorf("Board.FEN() = %q, want %q", got, want)
	}
}

func TestBoardFENAfterE4(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"
	game := ReadChessGame(fen, nil, nil)
	want := "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR"
	if got := game.Board.FEN(); got != want {
		t.Errorf("Board.FEN() = %q, want %q", got, want)
	}
}

// --- Move struct ---

func TestMoveIsEqual(t *testing.T) {
	base := Move{SR: 6, SF: 4, TR: 4, TF: 4}

	if !base.IsEqual(base) {
		t.Error("identical moves should be equal")
	}

	cases := []struct {
		name string
		mod  Move
	}{
		{"different SR", Move{SR: 5, SF: 4, TR: 4, TF: 4}},
		{"different SF", Move{SR: 6, SF: 3, TR: 4, TF: 4}},
		{"different TR", Move{SR: 6, SF: 4, TR: 3, TF: 4}},
		{"different TF", Move{SR: 6, SF: 4, TR: 4, TF: 5}},
		{"capture flag", Move{SR: 6, SF: 4, TR: 4, TF: 4, Capture: true}},
		{"castling flag", Move{SR: 6, SF: 4, TR: 4, TF: 4, Castling: true}},
		{"promotion", Move{SR: 6, SF: 4, TR: 4, TF: 4, Promotion: "Q"}},
	}
	for _, tc := range cases {
		if base.IsEqual(tc.mod) {
			t.Errorf("IsEqual should be false for %s", tc.name)
		}
	}
}

func TestMoveCopy(t *testing.T) {
	orig := Move{SR: 1, SF: 2, TR: 3, TF: 4, Capture: true, Promotion: "Q"}
	cp := orig.Copy()

	if !orig.IsEqual(cp) {
		t.Error("copy should equal original")
	}

	cp.SR = 99
	cp.Promotion = "R"
	if orig.SR == 99 || orig.Promotion == "R" {
		t.Error("mutating copy should not affect original")
	}
}

func TestParseMoveStr(t *testing.T) {
	cases := []struct {
		s        string
		sr, sf   int
		tr, tf   int
		capture  bool
		castling bool
		promo    string
	}{
		{"e2e4", 6, 4, 4, 4, false, false, ""},
		{"a1a8", 7, 0, 0, 0, false, false, ""},
		{"e4xd5", 4, 4, 3, 3, true, false, ""},
		{"e1g1", 7, 4, 7, 6, false, true, ""},
		{"e1c1", 7, 4, 7, 2, false, true, ""},
		{"e8g8", 0, 4, 0, 6, false, true, ""},
		{"e8c8", 0, 4, 0, 2, false, true, ""},
		{"e7e8Q", 1, 4, 0, 4, false, false, "Q"},
		{"a7a8n", 1, 0, 0, 0, false, false, "n"},
		{"e7xd8q", 1, 4, 0, 3, true, false, "q"},
	}
	for _, tc := range cases {
		m := ParseMoveStr(tc.s)
		if m.SR != tc.sr || m.SF != tc.sf || m.TR != tc.tr || m.TF != tc.tf ||
			m.Capture != tc.capture || m.Castling != tc.castling || m.Promotion != tc.promo {
			t.Errorf("ParseMoveStr(%q) = %+v, want SR=%d SF=%d TR=%d TF=%d cap=%v cas=%v promo=%q",
				tc.s, m, tc.sr, tc.sf, tc.tr, tc.tf, tc.capture, tc.castling, tc.promo)
		}
	}
}

func TestWriteMoveStrRoundTrip(t *testing.T) {
	for _, s := range []string{"e2e4", "e4xd5", "e1g1", "e1c1", "e8g8", "e8c8", "e7e8Q", "e7xd8q"} {
		if got := ParseMoveStr(s).WriteMoveStr(); got != s {
			t.Errorf("WriteMoveStr(ParseMoveStr(%q)) = %q", s, got)
		}
	}
}

// --- ReadChessGame ---

func TestReadStartingPosition(t *testing.T) {
	game := ReadChessGame(StartFEN, nil, nil)

	if game.Status != NotStarted {
		t.Errorf("status = %s, want NotStarted", game.Status)
	}
	if game.Side != "w" {
		t.Errorf("side = %q, want \"w\"", game.Side)
	}
	if game.Castle != "KQkq" {
		t.Errorf("castle = %q, want \"KQkq\"", game.Castle)
	}
	if game.EPTS != "-" {
		t.Errorf("epts = %q, want \"-\"", game.EPTS)
	}
	if game.HalfmoveClock != 0 {
		t.Errorf("halfmove = %d, want 0", game.HalfmoveClock)
	}
	if game.FullmoveCount != 1 {
		t.Errorf("fullmove = %d, want 1", game.FullmoveCount)
	}

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
	}
}

func TestReadAfterE4(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"
	game := ReadChessGame(fen, nil, nil)

	if game.Status != InProgress {
		t.Errorf("status = %s, want InProgress", game.Status)
	}
	if game.Side != "b" {
		t.Errorf("side = %q, want \"b\"", game.Side)
	}
	if game.EPTS != "e3" {
		t.Errorf("epts = %q, want \"e3\"", game.EPTS)
	}
	if game.Board.grid[4][4] != "P" {
		t.Errorf("e4 = %q, want \"P\"", game.Board.grid[4][4])
	}
	if game.Board.grid[6][4] != "." {
		t.Errorf("e2 = %q, want \".\"", game.Board.grid[6][4])
	}
}

func TestReadCastledPosition(t *testing.T) {
	// White has castled kingside, only black can still castle.
	fen := "r1bqkb1r/pppp1ppp/2n2n2/4p3/2B1P3/5N2/PPPP1PPP/RNBQ1RK1 b kq - 5 4"
	game := ReadChessGame(fen, nil, nil)

	if game.Castle != "kq" {
		t.Errorf("castle = %q, want \"kq\"", game.Castle)
	}
	if game.Board.grid[7][6] != "K" {
		t.Errorf("g1 = %q, want \"K\"", game.Board.grid[7][6])
	}
	if game.Board.grid[7][5] != "R" {
		t.Errorf("f1 = %q, want \"R\"", game.Board.grid[7][5])
	}
}

func TestReadNoCastlingRights(t *testing.T) {
	fen := "r4rk1/pp2ppbp/2np1np1/q7/3PP1b1/2N1BN2/PPP1BPPP/R2Q1RK1 b - - 7 10"
	game := ReadChessGame(fen, nil, nil)

	if game.Castle != "-" {
		t.Errorf("castle = %q, want \"-\"", game.Castle)
	}
	if game.HalfmoveClock != 7 {
		t.Errorf("halfmove = %d, want 7", game.HalfmoveClock)
	}
	if game.FullmoveCount != 10 {
		t.Errorf("fullmove = %d, want 10", game.FullmoveCount)
	}
}

func TestReadEndgamePosition(t *testing.T) {
	fen := "8/8/8/4k3/8/8/2K5/7R w - - 0 50"
	game := ReadChessGame(fen, nil, nil)

	if game.Side != "w" {
		t.Errorf("side = %q, want \"w\"", game.Side)
	}
	if game.FullmoveCount != 50 {
		t.Errorf("fullmove = %d, want 50", game.FullmoveCount)
	}

	// 3 pieces means 61 empty squares
	empty := 0
	for r := 0; r < 8; r++ {
		for f := 0; f < 8; f++ {
			if game.Board.grid[r][f] == "." {
				empty++
			}
		}
	}
	if empty != 61 {
		t.Errorf("empty squares = %d, want 61", empty)
	}
}

func TestBoardAllEmpty(t *testing.T) {
	fen := "4k3/8/8/8/8/8/8/4K3 w - - 0 1"
	game := ReadChessGame(fen, nil, nil)

	if game.Board.grid[0][4] != "k" {
		t.Errorf("e8 = %q, want \"k\"", game.Board.grid[0][4])
	}
	if game.Board.grid[7][4] != "K" {
		t.Errorf("e1 = %q, want \"K\"", game.Board.grid[7][4])
	}

	pieces := 0
	for r := 0; r < 8; r++ {
		for f := 0; f < 8; f++ {
			if game.Board.grid[r][f] != "." {
				pieces++
			}
		}
	}
	if pieces != 2 {
		t.Errorf("piece count = %d, want 2", pieces)
	}
}

func TestBoardFullRank(t *testing.T) {
	game := ReadChessGame(StartFEN, nil, nil)

	for f := 0; f < 8; f++ {
		if game.Board.grid[0][f] == "." {
			t.Errorf("rank 8 file %d should not be empty", f)
		}
		if game.Board.grid[7][f] == "." {
			t.Errorf("rank 1 file %d should not be empty", f)
		}
	}
	// ranks 3–6 (index 2–5) are all empty in the start position
	for r := 2; r < 6; r++ {
		for f := 0; f < 8; f++ {
			if game.Board.grid[r][f] != "." {
				t.Errorf("rank %d file %d should be empty, got %q", 8-r, f, game.Board.grid[r][f])
			}
		}
	}
}

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
			side:     "w", castle: "KQkq", epts: "-",
			halfmove: 0, fullmove: 1, status: NotStarted,
		},
		{
			name:     "after 1.e4",
			fen:      "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
			side:     "b", castle: "KQkq", epts: "e3",
			halfmove: 0, fullmove: 1, status: InProgress,
		},
		{
			name:     "king and pawn endgame",
			fen:      "8/5k2/8/5P2/5K2/8/8/8 w - - 3 45",
			side:     "w", castle: "-", epts: "-",
			halfmove: 3, fullmove: 45, status: InProgress,
		},
		{
			name:     "white kingside castling only",
			fen:      "r3kb1r/pppq1ppp/5n2/3p4/3P4/5N2/PPP2PPP/R1BQ1RK1 b Kkq - 0 9",
			side:     "b", castle: "Kkq", epts: "-",
			halfmove: 0, fullmove: 9, status: InProgress,
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

// --- ChessGame.FEN / ChessGame.Copy ---

func TestChessGameFENRoundTrip(t *testing.T) {
	fens := []string{
		StartFEN,
		"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
		"r1bqkb1r/pppp1ppp/2n2n2/4p3/2B1P3/5N2/PPPP1PPP/RNBQ1RK1 b kq - 5 4",
		"8/8/8/4k3/8/8/2K5/7R w - - 0 50",
	}
	for _, fen := range fens {
		game := ReadChessGame(fen, nil, nil)
		if got := game.FEN(); got != fen {
			t.Errorf("FEN round-trip:\n  in:  %q\n  out: %q", fen, got)
		}
	}
}

func TestChessGameCopyIndependence(t *testing.T) {
	orig := ReadChessGame(StartFEN, nil, nil)
	orig.NextMoves = orig.LegalMoves()
	cp := orig.Copy()

	cp.Side = "b"
	cp.Castle = "-"
	cp.EPTS = "e3"
	cp.HalfmoveClock = 42
	cp.FullmoveCount = 99
	cp.Board.grid[0][0] = "Q"

	if orig.Side != "w" {
		t.Error("Side mutated through copy")
	}
	if orig.Castle != "KQkq" {
		t.Error("Castle mutated through copy")
	}
	if orig.EPTS != "-" {
		t.Error("EPTS mutated through copy")
	}
	if orig.HalfmoveClock != 0 {
		t.Error("HalfmoveClock mutated through copy")
	}
	if orig.FullmoveCount != 1 {
		t.Error("FullmoveCount mutated through copy")
	}
	if orig.Board.grid[0][0] != "r" {
		t.Error("Board mutated through copy")
	}
}

// --- KingCoords ---

func TestKingCoordsStartPosition(t *testing.T) {
	game := ReadChessGame(StartFEN, nil, nil)

	wr, wf := game.KingCoords(true)
	if wr != 7 || wf != 4 {
		t.Errorf("white king: want (7,4), got (%d,%d)", wr, wf)
	}
	br, bf := game.KingCoords(false)
	if br != 0 || bf != 4 {
		t.Errorf("black king: want (0,4), got (%d,%d)", br, bf)
	}
}

func TestKingCoordsMissingKing(t *testing.T) {
	game := ReadChessGame("8/8/8/8/8/8/8/8 w - - 0 1", nil, nil)
	r, f := game.KingCoords(true)
	if r != -1 || f != -1 {
		t.Errorf("missing king: want (-1,-1), got (%d,%d)", r, f)
	}
}

func TestKingCoordsCustomPosition(t *testing.T) {
	// white on a1, black on h8
	game := ReadChessGame("7k/8/8/8/8/8/8/K7 w - - 0 1", nil, nil)

	wr, wf := game.KingCoords(true)
	if wr != 7 || wf != 0 {
		t.Errorf("white king: want (7,0), got (%d,%d)", wr, wf)
	}
	br, bf := game.KingCoords(false)
	if br != 0 || bf != 7 {
		t.Errorf("black king: want (0,7), got (%d,%d)", br, bf)
	}
}

// --- MakeMove ---

func TestMakeMoveUpdatesBoard(t *testing.T) {
	game := ReadChessGame(StartFEN, nil, nil)
	game.MakeMove(Move{SR: 6, SF: 4, TR: 4, TF: 4}) // e2e4

	if game.Board.grid[6][4] != EMPTY {
		t.Errorf("e2 should be empty after e4, got %q", game.Board.grid[6][4])
	}
	if game.Board.grid[4][4] != "P" {
		t.Errorf("e4 should be 'P', got %q", game.Board.grid[4][4])
	}
	if game.Side != "b" {
		t.Errorf("side should be 'b' after white move, got %q", game.Side)
	}
}

func TestMakeMoveCapture(t *testing.T) {
	// white knight c3, black pawn d5
	game := ReadChessGame("8/8/8/3p4/8/2N5/8/8 w - - 0 1", nil, nil)
	game.MakeMove(Move{SR: 5, SF: 2, TR: 3, TF: 3, Capture: true})

	if game.Board.grid[3][3] != "N" {
		t.Errorf("d5 should be 'N' after capture, got %q", game.Board.grid[3][3])
	}
	if game.Board.grid[5][2] != EMPTY {
		t.Errorf("c3 should be empty after move, got %q", game.Board.grid[5][2])
	}
	if game.HalfmoveClock != 1 {
		t.Errorf("halfmove should reset to 1 after capture, got %d", game.HalfmoveClock)
	}
}

func TestMakeMoveHalfmoveClock(t *testing.T) {
	game := ReadChessGame(StartFEN, nil, nil)

	game.MakeMove(Move{SR: 7, SF: 1, TR: 5, TF: 2}) // Nb1c3
	if game.HalfmoveClock != 1 {
		t.Errorf("after knight move: halfmove = %d, want 1", game.HalfmoveClock)
	}

	game.MakeMove(Move{SR: 1, SF: 4, TR: 3, TF: 4}) // e7e5
	if game.HalfmoveClock != 1 {
		t.Errorf("after pawn move: halfmove = %d, want 1 (should reset)", game.HalfmoveClock)
	}

	game.MakeMove(Move{SR: 5, SF: 2, TR: 4, TF: 4}) // Nc3e4
	if game.HalfmoveClock != 2 {
		t.Errorf("after second piece move: halfmove = %d, want 2", game.HalfmoveClock)
	}
}

func TestMakeMoveFullmoveCounter(t *testing.T) {
	game := ReadChessGame(StartFEN, nil, nil)

	game.MakeMove(Move{SR: 6, SF: 4, TR: 4, TF: 4}) // white
	if game.FullmoveCount != 1 {
		t.Errorf("after white move: fullmove = %d, want 1", game.FullmoveCount)
	}

	game.MakeMove(Move{SR: 1, SF: 4, TR: 3, TF: 4}) // black — counter increments
	if game.FullmoveCount != 2 {
		t.Errorf("after black move: fullmove = %d, want 2", game.FullmoveCount)
	}

	game.MakeMove(Move{SR: 7, SF: 6, TR: 5, TF: 5}) // white again — stays at 2
	if game.FullmoveCount != 2 {
		t.Errorf("after second white move: fullmove = %d, want 2", game.FullmoveCount)
	}
}

func TestMakeMoveSideToggle(t *testing.T) {
	game := ReadChessGame(StartFEN, nil, nil)

	game.MakeMove(Move{SR: 6, SF: 4, TR: 4, TF: 4})
	if game.Side != "b" {
		t.Errorf("after white move: side = %q, want \"b\"", game.Side)
	}
	game.MakeMove(Move{SR: 1, SF: 4, TR: 3, TF: 4})
	if game.Side != "w" {
		t.Errorf("after black move: side = %q, want \"w\"", game.Side)
	}
}

func TestMakeMoveEPTSSetOnDoublePush(t *testing.T) {
	game := ReadChessGame(StartFEN, nil, nil)
	game.MakeMove(Move{SR: 6, SF: 4, TR: 4, TF: 4}) // e2e4
	if game.EPTS != "e3" {
		t.Errorf("EPTS after e2e4 = %q, want \"e3\"", game.EPTS)
	}
}

func TestMakeMoveEPTSClearedAfterNonPawn(t *testing.T) {
	// EPTS is set from e2e4; black plays a knight — must clear to "-"
	game := ReadChessGame(
		"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
		nil, nil,
	)
	game.MakeMove(Move{SR: 0, SF: 1, TR: 2, TF: 2}) // Nb8c6
	if game.EPTS != "-" {
		t.Errorf("EPTS after non-pawn move = %q, want \"-\"", game.EPTS)
	}
}

func TestMakeMoveEnPassantWhite(t *testing.T) {
	// white pawn e5, black just played d7d5 — EPTS = d6
	fen := "rnbqkbnr/ppp1pppp/8/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 3"
	game := ReadChessGame(fen, nil, nil)
	game.MakeMove(Move{SR: 3, SF: 4, TR: 2, TF: 3, Capture: true}) // e5xd6

	if game.Board.grid[2][3] != "P" {
		t.Errorf("d6 should have 'P' after en passant, got %q", game.Board.grid[2][3])
	}
	// the captured pawn on d5 should disappear
	if game.Board.grid[3][3] != EMPTY {
		t.Errorf("d5 should be empty after en passant capture, got %q", game.Board.grid[3][3])
	}
	if game.Board.grid[3][4] != EMPTY {
		t.Errorf("e5 should be empty after pawn moved, got %q", game.Board.grid[3][4])
	}
}

func TestMakeMoveEnPassantBlack(t *testing.T) {
	// black pawn e4, white just played d2d4 — EPTS = d3
	fen := "rnbqkbnr/pppp1ppp/8/8/3Pp3/8/PPP2PPP/RNBQKBNR b KQkq d3 0 3"
	game := ReadChessGame(fen, nil, nil)
	game.MakeMove(Move{SR: 4, SF: 4, TR: 5, TF: 3, Capture: true}) // e4xd3

	if game.Board.grid[5][3] != "p" {
		t.Errorf("d3 should have 'p' after en passant, got %q", game.Board.grid[5][3])
	}
	if game.Board.grid[4][3] != EMPTY {
		t.Errorf("d4 should be empty after en passant capture, got %q", game.Board.grid[4][3])
	}
}

func TestMakeMoveCastlingKingsideWhite(t *testing.T) {
	fen := "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQkq - 0 1"
	game := ReadChessGame(fen, nil, nil)
	game.MakeMove(ParseMoveStr("e1g1"))

	if game.Board.grid[7][6] != "K" {
		t.Errorf("g1 should be 'K', got %q", game.Board.grid[7][6])
	}
	if game.Board.grid[7][5] != "R" {
		t.Errorf("f1 should be 'R', got %q", game.Board.grid[7][5])
	}
	if game.Board.grid[7][4] != EMPTY {
		t.Errorf("e1 should be empty, got %q", game.Board.grid[7][4])
	}
	if game.Board.grid[7][7] != EMPTY {
		t.Errorf("h1 should be empty after O-O, got %q", game.Board.grid[7][7])
	}
}

func TestMakeMoveCastlingQueensideWhite(t *testing.T) {
	fen := "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQkq - 0 1"
	game := ReadChessGame(fen, nil, nil)
	game.MakeMove(ParseMoveStr("e1c1"))

	if game.Board.grid[7][2] != "K" {
		t.Errorf("c1 should be 'K', got %q", game.Board.grid[7][2])
	}
	if game.Board.grid[7][3] != "R" {
		t.Errorf("d1 should be 'R', got %q", game.Board.grid[7][3])
	}
	if game.Board.grid[7][4] != EMPTY {
		t.Errorf("e1 should be empty, got %q", game.Board.grid[7][4])
	}
	if game.Board.grid[7][0] != EMPTY {
		t.Errorf("a1 should be empty after O-O-O, got %q", game.Board.grid[7][0])
	}
}

func TestMakeMoveCastlingKingsideBlack(t *testing.T) {
	fen := "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R b KQkq - 0 1"
	game := ReadChessGame(fen, nil, nil)
	game.MakeMove(ParseMoveStr("e8g8"))

	if game.Board.grid[0][6] != "k" {
		t.Errorf("g8 should be 'k', got %q", game.Board.grid[0][6])
	}
	if game.Board.grid[0][5] != "r" {
		t.Errorf("f8 should be 'r', got %q", game.Board.grid[0][5])
	}
	if game.Board.grid[0][4] != EMPTY {
		t.Errorf("e8 should be empty, got %q", game.Board.grid[0][4])
	}
	if game.Board.grid[0][7] != EMPTY {
		t.Errorf("h8 should be empty after o-o, got %q", game.Board.grid[0][7])
	}
}

func TestMakeMoveCastlingQueensideBlack(t *testing.T) {
	fen := "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R b KQkq - 0 1"
	game := ReadChessGame(fen, nil, nil)
	game.MakeMove(ParseMoveStr("e8c8"))

	if game.Board.grid[0][2] != "k" {
		t.Errorf("c8 should be 'k', got %q", game.Board.grid[0][2])
	}
	if game.Board.grid[0][3] != "r" {
		t.Errorf("d8 should be 'r', got %q", game.Board.grid[0][3])
	}
	if game.Board.grid[0][4] != EMPTY {
		t.Errorf("e8 should be empty, got %q", game.Board.grid[0][4])
	}
	if game.Board.grid[0][0] != EMPTY {
		t.Errorf("a8 should be empty after o-o-o, got %q", game.Board.grid[0][0])
	}
}

func TestMakeMovePromotion(t *testing.T) {
	// black king on h1 so e8 is clear
	fen := "8/4P3/8/8/8/8/8/4K2k w - - 0 1"
	game := ReadChessGame(fen, nil, nil)
	game.MakeMove(Move{SR: 1, SF: 4, TR: 0, TF: 4, Promotion: "Q"})

	if game.Board.grid[0][4] != "Q" {
		t.Errorf("e8 should be 'Q' after promotion, got %q", game.Board.grid[0][4])
	}
	if game.Board.grid[1][4] != EMPTY {
		t.Errorf("e7 should be empty after promotion, got %q", game.Board.grid[1][4])
	}
}

func TestMakeMoveCastleRightsKingMove(t *testing.T) {
	fen := "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1"
	game := ReadChessGame(fen, nil, nil)
	game.MakeMove(Move{SR: 7, SF: 4, TR: 7, TF: 3}) // Ke1d1

	if strings.ContainsAny(game.Castle, "KQ") {
		t.Errorf("white castling rights should be gone after king move; castle = %q", game.Castle)
	}
	if !strings.ContainsAny(game.Castle, "kq") {
		t.Errorf("black castling rights should be untouched; castle = %q", game.Castle)
	}
}

func TestMakeMoveCastleRightsRookKingside(t *testing.T) {
	fen := "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1"
	game := ReadChessGame(fen, nil, nil)
	game.MakeMove(Move{SR: 7, SF: 7, TR: 6, TF: 7}) // Rh1h2

	if strings.Contains(game.Castle, "K") {
		t.Errorf("white kingside castling should be removed; castle = %q", game.Castle)
	}
	if !strings.Contains(game.Castle, "Q") {
		t.Errorf("white queenside castling should remain; castle = %q", game.Castle)
	}
}

func TestMakeMoveCastleRightsRookQueenside(t *testing.T) {
	fen := "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1"
	game := ReadChessGame(fen, nil, nil)
	game.MakeMove(Move{SR: 7, SF: 0, TR: 6, TF: 0}) // Ra1a2

	if strings.Contains(game.Castle, "Q") {
		t.Errorf("white queenside castling should be removed; castle = %q", game.Castle)
	}
	if !strings.Contains(game.Castle, "K") {
		t.Errorf("white kingside castling should remain; castle = %q", game.Castle)
	}
}

// --- IsAttacked ---

func TestIsAttackedEmptySquare(t *testing.T) {
	game := ReadChessGame("4k3/8/8/8/8/8/8/4K3 w - - 0 1", nil, nil)
	if game.IsAttacked(4, 3) { // d4 is empty
		t.Error("empty square should not be reported as attacked")
	}
}

func TestIsAttackedByPawn(t *testing.T) {
	// black pawn on d5 attacks white king on e4 diagonally
	game := ReadChessGame("4k3/8/8/3p4/4K3/8/8/8 w - - 0 1", nil, nil)
	if !game.IsAttacked(4, 4) {
		t.Error("white king on e4 should be attacked by black pawn on d5")
	}

	// pawn directly in front (e5) does not attack e4 — only diagonals
	game2 := ReadChessGame("4k3/8/8/4p3/4K3/8/8/8 w - - 0 1", nil, nil)
	if game2.IsAttacked(4, 4) {
		t.Error("white king on e4 should NOT be attacked by pawn directly in front on e5")
	}
}

func TestIsAttackedByKnight(t *testing.T) {
	// black knight f5 reaches e3 via L-shape
	game := ReadChessGame("4k3/8/8/5n2/8/4K3/8/8 w - - 0 1", nil, nil)
	if !game.IsAttacked(5, 4) {
		t.Error("white king on e3 should be attacked by black knight on f5")
	}

	game2 := ReadChessGame("4k3/8/8/8/8/4K3/8/7n w - - 0 1", nil, nil)
	if game2.IsAttacked(5, 4) {
		t.Error("white king on e3 should NOT be attacked by knight on h1")
	}
}

func TestIsAttackedByBishop(t *testing.T) {
	// black bishop h7, white king e4 — clear diagonal
	game := ReadChessGame("4k3/7b/8/8/4K3/8/8/8 w - - 0 1", nil, nil)
	if !game.IsAttacked(4, 4) {
		t.Error("white king on e4 should be attacked by black bishop on h7")
	}
}

func TestIsAttackedByRook(t *testing.T) {
	// black rook e8, white king e4 — clear file
	game := ReadChessGame("4r2k/8/8/8/4K3/8/8/8 w - - 0 1", nil, nil)
	if !game.IsAttacked(4, 4) {
		t.Error("white king on e4 should be attacked by black rook on e8")
	}
}

func TestIsAttackedByQueen(t *testing.T) {
	// black queen a4, white king e4 — same rank
	game := ReadChessGame("4k3/8/8/8/q3K3/8/8/8 w - - 0 1", nil, nil)
	if !game.IsAttacked(4, 4) {
		t.Error("white king on e4 should be attacked by black queen on a4")
	}
}

func TestIsAttackedBlockedByPiece(t *testing.T) {
	// white pawn on e4 blocks the black rook on e8 from reaching white king on e3
	game := ReadChessGame("4r2k/8/8/8/4P3/4K3/8/8 w - - 0 1", nil, nil)
	if game.IsAttacked(5, 4) {
		t.Error("white king on e3 should NOT be attacked: pawn on e4 blocks rook on e8")
	}
}

func TestIsAttackedByKing(t *testing.T) {
	// white rook on g2, black king on g3 — king attacks adjacent piece
	game := ReadChessGame("4k3/8/8/8/8/6k1/6R1/8 b - - 0 1", nil, nil)
	if !game.IsAttacked(6, 6) {
		t.Error("white rook on g2 should be attacked by adjacent black king on g3")
	}
	// empty square next to the king — should still be false (no piece there)
	if game.IsAttacked(6, 7) {
		t.Error("empty h2 should not be reported as attacked")
	}

	// with white king g1 and black king g3, f2/g2/h2 are all adjacent to black king
	// so white has only f1 and h1 as safe squares
	game2 := ReadChessGame("8/8/8/8/8/6k1/8/6K1 w - - 0 1", nil, nil)
	if n := len(game2.LegalMoves()); n != 2 {
		t.Errorf("white king should have 2 legal moves (f1, h1), got %d", n)
	}
}

// --- PositionMoves ---

func TestPositionMovesStartCount(t *testing.T) {
	game := ReadChessGame(StartFEN, nil, nil)
	// 16 pawn moves (8 single + 8 double) + 4 knight moves
	if n := len(game.PositionMoves()); n != 20 {
		t.Errorf("PositionMoves from start = %d, want 20", n)
	}
}

func TestPositionMovesKnightCenter(t *testing.T) {
	game := ReadChessGame("8/8/8/8/3N4/8/8/4K2k w - - 0 1", nil, nil)
	var moves []Move
	for _, m := range game.PositionMoves() {
		if m.SR == 4 && m.SF == 3 { // d4
			moves = append(moves, m)
		}
	}
	if len(moves) != 8 {
		t.Errorf("knight on d4 should have 8 moves, got %d", len(moves))
	}
}

func TestPositionMovesKnightEdge(t *testing.T) {
	game := ReadChessGame("4k3/8/8/8/8/8/8/N3K3 w - - 0 1", nil, nil)
	var moves []Move
	for _, m := range game.PositionMoves() {
		if m.SR == 7 && m.SF == 0 { // a1
			moves = append(moves, m)
		}
	}
	if len(moves) != 2 {
		t.Errorf("knight on a1 should have 2 moves, got %d", len(moves))
	}
}

func TestPositionMovesRookOpen(t *testing.T) {
	game := ReadChessGame("4k3/8/8/8/3R4/8/8/4K3 w - - 0 1", nil, nil)
	var moves []Move
	for _, m := range game.PositionMoves() {
		if m.SR == 4 && m.SF == 3 { // d4
			moves = append(moves, m)
		}
	}
	// 3 left + 4 right + 4 up + 3 down = 14
	if len(moves) != 14 {
		t.Errorf("rook on d4 open board should have 14 moves, got %d", len(moves))
	}
}

func TestPositionMovesBishopOpen(t *testing.T) {
	game := ReadChessGame("4k3/8/8/8/3B4/8/8/4K3 w - - 0 1", nil, nil)
	var moves []Move
	for _, m := range game.PositionMoves() {
		if m.SR == 4 && m.SF == 3 { // d4
			moves = append(moves, m)
		}
	}
	// NE 4, NW 3, SE 3, SW 4 (board edges) = 13
	if len(moves) != 13 {
		t.Errorf("bishop on d4 open board should have 13 moves, got %d", len(moves))
	}
}

func TestPositionMovesQueenOpen(t *testing.T) {
	// queen = bishop (13) + rook (14) = 27
	game := ReadChessGame("4k3/8/8/8/3Q4/8/8/4K3 w - - 0 1", nil, nil)
	var moves []Move
	for _, m := range game.PositionMoves() {
		if m.SR == 4 && m.SF == 3 { // d4
			moves = append(moves, m)
		}
	}
	if len(moves) != 27 {
		t.Errorf("queen on d4 open board should have 27 moves, got %d", len(moves))
	}
}

func TestPositionMovesPawnPromotionSquares(t *testing.T) {
	// e8 is clear — one advance to promotion rank spawns 4 moves, one per piece
	game := ReadChessGame("8/4P3/8/8/8/8/8/4K2k w - - 0 1", nil, nil)
	var moves []Move
	for _, m := range game.PositionMoves() {
		if m.SR == 1 && m.SF == 4 { // e7
			moves = append(moves, m)
		}
	}
	if len(moves) != 4 {
		t.Errorf("pawn on e7 should generate 4 promotion moves, got %d", len(moves))
	}
	promos := map[string]bool{}
	for _, m := range moves {
		promos[m.Promotion] = true
	}
	for _, p := range []string{"Q", "R", "B", "N"} {
		if !promos[p] {
			t.Errorf("expected promotion to %q in moves", p)
		}
	}
}

// --- LegalMoves ---

func TestLegalMovesStartCount(t *testing.T) {
	game := ReadChessGame(StartFEN, nil, nil)
	if n := len(game.LegalMoves()); n != 20 {
		t.Errorf("LegalMoves from start = %d, want 20", n)
	}
}

func TestLegalMovesPinnedPiece(t *testing.T) {
	// white rook e2 is pinned to white king e1 by black rook e5 — can only move along the e-file
	fen := "4k3/8/8/4r3/8/8/4R3/4K3 w - - 0 1"
	game := ReadChessGame(fen, nil, nil)
	for _, m := range game.LegalMoves() {
		if m.SR == 6 && m.SF == 4 && m.TF != 4 { // rook on e2 moving off e-file
			t.Errorf("pinned rook moved off e-file to (%d,%d)", m.TR, m.TF)
		}
	}
}

func TestLegalMovesInCheck(t *testing.T) {
	// king e1 in check from rook e2 — can go to d1, f1, or capture
	fen := "4k3/8/8/8/8/8/4r3/4K3 w - - 0 1"
	game := ReadChessGame(fen, nil, nil)
	if n := len(game.LegalMoves()); n != 3 {
		t.Errorf("king in check should have 3 legal moves (d1, f1, xe2), got %d", n)
	}
}

func TestLegalMovesCheckmatePosition(t *testing.T) {
	// fool's mate final position — white is checkmated
	fen := "rnb1kbnr/pppp1ppp/8/4p3/6Pq/5P2/PPPPP2P/RNBQKBNR w KQkq - 1 3"
	game := ReadChessGame(fen, nil, nil)
	if n := len(game.LegalMoves()); n != 0 {
		t.Errorf("checkmate position should have 0 legal moves, got %d", n)
	}
}

func TestLegalMovesStalemate(t *testing.T) {
	// black king a8, white queen b6 and king h1 — black has no moves but is not in check
	fen := "k7/8/1Q6/8/8/8/8/7K b - - 0 1"
	game := ReadChessGame(fen, nil, nil)

	kr, kf := game.KingCoords(false)
	if game.IsAttacked(kr, kf) {
		t.Skip("position is check, not stalemate — adjust FEN")
	}
	if n := len(game.LegalMoves()); n != 0 {
		t.Errorf("stalemate position should have 0 legal moves, got %d", n)
	}
}

func TestLegalMovesEnPassantIncluded(t *testing.T) {
	fen := "rnbqkbnr/ppp1pppp/8/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 3"
	game := ReadChessGame(fen, nil, nil)
	if !containsMove(game.LegalMoves(), "e5xd6") {
		t.Error("e5xd6 (en passant) should be in legal moves")
	}
}

func TestLegalMovesCastlingIncluded(t *testing.T) {
	fen := "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQkq - 0 1"
	game := ReadChessGame(fen, nil, nil)
	moves := game.LegalMoves()

	if !containsMove(moves, "e1g1") {
		t.Error("e1g1 (O-O) should be legal when lane is clear")
	}
	if !containsMove(moves, "e1c1") {
		t.Error("e1c1 (O-O-O) should be legal when lane is clear")
	}
}

func TestLegalMovesCastlingBlockedByCheck(t *testing.T) {
	// black rook e8 puts white king in check — can't castle out of check
	fen := "4r3/8/8/8/8/8/8/R3K2R w KQ - 0 1"
	game := ReadChessGame(fen, nil, nil)
	moves := game.LegalMoves()

	if containsMove(moves, "e1g1") {
		t.Error("O-O should not be legal while king is in check")
	}
	if containsMove(moves, "e1c1") {
		t.Error("O-O-O should not be legal while king is in check")
	}
}

func TestLegalMovesCastlingThroughAttackedSquare(t *testing.T) {
	// black rook f8 attacks f1 — kingside castling passes through f1, must be illegal
	fen := "5r1k/8/8/8/8/8/8/R3K2R w KQ - 0 1"
	game := ReadChessGame(fen, nil, nil)
	moves := game.LegalMoves()

	if containsMove(moves, "e1g1") {
		t.Error("O-O should not be legal when f1 is attacked")
	}
	if !containsMove(moves, "e1c1") {
		t.Error("O-O-O should still be legal (queenside lane is clean)")
	}
}

// --- Move() ---

func TestMoveValidMove(t *testing.T) {
	game := ReadChessGame(StartFEN, nil, nil)
	game.NextMoves = game.LegalMoves()
	if err := game.Move("e2e4"); err != nil {
		t.Errorf("e2e4 should be legal from start; got error: %v", err)
	}
}

func TestMoveInvalidMove(t *testing.T) {
	game := ReadChessGame(StartFEN, nil, nil)
	game.NextMoves = game.LegalMoves()
	if err := game.Move("e2e5"); err == nil {
		t.Error("e2e5 should be illegal from start")
	}
}

func TestMoveStatusNotStartedToInProgress(t *testing.T) {
	game := ReadChessGame(StartFEN, nil, nil)
	game.NextMoves = game.LegalMoves()

	if game.Status != NotStarted {
		t.Fatalf("expected NotStarted before any moves, got %s", game.Status)
	}
	if err := game.Move("e2e4"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if game.Status != InProgress {
		t.Errorf("status after first move = %s, want InProgress", game.Status)
	}
}

func TestMovePrevMovesAppended(t *testing.T) {
	game := ReadChessGame(StartFEN, nil, nil)
	game.NextMoves = game.LegalMoves()
	_ = game.Move("e2e4")

	if len(game.PrevMoves) != 1 {
		t.Errorf("PrevMoves length = %d, want 1", len(game.PrevMoves))
	}
	if game.PrevMoves[0].WriteMoveStr() != "e2e4" {
		t.Errorf("PrevMoves[0] = %q, want \"e2e4\"", game.PrevMoves[0].WriteMoveStr())
	}
}

func TestMoveGameOverError(t *testing.T) {
	game := ReadChessGame(StartFEN, nil, nil)
	game.Status = WinWhite
	if err := game.Move("e2e4"); err == nil {
		t.Error("Move on a finished game should return an error")
	}
}

// plays the four-move fool's mate: 1. f3 e5 2. g4 Qh4#
func TestMoveFoolsMate(t *testing.T) {
	game := ReadChessGame(StartFEN, nil, nil)
	game.NextMoves = game.LegalMoves()

	for _, mv := range []string{"f2f3", "e7e5", "g2g4", "d8h4"} {
		if err := game.Move(mv); err != nil {
			t.Fatalf("unexpected error on %q: %v", mv, err)
		}
	}
	if game.Status != WinBlack {
		t.Errorf("after fool's mate: status = %s, want WinBlack", game.Status)
	}
}

func TestMove50MoveRule(t *testing.T) {
	// halfmove clock at 99 — one quiet rook move tips it to 100
	fen := "k6K/8/8/8/8/8/8/7R w - - 99 100"
	game := ReadChessGame(fen, nil, nil)
	game.NextMoves = game.LegalMoves()

	if err := game.Move("h1h2"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if game.Status != Draw {
		t.Errorf("50-move rule: status = %s, want Draw", game.Status)
	}
}

func TestMoveStalemateDetected(t *testing.T) {
	// after Qg2-c2, black king on a1 has no legal moves and is not in check
	fen := "8/8/8/8/8/8/6Q1/k1K5 w - - 0 1"
	game := ReadChessGame(fen, nil, nil)
	game.NextMoves = game.LegalMoves()

	if err := game.Move("g2c2"); err != nil {
		t.Fatalf("g2c2 should be legal: %v", err)
	}
	if game.Status != Draw {
		t.Errorf("stalemate should result in Draw, got %s", game.Status)
	}
}
