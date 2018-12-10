package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func problem10() {
	f, err := os.Open("10.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var points []point

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		input := scanner.Text()
		p := parsePoint(input)
		points = append(points, p)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	b := findBox(points)

	findMessage(points, b)
}

type point struct {
	x, y   int
	vx, vy int
}

func (p point) movePoint() point {
	p.x = p.x + p.vx
	p.y = p.y + p.vy
	return p
}

func parsePoint(input string) point {
	fields := strings.SplitAfter(input, "<")

	position := strings.Split(
		strings.TrimRight(fields[1], "> velocity=<"),
		",",
	)
	velocity := strings.Split(
		strings.TrimRight(fields[2], ">"),
		",",
	)

	var p point
	var err error
	p.x, err = strconv.Atoi(strings.TrimSpace(position[0]))
	if err != nil {
		log.Fatalln("Failed to convert position x of line:", input)
	}
	p.y, err = strconv.Atoi(strings.TrimSpace(position[1]))
	if err != nil {
		log.Fatalln("Failed to convert position y of line:", input)
	}
	p.vx, err = strconv.Atoi(strings.TrimSpace(velocity[0]))
	if err != nil {
		log.Fatalln("Failed to convert velocity x of line:", input)
	}
	p.vy, err = strconv.Atoi(strings.TrimSpace(velocity[1]))
	if err != nil {
		log.Fatalln("Failed to convert velocity y of line:", input)
	}

	return p
}

type box struct {
	xmin, xmax int
	ymin, ymax int
}

func (b box) size() int {
	return abs(b.xmin-b.xmax) + abs(b.ymin-b.ymax)
}

func findBox(points []point) box {
	b := box{
		xmin: points[0].x,
		xmax: points[0].x,
		ymin: points[0].y,
		ymax: points[0].y,
	}

	for _, p := range points {
		if p.x < b.xmin {
			b.xmin = p.x
		}
		if p.x > b.xmax {
			b.xmax = p.x
		}
		if p.y < b.ymin {
			b.ymin = p.y
		}
		if p.y > b.ymax {
			b.ymax = p.y
		}
	}
	return b
}

func findMessage(points []point, box box) {
	size := box.size()
	var newSize int

	seconds := 0
	for {
		// try next second
		next := make([]point, len(points))
		for i, p := range points {
			next[i] = p.movePoint()
		}
		newBox := findBox(next)
		newSize = newBox.size()

		// next wasn't better so quit
		if newSize > size {
			break
		}

		// bookkeeping before trying again
		size = newSize
		points = next
		box = newBox
		seconds += 1
	}

	printMessage(points, box)
	fmt.Println("Part 2:", seconds)
}

func printMessage(points []point, box box) {
	width := abs(box.xmin-box.xmax) + 1
	height := abs(box.ymin-box.ymax) + 1

	sky := make([][]string, height)
	for r := range sky {
		sky[r] = make([]string, width)
		for pos := range sky[r] {
			sky[r][pos] = "."
		}
	}

	for _, p := range points {

		xnew := p.x - box.xmin
		ynew := p.y - box.ymin

		sky[ynew][xnew] = "x"
	}

	fmt.Println("Part 1:")
	for r := range sky {
		fmt.Println(sky[r])
	}
}
