package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// --- helpers for unit tests (mirror logic from main where needed) ---

func resetBoards() {
	grid = [9][9]int{}
	solved = [9][9]int{}
}

// setGrid loads 9 strings of length 9 ('.' or '0' = empty, '1'-'9' = clue).
func setGrid(rows []string) {
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			ch := rows[r][c]
			if ch == '.' || ch == '0' {
				grid[r][c] = 0
			} else {
				grid[r][c] = int(ch - '0')
			}
		}
	}
}

// cluesAreValid uses the same check as main: each given digit must not conflict with others.
func cluesAreValid() bool {
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if grid[r][c] != 0 {
				num := grid[r][c]
				grid[r][c] = 0
				ok := isValid(r, c, num)
				grid[r][c] = num
				if !ok {
					return false
				}
			}
		}
	}
	return true
}

func countSolutions() int {
	solutions := 0
	solve(&solutions)
	return solutions
}

// boardToRows turns a 9×9 board into puzzle strings; listed cells are left empty.
func boardToRows(b [9][9]int, emptyCells [][2]int) []string {
	empty := make(map[[2]int]bool)
	for _, p := range emptyCells {
		empty[p] = true
	}
	rows := make([]string, 9)
	for r := 0; r < 9; r++ {
		chars := make([]byte, 9)
		for c := 0; c < 9; c++ {
			if empty[[2]int{r, c}] {
				chars[c] = '.'
			} else {
				chars[c] = byte('0' + b[r][c])
			}
		}
		rows[r] = string(chars)
	}
	return rows
}

