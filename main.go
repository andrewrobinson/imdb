package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

// go run main.go --titleType=short --primaryTitle=Conjuring --originalTitle=Escamotage
// go run main.go --originalTitle=Clown --genres=Comedy
// go run main.go --genres=Documentary

// TODO
// - maxApiRequests - maximum number of requests to be made to [omdbapi](https://www.omdbapi.com/)
// - maxRunTime - maximum run time of the application. Format is a `time.Duration` string see [here](https://godoc.org/time#ParseDuration)
// - maxRequests - maximum number of requests to send to [omdbapi](https://www.omdbapi.com/)
// - plotFilter - regex pattern to apply to the plot of a film retrieved from [omdbapi](https://www.omdbapi.com/)

func main() {

	filePathFlag := flag.String("filePath", "title.basics.truncated.tsv", "")
	titleTypeFlag := flag.String("titleType", "", "")
	primaryTitleFlag := flag.String("primaryTitle", "", "")
	originalTitleFlag := flag.String("originalTitle", "", "")
	startYearFlag := flag.String("startYear", "", "")
	endYearFlag := flag.String("endYear", "", "")
	runtimeMinutesFlag := flag.String("runtimeMinutes", "", "")
	genresFlag := flag.String("genres", "", "")

	flag.Parse()
	fmt.Print("\nMatches:\n")

	file, err := os.Open(*filePathFlag)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		line := scanner.Text()
		fields := strings.Split(line, "\t")

		//unlimited goroutines
		go printMatchingLines(fields, titleTypeFlag, primaryTitleFlag, originalTitleFlag,
			startYearFlag, endYearFlag, runtimeMinutesFlag, genresFlag)

	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

}

func printMatchingLines(fields []string, titleTypeFlag *string, primaryTitleFlag *string, originalTitleFlag *string, startYearFlag *string, endYearFlag *string, runtimeMinutesFlag *string, genresFlag *string) {

	titleType, primaryTitle, originalTitle := fields[1], fields[2], fields[3]
	startYear, endYear, runtimeMinutes, genres := fields[5], fields[6], fields[7], fields[8]

	titleTypeMatches := flagMatchesOrIsEmpty(*titleTypeFlag, titleType)
	primaryTitleMatches := flagMatchesOrIsEmpty(*primaryTitleFlag, primaryTitle)
	originalTitleMatches := flagMatchesOrIsEmpty(*originalTitleFlag, originalTitle)
	startYearMatches := flagMatchesOrIsEmpty(*startYearFlag, startYear)
	endYearMatches := flagMatchesOrIsEmpty(*endYearFlag, endYear)
	runtimeMinutesMatches := flagMatchesOrIsEmpty(*runtimeMinutesFlag, runtimeMinutes)
	genresMatches := flagMatchesOrIsEmpty(*genresFlag, genres)

	if titleTypeMatches && primaryTitleMatches && originalTitleMatches && startYearMatches && endYearMatches && runtimeMinutesMatches && genresMatches {
		fmt.Println(fields)
	}

}

func flagMatchesOrIsEmpty(filterValue string, columnValue string) bool {

	//if no flag value passed then don't filter, ie it passes
	if filterValue == "" {
		return true
	} else {
		return strings.Contains(columnValue, filterValue)
	}

}
