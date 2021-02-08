package game

import (
	"math"

	"github.com/dascr/dascr-board/player"
	"github.com/dascr/dascr-board/throw"
	"github.com/dascr/dascr-board/undo"
	"github.com/dascr/dascr-board/utils"
	"github.com/dascr/dascr-board/ws"
)

// This will check if an ongoing throw round
// is assosiated with the active player
// if not it will create one
func checkOngoingElseCreate(activePlayer *player.Player, base *BaseGame, sequence *undo.Sequence) {
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
		sequence.AddActionToSequence(undo.Action{
			Action:      "CREATETHROWROUND",
			RoundNumber: newRound.Round,
			Player:      activePlayer,
		})
	}
}

// This will add a provided throw to the last round
// assosiated with the active player
// It returns the throw round for further processing withing game logic
func addThrowToCurrentRound(activePlayer *player.Player, base *BaseGame, sequence *undo.Sequence, number, modifier int) *throw.Round {
	// Now add throw to that round or to the existing
	throwRound := &activePlayer.ThrowRounds[base.ThrowRound-1]
	newThrow := &throw.Throw{
		Number:   number,
		Modifier: modifier,
	}
	throwRound.Throws = append(throwRound.Throws, *newThrow)
	sequence.AddActionToSequence(undo.Action{
		Action:      "CREATETHROW",
		Player:      activePlayer,
		RoundNumber: throwRound.Round,
	})

	return throwRound
}

// This will check if player threw 3 darst
// and then close the round if so
func closePlayerRound(base *BaseGame, activePlayer *player.Player, throwRound *throw.Round, sequence *undo.Sequence, actions []undo.Action, previousState, previousMessage string) {
	throwRound.Done = true
	base.GameState = "NEXTPLAYER"
	base.Message = "Remove Darts!"

	for _, a := range actions {
		sequence.AddActionToSequence(a)
	}

	sequence.AddActionToSequence(undo.Action{
		Action:            "CLOSEPLAYERTHROWROUND",
		Player:            activePlayer,
		RoundNumber:       throwRound.Round,
		GameID:            base.UID,
		PreviousGameState: previousState,
		PreviousMessage:   previousMessage,
	})
}

