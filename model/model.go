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
	MaxApiRequestsFlag int
	MaxRunTimeFlag     int
	MaxRequestsFlag    int
	PlotFilterFlag     string
	ConcurrencyFactor  int
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
