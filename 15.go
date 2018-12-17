package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func problem15() {
	f, err := os.Open("15.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var board [][]byte

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		input := []byte(scanner.Text())
		board = append(board, input)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// part 1
	cave := newCave(board, 3)
	fmt.Println(cave)

}

type cave struct {
	board [][]byte
	units map[vec2]*unit
}

func (c *cave) String() string {
	var builder strings.Builder
	for y, line := range c.board {
		fmt.Fprintf(&builder, "%s  ", line)
		for x, sq := range line {
			if sq == 'G' || sq == 'E' {
				fmt.Fprintf(&builder, " %c(%d)", sq, c.units[vec2{x: x, y: y}].hp)
			}
		}
		fmt.Fprintln(&builder)
	}
	return builder.String()
}

func newCave(board [][]byte, elfAtk int) *cave {
	// copy board
	b := make([][]byte, len(board))
	for i, line := range board {
		b[i] = append([]byte(nil), line...)
	}

	// initialize units
	u := make(map[vec2]*unit)
	for y, line := range board {
		for x, sq := range line {
			if sq != '#' && sq != '.' {
				loc := vec2{x: x, y: y}
				unit := createUnit(sq, loc, elfAtk)
				u[loc] = unit
			}
		}
	}

	return &cave{
		board: b,
		units: u,
	}
}

func createUnit(sq byte, loc vec2, elfAtk int) *unit {
	var typ unitType
	var atk int

	switch sq {
	case 'G':
		typ = goblin
		atk = 3
	case 'E':
		typ = elf
		atk = elfAtk
	}

	return &unit{
		typ: typ,
		loc: loc,
		hp:  200,
		atk: atk,
	}
}

type unitType int

const (
	elf unitType = iota
	goblin
)

func (t unitType) byte() byte {
	switch t {
	case elf:
		return 'E'
	case goblin:
		return 'G'
	default:
		panic("Bad unit type")
	}
}

type unit struct {
	typ unitType
	loc vec2
	hp  int
	atk int
}
