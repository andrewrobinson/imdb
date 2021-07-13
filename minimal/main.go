package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func printResults(results []string) {
	fmt.Printf("full results:%v\n", results)
}

func printResultsSoFar(results []string) {
	fmt.Printf("results so far:%v\n", results)
}

func main() {

	var results []string

	strings := [5]string{"a", "b", "c", "d", "e"}

	resultsPipe := make(chan string, len(strings))

	//this is meant to signal "how execution ended", so we know if it is full or partial results being displayed
	//ie values of "done", "timedOut",
	finishedProcessingPipe := make(chan string)

	//produce to resultsPipe. Not sure how my actual code would interface with this.
	go func() {
		for _, n := range strings {
			resultsPipe <- n
			time.Sleep(600 * time.Millisecond)
		}
		//this does arrive - but before all the  result := <-resultsPipe cases have come in
		//with a sleep, or some time spent on real execution - no problem
		finishedProcessingPipe <- "done"
	}()

	shutdownSigTerm := make(chan os.Signal)
	signal.Notify(shutdownSigTerm, os.Interrupt, syscall.SIGTERM)

	shutdownSigInt := make(chan os.Signal)
	signal.Notify(shutdownSigInt, os.Interrupt, syscall.SIGINT)

LOOP:

	for {
		select {
		case <-time.Tick(time.Second * 30):
			fmt.Println("process timed out")

			//this doesn't arrive - why?
			// finishedProcessingPipe <- "timeOut"

			//b) this doesn't seemt to trigger the case sig. hence no results printed
			// shutdownSigTerm <- syscall.SIGTERM

			//so I added this instead. It is what happens at sigTerm
			printResultsSoFar(results)
			break LOOP
		case sigTerm := <-shutdownSigTerm:
			fmt.Printf("sigTerm signal %s received\n", sigTerm)
			//this doesn't arrive
			// finishedProcessingPipe <- "sigTerm"
			printResultsSoFar(results)
			break LOOP
		case sigInt := <-shutdownSigInt:
			fmt.Printf("sigInt signal %s received\n", sigInt)
			//this doesn't arrive
			// finishedProcessingPipe <- "sigInt"
			printResultsSoFar(results)
			break LOOP
		case message := <-finishedProcessingPipe:
			fmt.Printf("finishedProcessingPipe, message:%v\n", message)
			printResults(results)
		case result := <-resultsPipe:
			//this fires in an infinite loop if resultsPipe has been closed - why?
			fmt.Printf("intermediate result:%v\n", result)
			results = append(results, result)
		}
	}

}
