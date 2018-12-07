package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func problem07() {
	f, err := os.Open("07.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// child to parent mapping
	mapping := make(map[string]string)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		input := scanner.Text()

		s := parse(input)
		mapping[s.child] = mapping[s.child] + s.parent
		mapping[s.parent] = mapping[s.parent] + ""
	}

	copyMap := func() map[string]string {
		cpy := make(map[string]string)
		for k, v := range mapping {
			cpy[k] = v
		}
		return cpy
	}

	fmt.Println("Part 1:", findOrder(copyMap()))

	fmt.Println("Part 2:", timeWork(copyMap()))
}

type step struct {
	parent string
	child  string
}

func parse(input string) step {
	fields := strings.Fields(input)
	before := fields[1]
	after := fields[7]

	return step{parent: before, child: after}
}

func findOrder(m map[string]string) string {
	var instructions string
	var current string

	ready := make(map[string]struct{})

	for len(m) != 0 {
		// given current, find all ready children
		for child, parents := range m {
			if len(parents) == 0 || len(parents) == 1 && current == parents {
				ready[child] = struct{}{}
			}
		}

		// find next step
		current = ""
		for step := range ready {
			if len(current) == 0 {
				current = step
			}
			if step < current {
				current = step
			}
		}
		delete(ready, current)

		instructions = instructions + current
		delete(m, current)

		// remove current from parents
		for child, parents := range m {
			i := strings.Index(parents, current)
			if i != -1 && len(parents) != 1 {
				m[child] = parents[:i] + parents[i+1:]
			}
		}
	}

	return instructions
}

type worker struct {
	task   string
	points int
}

func timeWork(m map[string]string) int {
	var duration int

	workers := make([]worker, 5)

	var instructions string

	ready := make(map[string]struct{})

	for {
		// free up workers that are done
		for _, w := range workers {
			if w.points == 0 && w.task != "" {
				// remove finished task from parents
				for child, parents := range m {
					i := strings.Index(parents, w.task)
					if i != -1 {
						m[child] = parents[:i] + parents[i+1:]
					}
				}
				delete(m, w.task)
			}
		}

		// find all ready children
		for child, parents := range m {
			if len(parents) == 0 {
				ready[child] = struct{}{}
			}
		}

		allIdle := true
		// give work to idle workers
		for idx, w := range workers {
			if w.points <= 0 {
				if len(ready) > 0 {
					// try to find next step
					task := ""
					for step := range ready {
						if len(task) == 0 {
							task = step
						}
						if step < task {
							task = step
						}
					}
					// assign a job
					if task != "" {
						allIdle = false
						delete(ready, task)

						instructions = instructions + task
						m[task] = "*"

						workers[idx].task = task
						workers[idx].points = estimateTask(task)
					}
				}
			} else {
				allIdle = false
			}
			workers[idx].points--
		}
		// fmt.Printf("%3d %3s %3s %s\n", duration, workers[0].task, workers[1].task, instructions)
		if allIdle {
			break
		}
		duration++
	}

	return duration
}

func estimateTask(task string) int {
	return int(task[0]-'A') + 61
}
