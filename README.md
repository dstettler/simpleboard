# ♟️Simpleboard 

**Simpleboard** is a no-frills online multiplayer chess game with a focus on a clean user experience where users can jump straight into a game *with or without an account* to play their friends (or foes).

- Users *may login/register* to track games and identity or play instantly without any signup via an invite link.
- The game logic runs entirely on the server side ensuring fair play.

## Running the project
### Frontend

Requirements:
- Node.js

```bash
npm install -g @angular/cli
cd frontend
ng serve
```

Building:

``` bash
ng build
```

### Backend

Requirements:
- go go1.25+

Environment variables for the backend can be easily defined in an `env.sh` using the template:
``` bash
cd backend
cp env.sh.template env.sh
nano env.sh # edit values as needed
source ./env.sh
```

Build and run:
``` bash
cd backend/simpleboard
go build ./cmd/simpleboard
./simpleboard
```

Or simply run:

``` bash
go run ./cmd/simpleboard/
```

API documentation and more can be found in the [Backend Docs](https://github.com/dstettler/simpleboard/blob/main/backend/README.md)

## Features
- Real-time multiplayer chess with complete rule support
- Optional login / registration to track previous matches and personalize the user experience
- On-demand chess instances to support ephemeral game sessions
- Server side game validation

## Stack
- Frontend: Angular with TypeScript
- Backend: Implementation in Go and sqlite DB

## Team Members:
- Arunabho Basu (Backend)
- Sreeram Gangavarapu (Frontend)
- TJ Schultz (Backend)
- Devon Stettler (Frontend)
