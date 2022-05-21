package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"sudoku/sudoku"
)

var puzzleStrings = [][]string{
	// Easy (Valid Sudoku 1)
	{".96.4...1", "1...6...4", "5.481.39.", "..795..43", ".3..8....", "4.5.23.18", ".1.63..59", ".59.7.83.", "..359...7"},
	// Medium (Valid Sudoku 2)
	{"1.58.2...", ".9..764.5", "2..4..819", ".19..73.6", "762.83.9.", "....61.5.", "..76...3.", "43..2.5.1", "6..3.89.."},
	// Hard (Valid Sudoku 3)
	{"..5.3..81", "9.285..6.", "6....4.5.", "..74.283.", "34976...5", "..83..49.", "15..87..2", ".9....6..", ".26.495.3"},
}

func runBenchmark(solveFunc func(*[9][9]int) bool, puzzles [][9][9]int, N int) float64 {
	// Warm up
	for _, p := range puzzles {
		pCopy := p
		solveFunc(&pCopy)
	}

	start := time.Now()
	for i := 0; i < N; i++ {
		for _, p := range puzzles {
			pCopy := p
			solveFunc(&pCopy)
		}
	}
	elapsed := time.Since(start)
	totalRuns := N * len(puzzles)
	return float64(elapsed.Nanoseconds()) / 1000.0 / float64(totalRuns) // in microseconds
}