func TestAllScenarios(t *testing.T) {
	// 1. Build the binary
	tmpDir, err := os.MkdirTemp("", "sudoku-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	binPath := filepath.Join(tmpDir, "sudoku_bin")
	buildCmd := exec.Command("go", "build", "-o", binPath, ".")
	buildCmd.Dir = "."
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build sudoku binary: %v", err)
	}

	// Define the test scenarios
	type testCase struct {
		args     []string
		expected string
	}

	testCases := []testCase{
		// 1. Valid Sudoku 1
		{
			args: []string{".96.4...1", "1...6...4", "5.481.39.", "..795..43", ".3..8....", "4.5.23.18", ".1.63..59", ".59.7.83.", "..359...7"},
			expected: `3 9 6 2 4 5 7 8 1
1 7 8 3 6 9 5 2 4
5 2 4 8 1 7 3 9 6
2 8 7 9 5 1 6 4 3
9 3 1 4 8 6 2 7 5
4 6 5 7 2 3 9 1 8
7 1 2 6 3 8 4 5 9
6 5 9 1 7 4 8 3 2
8 4 3 5 9 2 1 6 7

`,
		},
		// 2. Valid Sudoku 2
		{
			args: []string{"1.58.2...", ".9..764.5", "2..4..819", ".19..73.6", "762.83.9.", "....61.5.", "..76...3.", "43..2.5.1", "6..3.89.."},
			expected: `1 4 5 8 9 2 6 7 3
8 9 3 1 7 6 4 2 5
2 7 6 4 3 5 8 1 9
5 1 9 2 4 7 3 8 6
7 6 2 5 8 3 1 9 4
3 8 4 9 6 1 7 5 2
9 5 7 6 1 4 2 3 8
4 3 8 7 2 9 5 6 1
6 2 1 3 5 8 9 4 7

`,
		},
		// 3. Valid Sudoku 3
		{
			args: []string{"..5.3..81", "9.285..6.", "6....4.5.", "..74.283.", "34976...5", "..83..49.", "15..87..2", ".9....6..", ".26.495.3"},
			expected: `4 7 5 9 3 6 2 8 1
9 3 2 8 5 1 7 6 4
6 8 1 2 7 4 3 5 9
5 1 7 4 9 2 8 3 6
3 4 9 7 6 8 1 2 5
2 6 8 3 1 5 4 9 7
1 5 3 6 8 7 9 4 2
7 9 4 5 2 3 6 1 8
8 2 6 1 4 9 5 7 3

`,
		},
		// 4. Valid Sudoku 4
		{
			args: []string{"34.91..2.", ".96.8..41", "..8.2..7.", ".6..57.39", "1.2.6.7..", "97..3..64", "45.2.8..6", ".8..9..5.", "6.3..189."},
			expected: `3 4 7 9 1 5 6 2 8
2 9 6 7 8 3 5 4 1
5 1 8 6 2 4 9 7 3
8 6 4 1 5 7 2 3 9
1 3 2 4 6 9 7 8 5
9 7 5 8 3 2 1 6 4
4 5 9 2 7 8 3 1 6
7 8 1 3 9 6 4 5 2
6 2 3 5 4 1 8 9 7

`,
		},
		// 5. Valid Sudoku 5
		{
			args: []string{"..73..4.5", "....2.9..", "253.6487.", ".9.74.36.", "....3..8.", "8362.9.47", "1..8.26.3", "6......18", ".8261...4"},
			expected: `9 6 7 3 8 1 4 2 5
4 1 8 5 2 7 9 3 6
2 5 3 9 6 4 8 7 1
5 9 1 7 4 8 3 6 2
7 2 4 1 3 6 5 8 9
8 3 6 2 5 9 1 4 7
1 4 9 8 7 2 6 5 3
6 7 5 4 9 3 2 1 8
3 8 2 6 1 5 7 9 4

`,
		},
		// 6. Valid Sudoku 6
		{
			args: []string{"935..7..8", "...3.8.7.", "6..5..49.", ".73..4...", "4..175.8.", ".618..247", ".187.....", "..6.8.75.", "75.4.3862"},
			expected: `9 3 5 6 4 7 1 2 8
1 2 4 3 9 8 6 7 5
6 8 7 5 2 1 4 9 3
8 7 3 2 6 4 5 1 9
4 9 2 1 7 5 3 8 6
5 6 1 8 3 9 2 4 7
2 1 8 7 5 6 9 3 4
3 4 6 9 8 2 7 5 1
7 5 9 4 1 3 8 6 2

`,
		},
		// 7. Valid Sudoku 7
		{
			args: []string{"..5.2...1", ".8735..46", "4...6.5..", ".5.9.....", ".7..3541.", "69314.857", "7415..6.8", "...284..5", "5.....3.4"},
			expected: `3 6 5 4 2 9 7 8 1
2 8 7 3 5 1 9 4 6
4 1 9 8 6 7 5 3 2
1 5 4 9 7 8 2 6 3
8 7 2 6 3 5 4 1 9
6 9 3 1 4 2 8 5 7
7 4 1 5 9 3 6 2 8
9 3 6 2 8 4 1 7 5
5 2 8 7 1 6 3 9 4

`,
		},
		// 8. Valid Sudoku 8
		{
			args: []string{"..75...3.", "8..23...9", ".3479.86.", "..3..4198", ".4815...3", "..6.23..7", "351.6.78.", "4..31...6", ".7...5..2"},
			expected: `9 6 7 5 4 8 2 3 1
8 1 5 2 3 6 4 7 9
2 3 4 7 9 1 8 6 5
5 2 3 6 7 4 1 9 8
7 4 8 1 5 9 6 2 3
1 9 6 8 2 3 5 4 7
3 5 1 9 6 2 7 8 4
4 8 2 3 1 7 9 5 6
6 7 9 4 8 5 3 1 2

`,
		},
		// 9. Valid Sudoku 9
		{
			args: []string{"53..7....", "6..195...", ".98....6.", "8...6...3", "4..8.3..1", "7...2...6", ".6....28.", "...419..5", "....8..79"},
			expected: `5 3 4 6 7 8 9 1 2
6 7 2 1 9 5 3 4 8
1 9 8 3 4 2 5 6 7
8 5 9 7 6 1 4 2 3
4 2 6 8 5 3 7 9 1
7 1 3 9 2 4 8 5 6
9 6 1 5 3 7 2 8 4
2 8 7 4 1 9 6 3 5
3 4 5 2 8 6 1 7 9

`,
		},
		// 10. Valid Sudoku 10
		{
			args: []string{".58..4.21", ".6.853..7", ".39.2...5", "8....1..6", "..37..21.", "1.6.825..", "67.2..18.", "9..4...5.", ".8.9167.2"},
			expected: `7 5 8 6 9 4 3 2 1
2 6 1 8 5 3 4 9 7
4 3 9 1 2 7 8 6 5
8 2 7 5 4 1 9 3 6
5 4 3 7 6 9 2 1 8
1 9 6 3 8 2 5 7 4
6 7 4 2 3 5 1 8 9
9 1 2 4 7 8 6 5 3
3 8 5 9 1 6 7 4 2

`,
		},
		// 11. Valid Sudoku 11
		{
			args: []string{"71.4.9..2", ".8.5.....", "9...3..1.", "839..21.4", "..7.4.2..", "4.13..795", ".5..7...8", ".....5.3.", "6..1.3.47"},
			expected: `7 1 6 4 8 9 3 5 2
3 8 2 5 1 6 4 7 9
9 4 5 2 3 7 8 1 6
8 3 9 7 5 2 1 6 4
5 6 7 9 4 1 2 8 3
4 2 1 3 6 8 7 9 5
1 5 3 6 7 4 9 2 8
2 7 4 8 9 5 6 3 1
6 9 8 1 2 3 5 4 7

`,
		},
		// 12. Invalid Sudoku 1
		{
			args:     []string{".932..8.", "27.3.85..", ".8.73.254", "9758...31", "....74.6.", "6.45.38.7", "7....2.48", "32.4...7.", "..8.579.."},
			expected: "Error\n",
		},
		// 13. Invalid Sudoku 2
		{
			args:     []string{".867.2..4", ".2.5..8..", "154.9.237", ".7.9.5..1", ".29..4.18", "51.6...42", "2.5.7..83", "...153...", "39...8.75"},
			expected: "Error\n",
		},
		// 14. Invalid Sudoku 3
		{
			args:     []string{".7....28.", ".2...6.57", "8654729..", "..925..64", ".4..19.7.", "7.8..4..9", "3..7..698", "..79.1...", "59..28.39"},
			expected: "Error\n",
		},
		// 15. Invalid Sudoku 4
		{
			args:     []string{"..213.748", "8.4.....2", ".178.26..", ".68.9.27.", ".932....4", "5..46.3..", "..9.24.23", "..63..19.", "385..1.2."},
			expected: "Error\n",
		},
		// 16. Invalid Sudoku 5
		{
			args:     []string{"9.46.3..1", "37.1..2.6", "..6..93.4", "..13..9.5", "56..91...", "82...461.", "..79...4.", "425.167..", "1.2..75.8"},
			expected: "Error\n",
		},
		// 17. Invalid arguments (not a sudoku)
		{
			args:     []string{"not", "a", "sudoku"},
			expected: "Error\n",
		},
		// 18. Invalid Sudoku 6 (incorrect row count)
		{
			args:     []string{"53..8294.", "8..34...5", "3542761..", "..6.3...4", "9....162.", ".9...3.78", "7438.9...", "..5..43.1"},
			expected: "Error\n",
		},
	}

	for i, tc := range testCases {
		cmd := exec.Command(binPath, tc.args...)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		_ = cmd.Run()

		outStr := stdout.String()
		outStr = strings.ReplaceAll(outStr, "\r\n", "\n")
		expected := strings.ReplaceAll(tc.expected, "\r\n", "\n")
		// Ignore trailing blank line (old solver printed an extra newline after the grid)
		outStr = strings.TrimRight(outStr, "\n")
		expected = strings.TrimRight(expected, "\n")

		if outStr != expected {
			t.Errorf("Case %d failed.\nArgs: %v\nExpected:\n%q\nGot:\n%q\nStderr: %s", i+1, tc.args, expected, outStr, stderr.String())
		}
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		name     string
		rows     []string
		row, col int
		num      int
		want     bool
	}{
		{
			name: "empty board allows digit in corner",
			rows: stringsRepeat(".........", 9),
			row:  0, col: 0, num: 5,
			want: true,
		},
		{
			name: "rejects duplicate in same row",
			rows: []string{
				"5........",
				".........", ".........", ".........", ".........",
				".........", ".........", ".........", ".........",
			},
			row: 0, col: 1, num: 5,
			want: false,
		},
		{
			name: "rejects duplicate in same column",
			rows: []string{
				"5........",
				".........", ".........", ".........", ".........",
				".........", ".........", ".........", ".........",
			},
			row: 1, col: 0, num: 5,
			want: false,
		},
		{
			name: "rejects duplicate in same 3x3 box",
			rows: []string{
				"5........",
				".5.......",
				".........", ".........", ".........", ".........",
				".........", ".........", ".........", ".........",
			},
			row: 2, col: 2, num: 5,
			want: false,
		},
		{
			name: "allows digit when row col and box are clear",
			rows: []string{
				".96.4...1",
				"1...6...4",
				"5.481.39.",
				"..795..43",
				".3..8....",
				"4.5.23.18",
				".1.63..59",
				".59.7.83.",
				"..359...7",
			},
			row: 0, col: 0, num: 3,
			want: true,
		},
		{
			name: "rejects 3 at top-left when puzzle clues block it",
			rows: []string{
				".96.4...1",
				"1...6...4",
				"5.481.39.",
				"..795..43",
				".3..8....",
				"4.5.23.18",
				".1.63..59",
				".59.7.83.",
				"..359...7",
			},
			row: 0, col: 0, num: 9,
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resetBoards()
			setGrid(tc.rows)
			got := isValid(tc.row, tc.col, tc.num)
			if got != tc.want {
				t.Errorf("isValid(%d, %d, %d) = %v, want %v", tc.row, tc.col, tc.num, got, tc.want)
			}
		})
	}
}

