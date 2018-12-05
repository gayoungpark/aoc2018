package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func problem04() {
	f, err := os.Open("04.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var records []record

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		input := scanner.Text()
		r := parseRecord(input)
		records = append(records, r)
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i].datetime.Before(records[j].datetime)
	})

	sleepyhead := findSleepyhead(records)
	fmt.Println("Part 1:", sleepyhead)

	// for _, r := range records {
	// 	fmt.Println(r.datetime, r.datetime.Minute(), r.activity)
	// }
}

type record struct {
	datetime time.Time
	activity string
}

func parseRecord(input string) record {
	fields := strings.Fields(input)

	var r record
	var err error

	r.datetime, err = time.Parse("2006-01-02 15:04", fields[0][1:]+" "+fields[1][:5])
	if err != nil {
		log.Fatalln("Could not parse datetime of record:", input)
	}

	switch fields[2] {
	case "Guard":
		r.activity = fields[3]
	case "falls":
		r.activity = "asleep"
	case "wakes":
		r.activity = "awake"
	}

	return r
}

func findSleepyhead(records []record) int {
	guardSchedule := make(map[string][]int)

	var guard string
	for _, r := range records {
		switch r.activity[0] {
		case '#':
			guard = r.activity
		default:
			guardSchedule[guard] = append(guardSchedule[guard], r.datetime.Minute())
		}
	}

	minuteDist := make(map[string][]int)
	minuteAsleep := make(map[string]int)
	for guard, timestamps := range guardSchedule {
		for i := 0; i < len(timestamps); i = i + 2 {
			asleep := timestamps[i]
			awake := timestamps[i+1]
			for j := asleep; j < awake; j++ {
				if len(minuteDist[guard]) == 0 {
					minuteDist[guard] = make([]int, 60)
				}
				minuteDist[guard][j]++
				minuteAsleep[guard]++
			}
		}
	}

	var sleepyhead string
	max := 0
	for guard, val := range minuteAsleep {
		if val > max {
			max = val
			sleepyhead = guard
		}
	}

	var sleepiestTime int
	max = 0
	for t, freq := range minuteDist[sleepyhead] {
		if freq > max {
			sleepiestTime = t
			max = freq
		}
	}

	sleepyheadId, _ := strconv.Atoi(sleepyhead[1:])

	// PART 2
	var g string
	var m int
	f := 0
	for guard, dist := range minuteDist {
		for t, freq := range dist {
			if freq > f {
				g = guard
				m = t
				f = freq
			}
		}
	}

	fmt.Println(g, m, f)
	fmt.Println(minuteDist)

	return sleepyheadId * sleepiestTime
}
