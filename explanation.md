# Step-by-Step Walkthrough: Dancing Links on a Real Test Puzzle

This document is a **companion to [README.md](README.md)**. The README explains concepts; this file shows **what actually happens in code** when the solver runs — using **Test Case 1** from `main_test.go`.

We stop as soon as the solver places a digit in the **first empty cell** of the input (top-left, row 0 column 0). We do **not** solve the whole board here.

---

## The example puzzle

From `TestAllScenarios` / `TestSolveExactCover` in `main_test.go`:

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

**First empty cell (reading left-to-right, top-to-bottom):** row **0**, column **0** — the top-left corner.

**Final answer for that cell:** **3** (from the full solution in the test).

**Important surprise for beginners:** The solver does **not** fill cells in reading order. It picks the **hardest constraint first** (column with fewest options). Cell (0,0) is chosen **last** among the first 10 decisions — but this walkthrough still shows exactly **how** it gets the value **3**.

---

## Big picture: five phases in `main.go`

| Phase | Function | What happens |
|-------|----------|--------------|
| 1. Parse input | `createBoard` | Turn 9 strings into a `9×9` int array (`0` = empty) |
| 2. Validate | `startValid` | Check clue count and no duplicate givens |
| 3. Build matrix | `solveExactCover` (lines 57–124) | Create 324 columns, 729 candidates, link nodes |
| 4. Lock givens | `solveExactCover` (lines 126–137) | `cover` every clue already on the board |
| 5. Search | `search()` (lines 141–166) | Recursively pick columns, try rows, backtrack |

---

## Phase 1 — Parse input (`createBoard`)

Each CLI string becomes one row. `.` or `0` → `0`. Digits `'1'`–`'9'` → `1`–`9`.

After parsing, row 0 looks like this in memory:

```text
board[0] = [0, 9, 6, 0, 4, 0, 0, 0, 1]
            ↑
         empty → first cell to solve visually
```

**Code reference:**

```go
// main.go — createBoard
if startCondition[i][j] >= '1' && startCondition[i][j] <= '9' {
    startBoard[i][j] = int(startCondition[i][j] - 48)
}
```

---

## Phase 2 — Validate (`startValid`)

Briefly: the board must have 17–77 givens, at least 8 different digits, and no row/column/box duplicates among **given** cells. Our puzzle passes.

If validation fails, `main` prints `Error` and never builds the matrix.

---

## Phase 3 — Build the DLX matrix

### One candidate = one Sudoku guess

Every guess is: **“put digit `v` in cell `(r, c)`.”**

There are **729** guesses (81 cells × 9 digits). Each guess becomes **one row** in the exact-cover matrix — a horizontal ring of **4 nodes**:

| Node links to column | Meaning |
|---------------------|---------|
| `c1 = r*9 + c` | Cell `(r,c)` is filled |
| `c2 = 81 + r*9 + (v-1)` | Row `r` contains digit `v` |
| `c3 = 162 + c*9 + (v-1)` | Column `c` contains digit `v` |
| `c4 = 243 + box*9 + (v-1)` | Box contains digit `v` |

### Mini example: guess “put **3** at **(0,0)**”

```go
// main.go — inside the r,c,v loops
c1 := r*9 + c                                    // 0   → cell(0,0)
c2 := 81 + r*9 + (v - 1)                         // 83  → row0-digit3
c3 := 162 + c*9 + (v - 1)                        // 164 → col0-digit3
c4 := 243 + ((r/3)*3+c/3)*9 + (v - 1)            // 245 → box0-digit3
```

For `(r,c,v) = (0,0,3)` the four column IDs are **0, 83, 164, 245**.

Four `Node` structs are created, inserted vertically into those columns, and linked left-right in a ring. The first node is stored in `rowNodes[0][0][3]`.

```go
// main.go — linking the horizontal ring
nodes[0].Right = nodes[1]
nodes[1].Left = nodes[0]
// ... nodes[1]↔nodes[2]↔nodes[3]↔nodes[0]
rowNodes[r][c][v] = nodes[0]
```

### Column headers and `root`

