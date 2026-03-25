package chess

import (
	"log"
	"simpleboard/internal/utils"
	"strings"
)

// Initial game state FEN
const StartFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// Vectors for move patterns
var knightVectors = [][2]int{{1, 2}, {-1, 2}, {1, -2}, {-1, -2}, {2, 1}, {-2, 1}, {2, -1}, {-2, -1}}
var bishopVectors = [][2]int{{1, 1}, {-1, 1}, {1, -1}, {-1, -1}}
var rookVectors = [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
var kingVectors = [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}}

// Generates legal moves from an array of possible move patterns
func (g *ChessGame) LegalMoves() []Move {
	plm := g.PositionMoves() // get possibly legal moves
	moves := []Move{}

	kr, kf := g.KingCoords(g.Side == "w")

	for _, m := range plm {
		copy := g.Copy()
		copy.MakeMove(m) // make the possible move

		// check king
		if !copy.IsAttacked(kr, kf) {
			moves = append(moves, m)
		}
	}

	return moves
}

// Generates all move patterns from the board state
// before checking the moving side's king
func (g *ChessGame) PositionMoves() []Move {
	moves := []Move{}
	white := bool(g.Side == "w")

	for r := 0; r < 8; r++ {
		for f := 0; f < 8; f++ {
			p := g.Board.grid[r][f] // get piece
			if p == EMPTY {
				continue
			}

			// generate moves by moving side
			// get white moves
			if white && utils.IsUpper(p) {
				moves = append(moves, g.generateMoves(p, r, f, true)...)
			}

			// get black moves
			if !white && utils.IsLower(p) {
				moves = append(moves, g.generateMoves(p, r, f, false)...)
			}
		}
	}
	return moves
}

// Generates possible moves from only patterns for a piece type (not legal)
func (g *ChessGame) generateMoves(p string, r, f int, white bool) []Move {
	switch strings.ToLower(p) {
	case "p":
		return g.generatePawnMoves(r, f, white)
	case "b":
		return g.generateBishopMoves(r, f, white)
	case "n":
		return g.generateKnightMoves(r, f, white)
	case "r":
		return g.generateRookMoves(r, f, white)
	case "q":
		return g.generateQueenMoves(r, f, white)
	case "k":
		return g.generateKingMoves(r, f, white)
	}
	return nil
}

// Generates possible moves from only patterns for pawns
func (g *ChessGame) generatePawnMoves(r, f int, white bool) []Move {
	moves := []Move{}
	b := g.Board.grid

	if white {
		// Single square
		if r > 0 && b[r-1][f] == EMPTY {
			moves = append(moves, Move{r, f, r - 1, f, false, false, ""})

			// Double square
			if r == 6 && b[5][f] == EMPTY && b[4][f] == EMPTY {
				moves = append(moves, Move{r, f, r - 2, f, false, false, ""})
			}
		}

		// Captures
		for _, dir := range []int{-1, 1} {
			fcap := f + dir
			if r > 0 && fcap >= 0 && fcap < 8 {
				if utils.IsLower(b[r-1][fcap]) {
					moves = append(moves, Move{r, f, r - 1, fcap, true, false, ""})
				}
			}
		}
	} else {
		// Single square
		if r < 7 && b[r+1][f] == EMPTY {
			moves = append(moves, Move{r, f, r + 1, f, false, false, ""})

			// Double square
			if r == 1 && b[2][f] == EMPTY && b[3][f] == EMPTY {
				moves = append(moves, Move{r, f, r + 2, f, false, false, ""})
			}
		}

		// Captures
		for _, dir := range []int{-1, 1} {
			fcap := f + dir
			if r < 7 && fcap >= 0 && fcap < 8 {
				if utils.IsLower(b[r+1][fcap]) {
					moves = append(moves, Move{r, f, r + 1, fcap, true, false, ""})
				}
			}
		}
		// TODO: add enpassant logic
		// TODO: add promotion logic
	}
	return moves
}

// Generates possible moves from only patterns for knights
func (g *ChessGame) generateKnightMoves(r, f int, white bool) []Move {
	moves := []Move{}
	b := g.Board.grid

	// define knight vectors
	vecs := knightVectors

	// iterate through the applied vector sums to the position and check bounds
	for _, v := range vecs {
		nr, nf := r+v[0], f+v[1]
		if nr < 0 || nr > 7 || nf < 0 || nf > 7 {
			continue
		} // out of bounds
		t := b[nr][nf]

		// determine same color / capture of target square piece / empty space
		if t == EMPTY {
			moves = append(moves, Move{r, f, nr, nf, false, false, ""})
			continue
		}
		if (white && utils.IsLower(t)) || (!white && utils.IsUpper(t)) {
			moves = append(moves, Move{r, f, nr, nf, true, false, ""})
		}
	}

	return moves
}

// Generic function to generate possible moves from only patterns for sliding pieces,
// i.e. patterns with direction and variable length (bishop, rook, queen)
func (g *ChessGame) generateVarLenMoves(r, f int, white bool, vecs [][2]int) []Move {
	moves := []Move{}
	b := g.Board.grid

	for _, v := range vecs {

		// unit scale applied
		s := 1
		nr, nf := int(r+(s*v[0])), int(f+(s*v[1]))

		// within this direction, apply the scalar to the vector, sum
		// increase the scalar for the next iteration
		// stop when out of bounds
		for nr >= 0 && nr <= 7 && nf >= 0 && nf <= 7 {
			t := b[nr][nf]

			// determine same color / capture of target square piece / empty space
			if t == EMPTY {
				moves = append(moves, Move{r, f, nr, nf, false, false, ""})
				s += 1
				nr, nf = int(r+(s*v[0])), int(f+(s*v[1]))
			}

			if (white && utils.IsLower(t)) || (!white && utils.IsUpper(t)) {
				moves = append(moves, Move{r, f, nr, nf, true, false, ""})
				break // capture breaks path
			} else {
				break // same color piece
			}
		}
	}

	return moves
}

