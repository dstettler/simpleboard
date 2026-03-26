# Sprint 2
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
- Built foundational chess engine
  - Legal move generation
  - Game state validator
  - Checkmate detection
- Added state elements to support engine
  - Move list format for requests
  - Pre-generating possible moves for frontend comparison
- Added game state unit tests
- Added password hashing
- Added login API placeholder (for real auth changes)

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
/internal/chess/
=== RUN   TestKingCoordsStartPosition
--- PASS: TestKingCoordsStartPosition (0.00s)
=== RUN   TestKingCoordsMissingKing
--- PASS: TestKingCoordsMissingKing (0.00s)
=== RUN   TestMakeMoveUpdatesBoard
--- PASS: TestMakeMoveUpdatesBoard (0.00s)
=== RUN   TestMakeMoveHalfmoveClock
--- PASS: TestMakeMoveHalfmoveClock (0.00s)
=== RUN   TestMakeMoveFullmoveCounter
--- PASS: TestMakeMoveFullmoveCounter (0.00s)
=== RUN   TestMakeMoveCapture
--- PASS: TestMakeMoveCapture (0.00s)
=== RUN   TestMakeMoveSideToggle
--- PASS: TestMakeMoveSideToggle (0.00s)
=== RUN   TestPositionMovesStartCount
--- PASS: TestPositionMovesStartCount (0.00s)
=== RUN   TestLegalMovesStartCount
--- PASS: TestLegalMovesStartCount (0.00s)
=== RUN   TestStatusString
--- PASS: TestStatusString (0.00s)
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
PASS
```

## Backend API Docs
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
        "user": {
            "user_id":  0,
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

#### Example Body
```
{
  "action": "create",
  "player_id": 0
}
```

#### Response
```
{
    "message":"state",
        "user": {
            "black_player_id":0,
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
    "message":"game created",
        "user": {
            "black_player_id":0,
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
  "game_id": 1,
  "player_id": 1
}
```

#### Response
```
{
    "message":"move applied",
        "user": {
            "black_player_id":0,
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
