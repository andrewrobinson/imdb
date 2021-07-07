package filter

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/andrewrobinson/imdb/model"
)

func BenchmarkRunFiltersAndPrint(b *testing.B) {

	flags := model.ProgramFlags{TitleTypeFlag: "short", PrimaryTitleFlag: "Conjuring", OriginalTitleFlag: "Escamotage"}

	file, err := os.Open("../title.basics.truncated.tsv")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	for i := 0; i < b.N; i++ {
		scanner := bufio.NewScanner(file)
		RunFiltersAndPrint(scanner, flags, false)
	}
}

func TestRunFiltersAndPrint(t *testing.T) {

	//running against ../title.basics.truncated.tsv, assert on matches and total lines for various filters
	//I got these numbers using text editor find counts / by eyeballing the data

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
