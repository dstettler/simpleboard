# Simpleboard Backend

Go backend server for the Simpleboard chess application.

## Directory Structure

```
backend/
├── api/                    # HTTP router and route definitions
├── cmd/
│   └── simpleboard/        # Application entry point (main.go)
├── internal/               # Private application packages
│   ├── chess/              # Chess rule logic and FEN handling
│   ├── domain/             # Domain type definitions
│   ├── handler/            # HTTP request handlers
│   ├── auth/               # Authentication
│   ├── repository/         # GORM model structs
│   └── service/            # Business logic services
├── pkg/                    # Shared utility packages
│   ├── config/             # Runtime configuration loader
│   ├── db/                 # Database connection wrapper
│   └── response/           # JSON response and error helpers
└── simpleboard.db          # Database instance
```

## Getting Started

```bash
go build ./cmd/simpleboard
./simpleboard.exe
```

The server starts on port **8080** by default.

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
