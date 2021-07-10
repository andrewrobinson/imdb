package plot

import (
	"fmt"
	"regexp"

	"github.com/andrewrobinson/imdb/model"
)

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

func LookupPlotsInParallel(filteredFileRows []model.FileRow, flags model.ProgramFlags, plotLookuptemplate string) map[string]string {
	urls := buildMapOfTconstToUrl(filteredFileRows, plotLookuptemplate)

	// fmt.Printf("urls:%+v", urls)
	// os.Exit(1)

	fmt.Printf("LookupPlotsInParallel using ConcurrencyFactor:%+v and RateLimitPerSecond:%v\n", flags.ConcurrencyFactorFlag, flags.RateLimitPerSecondFlag)
	parallelGetResults := BoundedParallelGet(urls, flags.ConcurrencyFactorFlag, flags.RateLimitPerSecondFlag)
	plots := buildMapOfTconstToPlot(parallelGetResults)
	return plots
}

func buildMapOfTconstToUrl(filteredFileRows []model.FileRow, plotLookuptemplate string) map[string]string {

	var urls map[string]string = make(map[string]string)

	for _, fileRow := range filteredFileRows {
		urls[fileRow.Tconst] = fmt.Sprintf(plotLookuptemplate, fileRow.Tconst)
	}

	return urls

}

func buildMapOfTconstToPlot(parallelGetResults []PlotLookupResult) map[string]string {

	var plots = make(map[string]string)

	for _, result := range parallelGetResults {
		plots[result.Tconst] = result.Plot
	}

	return plots

}
