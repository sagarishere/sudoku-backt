package sudoku

// Creates a slice of slices (integers), where each slice is a row
func CreateBoard(startCondition []string) ([9][9]int, bool) {

	sudokuSize := len(startCondition)
	validCreate := true

	// Initialise board with empty row slices
	var startBoard = [9][9]int{}
	if sudokuSize != 9 {
		validCreate = false
	}

	// Fill row slices with columns and fill with start values
	if validCreate == true {
		for i := 0; i < sudokuSize; i++ {
			if len(startCondition[i]) != 9 {
				validCreate = false
				break
			}
			for j := 0; j < sudokuSize; j++ {
				startBoard[i][j] = 0
				if startCondition[i][j] >= '1' && startCondition[i][j] <= '9' {
					startBoard[i][j] = int(startCondition[i][j] - 48)
				}
			}
		}
	}
	return startBoard, validCreate
}
