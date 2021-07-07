package filter

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/andrewrobinson/imdb/common"
	"github.com/andrewrobinson/imdb/model"
)

func RunFilters(scanner *bufio.Scanner, flags model.ProgramFlags, printRows bool) (int, int) {

	lineNumber := 0
	matches := 0
	for scanner.Scan() {

		lineNumber++
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		fileRow := common.BuildFileRow(fields)

		if rowMatchesFlags(fileRow, flags) {

			plot, err := common.LookupPlot(fileRow.Tconst)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fileRow.Plot = plot

			if flags.PlotFilterFlag != "" {

				match, _ := regexp.MatchString(flags.PlotFilterFlag, fileRow.Plot)
				// fmt.Printf("\nmatch:%v\n", match)

				if match {
					matches++
					if printRows {
						common.PrintFields(fileRow)
					}
				}

			} else {
				matches++
				if printRows {
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
