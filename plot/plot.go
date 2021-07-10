package plot

import (
	"encoding/json"
	"log"
	"regexp"

	"github.com/andrewrobinson/imdb/model"
)

type IMDBResponse struct {
	Plot string
}

func LookupPlotsAndMaybeRegexThem(filteredFileRows []model.FileRow, flags model.ProgramFlags) []model.FileRow {

	//This is something I would have done in parallel using a buffered channel and possibly
	//a ratelimiter like https://github.com/uber-go/ratelimit

	var urls map[string]string = make(map[string]string)

	for _, fileRow := range filteredFileRows {
		urls[fileRow.Tconst] = "http://localhost:3000/static/tt0000075.json"
	}

	var plots map[string]string = make(map[string]string)

	results := BoundedParallelGet(urls, 10)

	// fmt.Printf("results from BoundedParallelGet:%+v", results)

	for _, result := range results {

		var p IMDBResponse

		err := json.NewDecoder(result.Res.Body).Decode(&p)
		if err != nil {
			log.Fatalln(err)

		}
		plots[result.Tconst] = p.Plot

	}

	// fmt.Printf("plots:%+v", plots)

	var rowsWithPlots []model.FileRow

	for _, fileRow := range filteredFileRows {

		fileRow.Plot = plots[fileRow.Tconst]

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
