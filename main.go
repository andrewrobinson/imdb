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

/*
Problems with this code:

a) You only see the results upon timeout, or if you go Ctrl-C. Why can't you see results sooner?
b) shutdownSigTerm <- syscall.SIGTERM doesn't calls case sigTerm
c) I had to add L onto the break statement
d) case result := <-resultsPipe fires in an infinite loop if resultsPipe is closed. Why can't it be closed?
if it is closed, why can't that case not fire?

*/

func main() {

	var results []string

	strings := [5]string{"a", "b", "c", "d", "e"}

	resultsPipe := make(chan string, len(strings))

	//produce to resultsPipe
	go func() {
		for _, n := range strings {
			resultsPipe <- n
		}
	}()

	shutdownSigTerm := make(chan os.Signal)
	signal.Notify(shutdownSigTerm, os.Interrupt, syscall.SIGTERM)

	shutdownSigInt := make(chan os.Signal)
	signal.Notify(shutdownSigInt, os.Interrupt, syscall.SIGINT)

L:

	for {
		select {
		case <-time.Tick(time.Second * 30):
			fmt.Println("process timed out")
			//we time out via sending a sigterm, but could also do it via sigint

			//this doesn't seemt to trigger the case sig. hence no results printed
			// shutdownSigTerm <- syscall.SIGTERM

			//so I added this instead. It is what happens at sigTerm
			break L
		case sigTerm := <-shutdownSigTerm:
			fmt.Printf("sigTerm signal %s received\n", sigTerm)
			//I had to add the L
			break L
		case sigInt := <-shutdownSigInt:
			fmt.Printf("sigInt signal %s received\n", sigInt)
			//I had to add the L
			break L
		case result := <-resultsPipe:
			//this fires in an infinite loop if resultsPipe has been closed - why?
			results = append(results, result)
		}
	}

	fmt.Printf("results:%v\n", results)

}

func main4() {
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	var results []string
	resultsPipe := make(chan string)

	strings := [5]string{"a", "b", "c", "d"}

	go func() {
		for _, n := range strings {
			resultsPipe <- n
		}
		close(resultsPipe)
	}()

L:

	for {
		select {
		case <-time.Tick(time.Second * 30):
			fmt.Println("process timed out")
			shutdown <- syscall.SIGTERM
			// break L
		case sig := <-shutdown:
			fmt.Printf("shutdown signal %s recieved", sig)
			break L
		case result := <-resultsPipe:
			results = append(results, result)
		}
	}

	fmt.Printf("results:%v", results)
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

	processFile(flags, printRows, printMatches, plotLookuptemplate)

	elapsed := time.Since(start)
	if printDuration {
		fmt.Printf("finished, elapsed time:%v\n", elapsed)
	}

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
