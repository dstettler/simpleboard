# Simpleboard Backend

Go backend server for the Simpleboard chess application.

## Directory Structure

```
backend/simpleboard/
├── api/                    # HTTP router and route definitions
├── cmd/
│   └── simpleboard/        # Application entry point (main.go)
├── internal/               # Private application packages
|   ├── auth/               # Authentication
│   ├── chess/              # Chess rule logic and FEN handling
│   ├── domain/             # Domain type definitions
│   ├── handler/            # HTTP request handlers
│   ├── repository/         # GORM model structs
|   ├── timer/              # Game timer functionality
|   └── utils/              # Backend utils functions
├── pkg/                    # Shared utility packages
│   ├── config/             # Runtime configuration loader
│   ├── db/                 # Database connection wrapper
│   └── response/           # JSON response and error helpers
└── simpleboard.db          # Database instance (example path)
```

## Getting Started

```bash
cd simpleboard
go build ./cmd/simpleboard
./simpleboard.exe
```

The server starts on port **8080** by default.

## Environment Variables
Environment variables for the backend can be easily defined in an `env.sh` using the template:
``` bash
cp env.sh.template env.sh
nano env.sh # edit values as needed
source ./env.sh
```

| Variable                       | Default                  | Description                              |
|--------------------------------|--------------------------|------------------------------------------|
| `PORT`                         | `8080`                   | HTTP server port                         |
| `DB_PATH`                      | `./simpleboard.db`       | SQLite database file path                |
| `CORS_ORIGINS`                 | `http://localhost:4200`  | Comma-separated allowed origins          |
| `JWT_SECRET`                   | `no-secret`              | JWT Auth Secret Key                      |
| `DEFAULT_TIME_CONTROL_SECONDS` | `600`                    | Default per-side clock for new games (s) |
| `SWEEP_INTERVAL_SECONDS`       | `30`                     | Background flag-fall sweep interval (s)  |

## API Endpoints

| Method | Path             | Auth required | Description                      |
|--------|------------------|---------------|----------------------------------|
| GET    | `/api/health`    | No            | Health check                     |
| GET    | `/api/guest`     | No            | Generate a guest token           |
| POST   | `/api/register`  | No            | Register a new account           |
| POST   | `/api/login`     | No            | Login, returns token + streak    |
| POST   | `/api/game`      | Yes           | Create / join / poll / move      |
| GET    | `/api/dashboard` | Yes (user)    | Lifetime stats for current user  |
| GET    | `/api/games`     | Yes (user)    | Game history for current user    |

## Usage

### GET `/api/health` -> `200`
#### Response
```
{
  "status": "ok"
}
```

### GET `/api/guest` -> `200`
`/api/guest` will serve a new guest token if the request is made with no `Authorization` header.
- Used to generate a `guest_id` and required auth token for creating / joining ephemeral game sessions

#### Response
```
{
    "message":"guest creation successful",
    "token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWVzdF9pZCI6ImViOGIzMWQ1LTMyY2QtNGQ0NS05NTZkLTBkZGU1MzNlM2M0ZCIsImV4cCI6MTc3NzQzNjQxOSwiaWF0IjoxNzc3MzUwMDE5fQ.8VUlb9c0Jsfgh0fFlA3Tymz3ZVVf45rwSJwwTyZM_6k",
    "user":{
        "guest_id":"eb8b31d5-32cd-4d45-956d-0dde533e3c4d"
    }
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
        "token": <new-jwt-token>,
        "user": {
            "user_id": 1,
            "username": "example",
            "email": "example@example.com",
            "current_streak": 4,
            "longest_streak": 12
        }
}
```

### POST `/api/game` -> `200`
`/api/game` has 4 `"action"` field values that direct it's interaction with the game state:
- `"create"` - Creates new game
- `"join"` - Joins a game in queue (can be done via an invite link)
- `"state"` - Replies with the current game state
- `"move"` - Apply a user move to the game and get the result

All requests must have a valid `Authorization` header of the form:
```
... "Authorization: Bearer <YOUR_JWT_TOKEN_HERE>" ...
```
The JWT token is given upon a successful login or guest user creation, and expires in 24 hours.

