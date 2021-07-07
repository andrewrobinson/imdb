package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/andrewrobinson/imdb/common"
	"github.com/andrewrobinson/imdb/filter"
	"github.com/andrewrobinson/imdb/model"
)

//go run main.go --titleType=short --primaryTitle=Conjuring --originalTitle=Escamotage --plotFilter=female
// processed ok, matches:1 from lines processed:75
// finished, elapsed time:22.457956ms

//go run main.go --filePath=../title.basics.tsv --titleType=short --primaryTitle=Conjuring --originalTitle=Escamotage --plotFilter=female
// processed ok, matches:1 from lines processed:8061101
// finished, elapsed time:5.437130106s

//TODO - memory profile my current implementation
//I will have to use some buffer of memory if I want to improve performance?

func main() {

	printFlags := false
	printRows := false
	printMatches := true
	printDuration := true

	start := time.Now()

	//fmt.Printf("%v - main() invoked\n", start)
	//maxRunTime := 30 * time.Second

	flags := common.BuildProgramFlags()
	if printFlags {
		fmt.Printf("Flags passed: %+v\n", flags)
	}

	// go processFile(flags, false)
	processFile(flags, printRows, printMatches)

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

func processFile(flags model.ProgramFlags, printRows bool, printMatches bool) {

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
