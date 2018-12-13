package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

func problem13() {
	f, err := os.Open("13.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var tracks [][]byte
	var carts []*cart

	var y int
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		input := []byte(scanner.Text())
		input, carts = findCart(input, y, carts)
		tracks = append(tracks, input)
		y++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Part 1:", locateCollision(tracks, copyCarts(carts)))

	fmt.Println("Part 2:", locateSurvivor(tracks, copyCarts(carts)))
}

type cart struct {
	x, y      int
	direction byte
	nextTurn  byte
}

func (c *cart) String() string {
	return fmt.Sprintf(
		"(%3d, %3d): currently %c where next is %c",
		c.x, c.y, c.direction, c.nextTurn,
	)
}

func (c *cart) move(tracks [][]byte) *cart {
	switch c.direction {
	case '^':
		switch tracks[c.y-1][c.x] {
		case '\\':
			c.direction = '<'
		case '/':
			c.direction = '>'
		case '-':
			log.Fatal("Cannot move up into a horizontal track")
		case '+':
			c = c.intersection()
		}
		c.y--
	case 'v':
		switch tracks[c.y+1][c.x] {
		case '\\':
			c.direction = '>'
		case '/':
			c.direction = '<'
		case '-':
			log.Fatal("Cannot move down into a horizontal track")
		case '+':
			c = c.intersection()
		}
		c.y++
	case '<':
		switch tracks[c.y][c.x-1] {
		case '\\':
			c.direction = '^'
		case '/':
			c.direction = 'v'
		case '|':
			log.Fatal("Cannot move left into a vertical track")
		case '+':
			c = c.intersection()
		}
		c.x--
	case '>':
		switch tracks[c.y][c.x+1] {
		case '\\':
			c.direction = 'v'
		case '/':
			c.direction = '^'
		case '|':
			log.Fatal("Cannot move right into a vertical track")
		case '+':
			c = c.intersection()
		}
		c.x++
	}
	return c
}

func (c *cart) intersection() *cart {
	switch c.direction {
	case '^':
		switch c.nextTurn {
		case '<', '>':
			c.direction = c.nextTurn
		}
	case 'v':
		switch c.nextTurn {
		case '<':
			c.direction = '>'
		case '>':
			c.direction = '<'
		}
	case '<':
		switch c.nextTurn {
		case '<':
			c.direction = 'v'
		case '>':
			c.direction = '^'
		}
	case '>':
		switch c.nextTurn {
		case '<':
			c.direction = '^'
		case '>':
			c.direction = 'v'
		}
	}
	c.updateNextTurn()
	return c
}

func (c *cart) updateNextTurn() *cart {
	switch c.nextTurn {
	case '<':
		c.nextTurn = '-'
	case '|', '-':
		c.nextTurn = '>'
	case '>':
		c.nextTurn = '<'
	}
	return c
}

func findCart(input []byte, y int, carts []*cart) ([]byte, []*cart) {
	for x, c := range input {
		switch c {
		case '^', 'v', '<', '>':
			carts = append(carts, &cart{
				x:         x,
				y:         y,
				direction: c,
				nextTurn:  '<',
			})
			carts = carts
		}

		switch c {
		case '^', 'v':
			input[x] = '|'
		case '<', '>':
			input[x] = '-'
		}

	}
	return input, carts
}

type vec2 struct {
	x, y int
}

func (v vec2) String() string {
	return fmt.Sprintf("%d,%d", v.x, v.y)
}

func locateCollision(tracks [][]byte, carts []*cart) vec2 {
	for {
		sort.Slice(carts, func(i, j int) bool {
			if carts[i].y == carts[j].y {
				return carts[i].x < carts[j].x
			}
			return carts[i].y < carts[j].y
		})

		for _, c := range carts {
			c.move(tracks)
			if v, ok := detectCollision(carts); ok {
				return v
			}
		}
	}
	return vec2{}
}

func detectCollision(carts []*cart) (vec2, bool) {
	seen := make(map[vec2]struct{})
	for _, c := range carts {
		v := vec2{x: c.x, y: c.y}
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
		} else {
			return v, true
		}
	}
	return vec2{}, false
}

func locateSurvivor(tracks [][]byte, carts []*cart) vec2 {
	seen := make(map[*cart]struct{})
	for {
		sort.Slice(carts, func(i, j int) bool {
			if carts[i].y == carts[j].y {
				return carts[i].x < carts[j].x
			}
			return carts[i].y < carts[j].y
		})

		for _, c := range carts {
			// move all carts that aren't "removed"
			if _, ok := seen[c]; !ok {
				c.move(tracks)
			}
			// look for collision among carts not "removed"
			notSeen := removeCarts(carts, seen)
			if v, ok := detectCollision(notSeen); ok {
				// "remove" both carts by adding to seen set
				toRemove := locateCarts(v, notSeen)
				for _, c := range toRemove {
					seen[c] = struct{}{}
				}
			}
		}

		carts := removeCarts(carts, seen)

		if len(carts) == 1 {
			survivor := carts[0]
			return vec2{x: survivor.x, y: survivor.y}
		}
	}
	return vec2{}
}

func locateCarts(loc vec2, carts []*cart) []*cart {
	var found []*cart
	for _, c := range carts {
		if c.x == loc.x && c.y == loc.y {
			found = append(found, c)
		}
	}
	return found
}

func removeCarts(carts []*cart, toRemove map[*cart]struct{}) []*cart {
	var removed []*cart
	for _, c := range carts {
		if _, ok := toRemove[c]; !ok {
			removed = append(removed, c)
		}

	}
	return removed
}

func printState(tracks [][]byte, carts []*cart) {
	grid := make([][]byte, len(tracks))
	for i, line := range tracks {
		grid[i] = append([]byte(nil), line...)
	}

	for _, c := range carts {
		grid[c.y][c.x] = c.direction
	}

	for _, line := range grid {
		fmt.Println(string(line))
	}
	fmt.Println()
}

func copyCarts(carts []*cart) []*cart {
	cpy := make([]*cart, len(carts))
	for i, c := range carts {
		cpy[i] = &cart{
			x:         c.x,
			y:         c.y,
			direction: c.direction,
			nextTurn:  c.nextTurn,
		}
	}
	return cpy
}
