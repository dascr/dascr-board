package game

import (
	"errors"
	"math"

	"github.com/dascr/dascr-board/logger"
	"github.com/dascr/dascr-board/player"
	"github.com/dascr/dascr-board/throw"
	"github.com/dascr/dascr-board/undo"
	"github.com/dascr/dascr-board/utils"
	"github.com/dascr/dascr-board/ws"
)

// This will be used by setFrontendAssets
func calculateAverageAndTotalThrowCount(player *player.Player) *player.Player {
	// Set average
	allTotalThrowCount := 0
	totalSumCount := 0
	divider := float64(0)
	for _, rnd := range player.ThrowRounds {
		for _, thr := range rnd.Throws {
			totalSumCount += thr.Number * thr.Modifier
			divider += float64(1 / 3.0)
			allTotalThrowCount++
		}
	}
	avg := float64(totalSumCount) / divider
	avg = math.Round(avg*100) / 100
	if math.IsNaN(avg) {
		avg = 0
	}
	player.Average = avg
	player.TotalThrowCount = allTotalThrowCount

	return player
}

// This will be used by setFrontendAssets
func getLastThreeThrows(player *player.Player, base *BaseGame) *player.Player {
	// Set last 3 throws
	// and sum of it
	player.LastThrows = make([]throw.Throw, 0)
	lastThreeSum := 0
	if len(player.ThrowRounds) > 0 {
		if len(player.ThrowRounds[base.ThrowRound-1].Throws) != 0 {
			for _, thr := range player.ThrowRounds[base.ThrowRound-1].Throws {
				player.LastThrows = append(player.LastThrows, thr)
				lastThreeSum += thr.Number * thr.Modifier
			}
			player.ThrowSum = lastThreeSum
		}
	} else {
		player.LastThrows = make([]throw.Throw, 3)
	}

	return player
}

// This will take care of the stats to display at frontend
func setFrontendAssets(player *player.Player, base *BaseGame) {
	player = calculateAverageAndTotalThrowCount(player)
	getLastThreeThrows(player, base)
}

// This will set game state for podium case
func doPodium(base *BaseGame, activePlayer *player.Player, sequence int) {
	previousMessage := base.Message
	previousState := base.GameState
	playerLeft := len(base.Player) - len(base.Podium)

	// More than 2 player left to go to the podium
	if playerLeft > 2 {
		// Push active players uid to podium
		base.Podium = append(base.Podium, activePlayer.UID)
		base.GameState = "NEXTPLAYERWON"

		if len(base.Podium) == 1 {
			base.Message = "Winner!"
			base.UndoLog = append(base.UndoLog, undo.Undo{Sequence: sequence, Action: "DOPODIUM", Player: activePlayer, PreviousGameState: previousState, PreviousMessage: previousMessage, GameID: base.UID})
			return
		}

		base.Message = "Next Winner!"
		base.UndoLog = append(base.UndoLog, undo.Undo{Sequence: sequence, Action: "DOPODIUM", Player: activePlayer, PreviousGameState: previousState, PreviousMessage: previousMessage, GameID: base.UID})
		return
	}

	// Otherwise exactly 2 left
	base.Podium = append(base.Podium, activePlayer.UID)
	base.UndoLog = append(base.UndoLog, undo.Undo{Sequence: sequence, Action: "DOPODIUM", Player: activePlayer, PreviousGameState: previousState, PreviousMessage: previousMessage, GameID: base.UID})
	for i := range base.Player {
		var contained = false
		// When player UID is not in base.Podium
		for j := range base.Podium {
			if base.Player[i].UID == base.Podium[j] {
				contained = true
			}

		}

		if !contained {
			base.Podium = append(base.Podium, base.Player[i].UID)
			base.UndoLog = append(base.UndoLog, undo.Undo{Sequence: sequence, Action: "DOPODIUM", Player: &base.Player[i], PreviousGameState: previousState, PreviousMessage: previousMessage, GameID: base.UID})

			doWin(base)
			throwRound := &activePlayer.ThrowRounds[base.ThrowRound-1]
			throwRound.Done = true
			activePlayer.Score.Score = 0
			activePlayer.Score.ParkScore = 0
			previousScore := activePlayer.Score.Score
			previousParkScore := activePlayer.Score.ParkScore
			base.UndoLog = append(base.UndoLog, undo.Undo{Sequence: sequence, Action: "DOWIN", GameID: base.UID, PreviousGameState: previousState, PreviousMessage: previousMessage, Player: activePlayer, PreviousScore: previousScore, PreviousParkScore: previousParkScore})

			return
		}
	}
}

// This will set game state for win case
func doWin(base *BaseGame) {
	// Set game state
	base.GameState = "WON"
	base.Message = "Game shot!"
}

// This will strip response for Display function for FrontEnd
func stripDisplay(base *BaseGame) BaseGame {
	var returnBase BaseGame
	var returnPlayer = make([]player.Player, 0)
	returnBase = *base

	for _, p := range returnBase.Player {
		var tr = make([]throw.Round, 0)
		p.ThrowRounds = tr
		returnPlayer = append(returnPlayer, p)
	}

	returnBase.UndoLog = nil
	returnBase.Player = returnPlayer

	return returnBase
}

