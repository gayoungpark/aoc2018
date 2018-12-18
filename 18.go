package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func problem18() {
	f, err := os.Open("18.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var grid [][]byte
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		input := []byte(scanner.Text())
		grid = append(grid, input)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	l := stepMany(grid, 10)
	fmt.Println("Part 1:", l.trees*l.lumberyard)

	l = stepMany(grid, 1e9)
	fmt.Println(1e9)
	fmt.Println("Part 2:", l.trees*l.lumberyard)
}

func stepMany(grid [][]byte, steps int) landscape {
	if steps > 600 {
		steps = (steps-550)%28 + 550
	}

	l := newLandscape(grid)
	for i := 1; i <= steps; i++ {
		l = l.step()
	}
	return l
}

type landscape struct {
	area       [][]byte
	open       int
	trees      int
	lumberyard int
}

func newLandscape(grid [][]byte) landscape {
	g := make([][]byte, len(grid))
	for i, line := range grid {
		g[i] = append([]byte(nil), line...)
	}
	return landscape{area: g}
}

func (l landscape) String() string {
	var builder strings.Builder
	for _, line := range l.area {
		fmt.Fprintf(&builder, "%s\n", line)
	}
	return builder.String()
}

func (l landscape) step() landscape {
	cpy := newLandscape(l.area)
	for y, line := range l.area {
		for x, acre := range line {
			var next byte
			switch acre {
			case '.':
				next = l.stepOpen(x, y)
			case '|':
				next = l.stepTrees(x, y)
			case '#':
				next = l.stepLumberyard(x, y)
			}
			cpy.area[y][x] = next

			switch next {
			case '.':
				cpy.open++
			case '|':
				cpy.trees++
			case '#':
				cpy.lumberyard++
			}
		}
	}
	return cpy
}

func (l landscape) stepOpen(x, y int) byte {
	n := l.neighbors(x, y)
	if n.trees >= 3 {
		return '|'
	}
	return '.'
}

func (l landscape) stepTrees(x, y int) byte {
	n := l.neighbors(x, y)
	if n.lumberyard >= 3 {
		return '#'
	}
	return '|'
}

func (l landscape) stepLumberyard(x, y int) byte {
	n := l.neighbors(x, y)
	if n.lumberyard >= 1 && n.trees >= 1 {
		return '#'
	}
	return '.'
}

type neighborInfo struct {
	x, y       int
	neighbors  []byte
	open       int
	trees      int
	lumberyard int
}

func (l landscape) neighbors(x, y int) neighborInfo {
	n := neighborInfo{
		x: x,
		y: y,
		neighbors: []byte{
			l.safeGet(x, y-1),
			l.safeGet(x, y+1),
			l.safeGet(x-1, y),
			l.safeGet(x+1, y),
			l.safeGet(x-1, y-1),
			l.safeGet(x+1, y-1),
			l.safeGet(x-1, y+1),
			l.safeGet(x+1, y+1),
		},
	}
	for _, acre := range n.neighbors {
		switch acre {
		case '.':
			n.open++
		case '|':
			n.trees++
		case '#':
			n.lumberyard++
		}
	}
	return n
}

func (l landscape) safeGet(x, y int) byte {
	if y >= 0 && y < len(l.area) {
		if x >= 0 && x < len(l.area[0]) {
			return l.area[y][x]
		}
		return 'x'
	}
	return 'x'
}
