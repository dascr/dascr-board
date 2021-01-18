package undo

import (
	"github.com/dascr/dascr-board/player"
	"github.com/dascr/dascr-board/throw"
)

// Undo will hold the sequence and action of what to undo
type Undo struct {
	Sequence            int
	Action              string
	Points              int
	RoundNumber         int
	Player              *player.Player
	GameID              string
	NumberIndex         int
	Modifier            int
	PreviousGameState   string
	PreviousScore       int
	PreviousParkScore   int
	PreviousPlayerIndex int
	PreviousThrowSum    int
	PreviousAverage     float64
	PreviousLastThree   []throw.Throw
	PreviousMessage     string
}
