package sudoku

// Returns the indices of the next sudoku cell
// Moving first left to right (x-direction, on top-most row)
// Then top to bottom (y-direction)
func NextCell(currentY, currentX int) (int, int) {

	// Use %9 to handle new row changes
	nextY, nextX := currentY, (currentX+1)%9
	if nextX == 0 {
		nextY = currentY + 1
	}
	return nextY, nextX
}
