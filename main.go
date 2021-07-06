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

func main() {

	filePathFlag := flag.String("filePath", "title.basics.truncated.tsv", "")
	titleTypeFlag := flag.String("titleType", "", "")
	primaryTitleFlag := flag.String("primaryTitle", "Conjuring", "")
	originalTitleFlag := flag.String("originalTitle", "Escamotage", "")

	// numbPtr := flag.Int("numb", 42, "an int")
	// boolPtr := flag.Bool("fork", false, "a bool")

	flag.Parse()
	fmt.Println("filePathFlag:", *filePathFlag)
	fmt.Println("titleType:", *titleTypeFlag)
	fmt.Println("primaryTitleFlag:", *primaryTitleFlag)
	fmt.Println("originalTitleFlag:", *originalTitleFlag)

	file, err := os.Open(*filePathFlag)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		//now filter by all the filters

		titleType, primaryTitle, originalTitle := fields[1], fields[2], fields[3]

		if *titleTypeFlag != "" {
			if strings.Contains(titleType, *titleTypeFlag) {
				fmt.Println(fields)
			}

		}

		if *primaryTitleFlag != "" {
			if strings.Contains(primaryTitle, *primaryTitleFlag) {
				fmt.Println(fields)
			}

		}

		if *originalTitleFlag != "" {
			if strings.Contains(originalTitle, *originalTitleFlag) {
				fmt.Println(fields)
			}

		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

}
