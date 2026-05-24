# Go Sudoku Solver

A Sudoku solver implemented in Go using **Knuth's Algorithm X (Exact Cover)** with the **Dancing Links (DLX)** technique.

## Features

- **Algorithm X (Exact Cover)**: Formulates Sudoku as an exact cover problem and solves it with Dancing Links.
- **Grading Compliant**: Only prints the final solution or `Error` for invalid boards, with no extra debug lines.
- **Dependency-Free**: Uses only allowed Go built-ins (`os` and `fmt` in the entry point).
- **Robust Validation**: Pre-checks board dimensions, characters, row/column length, and minimum clues (minimum 17 numbers of which at least 8 must be unique).

---

## Architecture & Algorithm

The Sudoku grid is formulated as an **Exact Cover Problem**. A toroidal, circularly doubly-linked list (Dancing Links) manipulates columns and rows to find a solution.

*   **Detailed Guide**: See [Knuth's Algorithm X (Exact Cover) Detailed Explanation](docs/exact_cover.md).

---

## Usage

Run the program with exactly 9 arguments, each representing a row of the Sudoku board. Dots (`.`) or `0` denote empty cells.

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
├── main.go               # Entry point
├── main_test.go          # 18-case integration test suite
├── go.mod                # Module specification
└── sudoku/
    ├── algoX.go          # DLX matrix and Algorithm X solver
    ├── algoX_test.go     # Unit tests for the exact-cover solver
    ├── createBoard.go    # Argument-to-grid parsing
    ├── printBoard.go     # Output formatting
    └── startValid.go     # Board rules and clue pre-validation
```

---

## Testing

Integration tests in `main_test.go` compile the binary once and run all 18 subject-defined scenarios (valid and invalid layouts).

```bash
go test -v ./...
```

Unit tests for the DLX solver:

```bash
go test -v ./sudoku -run TestSolveExactCover
```
