# Go Sudoku Solver

A high-performance Sudoku solver implemented in Go. The application supports four distinct solving methodologies: **Traditional Grid Backtracking**, **Knuth's Algorithm X (Exact Cover)** using the **Dancing Links (DLX)** technique, **Bitmask Backtracking**, and **SIMD-Optimized (Tdoku-inspired)**.

## Features

- **Four-Solver Engine**: Easily toggle between standard Backtracking, Knuth's Algorithm X, Bitmask Backtracking, and SIMD-Optimized (Tdoku-inspired).
- **Grading Compliant**: Adheres to strict execution constraints (only prints the final solution or `Error` for invalid boards, with no extra debug lines).
- **Dependency-Free**: Developed strictly using allowed Go built-ins (`os` and `fmt` only).
- **Robust Validation**: Pre-checks board dimensions, characters, row/column length, and minimum clues (minimum 17 numbers of which at least 8 must be unique) to ensure valid Sudoku properties and prevent panics.

---

## Architecture & Algorithms

### 1. Traditional Grid Backtracking
An implementation of Depth-First Search (DFS) directly on the $9\times9$ grid array. It traverses the grid cell-by-cell, attempts digit placement from 1 to 9, and validates constraints (row, column, $3\times3$ box) directly on the grid before recurring. If a path fails, it resets the cell and backtracks.

### 2. Knuth's Algorithm X (Exact Cover via Dancing Links)
For advanced performance, the Sudoku grid is formulated as an **Exact Cover Problem**.
- **Constraint Matrix**: Maps the Sudoku rules to a binary matrix containing $729$ rows (all possible cell-row-value candidates) and $324$ columns (the constraint headers representing cell, row, column, and box constraints).
- **Dancing Links (DLX)**: Employs a toroidal, circularly doubly-linked list node grid. The recursive search chooses the constraint column with the fewest active rows (minimum size heuristic) and covers/uncovers rows using Knuth's pointer manipulation algorithm. Pre-existing clues on the board are pre-covered at startup to optimize execution.

### 3. Bitmask Backtracking
An optimized DFS backtracking engine. Instead of searching lists or maintaining dynamic doubly linked pointers, it stores the state of placed numbers in each row, column, and $3\times3$ box as bits inside $9$ integers (`rowsUsed`, `colsUsed`, `boxesUsed`). Checking validity resolves via single bitwise AND (`&`) operations, placing a digit updates via bitwise OR (`|`), and backtracking clears via bitwise AND NOT (`&^`). This completely eliminates memory allocation and pointer-chasing, keeping all state within ultra-fast CPU L1 cache or registers.

### 4. SIMD-Optimized (Tdoku-Inspired)
A highly optimized, hardware-friendly solver that simulates SIMD parallel vector architectures using native Go 64-bit integer registers:
- **Bitboard Candidate Representation**: The complete board is modeled as candidate bitboards (`Candidates [81]uint16`).
- **Parallel Constraint Propagation**: Integrates DPLL-style constraint satisfying. It recursively propagates Naked Singles and Hidden Singles across horizontal and vertical units.
- **Hardware Acceleration**: Employs Go's native, compiler-optimized `math/bits` package functions (e.g. `bits.OnesCount`, `bits.TrailingZeros`) which compile directly to hardware instructions (e.g., `POPCNT`, `TZCNT` / `BSF`) on both x86_64 and ARM64 (Apple Silicon) architectures.

---

## Configuration (`settings.txt`)

Since external CLI flag packages are not permitted, the solver's algorithm is configured via a local configuration file named `settings.txt` in the root directory.

Create or edit the `settings.txt` file and enter one of the following lines:

- `backtracking` *(Default)*: Uses the naive recursive backtracking solver.
- `exact-cover` (or `algo-x`): Uses Knuth's Algorithm X (DLX) solver.
- `bitmask`: Uses the ultra-fast Bitmask Backtracking solver.
- `tdoku` (or `simd-tdoku`): Uses the hardware-friendly SIMD-Optimized / Tdoku-inspired solver.

*Note: If `settings.txt` is missing, unreadable, or contains any other value, the application automatically defaults to `backtracking`.*

---

## Usage

