package main

import (
	"fmt"
	"os"
)

// grid is the board we fill while searching. 0 means an empty cell.
var grid [9][9]int

// solved holds the first complete answer we find.
var solved [9][9]int

// isValid checks whether num can go at (row, col) without breaking Sudoku rules.
// The cell at (row, col) should be empty (0) when this is called.
func isValid(row, col, num int) bool {
	// Same digit already in this row?
	for c := 0; c < 9; c++ {
		if grid[row][c] == num {
			return false
		}
	}
	// Same digit already in this column?
	for r := 0; r < 9; r++ {
		if grid[r][col] == num {
			return false
		}
	}
	// Same digit already in this 3×3 box?
	// Integer division maps row/col to the top-left corner of their box.
	startRow := (row / 3) * 3
	startCol := (col / 3) * 3
	for r := startRow; r < startRow+3; r++ {
		for c := startCol; c < startCol+3; c++ {
			if grid[r][c] == num {
				return false
			}
		}
	}
	return true
}

// solve fills empty cells using backtracking: try a digit, go deeper, undo if stuck.
// solutions counts how many complete grids we found (we only accept exactly one).
func solve(solutions *int) {
	if *solutions > 1 {
		return // not a unique puzzle; stop searching
	}

	// Find the next empty cell (scan left-to-right, top-to-bottom).
	row, col := -1, -1
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if grid[r][c] == 0 {
				row, col = r, c
				break
			}
		}
		if row != -1 {
			break
		}
	}

	// No empty cells left — the board is full.
	if row == -1 {
		*solutions++
		if *solutions == 1 {
			solved = grid // save the first (and hopefully only) answer
		}
		return
	}

	// Try digits 1–9 in this cell.
	for num := 1; num <= 9; num++ {
		if isValid(row, col, num) {
			grid[row][col] = num
			solve(solutions)
			grid[row][col] = 0 // undo (backtrack) before trying the next digit
		}
		if *solutions > 1 {
			return
		}
	}
}

func printError() {
	fmt.Println("Error")
	os.Exit(0)
}

func main() {
	args := os.Args[1:]
	if len(args) != 9 {
		printError()
	}

	// Read 9 rows from the command line into grid.
	for r, arg := range args {
		if len(arg) != 9 {
			printError()
		}
		for c, ch := range arg {
			if ch == '.' {
				grid[r][c] = 0
			} else if ch >= '1' && ch <= '9' {
				grid[r][c] = int(ch - '0') // '5' - '0' → 5
			} else {
				printError()
			}
		}
	}

	// Make sure the given clues do not already break any rule.
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if grid[r][c] != 0 {
				num := grid[r][c]
				grid[r][c] = 0 // pretend the cell is empty so isValid checks only other cells
				if !isValid(r, c, num) {
					printError()
				}
				grid[r][c] = num
			}
		}
	}

	solutions := 0
	solve(&solutions) // pass a pointer so every recursive call shares the same counter
	if solutions != 1 {
		printError() // 0 solutions or more than 1
	}

	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if c > 0 {
				fmt.Print(" ")
			}
			fmt.Print(solved[r][c])
		}
		fmt.Println()
	}
}
