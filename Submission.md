# Sprint 3
### Group : SimpleBoard
- TJ Schultz (Backend)
- Devon Stettler (Frontend)
- Sreeram Gangavarapu (Frontend)
- Arunabho Basu (Backend)

## Work Completed
### Frontend
- Added FEN string validation and unit tests
- Added Cypress e2e tests
- Added registration form
- Added homepage
- Fixed routing issues
### Backend
- Final chess engine w/ full rule support
- JWT authentication
  - Login authentication
  - Game endpoint authentication
- CORS fixes

## Tests (Frontend)
### Unit Tests
- Dashboard component instantiates successfully
- Game component instantiates successfully
- App component instantiates successfully
- Login component instantiates successfully
- Piece component instantiates successfully
- Board component instantiates successfully
- BoardLoadService instantiates successfully
- FEN string validation functions as expected
### Cypress Tests
- `cypress/e2e/home.cy.ts`  
  - Verifies the home page loads successfully
  - Verifies navbar links are visible
  - Verifies clicking the Login link from the navbar routes to the login page

- `cypress/e2e/login.cy.ts`  
  - Verifies the login page loads successfully
  - Verifies validation errors appear when login fields are touched and left invalid
  - Verifies the page can switch from login mode to register mode
  - Verifies registration form submission shows the success message

