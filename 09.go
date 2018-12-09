package main

import "fmt"

func problem09() {
	fmt.Println("Part 1:", scoreGame(413, 71082))
	fmt.Println("Part 2:", scoreGame(413, 7108200))
}

func scoreGame(numPlayers int, numMarbles int) int {
	scores := make([]int, numPlayers)

	current := &marble{id: 0}
	current.prev = current
	current.next = current

	for i := 1; i <= numMarbles; i++ {
		player := (i - 1) % numPlayers
		if i%23 == 0 {
			for j := 0; j < 7; j++ {
				current = current.prev
			}
			scores[player] += i + current.id
			current = current.remove()
		} else {
			current = current.next.addNext(i)
		}
	}

	return findMax(scores)
}

type marble struct {
	id   int
	prev *marble
	next *marble
}

func (m *marble) addNext(id int) *marble {
	next := &marble{
		id:   id,
		prev: m,
		next: m.next,
	}
	m.next.prev = next
	m.next = next

	return next
}

func (m *marble) remove() *marble {
	next := m.next
	m.prev.next = m.next
	m.next.prev = m.prev
	m.prev, m.next = nil, nil
	return next
}

func findMax(s []int) int {
	max := 0
	for _, n := range s {
		if n > max {
			max = n
		}
	}
	return max
}
