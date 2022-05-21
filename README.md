# Go Sudoku Solver

A high-performance Sudoku solver implemented in Go. The application supports three distinct solving methodologies: **Traditional Grid Backtracking**, **Knuth's Algorithm X (Exact Cover)** using the **Dancing Links (DLX)** technique, and **Bitmask Backtracking**.

## Features

- **Three-Solver Engine**: Easily toggle between standard Backtracking, Knuth's Algorithm X, and Bitmask Backtracking.
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

---

## Configuration (`settings.txt`)

Since external CLI flag packages are not permitted, the solver's algorithm is configured via a local configuration file named `settings.txt` in the root directory.

Create or edit the `settings.txt` file and enter one of the following lines:

- `backtracking` *(Default)*: Uses the naive recursive backtracking solver.
- `exact-cover` (or `algo-x`): Uses Knuth's Algorithm X (DLX) solver.
- `bitmask`: Uses the ultra-fast Bitmask Backtracking solver.

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
    ├── createBoard.go    # Argument-to-grid parsing with length checks
    ├── printBoard.go     # Output rendering format
    ├── checkValid.go     # Grid check helpers for backtracking
    ├── nextCell.go       # Traversal index helper for backtracking
    └── startValid.go     # Board rules and clue pre-validation
```

---

## Testing

A comprehensive, table-driven integration test suite is included in `main_test.go`. It compiles the binary once and executes all 18 subject-defined scenarios (including both valid and invalid layouts) under **all three** backtracking, exact-cover, and bitmask algorithm configurations.

Run all tests across the repository:
```bash
go test -v ./...
```
