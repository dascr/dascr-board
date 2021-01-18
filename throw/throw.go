package throw

// Round will be a throw round of 3 Throws
type Round struct {
	Round  int
	Done   bool
	Throws []Throw
}

// Throw will hold information about a single throw counted
type Throw struct {
	Number   int
	Modifier int
}
