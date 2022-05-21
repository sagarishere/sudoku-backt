package main

import (
	"fmt"
	"os"
	"sudoku/sudoku"
)

// Establish Global Variables
var inputBoard = os.Args[1:]
var board, validInput = sudoku.CreateBoard(inputBoard)

// An algorithm which solves a given sudoku puzzle using backtracking
func recursiveSolve(rowPosition, columnPosition int) bool {

	size := len(board)

	// End condition which should be recursively reached if solution found.
	// i.e. Finishes 9th row, moves to 10th row (non-existent)
	if rowPosition == 9 {
		return true
	}
	// Move to next cell if current cell already filled in
	if board[rowPosition][columnPosition] != 0 {
		return recursiveSolve(sudoku.NextCell(rowPosition, columnPosition))
	} else {
		for i := 1; i <= size; i++ {
			if sudoku.CheckValid(board, rowPosition, columnPosition, i) == true {
				board[rowPosition][columnPosition] = i
				if recursiveSolve(sudoku.NextCell(rowPosition, columnPosition)) {
					return true
				}
				board[rowPosition][columnPosition] = 0
			}
		}
		return false
	}
}

// See below for inspiration
// INSPIRATION: https://charltonaustin.com/posts/sudoku-using-go-lang/
// INSPIRATION: https://www.geeksforgeeks.org/sudoku-backtracking-7/
// INSPIRATION: https://www.5minsofcode.com/sodoku_solver.html
func main() {

	canProceed := true

	// Check starting board validity according to minimum number requirements
	if sudoku.StartValid(board) == false {
		canProceed = false
		fmt.Printf("Error: Input configuration is not valid.")
	} else if validInput == false {
		fmt.Printf("Error: Incorrect input - string cannot be read according to standard 9 x 9 dimensions")
	}
	// Recursively iterate through board and print results
	if canProceed {
		fmt.Println()
		fmt.Println("Initial sudoku board shown below:\n")
		sudoku.PrintBoard(board)
		if recursiveSolve(0, 0) {
			fmt.Println()
			fmt.Println("The following solution was found:\n")
			sudoku.PrintBoard(board)
			fmt.Println()
		} else {
			fmt.Println("\nA solution for this start configuration does not exist.")
			fmt.Println()
		}
	}
}
