package model

type ProgramFlags struct {
	FilePathFlag       string //the location of the imdb file to read
	TitleTypeFlag      string
	PrimaryTitleFlag   string
	OriginalTitleFlag  string
	StartYearFlag      string
	EndYearFlag        string
	RuntimeMinutesFlag string
	GenresFlag         string
	MaxApiRequestsFlag int    //the max number of search results the program will attempt to look plots up for
	MaxRunTimeFlag     int    //if the program hasn't finished by this time, print what it has and exit
	MaxRequestsFlag    int    //unused, not sure of the meaning of this one
	PlotFilterFlag     string //a regex filter run on the plot after it has been looked up
	//these 2 added by me
	ConcurrencyFactorFlag  int //parallel requests for plots are limited by this concurrency factor
	RateLimitPerSecondFlag int //parallel requests for plots use this as a rate limit per second
}

//matches the tsv
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
