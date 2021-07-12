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

//TODO - cancel this sigterm stuff and build it up slowly
//just get more channels in there 1st
//ie stop printing inside and get output out via channels

//The "kill" command in linux if specified without any signal number like -9, will send SIGTERM
//SIGINT is the signal generated when a user presses Ctrl+C

// Task 2.3: The program should exit when the `maxRunTime` is exceeded. At this point do not worry about a graceful exit or any resource cleanup.
// Task 2.4: The program should exit when stopped with `SIGTERM`. At this point do not worry about a graceful exit or any resource cleanup.

// Task 5.1: When your program exceeds the `maxRunTime`, it should gracefully exit. All resources should be released and any running goroutines should be stopped. Any results (if available should be printed)
// Task 5.2: When your program receives the `SIGTERM` signal, it should gracefully exit. All resources should be released and any running goroutines should be stopped. No results (if any) should be printed. A simple message such as "program exiting" is acceptable.

func main() {
	c := boring("Joe")
	timeout := time.After(5 * time.Second)
	for {
		select {
		case s := <-c:
			fmt.Println(s)
		case <-timeout:
			fmt.Println("You talk too much.")
			return
		}
	}
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

	//I'm imagining something along the following lines where resultsPipe is
	//some channel being fed with results which are being processed in a go func {}()
	//routine and your printout happens when you break out of the for loop.

	resultsPipe := make(chan string)

	go func() { resultsPipe <- "ping" }()

	go processFile(flags, printRows, printMatches, plotLookuptemplate, resultsPipe)

	var results []string

	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT)

	for {
		select {
		case <-time.Tick(time.Second * time.Duration(flags.MaxRunTimeFlag)):
			fmt.Println("process timed out")
			shutdown <- syscall.SIGINT
		case sig := <-shutdown:
			fmt.Printf("shutdown signal %s recieved", sig)
			break
		case result := <-resultsPipe:
			results = append(results, result)
		}
	}

	fmt.Println("outside loop, results:%v", results)

	// msg := <-resultsPipe
	// fmt.Println(msg)

	// I personally like using piping for processing multi-step processes.
	// I also find it allows you to use worker pools or rate limiters in conjunction with buffered channels
	// to achieve a certain amount of parallelism

	// https://blog.golang.org/pipelines

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

func processFile(flags model.ProgramFlags, printRows bool, printMatches bool, plotLookuptemplate string) <-chan string {

	resultsPipe := make(chan string)

	resultsPipe <- "ming"

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
		// fmt.Println("IMDB_ID	|	Title	|	Plot")
		resultsPipe <- "IMDB_ID	|	Title	|	Plot"
		for _, row := range filteredRowsWithPlots {
			// common.PrintRow(row)
			resultsPipe <- fmt.Sprintf("%v	|	%v	|	%v\n", row.Tconst, row.PrimaryTitle, row.Plot)
		}
	}

	if printMatches {
		// fmt.Printf("processed ok, matches:%v from lines processed:%v\n", len(filteredRowsWithPlots), highestLineNumber)
		resultsPipe <- fmt.Sprintf("processed ok, matches:%v from lines processed:%v\n", len(filteredRowsWithPlots), highestLineNumber)
	}

	return resultsPipe

}
