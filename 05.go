package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var letterPairs [256]byte

func init() {
	for c := byte('a'); c <= 'z'; c++ {
		letterPairs[c] = c + 'A' - 'a'
	}
	for c := byte('A'); c <= 'Z'; c++ {
		letterPairs[c] = c - 'A' + 'a'
	}
}

func problem05() {
	f, err := os.Open("05.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	polymer, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	polymer = bytes.TrimSpace(polymer)

	fmt.Println("Part 1:", reducePolymer(polymer))
	fmt.Println("Part 2:", findShortestPolymer(polymer))
}

func reducePolymer(polymer []byte) int {
	p := make([]byte, len(polymer))
	copy(p, polymer)

	for {
		n := len(p)
		p = reducePolymer1(p)
		if n == len(p) {
			return n
		}
	}
}

func reducePolymer1(polymer []byte) []byte {
	insert := 0
	for i := 0; i < len(polymer); i++ {
		if i == len(polymer)-1 || letterPairs[polymer[i]] != polymer[i+1] {
			polymer[insert] = polymer[i]
			insert++
		} else {
			i++
		}
	}
	return polymer[:insert]
}

func findShortestPolymer(polymer []byte) int {
	shortest := len(polymer)
	for c := byte('a'); c <= 'z'; c++ {
		n := reducePolymer(removeUnit(c, polymer))
		if n < shortest {
			shortest = n
		}
	}
	return shortest
}

func removeUnit(u byte, polymer []byte) []byte {
	removed := make([]byte, len(polymer))

	insert := 0
	for _, p := range polymer {
		if p == u || p == u+'A'-'a' {
			continue
		}
		removed[insert] = p
		insert++
	}
	return removed[:insert]
}
