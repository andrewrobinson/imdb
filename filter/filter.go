package filter

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/andrewrobinson/imdb/common"
	"github.com/andrewrobinson/imdb/model"
)

func RunFilters(scanner *bufio.Scanner, flags model.ProgramFlags) ([]model.FileRow, int) {

	lineNumber := 0

	var filteredFileRows []model.FileRow

	for scanner.Scan() {

		lineNumber++
		line := scanner.Text()
		if lineNumber != 1 {

			fields := strings.Split(line, "\t")
			fileRow := common.BuildFileRow(fields)

			if rowMatchesFlags(fileRow, flags) {
				filteredFileRows = append(filteredFileRows, fileRow)
			}

		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error on line %v: %v", lineNumber, err)
	}

	return filteredFileRows, lineNumber

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