#### Example Body
```
{
  "action": "create",
  "player_id": 1,
  "other_id": 2, // optional - only for games w/ 2 known users to start
  "starting_side": "w"
  "time_control_seconds": 700
}
```
`time_control_seconds` is **optional**. Omit (or send `0`) to use the server default (`DEFAULT_TIME_CONTROL_SECONDS`, 10 min). Both sides get the same starting clock.

#### Response
```
{
    "message":"game created",
        "state": {
            "black_guest_id":"",
            "black_player_id":2,
            "created_at":"2026-03-26T01:27:40.740472882-04:00",
            "game_id":"f0e510f2-0d72-4ce2-ab38-025e224c55c0",
            "next_moves":["a2a3","a2a4","b2b3","b2b4","c2c3","c2c4","d2d3","d2d4","e2e3","e2e4","f2f3","f2f4","g2g3","g2g4","h2h3","h2h4","b1c3","b1a3","g1h3","g1f3"],
            "prev_moves":[],
            "side":"w",
            "state":"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
            "status":"InProgress",
            "time_control_seconds": 700,
            "white_remaining_ms": 694120,
            "black_remaining_ms": 700000,
            "last_move_at": "2026-03-26T01:27:40.740472882-04:00",
            "server_time": "2026-03-26T01:27:40.740472882-04:00",
            "updated_at":"2026-03-26T01:27:40.740472882-04:00",
            "white_guest_id":"",
            "white_player_id":1
        }
}
```

#### Example Body
```
{
  "action": "create",
  "guest_id": "eb8b31d5-32cd-4d45-956d-0dde533e3c4d",
  "starting_side": "w"
}
```

#### Response
```
{
    "message":"game created",
        "state": {
            "black_guest_id":"",
            "black_player_id":0,
            "created_at":"2026-03-26T01:27:40.740472882-04:00",
            "game_id":"f0e510f2-0d72-4ce2-ab38-025e224c55c0",
            "next_moves":["a2a3","a2a4","b2b3","b2b4","c2c3","c2c4","d2d3","d2d4","e2e3","e2e4","f2f3","f2f4","g2g3","g2g4","h2h3","h2h4","b1c3","b1a3","g1h3","g1f3"],
            "prev_moves":[],
            "side":"w",
            "state":"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
            "status":"NotStarted",
            "time_control_seconds": 600,
            "white_remaining_ms": 600000,
            "black_remaining_ms": 600000,
            "last_move_at": "2026-03-26T01:33:52.383454683-04:00",
            "server_time": "2026-03-26T01:33:52.391204000-04:00",
            "updated_at":"2026-03-26T01:27:40.740472882-04:00",
            "white_guest_id":"eb8b31d5-32cd-4d45-956d-0dde533e3c4d"
            "white_player_id":0
        }
}
```

#### Example Body
```
{
  "action": "create",
  "player_id": 1,
  "starting_side": "w"
}
```

#### Response
```
{
    "message":"game created",
        "state": {
            "black_guest_id":"",
            "black_player_id":0,
            "created_at":"2026-03-26T01:27:40.740472882-04:00",
            "game_id":"f0e510f2-0d72-4ce2-ab38-025e224c55c0",
            "next_moves":["a2a3","a2a4","b2b3","b2b4","c2c3","c2c4","d2d3","d2d4","e2e3","e2e4","f2f3","f2f4","g2g3","g2g4","h2h3","h2h4","b1c3","b1a3","g1h3","g1f3"],
            "prev_moves":[],
            "side":"w",
            "state":"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
            "status":"NotStarted",
            "time_control_seconds": 600,
            "white_remaining_ms": 600000,
            "black_remaining_ms": 600000,
            "last_move_at": "2026-03-26T01:33:52.383454683-04:00",
            "server_time": "2026-03-26T01:33:52.391204000-04:00",
            "updated_at":"2026-03-26T01:27:40.740472882-04:00",
            "white_guest_id":""
            "white_player_id":1
        }
}
```

#### Example Body
```
{
  "action": "join",
  "game_id": "f0e510f2-0d72-4ce2-ab38-025e224c55c0", // game with existing white player
  "guest_id": "eb8b31d5-32cd-4d45-956d-0dde533e3c4d",
}
```

