package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

/*

This serves up a static json file for when main.go is configured with

plotLookuptemplate := "http://localhost:3000/%s.json"

(imdb gives you only 1000 plot lookups a day)

*/

func main() {
	mux := http.NewServeMux()

	th := &timeHandler{format: time.RFC1123}
	// mux.Handle("/static/([^/]+)", th)
	mux.Handle("/", th)

	log.Println("Listening...")
	http.ListenAndServe(":3000", mux)
}

type timeHandler struct {
	format string
}

func (th *timeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// fmt.Printf("request in:%+v\n", r)

	start := time.Now()

	sleepFor(30, 10)

	ret := `{
		"Title": "Escamotage d'une dame au théâtre Robert Houdin",
		"Year": "1896",
		"Rated": "N/A",
		"Released": "01 Oct 1896",
		"Runtime": "1 min",
		"Genre": "Short, Horror",
		"Director": "Georges Méliès",
		"Writer": "N/A",
		"Actors": "Jehanne d'Alcy, Georges Méliès",
		"Plot": "As an elegant maestro of mirage and delusion drapes his beautiful female assistant with a gauzy textile, much to our amazement, the lady vanishes into thin air.",
		"Language": "None",
		"Country": "France",
		"Awards": "N/A",
		"Poster": "https://m.media-amazon.com/images/M/MV5BNGRhNTcxMDMtYTMyMi00ZTIxLThiOWUtMTgwZDA2Njk4YTFjXkEyXkFqcGdeQXVyNDE5MTU2MDE@._V1_SX300.jpg",
		"Ratings": [
		{
		"Source": "Internet Movie Database",
		"Value": "6.3/10"
		}
		],
		"Metascore": "N/A",
		"imdbRating": "6.3",
		"imdbVotes": "1,665",
		"imdbID": "tt0000075",
		"Type": "movie",
		"DVD": "N/A",
		"BoxOffice": "N/A",
		"Production": "N/A",
		"Website": "N/A",
		"Response": "True"
		}`

	w.Write([]byte(ret))

	elapsed := time.Since(start)

	fmt.Printf("file served, elapsed time:%v\n", elapsed)

}

func sleepFor(base int, random int) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(random)
	time.Sleep(time.Duration(base+n) * time.Millisecond)
}