func main() {
	var puzzles [][9][9]int
	for _, ps := range puzzleStrings {
		board, ok := sudoku.CreateBoard(ps)
		if !ok {
			fmt.Printf("Failed to create board for: %v\n", ps)
			os.Exit(1)
		}
		puzzles = append(puzzles, board)
	}

	fmt.Println("Running benchmark suite (500 runs per engine)...")

	tBacktracking := runBenchmark(sudoku.SolveBacktracking, puzzles, 500)
	tExactCover := runBenchmark(sudoku.SolveExactCover, puzzles, 500)
	tBitmask := runBenchmark(sudoku.SolveBitmask, puzzles, 500)
	tTdoku := runBenchmark(sudoku.SolveTdoku, puzzles, 500)

	fmt.Printf("Backtracking:         %.2f μs\n", tBacktracking)
	fmt.Printf("Exact Cover (DLX):    %.2f μs\n", tExactCover)
	fmt.Printf("Bitmask Backtracking:  %.2f μs\n", tBitmask)
	fmt.Printf("SIMD-Optimized Tdoku: %.2f μs\n", tTdoku)

	// Max time for scaling
	maxTime := tBacktracking
	if tExactCover > maxTime {
		maxTime = tExactCover
	}
	if tBitmask > maxTime {
		maxTime = tBitmask
	}
	if tTdoku > maxTime {
		maxTime = tTdoku
	}

	// Avoid division by zero
	if maxTime == 0 {
		maxTime = 1.0
	}

	// Calculate bar widths (maximum width = 400px)
	wBacktracking := int(math.Round((tBacktracking / maxTime) * 400.0))
	wExactCover := int(math.Round((tExactCover / maxTime) * 400.0))
	wBitmask := int(math.Round((tBitmask / maxTime) * 400.0))
	wTdoku := int(math.Round((tTdoku / maxTime) * 400.0))

	// Adjust minimum width of 4px to keep bars visible if extremely fast
	if wBacktracking < 4 {
		wBacktracking = 4
	}
	if wExactCover < 4 {
		wExactCover = 4
	}
	if wBitmask < 4 {
		wBitmask = 4
	}
	if wTdoku < 4 {
		wTdoku = 4
	}

	// Text positions
	valXBacktracking := 220 + wBacktracking + 10
	valXExactCover := 220 + wExactCover + 10
	valXBitmask := 220 + wBitmask + 10
	valXTdoku := 220 + wTdoku + 10

	// Generate SVG
	svgContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<svg width="700" height="350" viewBox="0 0 700 350" xmlns="http://www.w3.org/2000/svg">
  <defs>
    <linearGradient id="grad-backtracking" x1="0%%" y1="0%%" x2="100%%" y2="0%%">
      <stop offset="0%%" stop-color="#F38BA8" />
      <stop offset="100%%" stop-color="#E78284" />
    </linearGradient>
    <linearGradient id="grad-exactcover" x1="0%%" y1="0%%" x2="100%%" y2="0%%">
      <stop offset="0%%" stop-color="#89B4FA" />
      <stop offset="100%%" stop-color="#30D5C8" />
    </linearGradient>
    <linearGradient id="grad-bitmask" x1="0%%" y1="0%%" x2="100%%" y2="0%%">
      <stop offset="0%%" stop-color="#A6E3A1" />
      <stop offset="100%%" stop-color="#40A02B" />
    </linearGradient>
    <linearGradient id="grad-tdoku" x1="0%%" y1="0%%" x2="100%%" y2="0%%">
      <stop offset="0%%" stop-color="#CBA6F7" />
      <stop offset="100%%" stop-color="#8839EF" />
    </linearGradient>
  </defs>
  <style>
    .background { fill: #1e1e2e; }
    .border { stroke: #313244; stroke-width: 2; fill: none; }
    .bar { transition: width 0.8s ease-out; }
    .bar:hover { fill-opacity: 0.85; filter: drop-shadow(0px 0px 8px rgba(255, 255, 255, 0.15)); }
    .text-title { font-family: 'Outfit', 'Inter', system-ui, sans-serif; font-weight: 700; fill: #cdd6f4; font-size: 20px; }
    .text-subtitle { font-family: 'Outfit', 'Inter', system-ui, sans-serif; font-weight: 400; fill: #a6adc8; font-size: 13px; }
    .text-label { font-family: 'Outfit', 'Inter', system-ui, sans-serif; font-weight: 600; fill: #cdd6f4; font-size: 14px; }
    .text-desc { font-family: 'Outfit', 'Inter', system-ui, sans-serif; font-weight: 400; fill: #a6adc8; font-size: 11px; }
    .text-value { font-family: 'Outfit', 'Inter', system-ui, sans-serif; font-weight: 700; fill: #cdd6f4; font-size: 13px; }
    .grid-line { stroke: #313244; stroke-width: 1; stroke-dasharray: 4; }
  </style>
  <rect class="background" width="100%%" height="100%%" rx="12" />
  <rect class="border" width="100%%" height="100%%" rx="12" />

  <!-- Title -->
  <text x="30" y="40" class="text-title">Sudoku Solvers Performance Benchmark</text>
  <text x="30" y="60" class="text-subtitle">Average execution time in microseconds (μs) - Lower is better</text>

  <!-- Gridlines -->
  <line x1="220" y1="80" x2="220" y2="300" stroke="#45475a" stroke-width="2" />
  <line x1="320" y1="80" x2="320" y2="300" class="grid-line" />
  <line x1="420" y1="80" x2="420" y2="300" class="grid-line" />
  <line x1="520" y1="80" x2="520" y2="300" class="grid-line" />
  <line x1="620" y1="80" x2="620" y2="300" class="grid-line" />

  <!-- Bar 1: Backtracking -->
  <text x="30" y="112" class="text-label">Traditional Backtracking</text>
  <text x="30" y="128" class="text-desc">Grid DFS (Baseline)</text>
  <rect class="bar" x="220" y="95" width="%d" height="36" rx="6" fill="url(#grad-backtracking)" />
  <text x="%d" y="118" class="text-value">%.2f μs</text>

  <!-- Bar 2: Exact Cover -->
  <text x="30" y="167" class="text-label">Knuth Algorithm X</text>
  <text x="30" y="183" class="text-desc">Dancing Links (DLX)</text>
  <rect class="bar" x="220" y="150" width="%d" height="36" rx="6" fill="url(#grad-exactcover)" />
  <text x="%d" y="173" class="text-value">%.2f μs</text>

  <!-- Bar 3: Bitmask Backtracking -->
  <text x="30" y="222" class="text-label">Bitmask Backtracking</text>
  <text x="30" y="238" class="text-desc">CPU Register Masks</text>
  <rect class="bar" x="220" y="205" width="%d" height="36" rx="6" fill="url(#grad-bitmask)" />
  <text x="%d" y="228" class="text-value">%.2f μs</text>

  <!-- Bar 4: Tdoku -->
  <text x="30" y="277" class="text-label">SIMD-Optimized</text>
  <text x="30" y="293" class="text-desc">Tdoku-Inspired Bit-Triads</text>
  <rect class="bar" x="220" y="260" width="%d" height="36" rx="6" fill="url(#grad-tdoku)" />
  <text x="%d" y="283" class="text-value">%.2f μs</text>
</svg>
`,
		wBacktracking, valXBacktracking, tBacktracking,
		wExactCover, valXExactCover, tExactCover,
		wBitmask, valXBitmask, tBitmask,
		wTdoku, valXTdoku, tTdoku,
	)

	// Write SVG file
	err := os.WriteFile("benchmark.svg", []byte(svgContent), 0644)
	if err != nil {
		fmt.Printf("Failed to write benchmark.svg: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully wrote benchmark.svg!")

	// Update README.md
	readmeBytes, err := os.ReadFile("README.md")
	if err != nil {
		fmt.Printf("Failed to read README.md: %v\n", err)
		os.Exit(1)
	}

	readmeContent := string(readmeBytes)

	startTag := "<!-- BENCHMARK_START -->"
	endTag := "<!-- BENCHMARK_END -->"

	startIdx := strings.Index(readmeContent, startTag)
	endIdx := strings.Index(readmeContent, endTag)

	if startIdx == -1 || endIdx == -1 || startIdx > endIdx {
		fmt.Println("Warning: could not find benchmark template markers in README.md")
		os.Exit(1)
	}

	// Compute speedups
	sExactCover := tBacktracking / tExactCover
	sBitmask := tBacktracking / tBitmask
	sTdoku := tBacktracking / tTdoku

	newBenchmarkMarkdown := fmt.Sprintf(`<!-- BENCHMARK_START -->
### Benchmark Results

Here is a performance comparison of the four Sudoku solver engines compiled and executed in your local environment.

| Solver Engine | Average Solve Time (μs) | Speedup Factor |
| :--- | :---: | :---: |
| **Traditional Backtracking** | %.2f μs | Baseline (1.0x) |
| **Knuth's Algorithm X (DLX)** | %.2f μs | %.1fx faster |
| **Bitmask Backtracking** | %.2f μs | %.1fx faster |
| **SIMD-Optimized (Tdoku)** | %.2f μs | %.1fx faster |

![Benchmark Results](benchmark.svg)
<!-- BENCHMARK_END -->`,
		tBacktracking,
		tExactCover, sExactCover,
		tBitmask, sBitmask,
		tTdoku, sTdoku,
	)

	updatedReadme := readmeContent[:startIdx] + newBenchmarkMarkdown + readmeContent[endIdx+len(endTag):]

	err = os.WriteFile("README.md", []byte(updatedReadme), 0644)
	if err != nil {
		fmt.Printf("Failed to update README.md: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully updated README.md with benchmark results!")
}
