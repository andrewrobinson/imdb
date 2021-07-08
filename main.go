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

// --maxRunTime=30 --filePath=../title.basics.tsv

// go run main.go --titleType=short --primaryTitle=Conjuring --originalTitle=Escamotage --plotFilter=female

//go test . ./...

func main() {

	// https://www.yellowduck.be/posts/graceful-shutdown/
	// https://www.yellowduck.be/posts/waitgroup-channels/
	// https://medium.com/code-zen/concurrency-in-go-5fcba11acb0f
	// https://stackoverflow.com/questions/36056615/what-is-the-advantage-of-sync-waitgroup-over-channels

	printFlags := false
	printRows := true
	printMatches := true
	printDuration := true

	start := time.Now()

	flags := common.BuildProgramFlags()
	if printFlags {
		fmt.Printf("Flags passed: %+v\n", flags)
	}

	//fmt.Printf("%v - main() invoked\n", start)
	// maxRunTime := time.Duration(flags.MaxRunTimeFlag) * time.Second
	// fmt.Printf("XX maxRunTime:%v\n", maxRunTime)

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

	scanner := bufio.NewScanner(file)

	filteredFileRows, highestLineNumber := filter.RunFilters(scanner, flags)

	rowsWithPlots := lookupPlots(filteredFileRows)

	if printRows {
		fmt.Printf("filteredFileRows:%+v\n", rowsWithPlots)
	}

	if printMatches {
		fmt.Printf("processed ok, matches:%v from lines processed:%v\n", len(rowsWithPlots), highestLineNumber)
	}

}

func lookupPlots(filteredFileRows []model.FileRow) []model.FileRow {

	var rowsWithPlots []model.FileRow

	for _, fileRow := range filteredFileRows {

		plot, err := common.LookupPlot(fileRow.Tconst)

		if err != nil {
			fmt.Printf("error while looking up plots:%+v\n", err)
			os.Exit(1)
		}

		fileRow.Plot = plot

		rowsWithPlots = append(rowsWithPlots, fileRow)

	}

	return rowsWithPlots
}
