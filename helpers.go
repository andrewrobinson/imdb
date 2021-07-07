package main

import (
	"flag"
	"fmt"
)

func BuildProgramFlags() ProgramFlags {

	filePathFlag := flag.String("filePath", "title.basics.truncated.tsv", "")
	titleTypeFlag := flag.String("titleType", "", "")
	primaryTitleFlag := flag.String("primaryTitle", "", "")
	originalTitleFlag := flag.String("originalTitle", "", "")
	startYearFlag := flag.String("startYear", "", "")
	endYearFlag := flag.String("endYear", "", "")
	runtimeMinutesFlag := flag.String("runtimeMinutes", "", "")
	genresFlag := flag.String("genres", "", "")

	flag.Parse()

	flagStruct := ProgramFlags{*filePathFlag, *titleTypeFlag, *primaryTitleFlag, *originalTitleFlag, *startYearFlag, *endYearFlag, *runtimeMinutesFlag, *genresFlag}
	return flagStruct

}

func BuildFileRow(fields []string) FileRow {

	rowStruct := FileRow{
		tconst:         fields[0],
		titleType:      fields[1],
		primaryTitle:   fields[2],
		originalTitle:  fields[3],
		isAdult:        fields[4],
		startYear:      fields[5],
		endYear:        fields[6],
		runtimeMinutes: fields[7],
		genres:         fields[8],
	}
	return rowStruct

}

func PrintFields(row FileRow) {
	//For now just print out the fields, but later output must be
	// IMDB_ID     |   Title               |   Plot
	// tt0000005   |   Blacksmith Scene    |   Three men hammer on an anvil and pass a bottle of beer around.

	// fmt.Printf("%+v\n", row)
	fmt.Println(row)
}
