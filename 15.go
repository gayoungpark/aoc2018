package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
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
	for r := 0; ; r++ {
		done, _ := cave.step()
		if done {
			fmt.Println("Part 1:", r*cave.sumHP())
			break
		}
	}

	// part 2
	for atk := 4; ; atk++ {
		cave := newCave(board, atk)
		for r := 0; ; r++ {
			done, elfDied := cave.step()
			if elfDied {
				break
			}
			if done {
				fmt.Println("Part 2:", r*cave.sumHP())
				return
			}
		}
	}

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
				fmt.Fprintf(&builder, " %c(%d)", sq, c.units[vec2{x, y}].hp)
			}
		}
		fmt.Fprintln(&builder)
	}
	return builder.String()
}

func (c *cave) sumHP() int {
	var sum int
	for _, u := range c.units {
		sum += u.hp
	}
	return sum
}

func (c *cave) setBoard(loc vec2, sq byte) {
	c.board[loc.y][loc.x] = sq
}

func (c *cave) getSquare(loc vec2) byte {
	return c.board[loc.y][loc.x]
}

func (c *cave) occupied(loc vec2) bool {
	return c.getSquare(loc) != '.'
}

func (c *cave) orderUnits() []*unit {
	var ordered []*unit
	for _, u := range c.units {
		ordered = append(ordered, u)
	}
	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].loc.readingOrderLess(ordered[j].loc)
	})
	return ordered
}

type vec2set map[vec2]struct{}

func (s vec2set) add(v vec2) {
	s[v] = struct{}{}
}

func (s vec2set) contains(v vec2) bool {
	_, ok := s[v]
	return ok
}

func (c *cave) findEnemies(u *unit) vec2set {
	enemies := make(vec2set)
	for _, u2 := range c.units {
		if u.typ != u2.typ {
			enemies.add(u2.loc)
		}
	}
	return enemies
}

func (c *cave) step() (done, elfDied bool) {
	ordered := c.orderUnits()
	for _, u := range ordered {
		if u.hp <= 0 {
			continue
		}
		enemies := c.findEnemies(u)
		if len(enemies) == 0 {
			return true, elfDied
		}
		attacked, dead := c.attack(u)
		if dead {
			elfDied = true
		}
		if !attacked {
			dest := c.findDest(u, enemies)
			c.move(u, dest)
			if _, dead = c.attack(u); dead {
				elfDied = true
			}
		}
	}
	return false, elfDied
}

func (c *cave) findDest(u *unit, enemies vec2set) vec2 {
	frontier := vec2set{u.loc: {}}

	type locInfo struct {
		n       int
		parents vec2set
	}
	explored := map[vec2]*locInfo{
		u.loc: {n: 0},
	}
	found := make(vec2set)

	for dist := 1; len(frontier) > 0 && len(found) == 0; dist++ {
		nextFrontier := make(vec2set)
		for loc := range frontier {
			neighbors := loc.neighbors()
			for _, n := range neighbors {
				if frontier.contains(n) {
					continue
				}
				if c.occupied(n) {
					if enemies.contains(n) {
						found.add(n)
					} else {
						continue
					}
				}
				if info, ok := explored[n]; ok {
					if dist == info.n {
						info.parents.add(loc)
					}
					continue
				}
				explored[n] = &locInfo{
					n:       dist,
					parents: vec2set{loc: {}},
				}
				nextFrontier.add(n)
			}
		}
		frontier = nextFrontier
	}

	// no valid move found
	if len(found) == 0 {
		return u.loc
	}

	chosen := vec2{-1, -1}
	for loc := range found {
		if chosen.x < 0 || loc.readingOrderLess(chosen) {
			chosen = loc
		}
	}

	dest := vec2{-1, -1}
	var backtrack func(vec2)
	backtrack = func(v vec2) {
		for p := range explored[v].parents {
			if p == u.loc {
				if dest.x < 0 || v.readingOrderLess(dest) {
					dest = v
				}
			} else {
				backtrack(p)
			}
		}
	}
	backtrack(chosen)

	return dest
}

func (c *cave) move(u *unit, dest vec2) {
	if u.loc == dest {
		return
	}
	c.setBoard(u.loc, '.')
	c.setBoard(dest, u.typ.String()[0])
	delete(c.units, u.loc)
	u.loc = dest
	c.units[dest] = u
}

func (c *cave) attack(u *unit) (attacked, elfDied bool) {
	var inRange []*unit
	for _, n := range u.loc.neighbors() {
		if u2, ok := c.units[n]; ok && u2.typ != u.typ {
			inRange = append(inRange, u2)
		}
	}

	// no enemy within range
	if len(inRange) == 0 {
		return false, elfDied
	}

	sort.Slice(inRange, func(i, j int) bool {
		u1, u2 := inRange[i], inRange[j]
		if u1.hp == u2.hp {
			return u1.loc.readingOrderLess(u2.loc)
		}
		return u1.hp < u2.hp
	})
	enemy := inRange[0]
	enemy.hp -= u.atk
	if enemy.hp <= 0 {
		c.setBoard(enemy.loc, '.')
		delete(c.units, enemy.loc)
		elfDied = enemy.typ == elf
	}
	return true, elfDied
}

func (v vec2) neighbors() []vec2 {
	return []vec2{
		vec2{v.x, v.y - 1},
		vec2{v.x - 1, v.y},
		vec2{v.x + 1, v.y},
		vec2{v.x, v.y + 1},
	}
}

func (v vec2) readingOrderLess(v2 vec2) bool {
	if v.y == v2.y {
		return v.x < v2.x
	}
	return v.y < v2.y
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
				loc := vec2{x, y}
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

func (t unitType) String() string {
	switch t {
	case elf:
		return "E"
	case goblin:
		return "G"
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
