package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
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

	printFlags := false
	printRows := true
	printMatches := true
	printDuration := true

	// plotLookuptemplate := "http://www.omdbapi.com/?i=%s&apikey=591edae0"

	//needs go run webserver/webserver.go running
	plotLookuptemplate := "http://localhost:3000/%s.json"

	// plotLookuptemplate := "https://raw.githubusercontent.com/andrewrobinson/imdb/207ba5bd2727dfadb65a3faccd6786a099dce5ef/static/tt0000075.json"

	start := time.Now()

	flags := common.BuildProgramFlags()
	if printFlags {
		fmt.Printf("Flags passed: %+v\n", flags)
	}

	//fake results to test channels
	var results []string
	strings := [5]string{"a", "b", "c", "d", "e"}
	resultsPipe := make(chan string, len(strings))

	//this is meant to signal "how execution ended", so we know if it is full or partial results being displayed
	//ie values of "done", "timedOut", "sigTerm", "sigInt"
	finishedProcessingPipe := make(chan string)

	go func() {

		//produce fake results to resultsPipe.
		for _, n := range strings {
			resultsPipe <- n
		}

		//This is the real results. It just prints output as it gets it
		//I'm tempted to try and produce to resultsPipe but the size of the results is unknown
		processFile(flags, printRows, printMatches, plotLookuptemplate)

		elapsed := time.Since(start)
		if printDuration {
			fmt.Printf("processFile finished, elapsed time:%v\n", elapsed)
		}

		//this works as long as processFile takes a while.
		//Without processFile being called, this message can arrive before any case result := <-resultsPipe have fired
		finishedProcessingPipe <- "done"
	}()

	//For convenience, since I'm not using Docker here yet, I need the program to stop on Ctrl-C as well as SIGTERM
	shutdownSigTerm := make(chan os.Signal)
	signal.Notify(shutdownSigTerm, os.Interrupt, syscall.SIGTERM)

	shutdownSigInt := make(chan os.Signal)
	signal.Notify(shutdownSigInt, os.Interrupt, syscall.SIGINT)

Loop:

	for {
		select {
		case <-time.Tick(time.Second * time.Duration(flags.MaxRunTimeFlag)):
			fmt.Printf("process timed out after %v seconds\n", flags.MaxRunTimeFlag)

			// Sending messages to channels inside the case statement doesn't seem to work
			// shutdownSigTerm <- syscall.SIGTERM
			// finishedProcessingPipe <- "timeOut"

			//so I added this instead of Tyrone's shutdownSigTerm <- syscall.SIGTERM. It is what happens at sigTerm
			printResultsSoFar(results)
			break Loop
		case sigTerm := <-shutdownSigTerm:
			fmt.Printf("sigTerm signal %s received\n", sigTerm)
			//this doesn't arrive
			// finishedProcessingPipe <- "sigTerm"
			printResultsSoFar(results)
			break Loop
		case sigInt := <-shutdownSigInt:
			fmt.Printf("sigInt signal %s received\n", sigInt)
			//this doesn't arrive
			// finishedProcessingPipe <- "sigInt"
			printResultsSoFar(results)
			break Loop
		case message := <-finishedProcessingPipe:
			fmt.Printf("finishedProcessingPipe, message:%v\n", message)
			printResults(results)
		case result := <-resultsPipe:
			//this fires in an infinite loop if resultsPipe has been closed - why?
			fmt.Printf("intermediate fake result:%v\n", result)
			results = append(results, result)
		}
	}

}

func printResults(results []string) {
	fmt.Printf("full results:%v\n", results)
}

func printResultsSoFar(results []string) {
	fmt.Printf("results so far:%v\n", results)
}

func processFile(flags model.ProgramFlags, printRows bool, printMatches bool, plotLookuptemplate string) {

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
			common.PrintRow(row)
		}
	}

	if printMatches {
		fmt.Printf("processed ok, matches:%v from lines processed:%v\n", len(filteredRowsWithPlots), highestLineNumber)
	}

}
