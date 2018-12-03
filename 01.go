package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func problem01() {
	f, err := os.Open("01.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var deltas []int64
	var sum int64
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		n, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		deltas = append(deltas, n)
		sum += n
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Part 1:", sum)

	freqSeen := map[int64]bool{0: true}
	var currFreq int64
	for i := 0; ; i = (i + 1) % len(deltas) {
		currFreq += deltas[i]
		if freqSeen[currFreq] {
			fmt.Println("Part 2:", currFreq)
			break
		}
		freqSeen[currFreq] = true
	}
}
