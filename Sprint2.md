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
