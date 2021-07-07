package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/andrewrobinson/imdb/common"
	"github.com/andrewrobinson/imdb/model"
)

// go run main.go
//processed ok, matches:75 from lines processed:75

// go run main.go --genres=Comedy
// processed ok, matches:7 from lines processed:75

// go run main.go --genres=Short
// processed ok, matches:73 from lines processed:75

// go run main.go --genres=Comedy,Short
// processed ok, matches:4 from lines processed:75

// go run main.go --genres=Animation,Comedy,Romance
// processed ok, matches:1 from lines processed:75

// go run main.go --genres=Comedy,Romance
// processed ok, matches:1 from lines processed:75

// go run main.go --genres=Documentary
// processed ok, matches:37 from lines processed:75

// go run main.go --originalTitle=Clown
//processed ok, matches:1 from lines processed:75

// go run main.go --originalTitle=Clown --genres=Comedy
// processed ok, matches:1 from lines processed:75

// go run main.go --originalTitle=Clown --genres=medy
// processed ok, matches:1 from lines processed:75

// go run main.go --originalTitle=Clown --genres=Dramedy
// processed ok, matches:0 from lines processed:75

// go run main.go --titleType=short --primaryTitle=Conjuring --originalTitle=Escamotage
// processed ok, matches:1 from lines processed:75%

func main() {

	flags := common.BuildProgramFlags()
	fmt.Printf("Flags passed: %+v\n", flags)

	fmt.Print("\nMatches:\n")

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
	lineNumber := 0
	matches := 0
	for scanner.Scan() {

		lineNumber++
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		fileRow := common.BuildFileRow(fields)
		if rowMatchesFlags(fileRow, flags) {
			matches++
			common.PrintFields(fileRow)
		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error on line %v: %v", lineNumber, err)
	} else {
		fmt.Printf("processed ok, matches:%v from lines processed:%v", matches, lineNumber)
	}

}

func rowMatchesFlags(row model.FileRow, flags model.ProgramFlags) bool {

	// TODO - examine buffered output
	// https://stackoverflow.com/questions/64638136/performance-issues-while-reading-a-file-line-by-line-with-bufio-newscanner

	titleTypeMatches := flagMatchesOrIsEmpty(flags.TitleTypeFlag, row.TitleType)
	primaryTitleMatches := flagMatchesOrIsEmpty(flags.PrimaryTitleFlag, row.PrimaryTitle)
	originalTitleMatches := flagMatchesOrIsEmpty(flags.OriginalTitleFlag, row.OriginalTitle)
	startYearMatches := flagMatchesOrIsEmpty(flags.StartYearFlag, row.StartYear)
	endYearMatches := flagMatchesOrIsEmpty(flags.EndYearFlag, row.EndYear)
	runtimeMinutesMatches := flagMatchesOrIsEmpty(flags.RuntimeMinutesFlag, row.RuntimeMinutes)
	genresMatches := flagMatchesOrIsEmpty(flags.GenresFlag, row.Genres)

	if titleTypeMatches && primaryTitleMatches && originalTitleMatches && startYearMatches && endYearMatches && runtimeMinutesMatches && genresMatches {
		return true
	} else {
		return false
	}

	//simulate the unpredictable time taken for the plot lookup
	//sleepForRandomTime()

}

func flagMatchesOrIsEmpty(filterValue string, columnValue string) bool {

	//if no flag value passed then don't filter, ie it passes
	if filterValue == "" {
		return true
	} else {
		return strings.Contains(columnValue, filterValue)
	}

}
