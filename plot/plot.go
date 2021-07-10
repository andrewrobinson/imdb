package plot

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/andrewrobinson/imdb/model"
)

func LookupPlots(filteredFileRows []model.FileRow, flags model.ProgramFlags) []model.FileRow {

	var rowsWithPlots []model.FileRow

	for _, fileRow := range filteredFileRows {

		plot := lookupPlot(fileRow.Tconst)

		fileRow.Plot = plot

		if flags.PlotFilterFlag != "" {

			match, _ := regexp.MatchString(flags.PlotFilterFlag, fileRow.Plot)

			if match {
				rowsWithPlots = append(rowsWithPlots, fileRow)
			}

		} else {
			rowsWithPlots = append(rowsWithPlots, fileRow)
		}

	}

	return rowsWithPlots
}

type IMDBResponse struct {
	Plot string
}

func lookupPlot(tconst string) string {

	fmt.Printf("Looking up plot for tconst:%v\n", tconst)

	//live location
	// "https://www.omdbapi.com/?i=tt0000075&apikey=591edae0"

	//currently served with a sleep by ../static/webserver.go:
	// http://localhost:3000/static/tt0000075.json

	var p IMDBResponse

	resp, err := http.Get("http://localhost:3000/static/tt0000075.json")

	if err != nil {
		log.Fatalln(err)
	}

	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		log.Fatalln(err)

	}

	//return "As an elegant maestro of mirage and delusion drapes his beautiful female assistant with a gauzy textile, much to our amazement, the lady vanishes into thin air.", nil

	fmt.Printf("returning plot for tconst:%v\n", tconst)
	return p.Plot

}

// func sleepFor(base int, random int) {
// 	rand.Seed(time.Now().UnixNano())
// 	n := rand.Intn(random)
// 	time.Sleep(time.Duration(base+n) * time.Millisecond)
// }

// func sleepForRandomTime() {
// 	rand.Seed(time.Now().UnixNano())
// 	n := rand.Intn(10) // n will be between 0 and 10
// 	// fmt.Printf("Sleeping %d milliseconds...\n", 10+n)
// 	time.Sleep(time.Duration(10+n) * time.Millisecond)
// 	// fmt.Println("Done")
// }