// Generates possible moves from only patterns for bishops
func (g *ChessGame) generateBishopMoves(r, f int, white bool) []Move {
	// define bishop directions
	vecs := bishopVectors
	return g.generateVarLenMoves(r, f, white, vecs)
}

// Generates possible moves from only patterns for rooks
func (g *ChessGame) generateRookMoves(r, f int, white bool) []Move {
	// define rook directions
	vecs := rookVectors
	return g.generateVarLenMoves(r, f, white, vecs)
}

// Generates possible moves from only patterns for queens
func (g *ChessGame) generateQueenMoves(r, f int, white bool) []Move {
	// assume the queen is a bishop + rook
	moves := g.generateBishopMoves(r, f, white)
	moves = append(moves, g.generateRookMoves(r, f, white)...)
	return moves
}

// Generates possible moves from only patterns for the king
func (g *ChessGame) generateKingMoves(r, f int, white bool) []Move {
	moves := []Move{}
	b := g.Board.grid

	// define king vectors
	vecs := kingVectors

	// iterate through the applied vector sums to the position and check bounds
	for _, v := range vecs {
		nr, nf := r+v[0], f+v[1]
		if nr < 0 || nr > 7 || nf < 0 || nf > 7 {
			continue
		} // out of bounds
		t := b[nr][nf]

		// determine same color / capture of target square piece / empty space
		if t == EMPTY {
			moves = append(moves, Move{r, f, nr, nf, false, false, ""})
			continue
		}
		if (white && utils.IsLower(t)) || (!white && utils.IsUpper(t)) {
			moves = append(moves, Move{r, f, nr, nf, true, false, ""})
		}
	}
	// TODO: add castling logic? or in LegalMoves()

	return moves
}

// Gets the position coords of the king; simple linear search
func (g *ChessGame) KingCoords(white bool) (int, int) {

	b := g.Board.grid

	p := "K"
	if !white {
		p = "k"
	}

	for r := 0; r < 8; r++ {
		for f := 0; f < 8; f++ {
			if b[r][f] == p { return r, f }
		}
	}
	return -1, -1
}

// Makes a verified move to the ChessGame, updating the board
// and fields.
func (g *ChessGame) MakeMove(m Move) {
	b := &g.Board.grid
	p := b[m.SR][m.SF]

	// TODO: En Passant
	// TODO: Castling

	b[m.TR][m.TF] = p
	b[m.SR][m.SF] = EMPTY

	// TODO: Promotion

	// reset halfmove clock
	if m.Capture || strings.ToLower(p) == "p" {
		g.HalfmoveClock = 0
	} else {
		g.HalfmoveClock++
	}

	// switch side to move; increment fullmove counter
	if g.Side == "w" {
		g.Side = "b"
	} else {
		g.Side = "w"
		g.FullmoveCount++
	}
}

// function that checks if a particular piece is attacked
func (g *ChessGame) IsAttacked(r, f int) bool {

	if r < 0 || r > 7 || f < 0 || f > 7 {
		log.Fatalf("Invalid coordinates to check attacked: r=%d, f=%d", r, f)
	}

	b := g.Board.grid

	p := b[r][f]
	if p == EMPTY {
		return false
	}

	// define attacker color
	white := !utils.IsUpper(p)

	// check pawns
	if !white {
		if r > 0 && f > 0 && b[r-1][f-1] == "p" {
			return true
		} // left attack
		if r > 0 && f < 7 && b[r-1][f+1] == "p" {
			return true
		} // right
	} else {
		if r < 7 && f > 0 && b[r+1][f-1] == "P" {
			return true
		} // left
		if r < 7 && f < 7 && b[r+1][f+1] == "P" {
			return true
		} // right
	}

	// check knights
	for _, v := range knightVectors {
		nr, nf := r+v[0], f+v[1]
		if nr < 0 || nr > 7 || nf < 0 || nf > 7 {
			continue
		} // out of bounds

		if white && b[nr][nf] == "N" {
			return true
		}
		if !white && b[nr][nf] == "n" {
			return true
		}
	}

	// check variable length attacks (bishops, rooks, and queens)
	vecs := rookVectors
	vecs = append(vecs, bishopVectors...)

	for _, v := range vecs {

		s := 1
		nr, nf := int(r+(s*v[0])), int(f+(s*v[1]))

		for nr >= 0 && nr <= 7 && nf >= 0 && nf <= 7 {
			a := b[nr][nf]

			if a != EMPTY {

				// same color piece check
				if (!white && utils.IsUpper(a)) || (white && utils.IsLower(a)) {
					break
				}

				if white {
					if a == "Q" {
						return true
					} // queen
					if a == "B" && (v[0] != 0 && v[1] != 0) {
						return true
					} // diagonal path
					if a == "R" && (v[0] == 0 || v[1] == 0) {
						return true
					} // straight path
				} else {
					if a == "q" {
						return true
					} // queen
					if a == "b" && (v[0] != 0 && v[1] != 0) {
						return true
					} // diagonal path
					if a == "r" && (v[0] == 0 || v[1] == 0) {
						return true
					} // straight path
				}
			}

			s += 1
			nr, nf = int(r+(s*v[0])), int(f+(s*v[1]))
		}
	}

	return false
}
