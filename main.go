package main

import (
	"fmt"
	"os"
	"sudoku/sudoku"
)

// Establish Global Variables
var board [9][9]int
var validInput bool

// readSettings reads the solver configuration from settings.txt.
// Defaults to "backtracking" if settings.txt is missing or unreadable.
func readSettings() string {
	data, err := os.ReadFile("settings.txt")
	if err != nil {
		return "backtracking"
	}
	s := string(data)
	// Custom space-trimming to avoid importing the "strings" package
	start := 0
	for start < len(s) && (s[start] == ' ' || s[start] == '\n' || s[start] == '\r' || s[start] == '\t') {
		start++
	}
	end := len(s)
	for end > start && (s[end-1] == ' ' || s[end-1] == '\n' || s[end-1] == '\r' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}

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
	algo := readSettings()

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

	var solved bool
	if algo == "exact-cover" || algo == "algo-x" {
		solved = sudoku.SolveExactCover(&board)
	} else if algo == "bitmask" {
		solved = sudoku.SolveBitmask(&board)
	} else {
		solved = recursiveSolve(0, 0)
	}

	if solved {
		sudoku.PrintBoard(board)
		fmt.Println()
	} else {
		fmt.Println("Error")
	}
}
