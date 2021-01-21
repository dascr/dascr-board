package undo

import (
	"errors"

	"github.com/dascr/dascr-board/player"
	"github.com/dascr/dascr-board/throw"
)

// Sequence will be identified by a number and will hold actions to be undone
type Sequence struct {
	Sequence int
	Action   []Action
}

// Log will hold a sequences of actions to undo
type Log []*Sequence

// Action will hold the information to undo an action
type Action struct {
	Number              int
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
	PreviousNumberToHit int
}

/* Log functions */

// ClearLog will reset the log to empty
func (l *Log) ClearLog() {
	*l = Log{}
}

// CreateSequence will add an undo sequence to the undo log
func (l *Log) CreateSequence() *Sequence {
	number := 1

	if len(*l) > 0 {
		number = (*l)[len(*l)-1].Sequence + 1
	}

	newSequence := &Sequence{Sequence: number}
	*l = append(*l, newSequence)

	return (*l)[len(*l)-1]
}

// GetLastSequence will return the last sequence to be able to add actions to
func (l *Log) GetLastSequence() (*Sequence, error) {
	if len(*l) > 0 {
		return (*l)[len(*l)-1], nil
	}
	return &Sequence{}, errors.New("There are no sequences yet")
}

// GetSequenceByNumber will return a specific sequence
func (l *Log) GetSequenceByNumber(n int) (*Sequence, error) {
	for _, s := range *l {
		if s.Sequence == n {
			return s, nil
		}
	}

	return &Sequence{}, errors.New("Undo Sequence was not found")
}

// RemoveLastSequence will remove the last undo sequence from the undo log
func (l *Log) RemoveLastSequence() {
	*l = (*l)[:len(*l)-1]
}

/* Sequence functions */

// AddActionToSequence will add an Action to a sequence
func (s *Sequence) AddActionToSequence(a Action) error {
	number := 1
	if len(s.Action) > 0 {
		number = s.Action[len(s.Action)-1].Number + 1
	}
	a.Number = number
	s.Action = append(s.Action, a)
	return nil
}
