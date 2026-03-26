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
│   ├── middleware/         # HTTP middleware (CORS, logging, recovery)
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

## API Endpoints

| Method | Path           | Description          |
|--------|----------------|----------------------|
| GET    | `/api/health`  | Health check         |
| POST   | `/api/register`| Register account     |
| POST   | `/api/login`   | Login to account     |

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

### POST `/api/login` `200`

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
