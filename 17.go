package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/cespare/hasty"
)

func problem17() {
	f, err := os.Open("17.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	clays := make(vec2set)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Bytes()
		if line[0] == 'x' {
			var v vertVein
			if err := hasty.Parse(line, &v, vertVeinRegexp); err != nil {
				log.Fatal(err)
			}
			for y := v.YBegin; y <= v.YEnd; y++ {
				clays.add(vec2{v.X, y})
			}
		}
		if line[0] == 'y' {
			var v horizVein
			if err := hasty.Parse(line, &v, horizVeinRegexp); err != nil {
				log.Fatal(err)
			}
			for x := v.XBegin; x <= v.XEnd; x++ {
				clays.add(vec2{x, v.Y})
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	g := newGround(clays)
	// g.print = true
	g.fill()

	water := len(g.water)
	wetSand := len(g.wetSand)
	fmt.Println("Part 1:", water+wetSand)
	fmt.Println("Part 2:", water)
}

var (
	vertVeinRegexp  = regexp.MustCompile(`^x=(?P<X>\d+), y=(?P<YBegin>\d+)\.\.(?P<YEnd>\d+)$`)
	horizVeinRegexp = regexp.MustCompile(`^y=(?P<Y>\d+), x=(?P<XBegin>\d+)\.\.(?P<XEnd>\d+)$`)
)

type vertVein struct {
	X      int
	YBegin int
	YEnd   int
}

type horizVein struct {
	Y      int
	XBegin int
	XEnd   int
}

type ground struct {
	spring  vec2
	clays   vec2set
	water   vec2set
	wetSand vec2set
	min     vec2
	max     vec2
	print   bool
}

func newGround(clays vec2set) *ground {
	cpy := make(vec2set)
	for c := range clays {
		cpy[c] = struct{}{}
	}

	min, max := vec2{-1, -1}, vec2{-1, -1}
	for c := range cpy {
		if min.x == -1 || c.x < min.x {
			min.x = c.x
		}
		if min.y == -1 || c.y < min.y {
			min.y = c.y
		}
		if max.x == -1 || c.x > max.x {
			max.x = c.x
		}
		if max.y == -1 || c.y > max.y {
			max.y = c.y
		}
	}

	// add padding
	min.x--
	max.x++

	return &ground{
		spring:  vec2{500, 0},
		clays:   cpy,
		water:   make(vec2set),
		wetSand: make(vec2set),
		min:     min,
		max:     max,
	}
}

func (g *ground) debug() {
	if !g.print {
		return
	}
	fmt.Println(g)
	time.Sleep(500 * time.Millisecond)
}

func (g *ground) String() string {
	var builder strings.Builder
	for y := g.min.y; y <= g.max.y; y++ {
		for x := g.min.x; x <= g.max.x; x++ {
			loc := vec2{x, y}
			switch {
			case g.clays.contains(loc):
				builder.WriteByte('#')
			case g.water.contains(loc):
				builder.WriteByte('~')
			case g.wetSand.contains(loc):
				builder.WriteByte('|')
			default:
				builder.WriteByte(' ')
			}
		}
		builder.WriteByte('\n')
	}
	return builder.String()
}

func (g *ground) fill() {
	start := vec2{g.spring.x, g.min.y}
	g.fillDown(start)
	g.debug()
}

func (g *ground) fillDown(start vec2) bool {
	g.debug()

	curr := start
	for !g.occupied(curr) {
		if curr.y > g.max.y {
			return false
		}
		if g.wetSand.contains(curr) {
			return false
		}
		g.wetSand.add(curr)
		curr.y++
	}

	for curr.y--; curr.y >= start.y; curr.y-- {
		left := g.fillLeft(curr)
		right := g.fillRight(curr)
		if !left || !right {
			return false
		}
		g.fillRow(curr)
	}

	return true
}

func (g *ground) fillRow(loc vec2) {
	g.debug()
	for v := loc; !g.clays.contains(v); v.x-- {
		g.water.add(v)
		delete(g.wetSand, v)
	}
	for v := loc; !g.clays.contains(v); v.x++ {
		g.water.add(v)
		delete(g.wetSand, v)
	}
}

func (g *ground) fillLeft(loc vec2) bool {
	return g.fillLeftOrRight(loc, -1)
}

func (g *ground) fillRight(loc vec2) bool {
	return g.fillLeftOrRight(loc, 1)
}

func (g *ground) fillLeftOrRight(loc vec2, dir int) bool {
	g.debug()
	for v := loc; !g.clays.contains(v); v.x += dir {
		below := vec2{v.x, v.y + 1}
		if !g.occupied(below) {
			if !g.fillDown(v) {
				return false
			}
		}
		g.wetSand.add(v)
	}
	return true
}

func (g *ground) occupied(loc vec2) bool {
	clay := g.clays.contains(loc)
	water := g.water.contains(loc)
	return (clay || water)
}
