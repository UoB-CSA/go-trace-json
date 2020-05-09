package main

import (
	"./trace"
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

type gtype struct {
	ID       uint64 // Unique identifier (PC).
	Name     string // Start function.
	N        int    // Total number of goroutines in this group.
	ExecTime int64  // Total execution time of all goroutines in this group.
}

type jsonOutput struct {
	Goroutines []gtype `json:"groroutines"`
}

func main() {
	inputFile := flag.String("input-file", "", "The input file. Required")
	outputFile := flag.String("output-file", "trace.json", "The output file. (default: traces.json)")
	flag.Parse()

	if *inputFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	parsedTraces, err := parseTrace(*inputFile)
	if err != nil {
		os.Exit(1)
	}
	events := parsedTraces.Events

	// Get goroutines
	gs := trace.GoroutineStats(events)
	gss := make(map[uint64]gtype)
	for _, g := range gs {
		gs1 := gss[g.PC]
		gs1.ID = g.PC
		gs1.Name = g.Name
		gs1.N++
		gs1.ExecTime += g.ExecTime
		gss[g.PC] = gs1
	}
	var glist []gtype
	for k, v := range gss {
		v.ID = k
		glist = append(glist, v)
		fmt.Printf("%s [%d]\n", v.Name, v.N)
	}

	jsonData := &jsonOutput{
		Goroutines: glist,
	}
	jsonString, err := json.Marshal(jsonData)
	if err != nil {
		fmt.Errorf("Failed to parse to JSON")
		os.Exit(1)
	}
	err = ioutil.WriteFile(*outputFile, jsonString, 0644)
	if err != nil {
		fmt.Errorf("Failed to write to output file")
		os.Exit(1)
	}

}

var loader struct {
	once sync.Once
	res  trace.ParseResult
	err  error
}

func parseTrace(traceFile string) (trace.ParseResult, error) {
	loader.once.Do(func() {
		tracef, err := os.Open(traceFile)
		if err != nil {
			loader.err = fmt.Errorf("failed to open trace file: %v", err)
			return
		}
		defer tracef.Close()

		// Parse and symbolize.
		res, err := trace.Parse(bufio.NewReader(tracef), "")
		if err != nil {
			loader.err = fmt.Errorf("failed to parse trace: %v", err)
			return
		}
		loader.res = res
	})
	return loader.res, loader.err
}