// This will be used by setFrontendAssets
// It calculates player average and total throw count
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
// It sets the players field LastThrows
// To the last 3 throws of current throw round
// Ommiting empty throws
func getLastThreeThrows(player *player.Player, base *BaseGame) *player.Player {
	// Set last 3 throws
	// and sum of it
	player.LastThrows = make([]throw.Throw, 0)
	lastThreeSum := 0
	if len(player.ThrowRounds) > 0 {
		if len(player.ThrowRounds[len(player.ThrowRounds)-1].Throws) != 0 {
			for _, thr := range player.ThrowRounds[len(player.ThrowRounds)-1].Throws {
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
// Placement of player on podium
func doPodium(base *BaseGame, activePlayer *player.Player, sequence *undo.Sequence) {
	previousMessage := base.Message
	previousState := base.GameState
	playerLeft := len(base.Player) - base.Podium.GetPodiumLength()

	// More than 2 player left to go to the podium
	if playerLeft > 2 {
		// Push active players uid to podium
		base.Podium.AddPlayerToPodium(activePlayer)
		base.GameState = "NEXTPLAYERWON"
		base.Message = "Next Winner!"

		if base.Podium.GetPodiumLength() == 1 {
			base.Message = "Winner!"
		}

		sequence.AddActionToSequence(undo.Action{
			Action:            "DOPODIUM",
			Player:            activePlayer,
			PreviousGameState: previousState,
			PreviousMessage:   previousMessage,
			GameID:            base.UID,
		})
		return
	}

	// Otherwise exactly 2 left
	base.Podium.AddPlayerToPodium(activePlayer)
	sequence.AddActionToSequence(undo.Action{
		Action:            "DOPODIUM",
		Player:            activePlayer,
		PreviousGameState: previousState,
		PreviousMessage:   previousMessage,
		GameID:            base.UID,
	})
	for i := range base.Player {
		var contained = false
		// When player UID is not in base.Podium
		podium := *base.Podium.GetPodium()
		for j := range podium {
			if base.Player[i].UID == podium[j].Player.UID {
				contained = true
			}

		}

		if !contained {
			base.Podium.AddPlayerToPodium(&base.Player[i])
			sequence.AddActionToSequence(undo.Action{
				Action:            "DOPODIUM",
				Player:            activePlayer,
				PreviousGameState: previousState,
				PreviousMessage:   previousMessage,
				GameID:            base.UID,
			})

			doWin(base)
			throwRound := &activePlayer.ThrowRounds[base.ThrowRound-1]
			throwRound.Done = true
			activePlayer.Score.Score = 0
			activePlayer.Score.ParkScore = 0

			previousScore := activePlayer.Score.Score
			previousParkScore := activePlayer.Score.ParkScore
			sequence.AddActionToSequence(undo.Action{
				Action:            "DOWIN",
				GameID:            base.UID,
				PreviousGameState: previousState,
				PreviousMessage:   previousMessage,
				Player:            activePlayer,
				PreviousScore:     previousScore,
				PreviousParkScore: previousParkScore,
			})

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
// To minimize traffic size
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

// This will fill the players current throw round with 0s
// for the case someone pressed nextPlayer before
// entering all 3 throws
func fillWithZeros(base *BaseGame, sequence *undo.Sequence, previousState, PreviousMessage string) {
	// fill empty throws of active player with 0
	// before switching to next player
	activePlayer := &base.Player[base.ActivePlayer]
	// Check if ongoing round
	checkOngoingElseCreate(activePlayer, base, sequence)

	currentThrowRound := &activePlayer.ThrowRounds[base.ThrowRound-1]
	count := 3 - len(currentThrowRound.Throws)
	if count != 0 {
		for i := 0; i < count; i++ {
			newThrow := &throw.Throw{
				Number:   0,
				Modifier: 1,
			}
			currentThrowRound.Throws = append(currentThrowRound.Throws, *newThrow)
			sequence.AddActionToSequence(undo.Action{
				Action:      "CREATETHROW",
				Player:      activePlayer,
				RoundNumber: currentThrowRound.Round,
			})

		}
		currentThrowRound.Done = true

		// When X01 we need to set new parkScore
		if base.Game == "x01" {
			activePlayer.Score.ParkScore = activePlayer.Score.Score
			previousParkScore := activePlayer.Score.ParkScore
			previousScore := activePlayer.Score.Score
			previousAverage := activePlayer.Average
			previousThrowSum := activePlayer.ThrowSum
			previousLastThree := activePlayer.LastThrows

			sequence.AddActionToSequence(undo.Action{
				Action:            "UPDATESCOREANDPARK",
				Player:            activePlayer,
				PreviousScore:     previousScore,
				PreviousParkScore: previousParkScore,
				PreviousAverage:   previousAverage,
				PreviousThrowSum:  previousThrowSum,
				PreviousLastThree: previousLastThree,
			})
		}

		sequence.AddActionToSequence(undo.Action{
			Action:            "CLOSEPLAYERTHROWROUND",
			Player:            activePlayer,
			RoundNumber:       currentThrowRound.Round,
			GameID:            base.UID,
			PreviousGameState: previousState,
			PreviousMessage:   PreviousMessage,
		})
	}

}

// This will switch to next player
func switchToNextPlayer(base *BaseGame, h *ws.Hub) {
	sequence := base.UndoLog.CreateSequence()
	previousPlayerIndex := base.ActivePlayer
	previousState := base.GameState
	previousMessage := base.Message

	// Switch to next player
	if base.GameState != "WON" {
		activePlayer := &base.Player[base.ActivePlayer]
		fillWithZeros(base, sequence, previousState, previousMessage)

		// Switch player
		base.ActivePlayer = utils.ChooseNextPlayer(base.Player, base.ActivePlayer, base.Podium.GetPodium())
		// Reset gamestate
		base.GameState = "THROW"
		base.Message = "-"

		sequence.AddActionToSequence(undo.Action{
			Action:              "NEXTPLAYER",
			PreviousPlayerIndex: previousPlayerIndex,
			GameID:              base.UID,
		})

		// Set assets for Frontend
		setFrontendAssets(activePlayer, base)

		// Update scoreboard
		utils.WSSendUpdate(base.UID, h)

		// Check if throw round done by all players and increase
		checkRoundDone(base, sequence)
	}
}

// This will check if the round is done and increase overall
// Throw round
func checkRoundDone(base *BaseGame, sequence *undo.Sequence) {
	// Check if to increase game.ThrowRound
	if roundDone := utils.CheckRoundDone(base.Player, base.ThrowRound, base.Podium.GetPodium()); roundDone {
		base.ThrowRound++
		sequence.AddActionToSequence(undo.Action{
			Action: "INCREASETHROWROUND",
			GameID: base.UID,
		})
	}
}

// This will trigger undo
func triggerUndo(base *BaseGame, h *ws.Hub) error {
	// Check length and do not remove "CREATEGAME"
	if len(*base.UndoLog) > 1 {
		sequence, err := base.UndoLog.GetLastSequence()
		if err != nil {
			return err
		}

		for _, a := range sequence.Action {
			if a.Action == "CREATETHROWROUND" {
				// look if there is CREATETHROWROUND in sequence
				// if so resort and put it last as otherwise
				// undo will be rendered unusable

				// park CREATETHROWROUND
				parkAction := sequence.Action[0]

				// remove CREATETHROWROUND from sequence
				sequence.Action = sequence.Action[1:]

				// add CREATETHROWROUND at the end
				sequence.Action = append(sequence.Action, parkAction)

				break
			}
		}

		for _, a := range sequence.Action {
			switch a.Action {
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

		// Remove sequence from undo log
		base.UndoLog.RemoveLastSequence()

		// reset assets
		setFrontendAssets(&base.Player[base.ActivePlayer], base)

		// Update trigger via hub
		message := ws.Message{
			Room: base.UID,
			Data: []byte("update"),
		}

		h.Broadcast <- message
	}

	return nil
}