324 column headers are chained in a circle. `root` sits at the start of that circle. While columns remain linked to `root`, the puzzle is not finished.

---

## Phase 4 — Cover the given clues

For every filled cell on the input board, the solver **commits** to that digit being correct — it removes conflicting options from the matrix.

```go
// main.go — cover givens
if board[r][c] != 0 {
    v := board[r][c]
    node := rowNodes[r][c][v]
    cover(node.Col)
    for j := node.Right; j != node; j = j.Right {
        cover(j.Col)
    }
}
```

**What `cover` does (plain English):**

1. Remove this **column header** from the horizontal header ring (this rule is “satisfied” for now).
2. For every remaining candidate in that column, remove that entire candidate row from all other columns it touches.

So covering the clue **9 at (0,1)** removes all guesses that would violate “row 0 already has a 9” or “column 1 already has a 9”, etc.

### Pencil marks for cell (0,0) after covering all givens

Before search starts, which digits could still go at (0,0)?

| Check | Already used near (0,0) | Ruled out |
|-------|-------------------------|-----------|
| Row 0 | 9, 6, 4, 1 | 1, 4, 6, 9 |
| Column 0 | 1, 5, 4 (from rows 1, 2, 5) | 1, 4, 5 |
| Top-left box | 9, 6, 1, 5, 4 | 1, 4, 5, 6, 9 |

**Still possible:** **2, 3, 7, 8** (four options).

The DLX column `cell(0,0)` (column ID **0**) has **Size = 4** — exactly four linked candidates remain.

---

## Phase 5 — Search begins (`search`)

### Step A: Stopping condition

```go
if root.Right == &root {
    return true  // no columns left → solved!
}
```

Not true yet — many columns remain.

### Step B: Pick a column (`selectColumn`)

```go
col := selectColumn(&root)  // column with smallest Size (MRV heuristic)
cover(col)
```

`selectColumn` scans all active column headers and picks the one with the **fewest remaining candidates**. This is the **“most constrained variable first”** idea — try the slot with the fewest choices before guessing on wide-open cells.

**First column chosen for our puzzle:** `cell(5,1)` — row 5, column 1 — with **Size = 1** (only one option left: digit **6**).

The solver does **not** start at (0,0). It starts wherever the matrix is most forced.

### Step C: Try each row in that column

```go
for r := col.Head.Down; r != &col.Head; r = r.Down {
    solution = append(solution, r)
    for j := r.Right; j != r; j = j.Right {
        cover(j.Col)   // cover the other 3 columns this guess touches
    }
    if search() { return true }
    // backtrack on failure ...
}
```

Each `r` is one **Node** representing a full guess `(RowVal, ColIdx, Digit)`.

---

## The 10 decisions that lead to (0,0) = 3

Tracing the actual search on this puzzle (same logic as `search()` in `main.go`):

| Step | Column chosen | Only option | Meaning |
|------|---------------|-------------|---------|
| 1 | `cell(5,1)` | put **6** at (5,1) | Row 5, col 1 forced |
| 2 | `cell(5,3)` | put **7** at (5,3) | forced |
| 3 | `cell(5,6)` | put **9** at (5,6) | forced |
| 4 | `row1-digit3` | put **3** at (1,3) | Row 1 must get a 3 in col 3 |
| 5 | `cell(0,3)` | put **2** at (0,3) | Top row, col 3 forced |
| 6 | `cell(2,5)` | put **7** at (2,5) | forced |
| 7 | `cell(0,5)` | put **5** at (0,5) | Top row, col 5 forced |
| 8 | `cell(0,6)` | put **7** at (0,6) | Top row, col 6 forced |
| 9 | `cell(0,7)` | put **8** at (0,7) | Top row, col 7 forced |
| 10 | `cell(0,0)` | put **3** at (0,0) | **First empty cell resolved!** |

After step 9, row 0 looks like:

```text
. 9 6 2 4 5 7 8 1
```

Only one digit is missing from row 0: **3**. Only one empty cell remains in row 0: **(0,0)**. So column `cell(0,0)` drops to **Size = 1**, and the solver picks **3**.