// This will switch to next player
func switchToNextPlayer(base *BaseGame, h *ws.Hub) {
	sequence := base.UndoLog[len(base.UndoLog)-1].Sequence + 1
	previousPlayerIndex := base.ActivePlayer
	previousMessage := base.Message
	previousState := base.GameState

	// Switch to next player
	if base.GameState != "WON" {
		// fill empty throws of active player with 0
		// before switching to next player
		activePlayer := &base.Player[base.ActivePlayer]
		// Check if ongoing round
		ongoing := utils.CheckOngoingRound(activePlayer.ThrowRounds, base.ThrowRound)
		if !ongoing {
			// If there is no round associated with the current throw round of the game
			// create one
			newRound := &throw.Round{
				Round:  base.ThrowRound,
				Throws: []throw.Throw{},
			}
			activePlayer.ThrowRounds = append(activePlayer.ThrowRounds, *newRound)
			base.UndoLog = append(base.UndoLog, undo.Undo{Sequence: sequence, Action: "CREATETHROWROUND", RoundNumber: newRound.Round, Player: activePlayer})
		}

		currentThrowRound := &activePlayer.ThrowRounds[base.ThrowRound-1]
		count := 3 - len(currentThrowRound.Throws)
		if count != 0 {
			for i := 0; i < count; i++ {
				newThrow := &throw.Throw{
					Number:   0,
					Modifier: 1,
				}
				currentThrowRound.Throws = append(currentThrowRound.Throws, *newThrow)
				base.UndoLog = append(base.UndoLog, undo.Undo{Sequence: sequence, Action: "CREATETHROW", Player: activePlayer, RoundNumber: currentThrowRound.Round})

			}
			currentThrowRound.Done = true
			base.UndoLog = append(base.UndoLog, undo.Undo{Sequence: sequence, Action: "CLOSEPLAYERTHROWROUND", Player: activePlayer, RoundNumber: currentThrowRound.Round, GameID: base.UID, PreviousGameState: previousState, PreviousMessage: previousMessage})
		}

		// Switch player
		base.ActivePlayer = utils.ChooseNextPlayer(base.Player, base.ActivePlayer, base.Podium)
		// Reset gamestate
		base.GameState = "THROW"
		base.Message = "-"

		base.UndoLog = append(base.UndoLog, undo.Undo{Sequence: sequence, Action: "NEXTPLAYER", PreviousPlayerIndex: previousPlayerIndex, GameID: base.UID})

		// Set assets for Frontend
		setFrontendAssets(activePlayer, base)

		// Update scoreboard
		utils.WSSendUpdate(base.UID, h)

		// Check if to increase game.ThrowRound
		if roundDone := utils.CheckRoundDone(base.Player, base.ThrowRound, base.Podium); roundDone {
			base.ThrowRound++
			base.UndoLog = append(base.UndoLog, undo.Undo{Sequence: sequence, Action: "INCREASETHROWROUND", GameID: base.UID})
		}
	}
}

// This will trigger undo
func triggerUndo(base *BaseGame, h *ws.Hub) error {
	lastSequence := base.UndoLog[len(base.UndoLog)-1].Sequence
	if lastSequence == 0 {
		return errors.New("there is nothing to undo")
	}
	var sequenceActions []undo.Undo
	var parkAction undo.Undo
	for _, s := range base.UndoLog {
		if s.Sequence == lastSequence {
			if s.Action == "CREATETHROWROUND" {
				parkAction = s
				continue
			}

			sequenceActions = append(sequenceActions, s)
		}

		if parkAction.Action == "CREATETHROWROUND" {
			sequenceActions = append(sequenceActions, parkAction)
		}
	}

	logger.Debugf("Sequence is: %+v", sequenceActions)

	for _, a := range sequenceActions {
		switch a.Action {
		// Common
		case "CREATEGAME":
			// Do nothing
			break
		case "CREATETHROWROUND":
			UndoCreateThrowRound(a)
		case "CREATETHROW":
			UndoCreateThrow(a)
		case "DOWIN":
			UndoWin(a, base)
		case "DOPODIUM":
			UndoDoPodium(a, base)
		case "NEXTPLAYER":
			UndoNextPlayer(a, base)
		case "CLOSEPLAYERTHROWROUND":
			UndoClosePlayerThrowRound(a, base)
		case "INCREASETHROWROUND":
			UndoIncreaseThrowRound(a, base)
		// X01
		case "UPDATESCORE":
			UndoScore(a)
		case "UPDATESCOREANDPARK":
			UndoScoreAndPark(a)
		case "X01BUST":
			UndoBustAndWin(a, base)
		// Cricket
		case "REVEALNUMBER":
			UndoRevealNumber(a, base)
		case "INCREASEHITCOUNT":
			UndoIncreaseHitCount(a)
		case "CLOSEPLAYERNUMBER":
			UndoClosePlayerNumber(a)
		case "CLOSECONTROLLERNUMBER":
			UndoCloseControllerNumber(a, base)
		case "GAINPOINTS":
			UndoGainPoints(a)
		// ATC
		case "ATCINCREASENUMBER":
			UndoATCIncreaseNumber(a)
		// Split
		case "UPDATESPLITSCORE":
			UndoUpdateSplitScore(a)
		default:
			break
		}
	}

	var newUndoLog []undo.Undo
	// Remove complete sequence from undoLog
	for i, entry := range base.UndoLog {
		if entry.Sequence != lastSequence {
			newUndoLog = append(newUndoLog, base.UndoLog[i])
		}
	}
	base.UndoLog = newUndoLog

	// reset assets
	setFrontendAssets(&base.Player[base.ActivePlayer], base)

	// Update trigger via hub
	message := ws.Message{
		Room: base.UID,
		Data: []byte("update"),
	}

	h.Broadcast <- message

	return nil
}
