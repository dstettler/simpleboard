package chess

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

const EMPTY = "."

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

// Type representing game board in a `[In rank][In file]` index schema
type Board struct {
	grid [8][8]string // 8x8
}

// Chess algebra '[In rank][In file]' to coords converter from ascii
func ParseCoords(s string) (int, int) {
	if len(s) != 2 {
		log.Fatalf("Invalid coordinates length: %d", len(s))
	}
	file := int(s[0] - 'a')            // parse first token as file
	rank := (int(s[1]-'0') - 8) * (-1) // parse second as rank

	return rank, file
}

// Coords to Chess algebra converter from ints
func ParseAlg(r, f int) string {
	files := []string{"a", "b", "c", "d", "e", "f", "g", "h"}

	if r < 0 || r > 7 || f < 0 || f > 7 {
		log.Fatalf("Invalid coordinates: %d, %d", r, f)
	}

	return files[f] + strconv.Itoa((r-8)*-1)
}

// Board to FEN string
func (b Board) FEN() string {
	fen := ""
	for r := 0; r < 8; r++ {
		empty := 0
		for f := 0; f < 8; f++ {
			sq := b.grid[r][f]
			if sq == EMPTY {
				empty++
			} else {
				if empty > 0 {
					fen += fmt.Sprintf("%d", empty) // format the empty space number
					empty = 0
				}
				fen += sq
			}
		}
		if empty > 0 {
			fen += fmt.Sprintf("%d", empty)
		}
		if r < 7 {
			fen += "/" // add rank delim
		}
	}
	return fen
}

// Board print
func (b Board) Print() {
	for _, r := range b.grid {
		fmt.Println(strings.Join(r[:], " "))
	}
}

// Go type for a single move on the board
type Move struct {
	SR, SF    int
	TR, TF    int
	Capture   bool
	Castling  bool
	Promotion string
}

// Parses a move string and returns a populated Move
func ParseMoveStr(s string) Move {
	length := len(s)
	if length < 4 || length > 6 {
		log.Fatalf("Invalid move string with length %d - '%s'", len(s), s)
	}

	capture := false
	castling := false
	promo := ""
	sr, sf := ParseCoords(s[0:2])
	tr, tf := -1, -1

	// only indicates the possibility of a castling move (matches the pattern).
	// If applied, must be able to castle
	if s == "e1g1" || s == "e1c1" || s == "e8g8" || s == "e8c8" { castling = true }

	if s[2:3] == "x" {
		capture = true
		tr, tf = ParseCoords(s[3:5])
	} else {
		tr, tf = ParseCoords(s[2:4])
	}

	if (length <= 4 && !capture) || (length > 5 && capture) {
		promo = s[length-1:]
	}

	return Move {
		SR: sr,
		SF: sf,
		TR: tr,
		TF: tf,
		Capture: capture,
		Castling: castling,
		Promotion: promo,
	}
}

// Writes move to a long-form notation string
func (m Move) WriteMoveStr() string {
	moveStr := ""
	s := ParseAlg(m.SR, m.SF)
	t := ParseAlg(m.TR, m.TF)
	moveStr += s
	if m.Capture {
		moveStr += "x"
	}
	moveStr += t

	return moveStr
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

// Create a deep copy of the ChessGame
func (orig ChessGame) Copy() ChessGame {
	cg := ChessGame{
		Board:         orig.Board,
		Side:          orig.Side,
		Castle:        orig.Castle,
		EPTS:          orig.EPTS,
		HalfmoveClock: orig.HalfmoveClock,
		FullmoveCount: orig.FullmoveCount,
		Status:        orig.Status,
	}
	return cg
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

// ChessGame to FEN string
// Returns the assembled FEN string
func (game ChessGame) FEN() string {
	fen := game.Board.FEN() // get FEN component for board
	components := []string{
		fen,
		game.Side,
		game.Castle,
		game.EPTS,
		fmt.Sprintf("%d", game.HalfmoveClock),
		fmt.Sprintf("%d", game.FullmoveCount),
	}
	return strings.Join(components, " ")
}

// Status to string function
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
					b.grid[rankindex][fileindex] = EMPTY
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
