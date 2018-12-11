package main

import "fmt"

func problem11() {
	grid := make([][]int, 300)
	for r := range grid {
		y := r + 1
		row := make([]int, 300)
		for c := range row {
			x := c + 1
			row[c] = score(x, y)
		}

		grid[r] = row
	}

	findLargest3x3(grid)
	findLargestSquare(grid)
}

func score(x, y int) int {
	serialNumber := 1955
	id := x + 10

	power := (id*y + serialNumber) * id
	power /= 100
	power = power % 10
	return power - 5
}

func findLargest3x3(grid [][]int) {
	var xTopLeft, yTopLeft int
	size := 3
	sum := 0

	for x := 1; x <= 298; x++ {
		for y := 1; y <= 298; y++ {
			windowSum := findWindowSum(grid, x, y, size)
			if windowSum > sum {
				sum = windowSum
				xTopLeft = x
				yTopLeft = y
			}
		}
	}

	fmt.Printf("Part 1: %d,%d\n", xTopLeft, yTopLeft)
}

func findLargestSquare(grid [][]int) {
	var xTopLeft, yTopLeft, sBest int
	sum := 0

	for s := 1; s <= 300; s++ {
		for x := 1; x <= 300-s+1; x++ {
			for y := 1; y <= 300-s+1; y++ {
				windowSum := findWindowSum(grid, x, y, s)
				if windowSum > sum {
					sum = windowSum
					xTopLeft = x
					yTopLeft = y
					sBest = s
				}
			}
		}
	}

	fmt.Printf("Part 2: %d,%d,%d\n", xTopLeft, yTopLeft, sBest)
}

func findWindowSum(grid [][]int, x, y, size int) int {
	sum := 0
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			sum += grid[y+i-1][x+j-1]
		}
	}
	return sum
}
