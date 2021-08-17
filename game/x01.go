package game

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/dascr/dascr-board/player"
	"github.com/dascr/dascr-board/podium"
	"github.com/dascr/dascr-board/score"
	"github.com/dascr/dascr-board/throw"
	"github.com/dascr/dascr-board/undo"
	"github.com/dascr/dascr-board/utils"
	"github.com/dascr/dascr-board/ws"
)

// X01Game will hold the X01 game information
type X01Game struct {
	Base BaseGame
}

// StartGame will satisfy interface Game for game X01
func (g *X01Game) StartGame() error {
	// Score to int
	setupscore, err := strconv.Atoi(g.Base.Variant)
	if err != nil {
		return err
	}

	// CreateScore for each player
	// and init empty throw splice
	for i := range g.Base.Player {
		score := score.BaseScore{
			Score:        setupscore,
			ParkScore:    setupscore,
			InitialScore: setupscore,
		}
		g.Base.Player[i].Score = score
		g.Base.Player[i].ThrowRounds = make([]throw.Round, 0)
		g.Base.Player[i].LastThrows = make([]throw.Throw, 3)
		g.Base.Player[i].ThrowSum = 0
		g.Base.Player[i].Average = 0
	}

	g.Base.Podium = &podium.Podium{}

	g.Base.UndoLog = &undo.Log{}
	sequence := g.Base.UndoLog.CreateSequence()
	sequence.AddActionToSequence(undo.Action{
		Action: "CREATEGAME",
	})

	return nil
}

// GetStatus will satisfy interface Game for game X01
func (g *X01Game) GetStatus() BaseGame {
	return g.Base
}

// GetStatusDisplay will satisfy interface Game for game X01
func (g *X01Game) GetStatusDisplay() BaseGame {
	return stripDisplay(&g.Base)
}

// NextPlayer will satisfy interface Game for game X01
func (g *X01Game) NextPlayer(h *ws.Hub) {
	switchToNextPlayer(&g.Base, h)
}

// RequestThrow will satisfy interface Game for game X01
func (g *X01Game) RequestThrow(number, modifier int, h *ws.Hub) error {
	sequence := g.Base.UndoLog.CreateSequence()

	points := number * modifier
	activePlayer := &g.Base.Player[g.Base.ActivePlayer]
	// Check game state
	if g.Base.GameState == "THROW" {
		// check if ongoing round else create
		checkOngoingElseCreate(activePlayer, &g.Base, sequence)

		// Add Throw to last round
		throwRound := addThrowToCurrentRound(activePlayer, &g.Base, sequence, number, modifier)

		// New score will be
		newScore := activePlayer.Score.Score - points

		// Add 100 if punisher enabled
		if g.Base.Punisher {
			if number == 0 || modifier == 0 {
				newScore = activePlayer.Score.Score + 100
			}
		}

		// Handle cases win, bust, normal throw
		switch {
		// BUST
		case newScore < 0:
			x01bust("BUST", g, throwRound, activePlayer, sequence, activePlayer.Score.Score)
		// WIN
		case newScore == 0:
			x01win(g, modifier, throwRound, activePlayer, sequence)
		// NORMAL THROW
		default:
			normalThrow(g, newScore, modifier, throwRound, activePlayer, sequence)
		}

		// Set assets for Frontend
		setFrontendAssets(activePlayer, &g.Base)

		// Update scoreboard
		utils.WSSendUpdate(g.Base.UID, h)

		return nil
	}

	return fmt.Errorf("game state is '%+v', so no throw accepted", g.Base.GameState)
}

// Undo will satisfy interface Game for game X01
func (g *X01Game) Undo(h *ws.Hub) error {
	if err := triggerUndo(&g.Base, h); err != nil {
		return err
	}

	return nil
}

// Rematch will satisfy interface Game for game X01
func (g *X01Game) Rematch(h *ws.Hub) error {
	// Score to int
	setupscore, err := strconv.Atoi(g.Base.Variant)
	if err != nil {
		return err
	}

	// Init random number generator
	s := rand.NewSource(time.Now().Unix())
	rg := rand.New(s)

	// Reset game state
	g.Base.Message = ""
	g.Base.GameState = "THROW"
	g.Base.Podium.ResetPodium()
	g.Base.UndoLog.ClearLog()
	g.Base.ActivePlayer = rg.Intn(len(g.Base.Player))
	g.Base.ThrowRound = 1

	// CreateScore for each player
	// and init empty throw splice
	for i := range g.Base.Player {
		score := score.BaseScore{
			Score:        setupscore,
			ParkScore:    setupscore,
			InitialScore: setupscore,
		}
		g.Base.Player[i].Score = score
		g.Base.Player[i].ThrowRounds = make([]throw.Round, 0)
		g.Base.Player[i].LastThrows = make([]throw.Throw, 3)
		g.Base.Player[i].TotalThrowCount = 0
		g.Base.Player[i].ThrowSum = 0
		g.Base.Player[i].Average = 0
	}

	sequence := g.Base.UndoLog.CreateSequence()
	sequence.AddActionToSequence(undo.Action{
		Action: "CREATEGAME",
	})

	// Update scoreboard
	utils.WSSendUpdate(g.Base.UID, h)

	return nil
}

