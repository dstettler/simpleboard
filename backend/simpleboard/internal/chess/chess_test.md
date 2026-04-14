
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
go test -v -run TestMoveFoolsMate ./internal/chess/

# run all tests for a specific function area
go test -v -run TestMakeMove ./internal/chess/
go test -v -run TestIsAttacked ./internal/chess/
```


### Status

| Test              | What it checks
|-------------------|----------------
| TestStatusString  | All 5 Status enum values + out-of-range returns "Invalid"


### ParseCoords / ParseAlg

| Test                        | What it checks
|-----------------------------|----------------
| TestParseCoords             | a1→(7,0), h8→(0,7), e4→(4,4), and others
| TestParseAlg                | (7,0)→"a1", (0,7)→"h8", (4,4)→"e4", and others
| TestParseAlgCoordsRoundTrip | ParseAlg(ParseCoords(sq)) == sq for 8 squares


### Board.FEN

| Test                  | What it checks
|-----------------------|----------------
| TestBoardFENStartPosition | Board.FEN() matches start position board string
| TestBoardFENAfterE4       | Board.FEN() matches board string after 1.e4


### Move struct

| Test                    | What it checks
|-------------------------|----------------
| TestMoveIsEqual         | Equal moves return true; 7 single-field differences each return false
| TestMoveCopy            | Copy matches original; mutating copy does not affect original
| TestParseMoveStr        | Table-driven: regular move, capture, all 4 castling patterns, promotion, capture+promotion
| TestWriteMoveStrRoundTrip | WriteMoveStr(ParseMoveStr(s)) == s for 8 move strings


### ReadChessGame

| Test                     | What it checks
|--------------------------|----------------
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


### ChessGame.FEN / ChessGame.Copy

| Test                       | What it checks
|----------------------------|----------------
| TestChessGameFENRoundTrip  | ReadChessGame(fen).FEN() == fen for 4 positions
| TestChessGameCopyIndependence | Mutating all fields on a copy does not affect the original


### KingCoords

| Test                       | What it checks
|----------------------------|----------------
| TestKingCoordsStartPosition | Both kings found at correct coordinates on start board
| TestKingCoordsMissingKing   | Returns (-1,-1) when king is absent from board
| TestKingCoordsCustomPosition | Kings-only board with non-standard placement


### MakeMove

| Test                               | What it checks
|------------------------------------|----------------
| TestMakeMoveUpdatesBoard           | Piece appears at target, source becomes empty, side switches
| TestMakeMoveCapture                | Captured piece replaced, halfmove clock resets to 1
| TestMakeMoveHalfmoveClock          | Knight move increments; pawn move resets; second piece move increments again
| TestMakeMoveFullmoveCounter        | Stays at 1 after white moves; increments to 2 after black moves
| TestMakeMoveSideToggle             | Side alternates w→b→w across two moves
| TestMakeMoveEPTSSetOnDoublePush    | EPTS = "e3" after white e2e4
| TestMakeMoveEPTSClearedAfterNonPawn | EPTS resets to "-" after a knight move
| TestMakeMoveEnPassantWhite         | d5 pawn removed, white pawn appears on d6 after e5xd6
| TestMakeMoveEnPassantBlack         | d4 pawn removed, black pawn appears on d3 after e4xd3
| TestMakeMoveCastlingKingsideWhite  | King on g1, rook on f1, e1 and h1 empty after O-O
| TestMakeMoveCastlingQueensideWhite | King on c1, rook on d1, e1 and a1 empty after O-O-O
| TestMakeMoveCastlingKingsideBlack  | King on g8, rook on f8, e8 and h8 empty after o-o
| TestMakeMoveCastlingQueensideBlack | King on c8, rook on d8, e8 and a8 empty after o-o-o
| TestMakeMovePromotion              | Pawn on e7 becomes queen on e8; e7 is empty
| TestMakeMoveCastleRightsKingMove   | Both white rights removed after king moves; black rights unchanged
| TestMakeMoveCastleRightsRookKingside   | "K" removed after h1 rook moves; "Q" unchanged
| TestMakeMoveCastleRightsRookQueenside  | "Q" removed after a1 rook moves; "K" unchanged


### IsAttacked

| Test                        | What it checks
|-----------------------------|----------------
| TestIsAttackedEmptySquare   | Returns false immediately for an empty square
| TestIsAttackedByPawn        | Black pawn on d5 attacks white king on e4; pawn directly in front does not
| TestIsAttackedByKnight      | Black knight on f5 attacks white king on e3; distant knight does not
| TestIsAttackedByBishop      | Black bishop on h7 attacks white king on e4 along clear diagonal
| TestIsAttackedByRook        | Black rook on e8 attacks white king on e4 along clear file
| TestIsAttackedByQueen       | Black queen on a4 attacks white king on e4 along rank
| TestIsAttackedBlockedByPiece | Rook attack on e3 blocked by own pawn on e4
| TestIsAttackedByKing        | Adjacent black king attacks white rook; LegalMoves excludes king-adjacent squares


### PositionMoves

| Test                             | What it checks
|----------------------------------|----------------
| TestPositionMovesStartCount      | 20 pseudo-legal moves (16 pawn + 4 knight) from start
| TestPositionMovesKnightCenter    | Knight on d4 generates exactly 8 moves
| TestPositionMovesKnightEdge      | Knight on a1 generates exactly 2 moves
| TestPositionMovesRookOpen        | Rook on d4 open board generates exactly 14 moves
| TestPositionMovesBishopOpen      | Bishop on d4 open board generates exactly 13 moves
| TestPositionMovesQueenOpen       | Queen on d4 open board generates exactly 27 moves
| TestPositionMovesPawnPromotionSquares | Pawn on e7 generates exactly 4 moves, one per promotion piece


### LegalMoves

| Test                                  | What it checks
|---------------------------------------|----------------
| TestLegalMovesStartCount              | 20 legal moves from start position
| TestLegalMovesPinnedPiece             | Rook pinned on e-file cannot move off the e-file
| TestLegalMovesInCheck                 | King in check from e2 rook has exactly 3 legal moves (d1, f1, xe2)
| TestLegalMovesCheckmatePosition       | Fool's mate final position returns 0 legal moves
| TestLegalMovesStalemate               | Stalemate position (king not in check) returns 0 legal moves
| TestLegalMovesEnPassantIncluded       | e5xd6 en passant appears in legal moves when EPTS = d6
| TestLegalMovesCastlingIncluded        | Both e1g1 and e1c1 present when lanes are clear
| TestLegalMovesCastlingBlockedByCheck  | Neither castling move legal when king is in check on e1
| TestLegalMovesCastlingThroughAttackedSquare | O-O illegal when f1 is attacked; O-O-O still legal


### Move()

| Test                            | What it checks
|---------------------------------|----------------
| TestMoveValidMove               | e2e4 from start returns no error
| TestMoveInvalidMove             | e2e5 from start returns an error
| TestMoveStatusNotStartedToInProgress | First move transitions status from NotStarted to InProgress
| TestMovePrevMovesAppended       | Played move is appended to PrevMoves
| TestMoveGameOverError           | Move on a finished game returns "Game is over."
| TestMoveFoolsMate               | Four-move fool's mate sequence ends with Status = WinBlack
| TestMove50MoveRule              | Halfmove clock at 99; one non-pawn move sets Status = Draw
| TestMoveStalemateDetected       | Stalemate-inducing queen move sets Status = Draw (not a win)
