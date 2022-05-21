package sudoku

import "fmt"

// Prints an input 'board' of integers (array)
// A space is printed between each value
func PrintBoard(board [9][9]int) {
	sudokuSize := len(board)

	// Assumes complete rectangular board, 9 x 9
	// Cycle through rows
	for i := 0; i < sudokuSize; i++ {
		// Cycle through columns
		for j := 0; j < sudokuSize; j++ {
			if board[i][j] < 1 && j != sudokuSize-1 {
				fmt.Print(board[i][j])
				fmt.Print(" ")
			} else if board[i][j] < 1 && j == sudokuSize-1 {
				fmt.Print(board[i][j])
				fmt.Print("\n")
			} else if board[i][j] >= 1 && j != sudokuSize-1 {
				fmt.Print(board[i][j])
				fmt.Print(" ")
			} else {
				fmt.Print(board[i][j])
				fmt.Print("\n")
			}
		}
	}
}
