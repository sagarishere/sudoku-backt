package main

import (
    "fmt"
    "os"
)

var grid [9][9]int
var solved [9][9]int

func isValid(row, col, num int) bool {
    for c := 0; c < 9; c++ {
        if grid[row][c] == num {
            return false
        }
    }
    for r := 0; r < 9; r++ {
        if grid[r][col] == num {
            return false
        }
    }
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

func solve(solutions *int) {
    if *solutions > 1 {
        return
    }
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
    if row == -1 {
        *solutions++
        if *solutions == 1 {
            solved = grid
        }
        return
}
    for num := 1; num <= 9; num++ {
        if isValid(row, col, num) {
            grid[row][col] = num
            solve(solutions)
            grid[row][col] = 0
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
    for r, arg := range args {
        if len(arg) != 9 {
            printError()
        }
        for c, ch := range arg {
            if ch == '.' {
                grid[r][c] = 0
            } else if ch >= '1' && ch <= '9' {
                grid[r][c] = int(ch - '0')
            } else {
                printError()
            }
        }
    }
    for r := 0; r < 9; r++ {
        for c := 0; c < 9; c++ {
            if grid[r][c] != 0 {
                num := grid[r][c]
                grid[r][c] = 0
                if !isValid(r, c, num) {
                    printError()
                }
                grid[r][c] = num
            }
        }
    }
    solutions := 0
    solve(&solutions)
    if solutions != 1 {
        printError()
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