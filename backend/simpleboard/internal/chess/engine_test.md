
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
go test -v -run TestMakeMoveCapture ./internal/chess/

# run all tests for a specific function
go test -v -run TestMakeMove ./internal/chess/
```

### KingCoords

| Test                         | What it checks
|------------------------------|----------------
| TestKingCoordsStartPosition  | Both kings on the starting board
| TestKingCoordsMissingKing    | Returns (-1,-1) when king is missing

### MakeMove

| Test                         | What it checks
|------------------------------|----------------
| TestMakeMoveUpdatesBoard     | Piece is moved and board is updated
| TestMakeMoveHalfmoveClock    | Halfmove clock resets on pawn move
| TestMakeMoveFullmoveCounter  | Fullmove counter increments after black moves
| TestMakeMoveCapture          | Capture replaces the target piece
| TestMakeMoveSideToggle       | Side toggles back to white after both sides move

### PositionMoves

| Test                         | What it checks
|------------------------------|----------------
| TestPositionMovesStartCount  | 20 pseudo-legal moves generated from start

### LegalMoves

| Test                         | What it checks
|------------------------------|----------------
| TestLegalMovesStartCount     | LegalMoves from start also returns 20
