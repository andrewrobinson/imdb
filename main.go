package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/andrewrobinson/imdb/common"
	"github.com/andrewrobinson/imdb/filter"
	"github.com/andrewrobinson/imdb/model"
	"github.com/andrewrobinson/imdb/plot"
)

// --maxRunTime=30 --filePath=../title.basics.tsv

// go run main.go --titleType=short --primaryTitle=Conjuring --originalTitle=Escamotage --plotFilter=female

//go test . ./...

//http://localhost:3000/static/tt0000075.json

func main() {

	printFlags := false
	printRows := false
	printMatches := true
	printDuration := true

	start := time.Now()

	flags := common.BuildProgramFlags()
	if printFlags {
		fmt.Printf("Flags passed: %+v\n", flags)
	}

	processFile(flags, printRows, printMatches)

	elapsed := time.Since(start)
	if printDuration {
		fmt.Printf("finished, elapsed time:%v\n", elapsed)
	}

}

func processFile(flags model.ProgramFlags, printRows bool, printMatches bool) {

	file, err := os.Open(flags.FilePathFlag)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	filteredFileRows, highestLineNumber := filter.RunFilters(scanner, flags)

	rowsWithPlots := plot.LookupPlotsAndMaybeRegexThem(filteredFileRows, flags)

	if printRows {
		fmt.Printf("filteredFileRows:%+v\n", rowsWithPlots)
	}

	if printMatches {
		fmt.Printf("processed ok, matches:%v from lines processed:%v\n", len(rowsWithPlots), highestLineNumber)
	}

}
