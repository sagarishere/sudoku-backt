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
// Defaults to "exact-cover" if settings.txt is missing, unreadable, or invalid.
func readSettings() string {
	data, err := os.ReadFile("settings.txt")
	if err != nil {
		return "exact-cover"
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
	algo := s[start:end]
	// Validate the algorithm. If it's invalid, default to exact-cover.
	if algo == "backtracking" || algo == "exact-cover" || algo == "algo-x" || algo == "bitmask" || algo == "tdoku" || algo == "simd-tdoku" {
		return algo
	}
	return "exact-cover"
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
	} else if algo == "tdoku" || algo == "simd-tdoku" {
		solved = sudoku.SolveTdoku(&board)
	} else {
		solved = sudoku.SolveBacktracking(&board)
	}

	if solved {
		sudoku.PrintBoard(board)
		fmt.Println()
	} else {
		fmt.Println("Error")
	}
}