Run the program by passing exactly 9 arguments, each representing a row of the Sudoku board. Dots (`.`) or `0` can be used to denote empty cells.

### Valid Sudoku Example
```bash
go run . ".96.4...1" "1...6...4" "5.481.39." "..795..43" ".3..8...." "4.5.23.18" ".1.63..59" ".59.7.83." "..359...7"
```

**Output:**
```
3 9 6 2 4 5 7 8 1
1 7 8 3 6 9 5 2 4
5 2 4 8 1 7 3 9 6
2 8 7 9 5 1 6 4 3
9 3 1 4 8 6 2 7 5
4 6 5 7 2 3 9 1 8
7 1 2 6 3 8 4 5 9
6 5 9 1 7 4 8 3 2
8 4 3 5 9 2 1 6 7

```

### Invalid Input Example
```bash
go run . "invalid" "args"
```

**Output:**
```
Error
```

---

## Directory Structure

```
├── main.go               # Entry point, settings parsing, and backtracking DFS solver
├── main_test.go          # Comprehensive 18-case integration test suite
├── settings.txt          # Solver configuration setting file (backtracking or exact-cover)
├── go.mod                # Module specification
└── sudoku/
    ├── algoX.go          # DLX Matrix and Algorithm X solver implementation
    ├── algoX_test.go     # Unit tests specifically for the exact-cover solver
    ├── bitmask.go        # High-performance Bitmask Backtracking solver
    ├── tdoku.go          # SIMD-Optimized (Tdoku-inspired) solver
    ├── tdoku_test.go     # Unit tests for the tdoku solver
    ├── createBoard.go    # Argument-to-grid parsing with length checks
    ├── printBoard.go     # Output rendering format
    ├── checkValid.go     # Grid check helpers for backtracking
    ├── nextCell.go       # Traversal index helper for backtracking
    └── startValid.go     # Board rules and clue pre-validation
```

---

## Performance Benchmarks

<!-- BENCHMARK_START -->
### Benchmark Results

Here is a performance comparison of the four Sudoku solver engines compiled and executed in your local environment.

| Solver Engine | Average Solve Time (μs) | Speedup Factor |
| :--- | :---: | :---: |
| **Traditional Backtracking** | 19.07 μs | Baseline (1.0x) |
| **Knuth's Algorithm X (DLX)** | 78.02 μs | 0.2x faster |
| **Bitmask Backtracking** | 1.31 μs | 14.6x faster |
| **SIMD-Optimized (Tdoku)** | 5.08 μs | 3.8x faster |

![Benchmark Results](benchmark.svg)
<!-- BENCHMARK_END -->

---

## Testing

A comprehensive, table-driven integration test suite is included in `main_test.go`. It compiles the binary once and executes all 18 subject-defined scenarios (including both valid and invalid layouts) under **all four** backtracking, exact-cover, bitmask, and tdoku/simd-tdoku algorithm configurations.

Run all tests across the repository:
```bash
go test -v ./...
```

---

## Faulty Settings Fallback & Resolution

If the `settings.txt` file is missing, unreadable, or contains an invalid/unsupported algorithm name, the application automatically defaults to **Knuth's Algorithm X (Exact Cover)** instead of standard backtracking.

### Detailed Solution & Rationale
1. **Robustness by Default**: In grading/testing environments, configurations might be missing or set incorrectly. Falling back to the highly optimized, reliable `exact-cover` ensures that the program always resolves the Sudoku puzzle under the best performing baseline engine, guaranteeing correct outcomes even for hard puzzles that could otherwise timeout in naive backtracking.
2. **Implementation Details**:
   - The file reading utility `readSettings()` in [main.go](file:///Users/sagar/Downloads/sudoku/main.go) handles file access errors (like `ErrNotExist` or permission issues) by catching the `error` from `os.ReadFile` and directly returning `"exact-cover"`.
   - Before returning the configuration value, the string is trimmed of whitespaces and validated against the whitelist of supported engines: `backtracking`, `exact-cover`, `algo-x`, `bitmask`, `tdoku`, `simd-tdoku`.
   - If the input does not match any entry in this whitelist, it is treated as a configuration fault, and `"exact-cover"` is returned.
