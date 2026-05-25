# Go Sudoku Solver

A Sudoku solver implemented in Go using **backtracking**. All logic lives in a single `main.go` file.

## Features

- **Backtracking search**: Tries digits 1–9 in empty cells, undoing choices when they lead to a dead end.
- **Unique-solution check**: Stops after finding two solutions and reports `Error` if the puzzle is invalid or ambiguous.
- **Grading compliant**: Prints only the final solution or `Error` for invalid boards.
- **Dependency-free**: Uses only `os` and `fmt`.
- **Input validation**: Checks argument count, row length, characters, and that given clues do not already violate Sudoku rules.

---

## What is Sudoku?

**Sudoku** is a number-placement puzzle on a **9×9 grid**. Some cells are filled in at the start (the **clues** or **givens**); the rest are empty. Your goal is to fill every empty cell with a digit from **1 to 9** so that the completed grid satisfies **three rules at the same time**.

Think of it as a logic puzzle, not a math puzzle — you never multiply or add. You only place digits so that nothing repeats where it should not.

### Rule 1: Rows

Each of the **9 horizontal rows** must contain the digits **1 through 9 exactly once**. No digit may appear twice in the same row.

<img src="images/three_rules/1.%20rule%20of%20rows.jpg" alt="Rule of rows: each horizontal row must contain 1–9 with no repeats" width="33%">

When checking a row, scan left to right and make sure you see each number from 1 to 9 once — and only once.

### Rule 2: Columns

Each of the **9 vertical columns** must also contain **1 through 9 exactly once**. No digit may appear twice in the same column.

<img src="images/three_rules/2.%20rule%20of%20columns.jpg" alt="Rule of columns: each vertical column must contain 1–9 with no repeats" width="33%">

When checking a column, scan top to bottom the same way: every digit 1–9 appears once, with no duplicates.

### Rule 3: Blocks (3×3 boxes)

The 9×9 grid is divided into **nine 3×3 sub-grids** (often called **boxes** or **blocks**). Each block must contain **1 through 9 exactly once** as well.

<img src="images/three_rules/3.%20rule_of_blocks.jpg" alt="Rule of blocks: each 3×3 box must contain 1–9 with no repeats" width="33%">

Thicker lines on the grid mark the block boundaries. When you place a digit, it must be valid for its **row**, its **column**, and its **block** all at once.

### Putting the three rules together

A completed Sudoku is a grid where:

- Every row is a permutation of 1–9.
- Every column is a permutation of 1–9.
- Every 3×3 block is a permutation of 1–9.

This solver takes a partially filled grid (via command-line input), checks that the starting clues are legal, and fills in the rest automatically. The sections below explain **how** it does that with backtracking.

---

## How It Works

