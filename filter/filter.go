package filter

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"github.com/andrewrobinson/imdb/common"
	"github.com/andrewrobinson/imdb/model"
)

func RunFilters(scanner *bufio.Scanner, flags model.ProgramFlags, printOutput bool) (int, int) {

	lineNumber := 0
	matches := 0
	for scanner.Scan() {

		lineNumber++
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		fileRow := common.BuildFileRow(fields)

		if rowMatchesFlags(fileRow, flags) {

			//this is the most immediate place to do the plot lookup
			//it could be done as a separate step, but then this process would need to return data
			//as opposed to just printing it out while it has it
			//need to balance memory usage/performance/fault tolerance etc

			plot := common.LookupPlot(fileRow.Tconst)
			fileRow.Plot = plot

			if flags.PlotFilterFlag != "" {

				match, _ := regexp.MatchString(flags.PlotFilterFlag, fileRow.Plot)
				fmt.Printf("\nmatch:%v\n", match)

				if match {
					matches++
					if printOutput {
						common.PrintFields(fileRow)
					}
				}

			} else {
				matches++
				if printOutput {
					common.PrintFields(fileRow)
				}
			}

		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error on line %v: %v", lineNumber, err)
	}

	return matches, lineNumber

}

func rowMatchesFlags(row model.FileRow, flags model.ProgramFlags) bool {

	titleTypeMatches := flagMatchesOrIsEmpty(flags.TitleTypeFlag, row.TitleType)
	primaryTitleMatches := flagMatchesOrIsEmpty(flags.PrimaryTitleFlag, row.PrimaryTitle)
	originalTitleMatches := flagMatchesOrIsEmpty(flags.OriginalTitleFlag, row.OriginalTitle)
	startYearMatches := flagMatchesOrIsEmpty(flags.StartYearFlag, row.StartYear)
	endYearMatches := flagMatchesOrIsEmpty(flags.EndYearFlag, row.EndYear)
	runtimeMinutesMatches := flagMatchesOrIsEmpty(flags.RuntimeMinutesFlag, row.RuntimeMinutes)
	genresMatches := flagMatchesOrIsEmpty(flags.GenresFlag, row.Genres)

	if titleTypeMatches && primaryTitleMatches && originalTitleMatches && startYearMatches && endYearMatches && runtimeMinutesMatches && genresMatches {
		return true
	} else {
		return false
	}

}

func flagMatchesOrIsEmpty(filterValue string, columnValue string) bool {

	//if no flag value passed then don't filter, ie it passes
	if filterValue == "" {
		return true
	} else {
		return strings.Contains(columnValue, filterValue)
	}

}