// This will handle x01bust and reset game state
func x01bust(nextState string, game *X01Game, throwRound *throw.Round, activePlayer *player.Player, sequence *undo.Sequence, previousScore int) {
	oldState := game.Base.GameState
	previousMessage := game.Base.Message
	previousParkScore := activePlayer.Score.ParkScore
	game.Base.GameState = nextState
	switch nextState {
	case "BUST":
		game.Base.Message = "Bust. Remove Darts!"
	case "BUSTNOCHECKOUT":
		game.Base.Message = "Bust. No Checkout left. Remove Darts!"
	case "BUSTCONDITION":
		game.Base.Message = fmt.Sprintf("Bust. Failed Condition: %s in / %s out. Remove Darts!", game.Base.In, game.Base.Out)
	default:
		game.Base.Message = "-"
	}
	throwRound.Done = true
	activePlayer.Score.Score = activePlayer.Score.ParkScore
	sequence.AddActionToSequence(undo.Action{
		Action:            "X01BUST",
		GameID:            game.Base.UID,
		Player:            activePlayer,
		PreviousGameState: oldState,
		PreviousScore:     previousScore,
		PreviousParkScore: previousParkScore,
		PreviousMessage:   previousMessage,
	})
}

// This will handle x01win
func x01win(game *X01Game, modifier int, throwRound *throw.Round, activePlayer *player.Player, sequence *undo.Sequence) {
	previousMessage := game.Base.Message
	previousState := game.Base.GameState
	previousScore := activePlayer.Score.Score
	previousParkScore := activePlayer.Score.ParkScore
	// Check if double or master out are met
	if !checkoutMet(game, modifier) {
		x01bust("BUSTCONDITION", game, throwRound, activePlayer, sequence, activePlayer.Score.Score)
		return
	}
	if game.Base.Settings.Podium {
		// Do podium and continue game
		doPodium(&game.Base, activePlayer, sequence)
		return
	}
	// Win and end game
	doWin(&game.Base)
	throwRound.Done = true
	activePlayer.Score.Score = 0
	activePlayer.Score.ParkScore = 0

	sequence.AddActionToSequence(undo.Action{
		Action:            "DOWIN",
		GameID:            game.Base.UID,
		PreviousGameState: previousState,
		PreviousMessage:   previousMessage,
		Player:            activePlayer,
		PreviousScore:     previousScore,
		PreviousParkScore: previousParkScore,
	})
}

// This will handle the normal throw routine
func normalThrow(game *X01Game, newScore, modifier int, throwRound *throw.Round, activePlayer *player.Player, sequence *undo.Sequence) {
	// For undo function
	previousScore := activePlayer.Score.Score
	previousParkScore := activePlayer.Score.ParkScore
	previousAverage := activePlayer.Average
	previousThrowSum := activePlayer.ThrowSum
	previousLastThree := activePlayer.LastThrows
	previousMessage := game.Base.Message
	previousState := game.Base.GameState

	// Check if checkout left on double and master out
	if !checkoutPossible(game, newScore) {
		x01bust("BUSTNOCHECKOUT", game, throwRound, activePlayer, sequence, activePlayer.Score.Score)
		return
	}
	// Check if first throw
	if activePlayer.Score.Score == activePlayer.Score.InitialScore {
		// Check if in condition met
		if !checkinPossible(game, modifier) {
			// if last throw and condition not met - bust
			if len(throwRound.Throws) == 3 {
				x01bust("BUSTCONDITION", game, throwRound, activePlayer, sequence, activePlayer.Score.Score)
			}
			return
		}
	}
	// Update player score
	activePlayer.Score.Score = newScore
	sequence.AddActionToSequence(undo.Action{
		Action:            "UPDATESCORE",
		GameID:            game.Base.UID,
		PreviousGameState: previousState,
		PreviousMessage:   previousMessage,
		Player:            activePlayer,
		PreviousScore:     previousScore,
		PreviousParkScore: previousParkScore,
		PreviousAverage:   previousAverage,
		PreviousThrowSum:  previousThrowSum,
		PreviousLastThree: previousLastThree,
	})

	// Check if 3 throws in round and close round
	// Also set gameState and perhaps increase game.Base.ThrowRound
	// if everyone has already thrown to this round
	if len(throwRound.Throws) == 3 {
		activePlayer.Score.Score = newScore
		activePlayer.Score.ParkScore = newScore
		closePlayerRound(&game.Base, activePlayer, throwRound, sequence, []undo.Action{
			undo.Action{
				Action:            "UPDATESCOREANDPARK",
				Player:            activePlayer,
				PreviousScore:     previousScore,
				PreviousParkScore: previousParkScore,
				PreviousAverage:   previousAverage,
				PreviousThrowSum:  previousThrowSum,
				PreviousLastThree: previousLastThree,
			},
		}, previousState, previousMessage)
	}

}

// This will check if a checkout is even possible anymore and return a bool
func checkoutPossible(game *X01Game, newScore int) bool {
	if game.Base.Out != "straight" {
		return newScore != 1
	}
	return true
}

// This will check if the checkout condition is met and return bool
func checkoutMet(game *X01Game, modifier int) bool {
	if game.Base.Out == "double" {
		return modifier == 2
	} else if game.Base.Out == "master" {
		return modifier != 1
	}
	return true
}

// This will check if a checkin is possible and return a bool
func checkinPossible(game *X01Game, modifier int) bool {
	if game.Base.In == "double" {
		return modifier == 2
	} else if game.Base.In == "master" {
		return modifier != 1
	}
	return true
}
