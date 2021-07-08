package common

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/andrewrobinson/imdb/model"
)

func LookupPlot(tconst string) (string, error) {

	// 	./curltime "https://www.omdbapi.com/?i=tt0000075&apikey=591edae0"
	//     time_namelookup:  0.002087s
	//        time_connect:  0.008398s
	//     time_appconnect:  0.036539is
	//    time_pretransfer:  0.036783s
	//       time_redirect:  0.000000s
	//  time_starttransfer:  0.052927s
	//                     ----------
	//          time_total:  0.053134s

	//So 53ms per plot lookup

	//if you have an overall 30 sec timeout to return to the user
	//30 / (53/1000) = 566 rows max

	//sleep for between 10 and 20 milliseconds
	//sleepForRandomTime()

	return "As an elegant maestro of mirage and delusion drapes his beautiful female assistant with a gauzy textile, much to our amazement, the lady vanishes into thin air.", nil

}

func sleepForRandomTime() {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(10) // n will be between 0 and 10
	// fmt.Printf("Sleeping %d milliseconds...\n", 10+n)
	time.Sleep(time.Duration(10+n) * time.Millisecond)
	// fmt.Println("Done")
}

func BuildProgramFlags() model.ProgramFlags {

	filePathFlag := flag.String("filePath", "title.basics.truncated.tsv", "")
	titleTypeFlag := flag.String("titleType", "", "")
	primaryTitleFlag := flag.String("primaryTitle", "", "")
	originalTitleFlag := flag.String("originalTitle", "", "")
	startYearFlag := flag.String("startYear", "", "")
	endYearFlag := flag.String("endYear", "", "")
	runtimeMinutesFlag := flag.String("runtimeMinutes", "", "")
	genresFlag := flag.String("genres", "", "")
	plotFilterFlag := flag.String("plotFilter", "", "")
	processingTypeFlag := flag.String("processingType", "lowmem", "lowmem or highmem")

	flag.Parse()

	flagStruct := model.ProgramFlags{FilePathFlag: *filePathFlag,
		TitleTypeFlag:      *titleTypeFlag,
		PrimaryTitleFlag:   *primaryTitleFlag,
		OriginalTitleFlag:  *originalTitleFlag,
		StartYearFlag:      *startYearFlag,
		EndYearFlag:        *endYearFlag,
		RuntimeMinutesFlag: *runtimeMinutesFlag,
		GenresFlag:         *genresFlag,
		PlotFilterFlag:     *plotFilterFlag,
		ProcessingTypeFlag: *processingTypeFlag}
	return flagStruct

}

func BuildFileRow(fields []string) model.FileRow {

	rowStruct := model.FileRow{
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
	return rowStruct

}

func PrintFields(row model.FileRow) {

	// TODO - examine buffered output
	// https://stackoverflow.com/questions/64638136/performance-issues-while-reading-a-file-line-by-line-with-bufio-newscanner

	//For now just print out the fields, but later output must be
	// IMDB_ID     |   Title               |   Plot
	// tt0000005   |   Blacksmith Scene    |   Three men hammer on an anvil and pass a bottle of beer around.

	// fmt.Printf("%+v\n", row)
	fmt.Println(row)
}
