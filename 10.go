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

func (p point) moveForward() point {
	p.x = p.x + p.vx
	p.y = p.y + p.vy
	return p
}

func (p point) moveBackward() point {
	p.x = p.x - p.vx
	p.y = p.y - p.vy
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

func (b box) area() int {
	return (b.xmax - b.xmin) * (b.ymax - b.ymin)
}

func findBox(points []point) box {
	var xmin, xmax, ymin, ymax int

	for i, p := range points {
		if i == 0 || p.x < xmin {
			xmin = p.x
		}
		if i == 0 || p.x > xmax {
			xmax = p.x
		}
		if i == 0 || p.y < ymin {
			ymin = p.y
		}
		if i == 0 || p.y > ymax {
			ymax = p.y
		}
	}
	return box{xmin: xmin, xmax: xmax, ymin: ymin, ymax: ymax}
}

func findMessage(points []point, box box) {
	area := box.area()
	var newArea int

	seconds := 0
	for {
		for i, p := range points {
			points[i] = p.moveForward()
		}
		box = findBox(points)
		newArea = box.area()

		// next second wasn't better so backtrack
		if newArea > area {
			for i, p := range points {
				points[i] = p.moveBackward()
			}
			break
		}

		area = newArea
		seconds += 1
	}

	printMessage(points)
	fmt.Println("Part 2:", seconds)
}

func printMessage(points []point) {
	b := findBox(points)
	width := b.xmax - b.xmin + 1
	height := b.ymax - b.ymin + 1

	sky := make([][]string, height)
	for r := range sky {
		sky[r] = make([]string, width)
		for pos := range sky[r] {
			sky[r][pos] = "."
		}
	}

	for _, p := range points {
		sky[p.y-b.ymin][p.x-b.xmin] = "#"
	}

	fmt.Println("Part 1:")
	for r := range sky {
		fmt.Println(sky[r])
	}
}