### Board snapshot when (0,0) becomes 3

```text
3 9 6 2 4 5 7 8 1   ← row 0 complete!
1 . . 3 6 . . . 4
5 . 4 8 1 7 3 9 .
. . 7 9 5 . . 4 3
. 3 . . 8 . . . .
4 6 5 7 2 3 9 1 8
. 1 . 6 3 . . 5 9
. 5 9 . 7 . 8 3 .
. . 3 5 9 . . . 7
```

**Notice:** the first **visual** empty cell is filled only after nine other forced placements — but none of those were random guesses. Every step up to Size = 1 was **forced** (only one candidate left in that column).

---

## Zoom in: what happens at step 10 (`put 3 at (0,0)`)

### 1. `selectColumn` returns column 0 (`cell(0,0)`), Size = 1

### 2. `cover(col)` hides column 0 and all conflicting rows

### 3. The loop finds one node: `RowVal=0, ColIdx=0, Digit=3`

That node is appended to `solution`.

### 4. Cover the other three columns in its horizontal ring

Walking `j := r.Right; j != r; j = j.Right` covers:

| Column ID | Rule satisfied |
|-----------|----------------|
| 83 | Row 0 contains digit 3 |
| 164 | Column 0 contains digit 3 |
| 245 | Top-left box contains digit 3 |

This removes every other guess that would put **3** in row 0, column 0, or the top-left box — or that would fill cell (0,0) with a different digit.

### 5. Recursive `search()` continues

We stop our walkthrough here. In the real run, search keeps going until all 324 columns are covered, then writes every node in `solution` back to the board:

```go
for _, node := range solution {
    board[node.RowVal][node.ColIdx] = node.Digit
}
```

---

## What if step 10 were wrong? (Backtracking in one sentence)

If a later recursive call failed, the solver would:

1. Remove the last node from `solution`
2. `uncover` the three columns from step 4 (in reverse order)
3. Try the **next** row in column 0 — but Size was 1, so there is no next row
4. `uncover` column 0 and backtrack further up

That is Algorithm X: **choose → cover → recurse → undo on failure**.

For this puzzle, **3** at (0,0) is correct and the search continues successfully.

---

## Map: concepts → code

| Concept | Where in `main.go` |
|---------|-------------------|
| Parse `.` and digits | `createBoard` |
| 324 column headers | `columns := make([]*Column, 324)` |
| 729 candidates | triple loop `for r`, `for c`, `for v := 1; v <= 9` |
| Four nodes per candidate | `nodes := make([]*Node, 4)` + horizontal linking |
| Column index formulas | `c1`, `c2`, `c3`, `c4` |
| Remove givens | loop with `cover(node.Col)` |
| Pick hardest column | `selectColumn` |
| Hide column + conflicts | `cover` |
| Restore on backtrack | `uncover` |
| Recursive search | `search` closure |
| Write answer to board | loop over `solution` at end of `solveExactCover` |

---

## Key beginner takeaways

1. **You see 81 cells; the solver tracks 324 rules.** Each placement must satisfy cell + row + column + box constraints at once.

2. **Each guess is a row of 4 linked nodes**, not a single number in an array.

3. **Givens are pre-committed** by covering their rows before search starts — the solver will not contradict input clues.

4. **Search order ≠ reading order.** MRV picks the tightest column first (often “naked singles” elsewhere on the board).

5. **Cell (0,0) had four candidates (2,3,7,8) after covering givens**, but other forced moves in row 0 narrowed it to **3** by the time the solver reached column `cell(0,0)`.

6. **Dancing Links** means never copying a big matrix — only rewiring `Left`, `Right`, `Up`, `Down` pointers in `cover` / `uncover`.

---

## Try it yourself

Run the same puzzle:

```bash
go run . ".96.4...1" "1...6...4" "5.481.39." "..795..43" ".3..8...." "4.5.23.18" ".1.63..59" ".59.7.83." "..359...7"
```

Run tests:

```bash
go test -v -run TestSolveExactCover .
```

For the full concept guide (linked lists, nodes, pointers), see [README.md — How It Works](README.md#how-it-works).
