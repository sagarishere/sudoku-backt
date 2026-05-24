package main

import (
	"fmt"
	"os"
	"sudoku/sudoku"
)

// Establish Global Variables
var board [9][9]int
var validInput bool

// See below for inspiration
// INSPIRATION: https://charltonaustin.com/posts/sudoku-using-go-lang/
// INSPIRATION: https://www.geeksforgeeks.org/sudoku-backtracking-7/
// INSPIRATION: https://www.5minsofcode.com/sodoku_solver.html
func main() {
	inputBoard := os.Args[1:]
	var err bool
	board, err = sudoku.CreateBoard(inputBoard)
	validInput = err

	canProceed := true

	// Check starting board validity according to minimum number requirements
	if sudoku.StartValid(board) == false {
		canProceed = false
	} else if validInput == false {
		canProceed = false
	}

	if !canProceed {
		fmt.Println("Error")
		return
	}

	if sudoku.SolveExactCover(&board) {
		sudoku.PrintBoard(board)
		fmt.Println()
	} else {
		fmt.Println("Error")
	}
}
