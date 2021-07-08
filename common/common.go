package common

import (
	"flag"
	"fmt"

	"github.com/andrewrobinson/imdb/model"
)

func BuildProgramFlags() model.ProgramFlags {

	filePathFlag := flag.String("filePath", "title.basics.truncated.tsv", "")
	titleTypeFlag := flag.String("titleType", "", "")
	primaryTitleFlag := flag.String("primaryTitle", "", "")
	originalTitleFlag := flag.String("originalTitle", "", "")
	startYearFlag := flag.String("startYear", "", "")
	endYearFlag := flag.String("endYear", "", "")
	runtimeMinutesFlag := flag.String("runtimeMinutes", "", "")
	genresFlag := flag.String("genres", "", "")
	maxApiRequestsFlag := flag.Int("maxApiRequests", 300, "")
	maxRunTimeFlag := flag.Int("maxRunTime", 30, "")
	maxRequestsFlag := flag.Int("maxRequests", 300, "")
	plotFilterFlag := flag.String("plotFilter", "", "")

	flag.Parse()

	flagStruct := model.ProgramFlags{FilePathFlag: *filePathFlag,
		TitleTypeFlag:      *titleTypeFlag,
		PrimaryTitleFlag:   *primaryTitleFlag,
		OriginalTitleFlag:  *originalTitleFlag,
		StartYearFlag:      *startYearFlag,
		EndYearFlag:        *endYearFlag,
		RuntimeMinutesFlag: *runtimeMinutesFlag,
		GenresFlag:         *genresFlag,
		MaxApiRequestsFlag: *maxApiRequestsFlag,
		MaxRunTimeFlag:     *maxRunTimeFlag,
		MaxRequestsFlag:    *maxRequestsFlag,
		PlotFilterFlag:     *plotFilterFlag,
	}
	return flagStruct

}

func BuildFileRow(fields []string) model.FileRow {

	return model.FileRow{
		Tconst:         fields[0],
		TitleType:      fields[1],
		PrimaryTitle:   fields[2],
		OriginalTitle:  fields[3],
		IsAdult:        fields[4],
		StartYear:      fields[5],
		EndYear:        fields[6],
		RuntimeMinutes: fields[7],
		Genres:         fields[8],
	}

}

func PrintFields(row model.FileRow) {

	// TODO - examine buffered output
	// https://stackoverflow.com/questions/64638136/performance-issues-while-reading-a-file-line-by-line-with-bufio-newscanner

	//For now just print out the fields, but later output must be
	// IMDB_ID     |   Title               |   Plot
	// tt0000005   |   Blacksmith Scene    |   Three men hammer on an anvil and pass a bottle of beer around.

	fmt.Printf("%+v\n", row)
	// fmt.Println(row)
}
