package model

type ProgramFlags struct {
	FilePathFlag       string
	TitleTypeFlag      string
	PrimaryTitleFlag   string
	OriginalTitleFlag  string
	StartYearFlag      string
	EndYearFlag        string
	RuntimeMinutesFlag string
	GenresFlag         string
	PlotFilterFlag     string
	ProcessingTypeFlag string

	// TODO
	// maxApiRequests - maximum number of requests to be made to omdbapi
	// maxRunTime - maximum run time of the application. Format is a time.Duration string see here
	// maxRequests - maximum number of requests to send to omdbapi
	// plotFilter - regex pattern to apply to the plot of a film retrieved from omdbapi

}

type FileRow struct {
	Tconst         string
	TitleType      string
	PrimaryTitle   string
	OriginalTitle  string
	IsAdult        string
	StartYear      string
	EndYear        string
	RuntimeMinutes string
	Genres         string
	Plot           string
}
