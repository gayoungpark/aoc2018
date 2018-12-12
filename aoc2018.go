package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
	"strconv"
)

var funcs = []func(){
	problem01,
	problem02,
	problem03,
	problem04,
	problem05,
	problem06,
	problem07,
	problem08,
	problem09,
	problem10,
	problem11,
	problem12,
}

func main() {
	log.SetFlags(0)

	profileFilename := flag.String("cpuprofile", "", "write cpu profile to `file`")
	flag.Parse()

	_ = profileFilename

	if len(flag.Args()) < 1 {
		log.Fatalf("Usage: %s [flags] <problem>", os.Args[0])
	}
	n, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		log.Fatalln("Bad problem number:", err)
	}
	if n < 1 || n > len(funcs) {
		log.Fatalf("Problem number out of range: %d", n)
	}

	if *profileFilename != "" {
		f, err := os.Create(*profileFilename)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Fatal(err)
			}
		}()

		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatalln("Cannot start profile:", err)
		}
		defer pprof.StopCPUProfile()
	}

	funcs[n-1]()
}
