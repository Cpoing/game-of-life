package main

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

const (
	width  = 20
	height = 10
)

type Grid [][]bool

func newGrid() Grid {
	grid := make(Grid, height)
	for i := range grid {
		grid[i] = make([]bool, width)
	}
	return grid
}

func randomize(grid Grid) {
	for y := range grid {
		for x := range grid[y] {
			grid[y][x] = rand.Float64() < 0.3
		}
	}
}

func nextGeneration(grid Grid) Grid {
	next := newGrid()
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			n := countNeighbors(grid, x, y)
			if grid[y][x] {
				next[y][x] = n == 2 || n == 3
			} else {
				next[y][x] = n == 3
			}
		}
	}
	return next
}

func countNeighbors(grid Grid, x, y int) int {
	count := 0
	dirs := [8][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	for _, dir := range dirs {
		dx := x + dir[0]
		dy := y + dir[1]

		if dx >= 0 && dx < width && dy >= 0 && dy < height {
			if grid[dy][dx] {
				count++
			}
		}
	}
	return count
}

func gridToString(grid Grid) string {
	var sb strings.Builder
	for _, row := range grid {
		for _, cell := range row {
			if cell {
				sb.WriteRune('x')
			} else {
				sb.WriteRune(' ')
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server")
		return
	}
	defer conn.Close()

	grid := newGrid()
	randomize(grid)

	for {
		frame := gridToString(grid)
		fmt.Fprint(conn, frame)
		time.Sleep(200 * time.Millisecond)
		grid = nextGeneration(grid)
	}
}
