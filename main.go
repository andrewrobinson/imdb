package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/andrewrobinson/imdb/common"
	"github.com/andrewrobinson/imdb/filter"
	"github.com/andrewrobinson/imdb/model"
)

//go run main.go --titleType=short --primaryTitle=Conjuring --originalTitle=Escamotage --plotFilter=female
// processed ok, matches:1 from lines processed:75
// finished, elapsed time:22.457956ms
// So 22.4ms for an 6 805 bytes (8 KB on disk) file

//highmem gives
// len stringContent: 6805
// finished, elapsed time:444.098µs

//go run main.go --filePath=../title.basics.tsv --titleType=short --primaryTitle=Conjuring --originalTitle=Escamotage --plotFilter=female
// processed ok, matches:1 from lines processed:8061101
// finished, elapsed time:5.437130106s
// So 5.4s for an 689 049 864 bytes (689,1 MB on disk) file

//highmem gives
// len stringContent: 689049864
// finished, elapsed time:3.742361676s

//TODO - memory profile
//TODO - docker - for profiling environment too

//TODO - time some alternative ways of reading this file
// https://hackernoon.com/leveraging-multithreading-to-read-large-files-faster-in-go-lmn32t7
// https://blog.cloudboost.io/reading-humongous-files-in-go-c894b05ac020
// https://stackoverflow.com/questions/52154609/fastest-way-of-reading-huge-file-in-go-lang-with-small-ram/52154800

//TODO - try a memory hungry version, ie slurp the whole thing in.
//the current impl is a low memory version

func main() {

	printFlags := false
	printRows := true
	printMatches := true
	printDuration := true

	useLowMemBufioScanner := true
	useHighMem := false

	start := time.Now()

	//fmt.Printf("%v - main() invoked\n", start)
	//maxRunTime := 30 * time.Second

	flags := common.BuildProgramFlags()
	if printFlags {
		fmt.Printf("Flags passed: %+v\n", flags)
	}

	// go processFile(flags, false)
	if useLowMemBufioScanner {
		processFileWithBufioScanner(flags, printRows, printMatches)
	}

	if useHighMem {
		processFileHighMem(flags, printRows, printMatches)
	}

	//a
	// time.Sleep(maxRunTime)

	//b is equiv to a
	// sleep := time.After(maxRunTime)
	// <-sleep

	elapsed := time.Since(start)
	if printDuration {
		fmt.Printf("finished, elapsed time:%v\n", elapsed)
	}
	// os.Exit(1)

}

func processFileHighMem(flags model.ProgramFlags, printRows bool, printMatches bool) {

	content, err := ioutil.ReadFile(flags.FilePathFlag)
	if err != nil {
		fmt.Println("Err")
	}

	stringContent := string(content)

	lines := strings.Split(stringContent, "\n")

	fmt.Printf("len stringContent: %v\n", len(stringContent))
	fmt.Printf("len lines: %v\n", len(lines))

	matches, highestLineNumber := filter.RunFiltersHighMem(lines, flags, printRows)

	if printMatches {
		fmt.Printf("processed ok, matches:%v from lines processed:%v\n", matches, highestLineNumber)
	}

}

func processFileWithBufioScanner(flags model.ProgramFlags, printRows bool, printMatches bool) {

	file, err := os.Open(flags.FilePathFlag)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	//https://golangdocs.com/reading-files-in-golang
	//https://devmarkpro.com/working-big-files-golang
	//https://golang.org/pkg/bufio/#Scanner
	//https://stackoverflow.com/questions/64638136/performance-issues-while-reading-a-file-line-by-line-with-bufio-newscanner
	scanner := bufio.NewScanner(file)

	matches, highestLineNumber := filter.RunFilters(scanner, flags, printRows)

	if printMatches {
		fmt.Printf("processed ok, matches:%v from lines processed:%v\n", matches, highestLineNumber)
	}

}
