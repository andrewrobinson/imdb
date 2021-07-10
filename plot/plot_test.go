package plot

import (
	"testing"
)

func TestLookupPlots(t *testing.T) {

	//TODO - plot regex logic must be tested here now

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
