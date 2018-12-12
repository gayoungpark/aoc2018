package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func problem12() {
	f, err := os.Open("12.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	rules := make(map[string]string)

	scanner := bufio.NewScanner(f)

	scanner.Scan()
	initialPlants := scanner.Text()
	initialState := plantState{
		plants: strings.TrimLeft(initialPlants, "initial state: "),
		offset: 0,
	}

	scanner.Scan()
	for scanner.Scan() {
		input := scanner.Text()
		rule := parseRule(input)
		rules[rule.sequence] = rule.outcome
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	firstState := findStateAfter(200, initialState, rules)
	fmt.Println("Part 1:", firstState.score())

	secondState := findStateAfter150(50e9, initialState, rules)
	fmt.Println("Part 2:", secondState.score())
}

type rule struct {
	sequence string
	outcome  string
}

func parseRule(input string) rule {
	fields := strings.Split(input, " => ")
	return rule{sequence: fields[0], outcome: fields[1]}
}

type plantState struct {
	plants string
	offset int
}

func (s plantState) trim() plantState {
	start := strings.IndexByte(s.plants, '#') - 5
	end := strings.LastIndexByte(s.plants, '#') + 6

	s.plants = s.plants[start:end]
	s.offset += start

	return s
}

func (s plantState) score() int {
	var sum int
	for i, r := range s.plants {
		if r == '#' {
			sum += i + s.offset
		}
	}
	return sum
}

func findStateAfter(iter int, state plantState, rules map[string]string) plantState {
	for i := 1; i <= iter; i++ {
		state = findNextState(state, rules)
	}
	return state
}

func findStateAfter150(iter int, state plantState, rules map[string]string) plantState {
	s := findStateAfter(150, state, rules)
	s.offset = iter - 76
	return s
}

func findNextState(state plantState, rules map[string]string) plantState {
	state = padSequence(state)

	next := state.plants[0:2]

	for i := 0; i < len(state.plants)-5; i++ {
		result := rules[state.plants[i:i+5]]
		next += result
	}
	state.plants = next
	return state
}

func padSequence(state plantState) plantState {
	s := plantState{
		plants: "....." + state.plants + ".....",
		offset: state.offset - 5,
	}
	return s.trim()
}
