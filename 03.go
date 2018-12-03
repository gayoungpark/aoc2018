package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func problem03() {
	const size = 1000

	f, err := os.Open("03.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	grid := make([][][]string, size)
	for i := range grid {
		grid[i] = make([][]string, size)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		input := scanner.Text()
		c := parseClaim(input)
		processClaim(grid, c)
	}

	overlapped := countOverlap(grid)
	fmt.Println("Part 1:", overlapped)

	perfectClaim := findPerfectClaim(grid)
	fmt.Println("Part 2:", perfectClaim)
}

type claim struct {
	id   string
	l, t int // inches between the left and top edges of the fabric
	w, h int // width and height of the rectanges in inches
}

func parseClaim(input string) claim {
	fields := strings.Fields(input)

	c := claim{id: fields[0]}
	var err error

	edges := strings.Split(strings.TrimSuffix(fields[2], ":"), ",")
	c.l, err = strconv.Atoi(edges[0])
	if err != nil {
		log.Fatalln("Could not convert left edge to int:", edges[0])
	}
	c.t, err = strconv.Atoi(edges[1])
	if err != nil {
		log.Fatalln("Could not convert right edge to int:", edges[1])
	}

	size := strings.Split(fields[3], "x")
	c.w, err = strconv.Atoi(size[0])
	if err != nil {
		log.Fatalln("Could not convert width to int:", size[0])
	}
	c.h, err = strconv.Atoi(size[1])
	if err != nil {
		log.Fatalln("Could not convert height to int:", size[1])
	}

	return c
}

func processClaim(g [][][]string, c claim) {
	for i := 0; i < c.w; i++ {
		for j := 0; j < c.h; j++ {
			x := c.l + i
			y := c.t + j
			g[y][x] = append(g[y][x], c.id)
		}
	}
}

func countOverlap(g [][][]string) int {
	var count int
	for row := 0; row < len(g); row++ {
		for col := 0; col < len(g[0]); col++ {
			if len(g[col][row]) > 1 {
				count++
			}
		}
	}
	return count
}

func findPerfectClaim(g [][][]string) string {
	overlaps := make(map[string]bool)

	for row := 0; row < len(g); row++ {
		for col := 0; col < len(g[0]); col++ {
			ids := g[col][row]
			if len(ids) > 1 {
				for _, id := range ids {
					overlaps[id] = true
				}
			}
			if len(ids) == 1 {
				if !overlaps[ids[0]] {
					overlaps[ids[0]] = false
				}
			}
		}
	}

	for id, faulty := range overlaps {
		if !faulty {
			return id
		}
	}

	panic("Perfect claim not found")
}
