# Simpleboard Backend

Go backend server for the Simpleboard chess application, built with Chi router, GORM, and SQLite.

## Directory Structure

```
backend/
├── api/                    # HTTP router and route definitions
├── cmd/
│   └── server/             # Application entry point (main.go)
├── configs/                # Configuration files
├── internal/               # Private application packages
│   ├── chess/              # Chess rule logic and FEN handling
│   ├── handler/            # HTTP request handlers
│   ├── middleware/         # HTTP middleware (CORS, logging, recovery)
│   ├── model/              # GORM model structs
│   └── service/            # Business logic services
├── pkg/                    # Shared utility packages
│   ├── config/             # Runtime configuration loader
│   ├── db/                 # Database connection wrapper
│   └── response/           # JSON response and error helpers
├── scripts/                # Build and utility scripts
└── test/
    └── integration/        # Integration tests
```

## Getting Started

```bash
go build ./cmd/server
./server.exe
```

The server starts on port **8080** by default.

## Environment Variables

| Variable       | Default                  | Description                    |
|----------------|--------------------------|--------------------------------|
| `PORT`         | `8080`                   | HTTP server port               |
| `DB_PATH`      | `./simpleboard.db`       | SQLite database file path      |
| `CORS_ORIGINS` | `http://localhost:4200`   | Comma-separated allowed origins|

## API Endpoints

| Method | Path           | Description          |
|--------|----------------|----------------------|
| GET    | `/api/health`  | Health check         |
