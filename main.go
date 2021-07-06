package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

// - filePath - absolute path to the inflated `title.basics.tsv.gz` file
// - titleType - filter on `titleType` column
// - primaryTitle - filter on `primaryTitle` column
// - originalTitle - filter on `originalTitle` column

// - genre - filter on `genre` column
// - startYear - filter on `startYear` column
// - endYear - filter on `endYear` column
// - runtimeMinutes - filter on `runtimeMinutes` column
// - genres - filter on `genres` column
// - maxApiRequests - maximum number of requests to be made to [omdbapi](https://www.omdbapi.com/)
// - maxRunTime - maximum run time of the application. Format is a `time.Duration` string see [here](https://godoc.org/time#ParseDuration)
// - maxRequests - maximum number of requests to send to [omdbapi](https://www.omdbapi.com/)
// - plotFilter - regex pattern to apply to the plot of a film retrieved from [omdbapi](https://www.omdbapi.com/)

// go run main.go --titleType=short --primaryTitle=Conjuring --originalTitle=Escamotage

func main() {

	filePathFlag := flag.String("filePath", "title.basics.truncated.tsv", "")
	titleTypeFlag := flag.String("titleType", "", "filter on `titleType` column")
	primaryTitleFlag := flag.String("primaryTitle", "", "filter on `primaryTitle` column")
	originalTitleFlag := flag.String("originalTitle", "", "filter on `originalTitle` column")

	// numbPtr := flag.Int("numb", 42, "an int")
	// boolPtr := flag.Bool("fork", false, "a bool")

	flag.Parse()
	fmt.Println("\nflag values passed:")
	fmt.Println("filePathFlag:", *filePathFlag)
	fmt.Println("titleType:", *titleTypeFlag)
	fmt.Println("primaryTitleFlag:", *primaryTitleFlag)
	fmt.Println("originalTitleFlag:", *originalTitleFlag)
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

		titleType, primaryTitle, originalTitle := fields[1], fields[2], fields[3]

		titleTypeMatches := flagMatchesOrIsEmpty(*titleTypeFlag, titleType)
		primaryTitleMatches := flagMatchesOrIsEmpty(*primaryTitleFlag, primaryTitle)
		originalTitleMatches := flagMatchesOrIsEmpty(*originalTitleFlag, originalTitle)

		if titleTypeMatches && primaryTitleMatches && originalTitleMatches {
			fmt.Println(fields)
		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

}

func flagMatchesOrIsEmpty(filterValue string, columnValue string) bool {

	//if no flag value, then don't filter, ie it passes
	if filterValue == "" {
		return true
	}

	return strings.Contains(columnValue, filterValue)

}
