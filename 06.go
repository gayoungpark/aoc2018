package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func problem06() {
	f, err := os.Open("06.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	coords := make(map[int]coord)
	x_coords := make(map[int][]coord)
	y_coords := make(map[int][]coord)

	areas := make(map[int]int)

	id, max_x, max_y := 1, 0, 0
	min_x, min_y := 1000, 1000
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		input := scanner.Text()

		c := parseCoordinate(input, id)

		if c.x > max_x {
			max_x = c.x
		}
		if c.y > max_y {
			max_y = c.y
		}
		if c.x < min_x {
			min_x = c.x
		}
		if c.y < min_y {
			min_y = c.y
		}

		coords[c.id] = c
		x_coords[c.x] = append(x_coords[c.x], c)
		y_coords[c.y] = append(y_coords[c.y], c)
		id++
	}

	// PART 1
	padding := 2
	for x := min_x; x <= max_x+padding; x++ {
		for y := min_y; y <= max_y+padding; y++ {
			closest_id := findClosestId(x, y, coords)
			areas[closest_id]++
		}
	}

	var remove []coord
	remove = append(remove, x_coords[min_x]...)
	remove = append(remove, x_coords[max_x]...)
	remove = append(remove, y_coords[min_y]...)
	remove = append(remove, y_coords[max_y]...)
	for _, r := range remove {
		areas[r.id] = -1
	}

	max_area := 0
	for id, area := range areas {
		if id != -1 && area > max_area {
			max_area = area
		}
	}

	fmt.Println("Part 1:", max_area)

	// PART 2
	max_sum := 10000
	var region int
	for x := min_x; x <= max_x+padding; x++ {
		for y := min_y; y <= max_y+padding; y++ {
			sum := findSumDistances(x, y, coords)
			if sum < max_sum {
				region++
			}
		}
	}
	fmt.Println("Part 2:", region)
}

type coord struct {
	id   int
	x, y int
}

func parseCoordinate(input string, id int) coord {
	split := strings.Split(input, ", ")

	c := coord{id: id}
	var err error

	c.x, err = strconv.Atoi(split[0])
	if err != nil {
		log.Fatalln("Could not convert x coordinate from input into int:", err)
	}
	c.y, err = strconv.Atoi(split[1])
	if err != nil {
		log.Fatalln("Could not convert y coordinate from input into int:", err)
	}

	return c
}

func findClosestId(x int, y int, coords map[int]coord) int {
	shortest := 1000
	dist_ids := make(map[int][]int)

	for _, c := range coords {
		c_dist := findDistance(c, x, y)
		dist_ids[c_dist] = append(dist_ids[c_dist], c.id)
		if c_dist < shortest {
			shortest = c_dist
		}
	}

	if len(dist_ids[shortest]) > 1 {
		return -1
	}
	return dist_ids[shortest][0]
}

func findDistance(c coord, x int, y int) int {
	return abs(c.x-x) + abs(c.y-y)
}

func abs(d int) int {
	if d < 0 {
		d = -d
	}
	return d
}

func findSumDistances(x int, y int, coords map[int]coord) int {
	var sum int
	for _, c := range coords {
		sum = sum + findDistance(c, x, y)
	}
	return sum
}
