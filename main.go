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

	// a := make([]FileRow, 5) // len(a)=5

	// b := make([]FileRow, 0, 5) // len(b)=0, cap(b)=5
	// lines []string

	matches, highestLineNumber := filter.RunFilters(scanner, flags, printRows)

	if printMatches {
		fmt.Printf("processed ok, matches:%v from lines processed:%v\n", matches, highestLineNumber)
	}

}
