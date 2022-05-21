package sudoku

// SolveBitmask solves the Sudoku puzzle using high-performance Bitmask Backtracking.
func SolveBitmask(board *[9][9]int) bool {
	// bitmasks to track numbers used in rows, columns, and 3x3 boxes.
	// For each, rowsUsed[r] stores the used numbers where the d-th bit is set to 1 if digit d is used.
	var rowsUsed [9]int
	var colsUsed [9]int
	var boxesUsed [9]int

	// Initialize the bitmasks from the starting board state
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if board[r][c] != 0 {
				v := board[r][c]
				mask := 1 << v
				b := (r/3)*3 + c/3
				rowsUsed[r] |= mask
				colsUsed[c] |= mask
				boxesUsed[b] |= mask
			}
		}
	}

	// Recursive backtracking depth-first search
	var solve func(r, c int) bool
	solve = func(r, c int) bool {
		// If we've successfully filled the 9th row (0-indexed), the board is solved
		if r == 9 {
			return true
		}

		// Calculate the indices for the next cell
		nextR, nextC := r, c+1
		if nextC == 9 {
			nextR = r + 1
			nextC = 0
		}

		// If the current cell is already filled, skip to the next cell
		if board[r][c] != 0 {
			return solve(nextR, nextC)
		}

		b := (r/3)*3 + c/3

		// Try placing numbers from 1 to 9
		for v := 1; v <= 9; v++ {
			mask := 1 << v

			// A simple bitwise check: verify if 'v' is unused in the current row, column, and box
			if (rowsUsed[r]&mask == 0) && (colsUsed[c]&mask == 0) && (boxesUsed[b]&mask == 0) {
				// Place the digit and update the bitmasks (bitwise OR)
				rowsUsed[r] |= mask
				colsUsed[c] |= mask
				boxesUsed[b] |= mask
				board[r][c] = v

				if solve(nextR, nextC) {
					return true
				}

				// Backtrack: reset the cell and clear the bitmasks (bitwise AND NOT)
				rowsUsed[r] &= ^mask
				colsUsed[c] &= ^mask
				boxesUsed[b] &= ^mask
				board[r][c] = 0
			}
		}

		return false
	}

	return solve(0, 0)
}