## Tests (Backend)
```
=== RUN   TestStatusString
--- PASS: TestStatusString (0.00s)
=== RUN   TestParseCoords
--- PASS: TestParseCoords (0.00s)
=== RUN   TestParseAlg
--- PASS: TestParseAlg (0.00s)
=== RUN   TestParseAlgCoordsRoundTrip
--- PASS: TestParseAlgCoordsRoundTrip (0.00s)
=== RUN   TestBoardFENStartPosition
--- PASS: TestBoardFENStartPosition (0.00s)
=== RUN   TestBoardFENAfterE4
--- PASS: TestBoardFENAfterE4 (0.00s)
=== RUN   TestMoveIsEqual
--- PASS: TestMoveIsEqual (0.00s)
=== RUN   TestMoveCopy
--- PASS: TestMoveCopy (0.00s)
=== RUN   TestParseMoveStr
--- PASS: TestParseMoveStr (0.00s)
=== RUN   TestWriteMoveStrRoundTrip
--- PASS: TestWriteMoveStrRoundTrip (0.00s)
=== RUN   TestReadStartingPosition
--- PASS: TestReadStartingPosition (0.00s)
=== RUN   TestReadAfterE4
--- PASS: TestReadAfterE4 (0.00s)
=== RUN   TestReadCastledPosition
--- PASS: TestReadCastledPosition (0.00s)
=== RUN   TestReadNoCastlingRights
--- PASS: TestReadNoCastlingRights (0.00s)
=== RUN   TestReadEndgamePosition
--- PASS: TestReadEndgamePosition (0.00s)
=== RUN   TestBoardAllEmpty
--- PASS: TestBoardAllEmpty (0.00s)
=== RUN   TestBoardFullRank
--- PASS: TestBoardFullRank (0.00s)
=== RUN   TestReadChessGameTable
=== RUN   TestReadChessGameTable/start_position
=== RUN   TestReadChessGameTable/after_1.e4
=== RUN   TestReadChessGameTable/king_and_pawn_endgame
=== RUN   TestReadChessGameTable/white_kingside_castling_only
--- PASS: TestReadChessGameTable (0.00s)
    --- PASS: TestReadChessGameTable/start_position (0.00s)
    --- PASS: TestReadChessGameTable/after_1.e4 (0.00s)
    --- PASS: TestReadChessGameTable/king_and_pawn_endgame (0.00s)
    --- PASS: TestReadChessGameTable/white_kingside_castling_only (0.00s)
=== RUN   TestChessGameFENRoundTrip
--- PASS: TestChessGameFENRoundTrip (0.00s)
=== RUN   TestChessGameCopyIndependence
--- PASS: TestChessGameCopyIndependence (0.00s)
=== RUN   TestKingCoordsStartPosition
--- PASS: TestKingCoordsStartPosition (0.00s)
=== RUN   TestKingCoordsMissingKing
--- PASS: TestKingCoordsMissingKing (0.00s)
=== RUN   TestKingCoordsCustomPosition
--- PASS: TestKingCoordsCustomPosition (0.00s)
=== RUN   TestMakeMoveUpdatesBoard
--- PASS: TestMakeMoveUpdatesBoard (0.00s)
=== RUN   TestMakeMoveCapture
--- PASS: TestMakeMoveCapture (0.00s)
=== RUN   TestMakeMoveHalfmoveClock
--- PASS: TestMakeMoveHalfmoveClock (0.00s)
=== RUN   TestMakeMoveFullmoveCounter
--- PASS: TestMakeMoveFullmoveCounter (0.00s)
=== RUN   TestMakeMoveSideToggle
--- PASS: TestMakeMoveSideToggle (0.00s)
=== RUN   TestMakeMoveEPTSSetOnDoublePush
--- PASS: TestMakeMoveEPTSSetOnDoublePush (0.00s)
=== RUN   TestMakeMoveEPTSClearedAfterNonPawn
--- PASS: TestMakeMoveEPTSClearedAfterNonPawn (0.00s)
=== RUN   TestMakeMoveEnPassantWhite
--- PASS: TestMakeMoveEnPassantWhite (0.00s)
=== RUN   TestMakeMoveEnPassantBlack
--- PASS: TestMakeMoveEnPassantBlack (0.00s)
=== RUN   TestMakeMoveCastlingKingsideWhite
--- PASS: TestMakeMoveCastlingKingsideWhite (0.00s)
=== RUN   TestMakeMoveCastlingQueensideWhite
--- PASS: TestMakeMoveCastlingQueensideWhite (0.00s)
=== RUN   TestMakeMoveCastlingKingsideBlack
--- PASS: TestMakeMoveCastlingKingsideBlack (0.00s)
=== RUN   TestMakeMoveCastlingQueensideBlack
--- PASS: TestMakeMoveCastlingQueensideBlack (0.00s)
=== RUN   TestMakeMovePromotion
--- PASS: TestMakeMovePromotion (0.00s)
=== RUN   TestMakeMoveCastleRightsKingMove
--- PASS: TestMakeMoveCastleRightsKingMove (0.00s)
=== RUN   TestMakeMoveCastleRightsRookKingside
--- PASS: TestMakeMoveCastleRightsRookKingside (0.00s)
=== RUN   TestMakeMoveCastleRightsRookQueenside
--- PASS: TestMakeMoveCastleRightsRookQueenside (0.00s)
=== RUN   TestIsAttackedEmptySquare
--- PASS: TestIsAttackedEmptySquare (0.00s)
=== RUN   TestIsAttackedByPawn
--- PASS: TestIsAttackedByPawn (0.00s)
=== RUN   TestIsAttackedByKnight
--- PASS: TestIsAttackedByKnight (0.00s)
=== RUN   TestIsAttackedByBishop
--- PASS: TestIsAttackedByBishop (0.00s)
=== RUN   TestIsAttackedByRook
--- PASS: TestIsAttackedByRook (0.00s)
=== RUN   TestIsAttackedByQueen
--- PASS: TestIsAttackedByQueen (0.00s)
=== RUN   TestIsAttackedBlockedByPiece
--- PASS: TestIsAttackedBlockedByPiece (0.00s)
=== RUN   TestIsAttackedByKing
--- PASS: TestIsAttackedByKing (0.00s)
=== RUN   TestPositionMovesStartCount
--- PASS: TestPositionMovesStartCount (0.00s)
=== RUN   TestPositionMovesKnightCenter
--- PASS: TestPositionMovesKnightCenter (0.00s)
=== RUN   TestPositionMovesKnightEdge
--- PASS: TestPositionMovesKnightEdge (0.00s)
=== RUN   TestPositionMovesRookOpen
--- PASS: TestPositionMovesRookOpen (0.00s)
=== RUN   TestPositionMovesBishopOpen
--- PASS: TestPositionMovesBishopOpen (0.00s)
=== RUN   TestPositionMovesQueenOpen
--- PASS: TestPositionMovesQueenOpen (0.00s)
=== RUN   TestPositionMovesPawnPromotionSquares
--- PASS: TestPositionMovesPawnPromotionSquares (0.00s)
=== RUN   TestLegalMovesStartCount
--- PASS: TestLegalMovesStartCount (0.00s)
=== RUN   TestLegalMovesPinnedPiece
--- PASS: TestLegalMovesPinnedPiece (0.00s)
=== RUN   TestLegalMovesInCheck
--- PASS: TestLegalMovesInCheck (0.00s)
=== RUN   TestLegalMovesCheckmatePosition
--- PASS: TestLegalMovesCheckmatePosition (0.00s)
=== RUN   TestLegalMovesStalemate
--- PASS: TestLegalMovesStalemate (0.00s)
=== RUN   TestLegalMovesEnPassantIncluded
--- PASS: TestLegalMovesEnPassantIncluded (0.00s)
=== RUN   TestLegalMovesCastlingIncluded
--- PASS: TestLegalMovesCastlingIncluded (0.00s)
=== RUN   TestLegalMovesCastlingBlockedByCheck
--- PASS: TestLegalMovesCastlingBlockedByCheck (0.00s)
=== RUN   TestLegalMovesCastlingThroughAttackedSquare
--- PASS: TestLegalMovesCastlingThroughAttackedSquare (0.00s)
=== RUN   TestMoveValidMove
--- PASS: TestMoveValidMove (0.00s)
=== RUN   TestMoveInvalidMove
--- PASS: TestMoveInvalidMove (0.00s)
=== RUN   TestMoveStatusNotStartedToInProgress
--- PASS: TestMoveStatusNotStartedToInProgress (0.00s)
=== RUN   TestMovePrevMovesAppended
--- PASS: TestMovePrevMovesAppended (0.00s)
=== RUN   TestMoveGameOverError
--- PASS: TestMoveGameOverError (0.00s)
=== RUN   TestMoveFoolsMate
--- PASS: TestMoveFoolsMate (0.00s)
=== RUN   TestMove50MoveRule
--- PASS: TestMove50MoveRule (0.00s)
=== RUN   TestMoveStalemateDetected
--- PASS: TestMoveStalemateDetected (0.00s)
PASS
ok      simpleboard/internal/chess      0.003s
```

