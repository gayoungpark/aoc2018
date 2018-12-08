package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

// A<childCount> A<metadataLen>

func problem08() {
	f, err := os.Open("08.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var input []int

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		n, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatalln("Failed to convert input character to int", err)
		}
		input = append(input, n)

	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	tree, remaining := parseTree(input)
	if len(remaining) != 0 {
		log.Fatalln("Invalid tree - input remains after parsing:", remaining)
	}

	fmt.Println("Part 1:", tree.sumMetadata())

	fmt.Println("Part 2:", tree.computeValue())
}

type node struct {
	children []*node
	metadata []int
}

func (n *node) sumMetadata() int {
	var sum int
	for _, c := range n.children {
		sum += c.sumMetadata()
	}
	for _, m := range n.metadata {
		sum += m
	}
	return sum
}

func (n *node) computeValue() int {
	var value int

	if len(n.children) == 0 {
		for _, m := range n.metadata {
			value += m
		}
		return value
	}

	for _, m := range n.metadata {
		if m <= len(n.children) && m > 0 {
			value += n.children[m-1].computeValue()
		}
	}

	return value
}

func parseTree(slice []int) (*node, []int) {

	n := &node{
		children: make([]*node, slice[0]),
		metadata: make([]int, slice[1]),
	}

	slice = slice[2:]

	for i := range n.children {
		n.children[i], slice = parseTree(slice)
	}

	copy(n.metadata, slice)

	return n, slice[len(n.metadata):]
}
