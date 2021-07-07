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

	flag.Parse()

	flagStruct := model.ProgramFlags{FilePathFlag: *filePathFlag,
		TitleTypeFlag:      *titleTypeFlag,
		PrimaryTitleFlag:   *primaryTitleFlag,
		OriginalTitleFlag:  *originalTitleFlag,
		StartYearFlag:      *startYearFlag,
		EndYearFlag:        *endYearFlag,
		RuntimeMinutesFlag: *runtimeMinutesFlag,
		GenresFlag:         *genresFlag}
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
	//For now just print out the fields, but later output must be
	// IMDB_ID     |   Title               |   Plot
	// tt0000005   |   Blacksmith Scene    |   Three men hammer on an anvil and pass a bottle of beer around.

	// fmt.Printf("%+v\n", row)
	fmt.Println(row)
}

// func sleepForRandomTime() {
// 	rand.Seed(time.Now().UnixNano())
// 	n := rand.Intn(10) // n will be between 0 and 10
// 	//fmt.Printf("Sleeping %d seconds...\n", n)
// 	time.Sleep(time.Duration(n) * time.Second)
// 	//fmt.Println("Done")
// }