#### Response
```
{
    "message":"game joined",
        "state": {
            "black_guest_id":"eb8b31d5-32cd-4d45-956d-0dde533e3c4d",
            "black_player_id":0,
            "created_at":"2026-03-26T01:27:40.740472882-04:00",
            "game_id":"f0e510f2-0d72-4ce2-ab38-025e224c55c0",
            "next_moves":["a2a3","a2a4","b2b3","b2b4","c2c3","c2c4","d2d3","d2d4","e2e3","e2e4","f2f3","f2f4","g2g3","g2g4","h2h3","h2h4","b1c3","b1a3","g1h3","g1f3"],
            "prev_moves":[],
            "side":"w",
            "state":"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
            "status":"InProgress",
            "time_control_seconds": 600,
            "white_remaining_ms": 594120,
            "black_remaining_ms": 600000,
            "last_move_at": "2026-03-26T01:33:52.383454683-04:00",
            "server_time": "2026-03-26T01:33:52.391204000-04:00",
            "updated_at":"2026-03-26T01:27:40.740472882-04:00",
            "white_guest_id":"",
            "white_player_id":1
        }
}
```

#### Example Body
```
{
    "action": "join",
    "game_id": "f0e510f2-0d72-4ce2-ab38-025e224c55c0", // game with existing white player
    "player_id": 2,
}
```

#### Response
```
{
    "message":"game joined",
        "state": {
            "black_guest_id":"",
            "black_player_id":2,
            "created_at":"2026-03-26T01:27:40.740472882-04:00",
            "game_id":"f0e510f2-0d72-4ce2-ab38-025e224c55c0",
            "next_moves":["a2a3","a2a4","b2b3","b2b4","c2c3","c2c4","d2d3","d2d4","e2e3","e2e4","f2f3","f2f4","g2g3","g2g4","h2h3","h2h4","b1c3","b1a3","g1h3","g1f3"],
            "prev_moves":[],
            "side":"w",
            "state":"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
            "status":"InProgress",
            "time_control_seconds": 600,
            "white_remaining_ms": 594120,
            "black_remaining_ms": 600000,
            "last_move_at": "2026-03-26T01:33:52.383454683-04:00",
            "server_time": "2026-03-26T01:33:52.391204000-04:00",
            "updated_at":"2026-03-26T01:27:40.740472882-04:00",
            "white_guest_id":"",
            "white_player_id":1
        }
}
```

#### Example Body
```
{
    "action": "state",
    "game_id": "f0e510f2-0d72-4ce2-ab38-025e224c55c0",
    "player_id": 1
}
```

#### Response
```
{
    "message":"state",
        "state": {
            "black_guest_id":"",
            "black_player_id":2,
            "created_at":"2026-03-26T01:27:40.740472882-04:00",
            "game_id":"f0e510f2-0d72-4ce2-ab38-025e224c55c0",
            "next_moves":["a2a3","a2a4","b2b3","b2b4","c2c3","c2c4","d2d3","d2d4","e2e3","e2e4","f2f3","f2f4","g2g3","g2g4","h2h3","h2h4","b1c3","b1a3","g1h3","g1f3"],
            "prev_moves":[],
            "side":"w",
            "state":"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
            "status":"InProgress",
            "time_control_seconds": 600,
            "white_remaining_ms": 597230,
            "black_remaining_ms": 600000,
            "last_move_at": "2026-03-26T01:27:40.740472882-04:00",
            "server_time": "2026-03-26T01:27:43.973102000-04:00",
            "updated_at":"2026-03-26T01:27:40.740472882-04:00",
            "white_guest_id":"",
            "white_player_id":1
        }
}
```
The remaining-ms values returned for `state` are **live**: the active side's clock is computed as `stored - (server_time - last_move_at)`. Re-poll periodically and the active side's number will keep dropping. The inactive side's number is the stored value.

#### Example Body
```
{
    "action": "move",
    "player_id": 1,
    "game_id": "f0e510f2-0d72-4ce2-ab38-025e224c55c0",
    "move":"a2a3"
}
```

