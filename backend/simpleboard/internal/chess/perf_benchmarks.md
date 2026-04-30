# Chess Engine Benchmarks

The chess package ships with a small benchmark suite in perf_test.go .
It does two jobs: measure how fast move generation runs, and double-check that we
still generate the *correct* set of moves. 


**Benchmarks** (for `ns/op`, `B/op`, `allocs/op`):

| Benchmark              | What it measures                                                          |
|------------------------|---------------------------------------------------------------------------|
| `BenchmarkLegalMoves`  | `LegalMoves()` throughput across three positions (start, kiwipete, endgame) |
| `BenchmarkPositionMoves` | Pseudo-legal generation only — diff against `LegalMoves` shows the king-safety filter cost |
| `BenchmarkPerft`       | Full move-tree enumeration from the start position at depths 1–4. Also validates the number of legal moves to generate at each depth (20, 400, 8902, 197281). |

**Verbose tests** (run with `-v`):

| Test                  | What it prints                                                                    |
|-----------------------|------------------------------------------------------------------------------------|
| `TestPerfLegalMoves`  | 10,000 calls per position, reports calls/sec and bytes per call                    |
| `TestPerfPerft`       | Perft from the start position at depths 1–3, reports nodes/sec and bytes allocated |


## Running them

From the module root (`backend/simpleboard/`):

```bash
# all benchmarks, default duration 
go test -bench=. -benchmem -run=^$ ./internal/chess/

# single benchmark
go test -bench=BenchmarkLegalMoves -benchmem -run=^$ ./internal/chess/

# run benchmarks longer for more normalized results
go test -bench=. -benchtime=5s -benchmem -run=^$ ./internal/chess/

# or tweak the iteration count
go test -bench=. -benchtime=100x -benchmem -run=^$ ./internal/chess/
```

verbose timing tests:

```bash
# verbose perft + LegalMoves timing
go test -v -run TestPerf ./internal/chess/

# just the LegalMoves throughput test
go test -v -run TestPerfLegalMoves ./internal/chess/
```


