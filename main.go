package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/andrewrobinson/imdb/common"
	"github.com/andrewrobinson/imdb/filter"
	"github.com/andrewrobinson/imdb/model"
	"github.com/andrewrobinson/imdb/plot"
)

/*

LookupPlotsInParallel using ConcurrencyFactor:20 and RateLimitPerSecond:100
When hitting https://raw.githubusercontent.com/andrewrobinson/imdb/207ba5bd2727dfadb65a3faccd6786a099dce5ef/static/tt0000075.json


//1 row
// go run main.go --filePath=../title.basics.tsv --titleType=short --primaryTitle=Conjuring --originalTitle=Escamotage --plotFilter=female

//10 rows - 4.24 seconds
//go run main.go --primaryTitle=Almodovar --filePath=../title.basics.tsv

//106 rows - 5 seconds
//go run main.go --primaryTitle=Conjuring --filePath=../title.basics.tsv

//290 rows - 6.96 seconds
//go run main.go --primaryTitle=Xavier --filePath=../title.basics.tsv

//1027 rows - 14.54 seconds
//go run main.go --primaryTitle=Stewart --filePath=../title.basics.tsv

//2206 rows - 26.2 seconds
//go run main.go --primaryTitle=Andrew --filePath=../title.basics.tsv

//4534 rows - 50 seconds
//go run main.go --primaryTitle=Adam --filePath=../title.basics.tsv

//16385 rows - 2 minutes 53 seconds
//go run main.go --primaryTitle=John --maxApiRequests=17000 --filePath=../title.basics.tsv

*/

// --maxRunTime=30 --filePath=../title.basics.tsv --concurrencyFactor=20

// go run main.go --filePath=../title.basics.tsv --titleType=short --primaryTitle=Conjuring --originalTitle=Escamotage --plotFilter=female

// go run main.go --titleType=short --primaryTitle=Conjuring --originalTitle=Escamotage --plotFilter=female

//go test . ./...

//http://localhost:3000/static/tt0000075.json

func main() {

	printFlags := true
	printRows := true
	printMatches := true
	printDuration := true

	// plotLookuptemplate := "http://www.omdbapi.com/?i=%s&apikey=591edae0"

	//needs go run webserver.go
	plotLookuptemplate := "http://localhost:3000/%s.json"

	// plotLookuptemplate := "https://raw.githubusercontent.com/andrewrobinson/imdb/207ba5bd2727dfadb65a3faccd6786a099dce5ef/static/tt0000075.json"

	start := time.Now()

	flags := common.BuildProgramFlags()
	if printFlags {
		fmt.Printf("Flags passed: %+v\n", flags)
	}

	// c := moring("Joe")

	//erm you can't put go in front here?
	resultsPipe := processFile(flags, printRows, printMatches, plotLookuptemplate)

	timeout := time.After(5 * time.Second)

L:

	for {
		select {
		case s := <-resultsPipe:
			fmt.Println(s)
		case <-timeout:
			fmt.Println("You talk too much.")
			break L
		}
	}

	elapsed := time.Since(start)
	if printDuration {
		fmt.Printf("finished, elapsed time:%v\n", elapsed)
	}

}

// func boring(msg string) <-chan string { // Returns receive-only channel of strings.
// 	c := make(chan string)
// 	go func() { // We launch the goroutine from inside the function.
// 		for i := 0; ; i++ {
// 			c <- fmt.Sprintf("%s %d", msg, i)
// 			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
// 		}
// 	}()
// 	return c // Return the channel to the caller.
// }

func moring(msg string) <-chan string { // Returns receive-only channel of strings.
	resultsPipe := make(chan string)
	go func() { // We launch the goroutine from inside the function.
		for i := 0; ; i++ {
			resultsPipe <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()
	return resultsPipe // Return the channel to the caller.
}

func processFile(flags model.ProgramFlags, printRows bool, printMatches bool, plotLookuptemplate string) <-chan string {

	resultsPipe := make(chan string)

	file, err := os.Open(flags.FilePathFlag)
	if err != nil {
		resultsPipe <- fmt.Sprintf("%v", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	filteredRows, highestLineNumber := filter.RunFilters(scanner, flags)

	if len(filteredRows) > flags.MaxApiRequestsFlag {
		resultsPipe <- fmt.Sprintf("filteredRows size:%v larger than MaxApiRequestsFlag:%v, exiting.\n", len(filteredRows), flags.MaxApiRequestsFlag)
		os.Exit(1)
	} else {
		resultsPipe <- fmt.Sprintf("filteredRows size:%+v\n", len(filteredRows))
	}

	plotMap := plot.LookupPlotsInParallel(filteredRows, flags, plotLookuptemplate)

	filteredRowsWithPlots := plot.AddPlotsAndMaybeRegexFilter(filteredRows, plotMap, flags)

	if printRows {
		resultsPipe <- fmt.Sprintf("IMDB_ID	|	Title	|	Plot")
		for _, row := range filteredRowsWithPlots {
			resultsPipe <- fmt.Sprintf("%v	|	%v	|	%v\n", row.Tconst, row.PrimaryTitle, row.Plot)
		}
	}

	if printMatches {
		resultsPipe <- fmt.Sprintf("processed ok, matches:%v from lines processed:%v\n", len(filteredRowsWithPlots), highestLineNumber)
	}

	return resultsPipe

}