#### Response
```
{
    "message":"move applied",
        "state": {
            "black_guest_id":"",
            "black_player_id":2,
            "created_at":"2026-03-26T01:27:40.740472882-04:00",
            "game_id":"f0e510f2-0d72-4ce2-ab38-025e224c55c0",
            "next_moves":["b8c6","b8a6","g8h6","g8f6","a7a6","a7a5","b7b6","b7b5","c7c6","c7c5","d7d6","d7d5","e7e6","e7e5","f7f6","f7f5","g7g6","g7g5","h7h6","h7h5"],
            "prev_moves":["a2a3"],
            "side":"b",
            "state":"rnbqkbnr/pppppppp/8/8/8/P7/1PPPPPPP/RNBQKBNR b KQkq - 1 1",
            "status":"InProgress",
            "time_control_seconds": 600,
            "white_remaining_ms": 594120,
            "black_remaining_ms": 600000,
            "last_move_at": "2026-03-26T01:33:52.383454683-04:00",
            "server_time": "2026-03-26T01:33:52.391204000-04:00",
            "updated_at":"2026-03-26T01:33:52.383454683-04:00",
            "white_guest_id":"",
            "white_player_id":1
        }
}
```
On a successful move, the responding side has switched (`side` is now the opponent), the moving player's elapsed time has been deducted from their clock, and `last_move_at` jumps to the move time. If the moving player's clock had already run out, the response instead has `"message":"flag fall"` and `status` is set to `WinWhite` or `WinBlack` -- the move is **not** applied in that case.

## Game Timer Functionality

Every game now carries a per-side chess clock. The server is the source of truth -- it decides when a player has run out of time, so clients can never cheat by stalling. This section is everything the frontend needs to render and use it.

### State payload fields

Every `create` / `state` / `move` response includes these on the `state` object:

| Field                  | Type      | Meaning                                                                  |
|------------------------|-----------|--------------------------------------------------------------------------|
| `time_control_seconds` | int       | Starting clock per side (e.g. `600` = 10 min)                            |
| `white_remaining_ms`   | int       | White's remaining time, in milliseconds                                  |
| `black_remaining_ms`   | int       | Black's remaining time, in milliseconds                                  |
| `last_move_at`         | timestamp | When the active side's clock started ticking (game start, or last move)  |
| `server_time`          | timestamp | Server's current time at the moment of response (for drift correction)   |
| `side`                 | `"w"`/`"b"` | Whose turn it is -- this is the side whose clock is currently ticking  |
| `status`               | string    | `InProgress`, `Draw`, `WinWhite`, `WinBlack` (last two cover flag falls) |

The remaining-ms values are **already live** at the moment of response: for the side whose clock is ticking, the server has already subtracted `(server_time - last_move_at)` before sending. You don't need to re-do that math on receipt.

### Rendering a smooth countdown between polls

Polling every second would waste bandwidth. The recommended pattern:

1. On each response, capture `serverTime` and `lastMoveAt` from the payload, plus the local clock time `t0 = Date.now()` at receipt.
2. For the **active side** (`side` field), display:
   ```
   activeRemainingMs - (Date.now() - t0)
   ```
   (i.e. start a local 1-second tick and decrement smoothly)
