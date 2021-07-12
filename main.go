package main

import (
	"bufio"
	"fmt"
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
	output := make(chan string, 6)

	// messages := make(chan string, 6)

	// ret := fmt.Sprintf("s%v", 4)
	// fmt.Println(ret)

	go func() {
		//only the 1st one goes in .... ?
		output <- "andr"
		output <- "ping"
		output <- "ew"

		//only the 1st one goes in .... ?
		// By default sends and receives block until both the sender and receiver are ready.
		// This property allowed us to wait at the end of our program for the "ping" message
		// without having to use any other synchronization.
		// for i := 0; i < 5; i++ {
		// 	output <- fmt.Sprintf("s%v", i)
		// }

		//and this doesn't get printed
		// fmt.Println("output in goroutine %+v", output)

	}()
	msg0 := <-output
	msg1 := <-output
	msg2 := <-output

	//and it blocks here
	// msg3 := <-output

	fmt.Println(msg0)
	fmt.Println(msg1)
	fmt.Println(msg2)

	// fmt.Println(msg3)

	// fmt.Printf("output after make: %+v\n", output)

	// output <- "andr"

	// fmt.Printf("output after push1: %+v\n", output)

	// output <- "ew"

	// fmt.Printf("output after push2: %+v\n", output)

	// fmt.Printf("output: %+v\n", output)

}

func main2() {

	printFlags := false
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

	output := make(chan string)

	output <- "andr"

	// output <- processFile(flags, printRows, printMatches, plotLookuptemplate, output)

	processFile(flags, printRows, printMatches, plotLookuptemplate)

	output <- "ew"

	fmt.Printf("output: %+v\n", output)

	elapsed := time.Since(start)
	if printDuration {
		fmt.Printf("finished, elapsed time:%v\n", elapsed)
	}

}

func processFile(flags model.ProgramFlags, printRows bool, printMatches bool, plotLookuptemplate string) {
	// func processFile(flags model.ProgramFlags, printRows bool, printMatches bool, plotLookuptemplate string, output chan string) chan string {

	file, err := os.Open(flags.FilePathFlag)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	filteredRows, highestLineNumber := filter.RunFilters(scanner, flags)

	if len(filteredRows) > flags.MaxApiRequestsFlag {
		fmt.Printf("filteredRows size:%v larger than MaxApiRequestsFlag:%v, exiting.\n", len(filteredRows), flags.MaxApiRequestsFlag)
		os.Exit(1)
	} else {
		fmt.Printf("filteredRows size:%+v\n", len(filteredRows))
	}

	plotMap := plot.LookupPlotsInParallel(filteredRows, flags, plotLookuptemplate)

	filteredRowsWithPlots := plot.AddPlotsAndMaybeRegexFilter(filteredRows, plotMap, flags)

	if printRows {
		fmt.Println("IMDB_ID	|	Title	|	Plot")
		for _, row := range filteredRowsWithPlots {

			fmt.Printf("%v	|	%v	|	%v\n", row.Tconst, row.PrimaryTitle, row.Plot)
			// common.PrintRow(row)
		}
	}

	if printMatches {
		fmt.Printf("processed ok, matches:%v from lines processed:%v\n", len(filteredRowsWithPlots), highestLineNumber)
	}

	// output <- "andrew"
	// return output

}
