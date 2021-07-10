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
	//the max number of results the program will attempt to look plots up for
	maxApiRequestsFlag := flag.Int("maxApiRequests", 5000, "")
	maxRunTimeFlag := flag.Int("maxRunTime", 30, "")
	maxRequestsFlag := flag.Int("maxRequests", 300, "")
	plotFilterFlag := flag.String("plotFilter", "", "")
	concurrencyFactorFlag := flag.Int("concurrencyFactor", 20, "")
	rateLimitPerSecondFlag := flag.Int("rateLimitPerSecond", 100, "")

	flag.Parse()

	flagStruct := model.ProgramFlags{FilePathFlag: *filePathFlag,
		TitleTypeFlag:          *titleTypeFlag,
		PrimaryTitleFlag:       *primaryTitleFlag,
		OriginalTitleFlag:      *originalTitleFlag,
		StartYearFlag:          *startYearFlag,
		EndYearFlag:            *endYearFlag,
		RuntimeMinutesFlag:     *runtimeMinutesFlag,
		GenresFlag:             *genresFlag,
		MaxApiRequestsFlag:     *maxApiRequestsFlag,
		MaxRunTimeFlag:         *maxRunTimeFlag,
		MaxRequestsFlag:        *maxRequestsFlag,
		PlotFilterFlag:         *plotFilterFlag,
		ConcurrencyFactorFlag:  *concurrencyFactorFlag,
		RateLimitPerSecondFlag: *rateLimitPerSecondFlag,
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

func PrintRow(row model.FileRow) {

	// TODO - examine buffered output
	// https://stackoverflow.com/questions/64638136/performance-issues-while-reading-a-file-line-by-line-with-bufio-newscanner

	//For now just print out the fields, but later output must be
	// IMDB_ID     |   Title               |   Plot
	// tt0000005   |   Blacksmith Scene    |   Three men hammer on an anvil and pass a bottle of beer around.

	fmt.Printf("%v	|	%v	|	%v\n", row.Tconst, row.PrimaryTitle, row.Plot)
	// fmt.Println(row)
}