3. For the **inactive side**, just display the stored remaining ms (it's frozen).
4. Re-poll `state` every ~5-10 seconds to resync against drift, or after every move you receive.

`server_time` exists so you can detect clock drift between client and server -- if you want to be robust, anchor display time to `server_time` rather than `Date.now()`.

### Setting a custom time control

When calling `action: "create"`, optionally pass `time_control_seconds`. Examples:

| Mode    | Seconds |
|---------|---------|
| Bullet  | `60` or `120` |
| Blitz   | `180` or `300` |
| Rapid   | `600` (default) or `900` |
| Classical | `1800`+ |

Omit the field, send `0`, or send a negative value -> server falls back to `DEFAULT_TIME_CONTROL_SECONDS` (10 min).

### What happens when a clock runs out

1. **During a move**: if the moving player's flag has fallen, the server returns `"message":"flag fall"`, sets `status` to `WinWhite` or `WinBlack`, and the move is **rejected** (board state is unchanged).
2. **During a state poll**: if the active side has run out while idle, the server marks the game as ended and returns the final status. Clients see `status` flip from `InProgress` to `WinWhite` / `WinBlack` between polls.
3. **With nobody polling**: a background sweeper goroutine scans in-progress games every `SWEEP_INTERVAL_SECONDS` (default 30s) and ends any whose active side has flag-fallen. So a game can't sit alive forever just because both browsers were closed.

After flag fall, the loser's `*_remaining_ms` is `0`. The opponent's number is whatever was stored at the moment of the loss.

### Status values to watch for

| Status       | Meaning                                                               |
|--------------|-----------------------------------------------------------------------|
| `InProgress` | Game is live, clocks tick on the active side                          |
| `Draw`       | Stalemate / 50-move rule / etc.                                       |
| `WinWhite`   | White wins (checkmate, resignation, **or black flag-fell on time**)   |
| `WinBlack`   | Black wins (checkmate, resignation, **or white flag-fell on time**)   |

The status string alone doesn't tell you whether a win was by checkmate or by time -- if you need to differentiate in the UI, check whether the loser's `*_remaining_ms` is `0` at game end.

---

## User Stats, Game History & Daily Streak

These three endpoints back the dashboard and streak features. All require a registered user token — guests get `401`.

---

### POST `/api/login` response

The login response includes streak fields so the UI can show the streak immediately after sign-in without a second request.

```
{
    "message": "login successful",
        "token": <new-jwt-token>,
        "user": {
            "user_id": 1,
            "username": "example",
            "email": "example@example.com",
            "current_streak": 4,
            "longest_streak": 12
        }
}
```

---

### GET `/api/dashboard` → `200`

Returns lifetime stats for the currently authenticated user.

**Headers**
```
Authorization: Bearer <token>
```

**Response**
```
{
    "user_id":        1,
    "username":       "example",
    "total_games":    38,
    "wins":           20,
    "losses":         14,
    "win_rate":       0.526,
    "current_streak": 4,
    "longest_streak": 12
}
```

| Field            | Type    | Notes                                             |
|------------------|---------|---------------------------------------------------|
| `total_games`    | int     | All completed games (wins + losses + draws)       |
| `wins`           | int     | Games won                                         |
| `losses`         | int     | Games lost                                        |
| `win_rate`       | float64 | `wins / total_games`; `0.0` if no games yet       |
| `current_streak` | int     | Consecutive days logged in ending today           |
| `longest_streak` | int     | All-time best streak                              |

> **Note:** `total_games` counts draws too, so `wins + losses` may be less than `total_games`.

---

### GET `/api/games` → `200`

Returns the authenticated user's game history, newest first.

**Headers**
```
Authorization: Bearer <token>
```

**Response**
```
{
    "user_id": 1,
    "games": [
    {
        "game_id":     "f0e510f2-0d72-4ce2-ab38-025e224c55c0",
        "status":      "WinWhite",
        "played_as":   "w",
        "opponent_id": 2,
        "created_at":  "2026-04-28T14:00:00Z",
        "updated_at":  "2026-04-28T14:22:00Z"
    },
    {
        "game_id":     "a1b2c3d4-...",
        "status":      "WinBlack",
        "played_as":   "b",
        "opponent_id": 0,
        "created_at":  "2026-04-27T10:00:00Z",
        "updated_at":  "2026-04-27T10:31:00Z"
    }
  ]
}
```

| Field         | Type   | Notes                                                   |
|---------------|--------|---------------------------------------------------------|
| `game_id`     | string | UUID of the game                                        |
| `status`      | string | `NotStarted`, `InProgress`, `Draw`, `WinWhite`, `WinBlack` |
| `played_as`   | string | `"w"` or `"b"` — which side this user played            |
| `opponent_id` | uint   | The other player's user ID; `0` if they were a guest    |
| `created_at`  | string | ISO 8601 timestamp                                      |
| `updated_at`  | string | ISO 8601 timestamp — effectively when the game ended    |

---

### How the daily streak works

The streak is **login-based**: it increments once per calendar day when the user logs in.

`current_streak` and `longest_streak` are returned in both the `/api/login` response and `/api/dashboard` — show whichever fits the UI context.

**Game stats** (`wins`, `losses`, `total_games`) update automatically the moment a game ends — checkmate, stalemate, 50-move draw, or flag fall. Guest players don't accumulate stats; only registered accounts do.

---

