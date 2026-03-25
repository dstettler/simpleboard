
# Running Tests
From the module root (`backend/simpleboard/`):

```bash
# run all tests (not recommended rn)
go test ./...

# run only the chess package tests
go test ./internal/chess/

# verbose output
go test -v ./internal/chess/

# run a specific test 
go test -v -run TestReadStartingPosition ./internal/chess/

# run the table-driven test suite
go test -v -run TestReadChessGameTable ./internal/chess/
```


| Test                     | What it checks 
|--------------------------|----------------
| TestStatusString         | All 5 Status enum values + out-of-range returns "Invalid"
| TestReadStartingPosition | Full board layout, all metadata fields, NotStarted status
| TestReadAfterE4          | Pawn on e4, e2 empty, side is black, en passant set
| TestReadCastledPosition  | King on g1, rook on f1, only black can castle
| TestReadNoCastlingRights | Castle field is "-", halfmove and fullmove parsed
| TestReadEndgamePosition  | 3 pieces on board, 61 empty squares, piece placement
| TestBoardAllEmpty        | Kings-only board, exactly 2 pieces counted
| TestBoardFullRank        | Back rows full, middle rows empty in start position

| **TestReadChessGameTable** | **Table-driven subtests (4):** 
| — start position           | Starting FEN, NotStarted status 
| — after 1.e4               | First move, en passant target e3 
| — king and pawn endgame    | Sparse board, no castling, fullmove 45
| — white kingside castling only | Castle = "Kkq", fullmove 9
