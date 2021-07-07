package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/andrewrobinson/imdb/common"
	"github.com/andrewrobinson/imdb/model"
	// "PrintFields"
)

// go run . --titleType=short --primaryTitle=Conjuring --originalTitle=Escamotage
// go run . --titleType=short --primaryTitle=Conjuring --originalTitle=Escamotage | wc -l
// 4 (includes Flags printing lines etc - 1 result)go run

// go run . --originalTitle=Clown --genres=Comedy
// go run . --originalTitle=Clown --genres=Comedy | wc -l
// 4

// go run . --genres=Documentary
// go run . --genres=Documentary | wc -l
// 40

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
	for scanner.Scan() {

		lineNumber++
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		rowStruct := common.BuildFileRow(fields)
		printMatchingLines(rowStruct, flags)

	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("something bad happened in the line %v: %v", lineNumber, err)
	} else {
		fmt.Printf("processed ok and reached lineNumber: %v", lineNumber)
	}

}

func printMatchingLines(row model.FileRow, flags model.ProgramFlags) {

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
		common.PrintFields(row)
	}

	//simulate the unpredictable time taken for the plot lookup
	//sleepForRandomTime()

}

// func sleepForRandomTime() {
// 	rand.Seed(time.Now().UnixNano())
// 	n := rand.Intn(10) // n will be between 0 and 10
// 	//fmt.Printf("Sleeping %d seconds...\n", n)
// 	time.Sleep(time.Duration(n) * time.Second)
// 	//fmt.Println("Done")
// }

//moved to helpers.go
// func printFields(row FileRow) {
// 	//For now just print out the fields, but later output must be
// 	// IMDB_ID     |   Title               |   Plot
// 	// tt0000005   |   Blacksmith Scene    |   Three men hammer on an anvil and pass a bottle of beer around.

// 	// fmt.Printf("%+v\n", row)
// 	fmt.Println(row)
// }

func flagMatchesOrIsEmpty(filterValue string, columnValue string) bool {

	//if no flag value passed then don't filter, ie it passes
	if filterValue == "" {
		return true
	} else {
		return strings.Contains(columnValue, filterValue)
	}

}
