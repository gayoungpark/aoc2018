package main

import (
	"bytes"
	"fmt"
	"strconv"
)

func problem14() {
	iter := 330121

	sb := newScoreboard()
	for i := 0; i < iter+10; i++ {
		sb.step()
	}
	fmt.Println("Part 1:", bytesToString(sb.recipes[iter:iter+10]))

	sb = newScoreboard()
	sequence := []byte{3, 3, 0, 1, 2, 1}
	var start int
	for {
		sb.step()
		start = len(sb.recipes) - len(sequence) - 1
		if start >= 0 {
			if bytes.Equal(sb.recipes[start:len(sb.recipes)-1], sequence) {
				break
			}
			start++
			if bytes.Equal(sb.recipes[start:], sequence) {
				break
			}
		}
	}
	fmt.Println("Part 2:", start)
}

type scoreboard struct {
	recipes []byte
	elf1    int
	elf2    int
}

func newScoreboard() *scoreboard {
	return &scoreboard{
		recipes: []byte{3, 7},
		elf1:    0,
		elf2:    1,
	}
}

func (sb *scoreboard) String() string {
	return fmt.Sprintf("%v", sb.recipes)
}

func (sb *scoreboard) step() {
	newScore := sb.recipes[sb.elf1] + sb.recipes[sb.elf2]
	if newScore > 9 {
		sb.recipes = append(sb.recipes, 1, newScore-10)
	} else {
		sb.recipes = append(sb.recipes, newScore)
	}
	sb.elf1 = (sb.elf1 + int(sb.recipes[sb.elf1]) + 1) % len(sb.recipes)
	sb.elf2 = (sb.elf2 + int(sb.recipes[sb.elf2]) + 1) % len(sb.recipes)
}

func bytesToString(bytes []byte) string {
	var s string
	for _, b := range bytes {
		s += strconv.Itoa(int(b))
	}
	return s
}
