package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func problem02() {
	f, err := os.Open("02.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var ids []string
	var twice int64
	var thrice int64

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		id := scanner.Text()
		ids = append(ids, id)
		freqs := make(map[rune]int64)
		for _, r := range id {
			freqs[r]++
		}
		var twoFound, threeFound bool
		for _, c := range freqs {
			switch c {
			case 2:
				twoFound = true
			case 3:
				threeFound = true
			}
		}
		if twoFound {
			twice++
		}
		if threeFound {
			thrice++
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Part 1:", twice*thrice)

	combos := make(map[string]struct{})
idLoop:
	for _, id := range ids {
		runes := []rune(id)
		for i := range runes {
			combo := string(runes[:i]) + "*" + string(runes[i+1:])
			if _, ok := combos[combo]; ok {
				fmt.Println("Part 2:", string(runes[:i])+string(runes[i+1:]))
				break idLoop
			}
			combos[combo] = struct{}{}
		}
	}
}
