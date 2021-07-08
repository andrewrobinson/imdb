package filter

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/andrewrobinson/imdb/model"
)

func TestRunFilters(t *testing.T) {

	//running against ../title.basics.truncated.tsv, assert on matches and total lines for various filters
	//I got these numbers using text editor find counts / by eyeballing the data

	t.Run("no filters", func(t *testing.T) {
		emptyFlags := model.ProgramFlags{}
		genericTest(t, emptyFlags, 74, 75)
	})

	t.Run("--genres=Comedy", func(t *testing.T) {
		flags := model.ProgramFlags{GenresFlag: "Comedy"}
		genericTest(t, flags, 7, 75)
	})

	t.Run("--genres=Short", func(t *testing.T) {
		flags := model.ProgramFlags{GenresFlag: "Short"}
		genericTest(t, flags, 73, 75)
	})

	t.Run("--genres=Comedy,Short", func(t *testing.T) {
		flags := model.ProgramFlags{GenresFlag: "Comedy,Short"}
		genericTest(t, flags, 4, 75)
	})

	t.Run("--genres=Animation,Comedy,Romance", func(t *testing.T) {
		flags := model.ProgramFlags{GenresFlag: "Animation,Comedy,Romance"}
		genericTest(t, flags, 1, 75)
	})

	t.Run("--genres=Comedy,Romance", func(t *testing.T) {
		flags := model.ProgramFlags{GenresFlag: "Comedy,Romance"}
		genericTest(t, flags, 1, 75)
	})

	t.Run("--genres=Documentary", func(t *testing.T) {
		flags := model.ProgramFlags{GenresFlag: "Documentary"}
		genericTest(t, flags, 37, 75)
	})

	t.Run("--originalTitle=Clown", func(t *testing.T) {
		flags := model.ProgramFlags{OriginalTitleFlag: "Clown"}
		genericTest(t, flags, 1, 75)
	})

	t.Run("--originalTitle=Clown --genres=Documentary", func(t *testing.T) {
		flags := model.ProgramFlags{OriginalTitleFlag: "Clown", GenresFlag: "Documentary"}
		genericTest(t, flags, 0, 75)
	})

	t.Run("--originalTitle=Clown --genres=Comedy", func(t *testing.T) {
		flags := model.ProgramFlags{OriginalTitleFlag: "Clown", GenresFlag: "Comedy"}
		genericTest(t, flags, 1, 75)
	})

	t.Run("--originalTitle=Clown --genres=medy", func(t *testing.T) {
		flags := model.ProgramFlags{OriginalTitleFlag: "Clown", GenresFlag: "medy"}
		genericTest(t, flags, 1, 75)
	})

	t.Run("--originalTitle=Clown --genres=Dramedy", func(t *testing.T) {
		flags := model.ProgramFlags{OriginalTitleFlag: "Clown", GenresFlag: "Dramedy"}
		genericTest(t, flags, 0, 75)
	})

	t.Run("--genres=Dramedy", func(t *testing.T) {
		flags := model.ProgramFlags{GenresFlag: "Dramedy"}
		genericTest(t, flags, 0, 75)
	})

	t.Run("--titleType=short --primaryTitle=Conjuring --originalTitle=Escamotage", func(t *testing.T) {
		flags := model.ProgramFlags{TitleTypeFlag: "short", PrimaryTitleFlag: "Conjuring", OriginalTitleFlag: "Escamotage"}
		genericTest(t, flags, 1, 75)
	})

	//TODO - plot regex matching has moved to another function so it needs to be tested there now

	//all rows currently have the plot for the below film hardcoded to save on limited api requests allowed:

	//"As an elegant maestro of mirage and delusion drapes his beautiful female assistant with a gauzy textile,
	// much to our amazement, the lady vanishes into thin air."

	//female should match
	// t.Run("--primaryTitle=Conjuring --originalTitle=Escamotage --plotFilter=female", func(t *testing.T) {
	// 	flags := model.ProgramFlags{PrimaryTitleFlag: "Conjuring", OriginalTitleFlag: "Escamotage", PlotFilterFlag: "female"}
	// 	genericTest(t, flags, 1, 75)
	// })

	// //females should not regex match
	// t.Run("--primaryTitle=Conjuring --originalTitle=Escamotage --plotFilter=females", func(t *testing.T) {
	// 	flags := model.ProgramFlags{PrimaryTitleFlag: "Conjuring", OriginalTitleFlag: "Escamotage", PlotFilterFlag: "females"}
	// 	genericTest(t, flags, 0, 75)
	// })

}

func genericTest(t *testing.T, flags model.ProgramFlags, expectedMatches int, expectedHighestLineNumber int) {

	file, err := os.Open("../title.basics.truncated.tsv")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	matchingFileRows, highestLineNumber := RunFilters(scanner, flags)

	if len(matchingFileRows) != expectedMatches || highestLineNumber != expectedHighestLineNumber {
		t.Errorf("got (%d, %d); wanted (%d, %d)", len(matchingFileRows), highestLineNumber, expectedMatches, expectedHighestLineNumber)
	}

}

func BenchmarkRunFiltersAndPrint(b *testing.B) {

	flags := model.ProgramFlags{TitleTypeFlag: "short", PrimaryTitleFlag: "Conjuring", OriginalTitleFlag: "Escamotage"}

	file, err := os.Open("../title.basics.truncated.tsv")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	for i := 0; i < b.N; i++ {
		scanner := bufio.NewScanner(file)
		RunFilters(scanner, flags)
	}
}
