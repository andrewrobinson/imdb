package plot

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/andrewrobinson/imdb/model"
)

type IMDBResponse struct {
	Plot string
}

// a struct to hold the result from each request
type MontanaResult struct {
	Tconst string
	Res    http.Response
	Err    error
}

// plot.MaybeRegexFilter(filteredFileRows, mapOfTconstToPlot, flags)

func AddPlotsAndMaybeRegexFilter(filteredFileRows []model.FileRow, mapOfTconstToPlot map[string]string, flags model.ProgramFlags) []model.FileRow {

	var rowsWithPlots []model.FileRow

	for _, fileRow := range filteredFileRows {

		fileRow.Plot = mapOfTconstToPlot[fileRow.Tconst]

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

func LookupPlotsInParallel(filteredFileRows []model.FileRow, flags model.ProgramFlags) map[string]string {

	//This is something I would have done in parallel using a buffered channel and possibly
	//a ratelimiter like https://github.com/uber-go/ratelimit

	var urls map[string]string = buildMapOfTconstToUrl(filteredFileRows)

	fmt.Printf("LookupPlotsInParallel using flags.ConcurrencyFactor:%+v\n", flags.ConcurrencyFactor)

	//TODO - pull 10 from flags
	var parallelGetResults []MontanaResult = MontanaBoundedParallelGet(urls, flags.ConcurrencyFactor)

	var plots map[string]string = buildMapOfTconstToPlot(parallelGetResults)

	return plots

}

func buildMapOfTconstToUrl(filteredFileRows []model.FileRow) map[string]string {

	var urls map[string]string = make(map[string]string)

	for _, fileRow := range filteredFileRows {
		urls[fileRow.Tconst] = "http://localhost:3000/static/tt0000075.json"
	}

	return urls

}

func buildMapOfTconstToPlot(parallelGetResults []MontanaResult) map[string]string {

	var plots = make(map[string]string)

	// fmt.Printf("results from BoundedParallelGet:%+v", results)

	for _, result := range parallelGetResults {

		var p IMDBResponse

		err := json.NewDecoder(result.Res.Body).Decode(&p)
		if err != nil {
			log.Fatalln(err)

		}
		plots[result.Tconst] = p.Plot

	}

	return plots

}
