package filter

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/andrewrobinson/imdb/model"
)

func TestRunFiltersAndPrint(t *testing.T) {

	//running against ../title.basics.truncated.tsv, assert on matches and total lines for various filters

	t.Run("no filters", func(t *testing.T) {
		emptyFlags := model.ProgramFlags{}
		genericTest(t, emptyFlags, 75, 75)
	})

	t.Run("--genres=Comedy", func(t *testing.T) {
		flags := model.ProgramFlags{GenresFlag: "Comedy"}
		genericTest(t, flags, 7, 75)
	})

	t.Run("--genres=Short", func(t *testing.T) {
		flags := model.ProgramFlags{GenresFlag: "Short"}
		genericTest(t, flags, 73, 75)
	})

}

func genericTest(t *testing.T, flags model.ProgramFlags, expectedMatches int, expectedHighestLineNumber int) {

	file, err := os.Open("../title.basics.truncated.tsv")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	matches, highestLineNumber := RunFiltersAndPrint(scanner, flags, false)

	if matches != expectedMatches || highestLineNumber != expectedHighestLineNumber {
		t.Errorf("got (%d, %d); wanted (%d, %d)", matches, highestLineNumber, expectedMatches, expectedHighestLineNumber)
	}

}
