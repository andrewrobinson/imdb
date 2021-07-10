package common

import (
	"flag"
	"fmt"

	"github.com/andrewrobinson/imdb/model"
)

<<<<<<< HEAD
func LookupPlot(tconst string) (string, error) {

	//TODO - make a localhost call or something
	// https://raw.githubusercontent.com/andrewrobinson/imdb/main/tt0000075.json

	// the real one but limited to 1000 a day
	// "https://www.omdbapi.com/?i=tt0000075&apikey=591edae0"

	//this waits 10-20ms, the actual call is about 53ms
	sleepForRandomTime()
	return "As an elegant maestro of mirage and delusion drapes his beautiful female assistant with a gauzy textile, much to our amazement, the lady vanishes into thin air.", nil
}

func sleepForRandomTime() {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(10)
	// fmt.Printf("Sleeping %d milliseconds...\n", 10+n)
	time.Sleep(time.Duration(10+n) * time.Millisecond)
	// fmt.Println("Done")
}

=======
>>>>>>> lines-in-memory
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

	fmt.Printf("%v	%v	%v\n", row.Tconst, row.PrimaryTitle, row.Plot)
	// fmt.Println(row)
}
