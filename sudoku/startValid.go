package sudoku

// Checks if input has the minimum requirements to produce a unique solution
// - A minimum of 17 numbers
// - Of which, 8 must be unique
func StartValid(inputBoard [9][9]int) bool {
	uniqueNumSlice := make([]int, 9)
	uniqueNumCount := 0
	numberCount := 0
	canContinue := true

	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if inputBoard[i][j] >= 1 && inputBoard[i][j] <= 9 {
				numberCount++
				if uniqueNumSlice[inputBoard[i][j]-1] < 1 {
					uniqueNumSlice[inputBoard[i][j]-1]++
					uniqueNumCount++
				}
			}
		}
	}
	if uniqueNumCount < 8 || numberCount < 17 || numberCount > 77 {
		canContinue = false
	}

	if canContinue == true {
		for k := 0; k < 9; k++ {
			for l := 0; l < 9; l++ {
				if inputBoard[k][l] != 0 {
					// Row check
					for m := 0; m < 9; m++ {
						if inputBoard[k][m] == inputBoard[k][l] && m != l {
							return false
						}
					}
					// Column check
					for n := 0; n < 9; n++ {
						if inputBoard[n][l] == inputBoard[k][l] && n != k {
							return false
						}
					}
					// Box check
					boxStartRow := (k / 3) * 3 // Reduces to row multiple 3
					boxStartCol := (l / 3) * 3 // reduces to column multiple 3
					for p := boxStartRow; p < boxStartRow+3; p++ {
						for q := boxStartCol; q < boxStartCol+3; q++ {
							if inputBoard[p][q] == inputBoard[k][l] && p != k && q != l {
								return false
							}
						}
					}
				}
			}
		}
	}
	return canContinue
}
