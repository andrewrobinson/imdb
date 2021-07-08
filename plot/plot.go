package plot

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"time"

	"github.com/andrewrobinson/imdb/model"
)

func LookupPlots(filteredFileRows []model.FileRow, flags model.ProgramFlags) []model.FileRow {

	var rowsWithPlots []model.FileRow

	for _, fileRow := range filteredFileRows {

		plot, err := lookupPlot(fileRow.Tconst)

		if err != nil {
			fmt.Printf("error while looking up plots:%+v\n", err)
			os.Exit(1)
		}

		fileRow.Plot = plot

		if flags.PlotFilterFlag != "" {

			match, _ := regexp.MatchString(flags.PlotFilterFlag, fileRow.Plot)

			if match {
				rowsWithPlots = append(rowsWithPlots, fileRow)
			}

		} else {
			rowsWithPlots = append(rowsWithPlots, fileRow)
		}

	}

	return rowsWithPlots
}

func lookupPlot(tconst string) (string, error) {

	// https://raw.githubusercontent.com/andrewrobinson/imdb/main/tt0000075.json

	// "https://www.omdbapi.com/?i=tt0000075&apikey=591edae0"
	//TODO - make a localhost call
	return "As an elegant maestro of mirage and delusion drapes his beautiful female assistant with a gauzy textile, much to our amazement, the lady vanishes into thin air.", nil
}

func sleepForRandomTime() {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(10) // n will be between 0 and 10
	// fmt.Printf("Sleeping %d milliseconds...\n", 10+n)
	time.Sleep(time.Duration(10+n) * time.Millisecond)
	// fmt.Println("Done")
}
