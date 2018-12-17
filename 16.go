package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func problem16() {
	f, err := os.Open("16.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var samples []*sample
	var program []instruction

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		input := scanner.Text()
		if fields := strings.Split(input, "Before:"); len(fields) > 1 {
			s := &sample{}

			s.before = parseRegisters(fields[1])

			scanner.Scan()
			s.instruction = parseInstruction(scanner.Text())

			scanner.Scan()
			s.after = parseRegisters(strings.Split(scanner.Text(), "After:")[1])

			scanner.Scan()

			samples = append(samples, s)
		} else if len(input) > 0 {
			program = append(program, parseInstruction(input))
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	part1 := findOpcodeMapping(samples)
	fmt.Println("Part 1:", part1)

	r0 := runProgram(program)
	fmt.Println("Part 2:", r0)
}

type sample struct {
	before      registers
	instruction instruction
	after       registers
}

type registers [4]int

type instruction struct {
	opcode int
	input1 int
	input2 int
	output int
}

func parseRegisters(input string) registers {
	var r registers
	var err error

	input = strings.TrimSpace(input)
	input = strings.TrimLeft(input, "[")
	input = strings.TrimRight(input, "]")
	fields := strings.Split(input, ", ")
	for i, f := range fields {
		r[i], err = strconv.Atoi(f)
		if err != nil {
			log.Fatalf("Could not convert register %d to int: %s", i, f)
		}
	}
	return r
}

func parseInstruction(input string) instruction {
	var inst instruction
	var err error

	slice := strings.Fields(input)
	inst.opcode, err = strconv.Atoi(slice[0])
	if err != nil {
		log.Fatalln("Could not convert instruction opcode to int:", slice[0])
	}
	inst.input1, err = strconv.Atoi(slice[1])
	if err != nil {
		log.Fatalln("Could not convert instruction input1 to int:", slice[1])
	}
	inst.input2, err = strconv.Atoi(slice[2])
	if err != nil {
		log.Fatalln("Could not convert instruction input2 to int:", slice[2])
	}
	inst.output, err = strconv.Atoi(slice[3])
	if err != nil {
		log.Fatalln("Could not convert instruction output to int:", slice[3])
	}

	return inst
}

func (r registers) addr(a, b, c int) registers {
	r[c] = r[a] + r[b]
	return r
}

func (r registers) addi(a, b, c int) registers {
	r[c] = r[a] + b
	return r
}

func (r registers) mulr(a, b, c int) registers {
	r[c] = r[a] * r[b]
	return r
}

func (r registers) muli(a, b, c int) registers {
	r[c] = r[a] * b
	return r
}

func (r registers) banr(a, b, c int) registers {
	r[c] = r[a] & r[b]
	return r
}

func (r registers) bani(a, b, c int) registers {
	r[c] = r[a] & b
	return r
}

func (r registers) borr(a, b, c int) registers {
	r[c] = r[a] | r[b]
	return r
}

func (r registers) bori(a, b, c int) registers {
	r[c] = r[a] | b
	return r
}

func (r registers) setr(a, _, c int) registers {
	r[c] = r[a]
	return r
}

func (r registers) seti(a, _, c int) registers {
	r[c] = a
	return r
}

func boolToInt(a bool) int {
	if a {
		return 1
	}
	return 0
}

func (r registers) gtir(a, b, c int) registers {
	r[c] = boolToInt(a > r[b])
	return r
}

func (r registers) gtri(a, b, c int) registers {
	r[c] = boolToInt(r[a] > b)
	return r
}

func (r registers) gtrr(a, b, c int) registers {
	r[c] = boolToInt(r[a] > r[b])
	return r
}

func (r registers) eqir(a, b, c int) registers {
	r[c] = boolToInt(a == r[b])
	return r
}

func (r registers) eqri(a, b, c int) registers {
	r[c] = boolToInt(r[a] == b)
	return r
}

func (r registers) eqrr(a, b, c int) registers {
	r[c] = boolToInt(r[a] == r[b])
	return r
}

var opFuncs = []func(registers, int, int, int) registers{
	registers.addr,
	registers.addi,
	registers.mulr,
	registers.muli,
	registers.banr,
	registers.bani,
	registers.borr,
	registers.bori,
	registers.setr,
	registers.seti,
	registers.gtir,
	registers.gtri,
	registers.gtrr,
	registers.eqir,
	registers.eqri,
	registers.eqrr,
}

var ops = make([]func(registers, int, int, int) registers, len(opFuncs))

func findOpcodeMapping(samples []*sample) int {
	// opcode : set of opfunc indices
	possibilities := make(map[int]map[int]struct{})
	for opcode := 0; opcode < 16; opcode++ {
		possibilities[opcode] = make(map[int]struct{})
		for idx := 0; idx < len(opFuncs); idx++ {
			possibilities[opcode][idx] = struct{}{}
		}
	}

	var part1 int
	for _, s := range samples {
		var match int
		for j, opFunc := range opFuncs {
			if s.after == opFunc(
				s.before,
				s.instruction.input1,
				s.instruction.input2,
				s.instruction.output,
			) {
				match++
			} else {
				delete(possibilities[s.instruction.opcode], j)
			}
		}
		if match >= 3 {
			part1++
		}
	}

	for len(possibilities) > 0 {
		for opcode, indices := range possibilities {
			if len(indices) == 1 {
				var toRemove int
				for i := range indices {
					toRemove = i
					ops[opcode] = opFuncs[i]
				}
				for _, indices2 := range possibilities {
					delete(indices2, toRemove)
				}
				delete(possibilities, opcode)
			}
		}
	}

	return part1
}

func runProgram(program []instruction) int {
	var curr registers
	for _, line := range program {
		curr = ops[line.opcode](curr, line.input1, line.input2, line.output)
	}
	return curr[0]
}