func TestCluesAreValid(t *testing.T) {
	tests := []struct {
		name string
		rows []string
		want bool
	}{
		{
			name: "valid partial puzzle",
			rows: []string{
				".96.4...1",
				"1...6...4",
				"5.481.39.",
				"..795..43",
				".3..8....",
				"4.5.23.18",
				".1.63..59",
				".59.7.83.",
				"..359...7",
			},
			want: true,
		},
		{
			name: "two 5s in same row",
			rows: []string{
				"55.......",
				".........", ".........", ".........", ".........",
				".........", ".........", ".........", ".........",
			},
			want: false,
		},
		{
			name: "two 1s in same column",
			rows: []string{
				"1........",
				"1........",
				".........", ".........", ".........",
				".........", ".........", ".........", ".........",
			},
			want: false,
		},
		{
			name: "two 7s in same box",
			rows: []string{
				"7........",
				".7.......",
				".........", ".........", ".........",
				".........", ".........", ".........", ".........",
			},
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resetBoards()
			setGrid(tc.rows)
			if got := cluesAreValid(); got != tc.want {
				t.Errorf("cluesAreValid() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestSolve(t *testing.T) {
	expectedCase1 := [9][9]int{
		{3, 9, 6, 2, 4, 5, 7, 8, 1},
		{1, 7, 8, 3, 6, 9, 5, 2, 4},
		{5, 2, 4, 8, 1, 7, 3, 9, 6},
		{2, 8, 7, 9, 5, 1, 6, 4, 3},
		{9, 3, 1, 4, 8, 6, 2, 7, 5},
		{4, 6, 5, 7, 2, 3, 9, 1, 8},
		{7, 1, 2, 6, 3, 8, 4, 5, 9},
		{6, 5, 9, 1, 7, 4, 8, 3, 2},
		{8, 4, 3, 5, 9, 2, 1, 6, 7},
	}

	tests := []struct {
		name       string
		rows       []string
		wantCount  int
		wantSolved [9][9]int // zero value means do not check solved grid
	}{
		{
			name: "unique solution (integration case 1)",
			rows: []string{
				".96.4...1",
				"1...6...4",
				"5.481.39.",
				"..795..43",
				".3..8....",
				"4.5.23.18",
				".1.63..59",
				".59.7.83.",
				"..359...7",
			},
			wantCount:  1,
			wantSolved: expectedCase1,
		},
		{
			name:       "one empty cell in complete grid",
			rows:       boardToRows(expectedCase1, [][2]int{{0, 0}}),
			wantCount:  1,
			wantSolved: expectedCase1,
		},
		{
			name:      "empty grid has multiple solutions",
			rows:      stringsRepeat(".........", 9),
			wantCount: 2, // solve stops after finding a second solution
		},
		{
			name: "no solution (integration case 16)",
			rows: []string{
				"9.46.3..1",
				"37.1..2.6",
				"..6..93.4",
				"..13..9.5",
				"56..91...",
				"82...461.",
				"..79...4.",
				"425.167..",
				"1.2..75.8",
			},
			wantCount: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resetBoards()
			setGrid(tc.rows)
			got := countSolutions()
			if got != tc.wantCount {
				t.Fatalf("countSolutions() = %d, want %d", got, tc.wantCount)
			}
			if tc.wantCount == 1 && tc.wantSolved[0][0] != 0 {
				for r := 0; r < 9; r++ {
					for c := 0; c < 9; c++ {
						if solved[r][c] != tc.wantSolved[r][c] {
							t.Errorf("solved[%d][%d] = %d, want %d", r, c, solved[r][c], tc.wantSolved[r][c])
						}
					}
				}
			}
		})
	}
}

// stringsRepeat builds n copies of s (helper for empty-grid tests).
func stringsRepeat(s string, n int) []string {
	out := make([]string, n)
	for i := range out {
		out[i] = s
	}
	return out
}
