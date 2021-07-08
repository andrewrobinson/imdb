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
		// fmt.Printf("lowmem lineNumber:'%v'\n", lineNumber)
		line := scanner.Text()
		// fmt.Printf("lowmem line:'%v'\n", line)
		if lineNumber != 1 {
			matches = handleLine(line, flags, matches, printRows)
		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error on line %v: %v", lineNumber, err)
	}

	// fmt.Printf("\nRunFilters returning matches:%v, lineNumber:%v\n", matches, lineNumber)
	return matches, lineNumber

}

func handleLine(line string, flags model.ProgramFlags, matches int, printRows bool) int {
	fields := strings.Split(line, "\t")
	fileRow := common.BuildFileRow(fields)

	if rowMatchesFlags(fileRow, flags) {

		//this shifts out, and the line as well as the incremented matches must be returned
		plot, err := common.LookupPlot(fileRow.Tconst)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fileRow.Plot = plot

		if flags.PlotFilterFlag != "" {

			match, _ := regexp.MatchString(flags.PlotFilterFlag, fileRow.Plot)

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

	// fmt.Printf("\nhandleLine returning matches:%v\n", matches)
	return matches
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
