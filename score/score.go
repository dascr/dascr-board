package score

// BaseScore will keep track of the players score
// for all games
type BaseScore struct {
	Score         int    // universal
	ParkScore     int    // X01
	InitialScore  int    // X01
	Numbers       []int  // Cricket
	Closed        []bool // Cricket
	CurrentNumber int    // ATC
	Split         bool   // Split
	Hit           bool   // Split
}
