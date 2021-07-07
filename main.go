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

//TODO - memory profile my current implementation
//I will have to use some buffer of memory if I want to improve performance?

//DONE - crude timing for my unperformant version

func main() {

	// now := time.Now()

	// fmt.Printf("%v - main() invoked\n", now)

	// maxRunTime := 30 * time.Second

	flags := common.BuildProgramFlags()
	// fmt.Printf("Flags passed: %+v\n", flags)

	// go processFile(flags)
	processFile(flags)

	//a
	// time.Sleep(maxRunTime)

	//b is equiv to a
	// sleep := time.After(maxRunTime)
	// <-sleep

	// now = time.Now()
	// fmt.Printf("%v - timed out\n", now)
	// os.Exit(1)

}

func processFile(flags model.ProgramFlags) {

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

	matches, highestLineNumber := filter.RunFilters(scanner, flags, true)

	now := time.Now()

	fmt.Printf("%v: processed ok, matches:%v from lines processed:%v\n", now, matches, highestLineNumber)

}