## Backend API Docs
## API Endpoints

## Environment Variables

| Variable       | Default                  | Description                    |
|----------------|--------------------------|--------------------------------|
| `PORT`         | `8080`                   | HTTP server port               |
| `DB_PATH`      | `./simpleboard.db`       | SQLite database file path      |
| `CORS_ORIGINS` | `http://localhost:4200`  | Comma-separated allowed origins|
| `JWT_SECRET`   | `no-secret`              | JWT Auth Secret Key            |

## API Endpoints

| Method | Path           | Description          |
|--------|----------------|----------------------|
| GET    | `/api/health`  | Health check         |
| POST   | `/api/register`| Register account     |
| POST   | `/api/login`   | Login to account     |
| POST   | `/api/game`    | Game interaction     |

## Usage

### GET `/api/health` -> `200`

```
{
  "status": "ok"
}
```

### POST `/api/register` -> `201`

#### Example Body
```
{
  "username": "example",
  "email": "example@example.com",
  "password": "secretpassword"
}
```

#### Response
```
{
    "message": "user registered",
        "user": {
            "user_id":  0,
            "username": "example",
            "email":    "example@example.com"
            }
}
```

### POST `/api/login` -> `200`

#### Example Body
```
{
  "username": "example",
  "password": "secretpassword"
}
```

#### Response
```
{
    "message": "login successful",
        "token": "new-jwt-token",
        "user": {
            "user_id":  1,
            "username": "example",
            "email":    "example@example.com"
            }
}
```

### POST `/api/game` -> `200`
`/api/game` has 3 `"action"` field values that direct it's interaction with the game state:
- `"create"` - Creates new game
- `"state"` - Replies with the current game state
- `"move"` - Apply a user move to the game and get the result

All requests must have a valid `Authorization` header of the form:
```
... "Authorization: Bearer <YOUR_JWT_TOKEN_HERE>" ...
```
The JWT token is given upon a successful login, and expires in 24 hours.

#### Example Body
```
{
  "action": "create",
  "player_id": 1,
  "other_id": 2,
  "starting_side": "w"
}
```

#### Response
```
{
    "message":"game created",
        "state": {
            "black_player_id":2,
            "created_at":"2026-03-26T01:27:40.740472882-04:00",
            "game_id":1,
            "next_moves":["a2a3","a2a4","b2b3","b2b4","c2c3","c2c4","d2d3","d2d4","e2e3","e2e4","f2f3","f2f4","g2g3","g2g4","h2h3","h2h4","b1c3","b1a3","g1h3","g1f3"],
            "prev_moves":[],
            "side":"w",
            "state":"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
            "status":"InProgress",
            "updated_at":"2026-03-26T01:27:40.740472882-04:00",
            "white_player_id":1
        }
}
```

#### Example Body
```
{
  "action": "state",
  "game_id": 1,
  "player_id": 1
}
```

#### Response
```
{
    "message":"state",
        "state": {
            "black_player_id":2,
            "created_at":"2026-03-26T01:27:40.740472882-04:00",
            "game_id":1,
            "next_moves":["a2a3","a2a4","b2b3","b2b4","c2c3","c2c4","d2d3","d2d4","e2e3","e2e4","f2f3","f2f4","g2g3","g2g4","h2h3","h2h4","b1c3","b1a3","g1h3","g1f3"],
            "prev_moves":[],
            "side":"w",
            "state":"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
            "status":"InProgress",
            "updated_at":"2026-03-26T01:27:40.740472882-04:00",
            "white_player_id":1
        }
}
```

#### Example Body
```
{
  "action": "move",
  "player_id": 1,
  "game_id": 1,
  "move":"a2a3"
}
```

#### Response
```
{
    "message":"move applied",
        "user": {
            "black_player_id":2,
            "created_at":"2026-03-26T01:27:40.740472882-04:00",
            "game_id":1,
            "next_moves":["b8c6","b8a6","g8h6","g8f6","a7a6","a7a5","b7b6","b7b5","c7c6","c7c5","d7d6","d7d5","e7e6","e7e5","f7f6","f7f5","g7g6","g7g5","h7h6","h7h5"],
            "prev_moves":["a2a3"],
            "side":"b",
            "state":"rnbqkbnr/pppppppp/8/8/8/P7/1PPPPPPP/RNBQKBNR b KQkq - 1 1",
            "status":"InProgress",
            "updated_at":"2026-03-26T01:33:52.383454683-04:00",
            "white_player_id":1
        }
}
```
