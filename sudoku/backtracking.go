package sudoku

// SolveBacktracking solves the Sudoku puzzle using traditional grid backtracking.
func SolveBacktracking(board *[9][9]int) bool {
	var solve func(row, col int) bool
	solve = func(row, col int) bool {
		if row == 9 {
			return true
		}
		nextR, nextC := NextCell(row, col)
		if board[row][col] != 0 {
			return solve(nextR, nextC)
		}
		for i := 1; i <= 9; i++ {
			if CheckValid(*board, row, col, i) {
				board[row][col] = i
				if solve(nextR, nextC) {
					return true
				}
				board[row][col] = 0
			}
		}
		return false
	}
	return solve(0, 0)
}
