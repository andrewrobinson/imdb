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

func main() {

	maxRunTime := 30 * time.Second

	flags := common.BuildProgramFlags()
	//fmt.Printf("Flags passed: %+v\n", flags)

	go processFile(flags)

	time.Sleep(maxRunTime)
	fmt.Println("timed out")
	os.Exit(0)

}

func processFile(flags model.ProgramFlags) {

	file, err := os.Open(flags.FilePathFlag)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	//https://golangdocs.com/reading-files-in-golang
	// https://devmarkpro.com/working-big-files-golang
	//https://golang.org/pkg/bufio/#Scanner
	// https://stackoverflow.com/questions/64638136/performance-issues-while-reading-a-file-line-by-line-with-bufio-newscanner
	scanner := bufio.NewScanner(file)

	matches, highestLineNumber := filter.RunFiltersAndPrint(scanner, flags, true)

	fmt.Printf("processed ok, matches:%v from lines processed:%v\n", matches, highestLineNumber)

}