This section explains the **backtracking** approach used in `main.go`. If you already know Sudoku, you can skim [What is Sudoku?](#what-is-sudoku) above.

### What problem are we solving?

Given a 9×9 Sudoku board with some cells already filled, we must fill every empty cell so that **all three rules hold at once**. A well-formed puzzle has **exactly one** solution.

A natural human strategy is: pick an empty cell, try a digit that does not break any rule, and if you get stuck later, **undo** that digit and try another. That is exactly what this program does, recursively.

### The grid in memory

Two global 9×9 arrays hold state:

| Variable | Role |
|----------|------|
| `grid` | Working board while searching (cells change during backtracking) |
| `solved` | Copy of `grid` when the first complete solution is found |

Empty cells are stored as `0`. Given clues are digits `1`–`9`.

### `isValid` — can we place a digit here?

Before placing digit `num` at `(row, col)`, the solver checks three things (the cell itself must be empty when called):

1. **Row** — no other cell in that row already has `num`.
2. **Column** — no other cell in that column already has `num`.
3. **Block** — no other cell in the same 3×3 box already has `num`.

The block’s top-left corner is computed with integer division:

```go
startRow := (row / 3) * 3
startCol := (col / 3) * 3
```

If any check fails, that digit cannot go in this cell.

### `solve` — backtracking search

`solve` is a recursive function that:

1. **Stops early** if more than one solution has been found (puzzle is not unique).
2. **Finds the next empty cell** by scanning the grid top-to-bottom, left-to-right.
3. If **no empty cell** remains, the board is full — increment the solution counter and save the first solution into `solved`.
4. Otherwise, **try digits 1 through 9** in that cell:
   - If `isValid(row, col, num)`:
     - Place `num` in `grid[row][col]`.
     - Call `solve` recursively.
     - **Undo** the placement (`grid[row][col] = 0`) before trying the next digit.

This is classic **backtracking**: choose → explore → undo on failure.

### `main` — parse, validate, solve, print

1. **Parse** exactly 9 CLI arguments (one row each). `.` means empty; `1`–`9` are clues. Anything else → `Error`.
2. **Validate givens** — for each filled cell, temporarily clear it and run `isValid`. If a given digit already conflicts with another given in its row, column, or box → `Error`.
3. **Solve** — call `solve(&solutions)`. If `solutions != 1` (zero or multiple solutions) → `Error`.
4. **Print** the grid from `solved`, one row per line, digits separated by spaces.

---

## Code walkthrough: Test Case 1

This walkthrough follows the **first valid puzzle** in `main_test.go` — the same example used by `TestSolveBacktracking`. We trace what the code does until the **first empty cell** (row 0, column 0) receives its digit.

### The example puzzle

```text
Input (`.` = empty):

. 9 6 . 4 . . . 1
1 . . . 6 . . . 4
5 . 4 8 1 . 3 9 .
. . 7 9 5 . . 4 3
. 3 . . 8 . . . .
4 . 5 . 2 3 . 1 8
. 1 . 6 3 . . 5 9
. 5 9 . 7 . 8 3 .
. . 3 5 9 . . . 7
```

**First empty cell (reading order):** row **0**, column **0** — top-left corner.

**Final answer for that cell:** **3** (from the full solution in the test).

With backtracking, cells are filled in **scan order** (row 0 col 0, then row 0 col 1, …), not by “hardest cell first.” So (0,0) is actually the **first** cell the solver tries to fill — but many later cells get resolved before (0,0) succeeds, because earlier attempts at (0,0) may fail and backtrack.

### Phase 1 — Parse input (`main`)

Each CLI string becomes one row. `.` → `0`; digits `'1'`–`'9'` → `1`–`9`.

After parsing, row 0 in memory is:

```text
grid[0] = [0, 9, 6, 0, 4, 0, 0, 0, 1]
            ↑
         first empty cell in scan order
```

### Phase 2 — Validate givens

For each non-zero cell, the code clears that cell, calls `isValid`, then restores the value. Our puzzle has no conflicting givens, so validation passes.

If validation fails, `main` prints `Error` and exits before solving.

### Phase 3 — First visit to cell (0,0)

`solve` finds the first `0` at `(0, 0)` and tries digits **1** through **9**:

| Try | `isValid(0, 0, num)` | Why |
|-----|----------------------|-----|
| 1 | false | Row 0 already has 1; column 0 has 1 from row 1 |
| 2 | may pass row/col/box checks | Search continues deeper… |
| … | … | Other branches fill more cells or hit dead ends |
| 3 | true (on the successful path) | Eventually leads to the unique solution |

Backtracking explores failed branches by resetting `grid[0][0] = 0` and trying the next digit. When a full grid is found, `solutions` becomes 1 and `solved` stores that grid.

### Phase 4 — What the successful path looks like

On the path that yields the unique solution, the solver fills cells in scan order. By the time it commits **3** at (0,0), other cells on row 0 are already settled on that branch:

```text
3 9 6 2 4 5 7 8 1   ← row 0 complete
1 . . . 6 . . . 4
5 . 4 8 1 . 3 9 .
...
```

Row 0 needs **3** in the top-left: row 0 already contains 9, 6, 4, 1 and the top-left box cannot repeat those. Digits **2, 3, 7, 8** are the only candidates by pencil-mark logic; backtracking discovers that **3** is the only digit that completes the **entire** puzzle consistently.

### Phase 5 — Print result

When exactly one solution exists, `main` prints `solved`:

```text
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

### What if a guess is wrong?

If placing a digit leads to a dead end (no empty cell can accept any digit 1–9), recursion returns, the code sets `grid[row][col] = 0`, and the next digit is tried. If all digits fail at some cell, the branch fails and search backtracks further up.

The puzzle is rejected when:

- **No** complete grid exists (`solutions == 0`), or
- **More than one** complete grid exists (`solutions > 1`).

### Map: concepts → code

| Concept | Where in `main.go` |
|---------|-------------------|
| Parse `.` and digits | `main` (loop over `args`) |
| Check row/column/block rules | `isValid` |
| Find next empty cell | `solve` (nested `for` loops) |
| Try digit and recurse | `solve` (`for num := 1; num <= 9`) |
| Undo a bad guess | `grid[row][col] = 0` after recursive call |
| Count solutions | `solutions` pointer; stop if `> 1` |
| Save first solution | `solved = grid` when `*solutions == 1` |
| Print answer | `main` (loop over `solved`) |

### Key takeaways

1. **Backtracking** = try a legal digit, recurse, undo if the deeper search fails.
2. **Empty cells are visited in scan order** (row-major), not by a separate “hardest first” heuristic.
3. **Givens are fixed** before `solve` runs; validation ensures they do not already break Sudoku rules.
4. **Uniqueness** matters: the program only accepts puzzles with exactly one solution.

---

## Usage

Run with exactly 9 arguments (one row each). Use `.` for empty cells.

### Valid example

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

### Invalid example

```bash
go run . "invalid" "args"
```

**Output:**

```
Error
```

---

## Project layout

```
├── main.go          # Entry point, parser, validator, and backtracking solver
├── main_test.go     # Integration and unit tests
├── README.md        # Concepts, usage, and code walkthrough
├── images/          # Diagrams for Sudoku rules
└── go.mod
```

---

## Testing

```bash
go test -v .
```

- **`TestAllScenarios`** — builds the binary and runs 18 integration cases (11 valid puzzles, invalid puzzles, and bad arguments).
- **`TestSolveBacktracking`** — loads Test Case 1 into `grid`, calls `solve` directly, and checks `solved` against the expected 9×9 answer.
