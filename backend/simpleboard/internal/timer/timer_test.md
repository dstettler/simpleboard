
# Running Tests
From the module root (`backend/simpleboard/`):

```bash
# run all timer package tests
go test ./internal/timer/

# verbose output
go test -v ./internal/timer/

# run a specific test
go test -v -run TestApplyMove_FlagFall ./internal/timer/

# run all tests for a specific function area
go test -v -run TestApplyMove   ./internal/timer/
go test -v -run TestLiveRemaining ./internal/timer/
go test -v -run TestMarkIfTimedOut ./internal/timer/
```

The timer package is CGO-free: every test runs against in-memory structs, no SQLite driver required. `go test ./internal/timer/...` works on a clean Go install with no extra toolchain.


### InitGame

| Test                                            | What it checks
|-------------------------------------------------|----------------
| TestInitGame_UsesProvidedControl                | Provided control value sets TimeControlSeconds, both clocks = control * 1000 ms, LastMoveAt = `now`
| TestInitGame_FallsBackToDefaultWhenZeroOrNegative | Zero and negative inputs both fall back to the package default (600s)


### ApplyMove

| Test                                  | What it checks
|---------------------------------------|----------------
| TestApplyMove_DeductsActiveSideOnly   | White-to-move with 2s elapsed: white's clock drops by ~2000ms, black's clock untouched, LastMoveAt advances to `now`
| TestApplyMove_FlagFall                | Black has 500ms left and 2s have elapsed: returns `(true, "b")`, BlackRemainingMs clamped to 0
| TestApplyMove_ClampsBackwardClockSkew | LastMoveAt set in the future (clock skew): elapsed treated as 0, no negative deduction


### LiveRemaining

| Test                                | What it checks
|-------------------------------------|----------------
| TestLiveRemaining_DoesNotMutate     | Stored WhiteRemainingMs, BlackRemainingMs, and LastMoveAt are untouched after the call
| TestLiveRemaining_ComputesActiveSide | Black to move with 4s elapsed: returned black ms ≈ stored - 4000, white ms = stored (frozen)
| TestLiveRemaining_FrozenWhenGameFinished | Status = WinBlack: returned values match stored, timedOut is false even with stale clocks
| TestLiveRemaining_DetectsFlagFall   | White has 100ms left and 2s have elapsed: returns timedOut=true, loser="w", whiteMs=0


### FlagFallStatus

| Test               | What it checks
|--------------------|----------------
| TestFlagFallStatus | `"w"` -> `chess.WinBlack`, `"b"` -> `chess.WinWhite`


### markIfTimedOut (sweeper helper)

| Test                                  | What it checks
|---------------------------------------|----------------
| TestMarkIfTimedOut_FlagsWhiteOutOfTime | White-to-move with 100ms left and 5s elapsed: returns true, Status set to WinBlack, WhiteRemainingMs zeroed
| TestMarkIfTimedOut_LeavesActiveGameAlone | Plenty of time on both clocks, recent move: returns false, Status stays InProgress
| TestMarkIfTimedOut_IgnoresFinishedGames | Already-won game with stale clocks: returns false, Status stays WinWhite, no clock fields rewritten


## Test design notes

- The pure helpers (`InitGame`, `ApplyMove`, `LiveRemaining`, `FlagFallStatus`, `markIfTimedOut`) operate on `*repository.Game` structs in memory, so all 13 tests run without spinning up SQLite or any other dependency.
- `markIfTimedOut` was extracted from `sweepOnce` specifically so the per-game decision logic can be unit-tested without a real database. The DB-coupled `sweepOnce` itself is a thin loop (query -> apply -> save) and is left to integration testing.
- The `_ClampsBackwardClockSkew` and `_DoesNotMutate` cases guard properties that are easy to regress on -- negative durations and accidental writes through the live-clock helper. Worth keeping if anyone refactors the timer math.
