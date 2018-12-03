package main

import (
	"log"
	"os"
	"strconv"
)

var funcs = []func(){
	problem01,
	problem02,
	problem03,
}

func main() {
	log.SetFlags(0)
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <problem>", os.Args[0])
	}
	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalln("Bad problem number:", err)
	}
	if n < 1 || n > len(funcs) {
		log.Fatalf("Problem number out of range: %d", n)
	}
	funcs[n-1]()
}
