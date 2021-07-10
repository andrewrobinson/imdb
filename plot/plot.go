package plot

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	"github.com/andrewrobinson/imdb/model"
)

type IMDBResponse struct {
	Plot string
}

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

	fmt.Printf("LookupPlotsInParallel using ConcurrencyFactor:%+v and RateLimitPerSecond:%v\n", flags.ConcurrencyFactorFlag, flags.RateLimitPerSecondFlag)

	var parallelGetResults []MontanaResult = MontanaBoundedParallelGet(urls, flags.ConcurrencyFactorFlag)

	var plots map[string]string = buildMapOfTconstToPlot(parallelGetResults)

	return plots

}

func buildMapOfTconstToUrl(filteredFileRows []model.FileRow) map[string]string {

	var urls map[string]string = make(map[string]string)

	for _, fileRow := range filteredFileRows {

		urls[fileRow.Tconst] = "https://raw.githubusercontent.com/andrewrobinson/imdb/207ba5bd2727dfadb65a3faccd6786a099dce5ef/static/tt0000075.json"
		// urls[fileRow.Tconst] = "http://localhost:3000/static/tt0000075.json"
	}

	return urls

}

func buildMapOfTconstToPlot(parallelGetResults []MontanaResult) map[string]string {

	var plots = make(map[string]string)

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
