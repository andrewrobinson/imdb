package plot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"go.uber.org/ratelimit"
)

// a modified version of:
// https://gist.github.com/montanaflynn/ea4b92ed640f790c4b9cee36046a5383

// a struct to hold the result from each request
type MontanaResult struct {
	Tconst string
	Plot   string
	// Res    http.Response
	// Err    error
}

// Sends requests in parallel but only up to a certain
// limit, and furthermore it's only parallel up to the amount of CPUs but
// is always concurrent up to the concurrency limit
func MontanaBoundedParallelGet(urls map[string]string, concurrencyLimit int) []MontanaResult {

	//see if a rate limiter helps - was getting localhost timeouts
	rl := ratelimit.New(100) // per second

	// this buffered channel will block at the concurrency limit
	semaphoreChan := make(chan struct{}, concurrencyLimit)

	// this channel will not block and collect the http request results
	resultsChan := make(chan *MontanaResult)

	// make sure we close these channels when we're done with them
	defer func() {
		close(semaphoreChan)
		close(resultsChan)
	}()

	// keep an index and loop through every url we will send a request to
	for key, url := range urls {

		// fmt.Sprintf("key:%v, url:%v", key, url)

		// start a go routine with the key and url in a closure
		go func(key string, url string) {

			// this sends an empty struct into the semaphoreChan which
			// is basically saying add one to the limit, but when the
			// limit has been reached block until there is room
			semaphoreChan <- struct{}{}

			// send the request and put the response in a result struct
			// along with the key so we can sort them later along with
			// any error that might have occoured

			//see if a rate limiter helps
			rl.Take()

			res, err := http.Get(url)

			if err != nil {
				fmt.Printf("XXX ERR while getting url for key:%v, url:%v, err:%+v", key, url, err)
				os.Exit(1)
			}

			var p IMDBResponse

			err = json.NewDecoder(res.Body).Decode(&p)
			if err != nil {
				fmt.Printf("XXX ERR while decoding json for key:%v, url:%v, err:%+v", key, url, err)
				os.Exit(1)
			}

			result := &MontanaResult{key, p.Plot}

			// now we can send the result struct through the resultsChan
			resultsChan <- result

			// once we're done it's we read from the semaphoreChan which
			// has the effect of removing one from the limit and allowing
			// another goroutine to start
			<-semaphoreChan

		}(key, url)
	}

	// make a slice to hold the results we're expecting
	var results []MontanaResult

	// start listening for any results over the resultsChan
	// once we get a result append it to the result slice
	for {
		result := <-resultsChan
		results = append(results, *result)

		// if we've reached the expected amount of urls then stop
		if len(results) == len(urls) {
			break
		}
	}

	// let's sort these results real quick
	// sort.Slice(results, func(i, j string) bool {
	// 	return results[i].tconst < results[j].tconst
	// })

	// now we're done we return the results
	return results
}
