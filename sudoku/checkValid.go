package sudoku

// Takes an input board, along with row and column position, as well as a cell value
// Performs a check of a newly placed number
// Checks: Row; Column; Box
func CheckValid(board [9][9]int, rowPosition int, columnPosition int, value int) bool {

	sizeSudoku := len(board)

	board[rowPosition][columnPosition] = value

	// Row check
	for i := 0; i < sizeSudoku; i++ {
		if board[rowPosition][i] == value && i != columnPosition {
			return false
		}
	}

	// Column check
	for j := 0; j < sizeSudoku; j++ {
		if board[j][columnPosition] == value && j != rowPosition {
			return false
		}
	}

	// Box check
	boxStartRow := (rowPosition / 3) * 3    // Reduces to row multiple 3
	boxStartCol := (columnPosition / 3) * 3 // reduces to column multiple 3
	for k := boxStartRow; k < boxStartRow+3; k++ {
		for l := boxStartCol; l < boxStartCol+3; l++ {
			if board[k][l] == value && k != rowPosition && l != columnPosition {
				return false
			}
		}
	}
	return true
}
