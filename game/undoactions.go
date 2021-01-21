package game

import (
	"github.com/dascr/dascr-board/throw"
	"github.com/dascr/dascr-board/undo"
	"github.com/dascr/dascr-board/utils"
)

// UndoCreateThrowRound is the undo action for CREATETHROWROUND
func UndoCreateThrowRound(action undo.Action) {
	for i, rnd := range action.Player.ThrowRounds {
		if rnd.Round == action.RoundNumber {
			copy(action.Player.ThrowRounds[i:], action.Player.ThrowRounds[i+1:])
			action.Player.ThrowRounds[len(action.Player.ThrowRounds)-1] = throw.Round{}
			action.Player.ThrowRounds = action.Player.ThrowRounds[:len(action.Player.ThrowRounds)-1]
		}
	}
}

// UndoCreateThrow is the undo action for CREATETHROW
func UndoCreateThrow(action undo.Action) {
	if len(action.Player.ThrowRounds) > 0 {
		throwRound := &action.Player.ThrowRounds[action.RoundNumber-1]
		if len(throwRound.Throws) > 0 {
			throwRound.Throws = throwRound.Throws[:len(throwRound.Throws)-1]
		}
		action.Player.TotalThrowCount--
	}
}

// UndoScore is the undo action for UPDATESCORE
func UndoScore(action undo.Action) {
	action.Player.Score.Score = action.PreviousScore
	action.Player.Average = action.PreviousAverage
	action.Player.ThrowSum = action.PreviousThrowSum
	action.Player.LastThrows = action.PreviousLastThree
}

// UndoScoreAndPark is the undo action for UPDATESCOREANDPARK
func UndoScoreAndPark(action undo.Action) {
	action.Player.Score.Score = action.PreviousScore
	action.Player.Score.ParkScore = action.PreviousParkScore
	action.Player.Average = action.PreviousAverage
	action.Player.ThrowSum = action.PreviousThrowSum
	action.Player.LastThrows = action.PreviousLastThree
}

// UndoBustAndWin is the undo action for X01BUST and DOWIN
func UndoBustAndWin(action undo.Action, base *BaseGame) {
	base.GameState = action.PreviousGameState
	base.Message = action.PreviousMessage
	action.Player.Score.Score = action.PreviousScore
	action.Player.Score.ParkScore = action.PreviousParkScore
	action.Player.ThrowRounds[len(action.Player.ThrowRounds)-1].Done = false
}

// UndoDoPodium is the undo action for DOPODIUM
func UndoDoPodium(action undo.Action, base *BaseGame) {
	base.GameState = action.PreviousGameState
	base.Message = action.PreviousMessage
	if base.Podium.GetPodiumLength() > 0 {
		base.Podium.RemoveLastPlacement()
	}
}

// UndoClosePlayerThrowRound is the undo action for CLOSEPLAYERTHROWROUND
func UndoClosePlayerThrowRound(action undo.Action, base *BaseGame) {
	base.GameState = action.PreviousGameState
	base.Message = action.PreviousMessage
	for _, rnd := range action.Player.ThrowRounds {
		if rnd.Round == action.RoundNumber {
			rnd.Done = false
		}
	}
}

// UndoIncreaseThrowRound is the undo action for INCREASETHROWROUND
func UndoIncreaseThrowRound(action undo.Action, base *BaseGame) {
	base.ThrowRound--
}

// UndoNextPlayer is the undo action for NEXTPLAYER
func UndoNextPlayer(action undo.Action, base *BaseGame) {
	base.ActivePlayer = action.PreviousPlayerIndex
}

// UndoCloseControllerNumber is the undo action for CLOSECONTROLLERNUMBER
func UndoCloseControllerNumber(action undo.Action, base *BaseGame) {
	base.CricketController.NumberClosed[action.NumberIndex] = false
}

// UndoClosePlayerNumber is the undo action for CLOSEPLAYERNUMBER
func UndoClosePlayerNumber(action undo.Action) {
	action.Player.Score.Closed[action.NumberIndex] = false
}

// UndoRevealNumber is the undo action for REVEALNUMBER
func UndoRevealNumber(action undo.Action, base *BaseGame) {
	// Hide it again
	base.CricketController.NumberRevealed[action.NumberIndex] = false
	// Choose new random number which is not one of the 7 existing
	for {
		var newNumber = utils.GetSingleRandomCricketNumber()
		var contained = false
		for _, n := range base.CricketController.Numbers {
			if newNumber == n {
				contained = true
			}
		}
		if !contained {
			// substitute the revealed number with the new one
			base.CricketController.Numbers[action.NumberIndex] = newNumber
			break
		}
	}
}

// UndoIncreaseHitCount is the undo action for INCREASEHITCOUNT
func UndoIncreaseHitCount(action undo.Action) {
	action.Player.Score.Numbers[action.NumberIndex] -= action.Modifier
}

// UndoGainPoints is the undo action for GAINPOINTS
func UndoGainPoints(action undo.Action) {
	action.Player.Score.Score -= action.Points
}

// UndoWin is the undo action for DOWIN
func UndoWin(action undo.Action, base *BaseGame) {
	base.GameState = action.PreviousGameState
	base.Message = action.PreviousMessage
}

// UndoATCIncreaseNumber is the undo action for ATCINCREASENUMBER
func UndoATCIncreaseNumber(action undo.Action) {
	action.Player.Score.CurrentNumber = action.PreviousNumberToHit
}

// UndoUpdateSplitScore is the undo action for UPDATESPLITSCORE
func UndoUpdateSplitScore(action undo.Action) {
	action.Player.Score.Score = action.PreviousScore
}
