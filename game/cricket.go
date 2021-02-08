package game

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/dascr/dascr-board/player"
	"github.com/dascr/dascr-board/podium"
	"github.com/dascr/dascr-board/score"
	"github.com/dascr/dascr-board/throw"
	"github.com/dascr/dascr-board/undo"
	"github.com/dascr/dascr-board/utils"
	"github.com/dascr/dascr-board/ws"
)

// CricketGameController will hold the status of different numbers
// This will hold information if the number is closed in game and can be hit or not
// Also it maps the 7 numbers to actual dart values
// Classic is 15-20 + 25
// Otherwise random or ghost can be used
type CricketGameController struct {
	NumberClosed   []bool
	Numbers        []int
	NumberRevealed []bool // For Ghost Cricket
	Random         bool
	Ghost          bool
}

// CricketGame will hold the Cricket game information
type CricketGame struct {
	Base BaseGame
}

// StartGame will satisfy interface Game for game Cricket
func (g *CricketGame) StartGame() error {
	// Different cases of game
	var mode string

	if g.Base.CricketController.Ghost {
		mode = "ghost"
	} else if g.Base.CricketController.Random {
		mode = "random"
	} else {
		mode = "default"
	}

	// fill controller accordingly
	switch mode {
	case "ghost":
		// Random numbers + ghost
		g.Base.CricketController.initGhost()
	case "random":
		// Random numbers
		g.Base.CricketController.initRandom()
	case "default":
		// Default numbers
		g.Base.CricketController.initDefault()
	}

	// set players scores
	for i := range g.Base.Player {
		score := score.BaseScore{
			Score:   0,
			Numbers: []int{0, 0, 0, 0, 0, 0, 0},
			Closed:  []bool{false, false, false, false, false, false, false},
		}
		g.Base.Player[i].Score = score
		g.Base.Player[i].ThrowRounds = make([]throw.Round, 0)
		g.Base.Player[i].LastThrows = make([]throw.Throw, 3)
	}

	g.Base.Podium = &podium.Podium{}

	g.Base.UndoLog = &undo.Log{}
	sequence := g.Base.UndoLog.CreateSequence()
	sequence.AddActionToSequence(undo.Action{
		Action: "CREATEGAME",
	})

	return nil
}

// GetStatus will satisfy interface Game for game Cricket
func (g *CricketGame) GetStatus() BaseGame {
	return g.Base
}

// GetStatusDisplay will satisfy interface Game for game Cricket
func (g *CricketGame) GetStatusDisplay() BaseGame {
	return stripDisplay(&g.Base)
}

// NextPlayer will satisfy interface Game for game Cricket
func (g *CricketGame) NextPlayer(h *ws.Hub) {
	switchToNextPlayer(&g.Base, h)
}

// RequestThrow will satisfy interface Game for game Cricket
func (g *CricketGame) RequestThrow(number, modifier int, h *ws.Hub) error {
	sequence := g.Base.UndoLog.CreateSequence()
	previousMessage := g.Base.Message
	previousState := g.Base.GameState

	activePlayer := &g.Base.Player[g.Base.ActivePlayer]
	allPlayer := g.Base.Player

	// Check game state
	if g.Base.GameState == "THROW" {
		// check if ongoing round else create
		checkOngoingElseCreate(activePlayer, &g.Base, sequence)

		// Add Throw to last round
		throwRound := addThrowToCurrentRound(activePlayer, &g.Base, sequence, number, modifier)

		// Filter if throw is relevant at all
		// and assign index if applicable
		var relevant = false
		var index = -1
		for i, n := range g.Base.CricketController.Numbers {
			if number == n {
				relevant = true
				index = i
			}
		}

		if relevant {
			// Throw is relevant, continue
			// In Ghost reveal if not revealed yet
			if g.Base.CricketController.Ghost {
				if !g.Base.CricketController.NumberRevealed[index] {
					g.Base.CricketController.NumberRevealed[index] = true
					sequence.AddActionToSequence(undo.Action{
						Action:      "REVEALNUMBER",
						NumberIndex: index,
						GameID:      g.Base.UID,
					})
				}
			}

			// Open Close handling
			// Check if already closed
			if !g.Base.CricketController.NumberClosed[index] {
				// Continue with scoring and such
				relevantMod := increaseHitCount(activePlayer.Score, index, modifier)
				// calculate relevant mod

				points := number * relevantMod
				sequence.AddActionToSequence(undo.Action{
					Action:      "INCREASEHITCOUNT",
					NumberIndex: index,
					Modifier:    modifier,
					Player:      activePlayer,
				})

				// Check if number has to be closed for player
				if !checkClosed(activePlayer.Score, index) && checkToClose(activePlayer.Score, index) {
					activePlayer.Score.Closed[index] = true
					sequence.AddActionToSequence(undo.Action{
						Action:      "CLOSEPLAYERNUMBER",
						NumberIndex: index,
						Player:      activePlayer,
					})
				}

				// Dispatch to different scoring modes
				switch g.Base.Variant {
				case "cut":
					if checkClosed(activePlayer.Score, index) {
						scoreCut(g, points, index, activePlayer.UID, allPlayer, sequence)
					}
				case "normal":
					scoreNormal(g, points, index, *activePlayer, sequence)
				case "no":
				}

				// Now check again if the number has to be closed
				var done = true
				for _, p := range allPlayer {
					if !checkClosed(p.Score, index) {
						done = false
					}
				}
				if done {
					g.Base.CricketController.NumberClosed[index] = true
					sequence.AddActionToSequence(undo.Action{
						Action:      "CLOSECONTROLLERNUMBER",
						NumberIndex: index,
						GameID:      g.Base.UID,
					})
				}
			}
		}

		// Check if 3 throws in round and close round
		// Also set gameState and perhaps increase g.Base.ThrowRound
		// if everyone has already thrown to this round
		if len(throwRound.Throws) == 3 {
			closePlayerRound(&g.Base, activePlayer, throwRound, sequence, []undo.Action{}, previousState, previousMessage)
		}

		// Check if game shot and handle win
		if checkWin(g) {
			cricketWin(g, activePlayer, sequence, throwRound)
		}

		// Set assets for Frontend
		setFrontendAssets(&g.Base.Player[g.Base.ActivePlayer], &g.Base)

		// Update scoreboard
		utils.WSSendUpdate(g.Base.UID, h)
		return nil
	}

	return fmt.Errorf("game state is '%+v', so no throw accepted", g.Base.GameState)
}

// Undo will satisfy interface Game for game Cricket
func (g *CricketGame) Undo(h *ws.Hub) error {
	if err := triggerUndo(&g.Base, h); err != nil {
		return err
	}

	return nil
}

// Rematch will satisfy interface Game for game Cricket
func (g *CricketGame) Rematch(h *ws.Hub) error {
	// Init random number generator
	s := rand.NewSource(time.Now().Unix())
	rg := rand.New(s)

	// Different cases of game
	var mode string

	if g.Base.CricketController.Ghost {
		mode = "ghost"
	} else if g.Base.CricketController.Random {
		mode = "random"
	} else {
		mode = "default"
	}

	// fill controller accordingly
	switch mode {
	case "ghost":
		// Random numbers + ghost
		g.Base.CricketController.initGhost()
	case "random":
		// Random numbers
		g.Base.CricketController.initRandom()
	case "default":
		// Default numbers
		g.Base.CricketController.initDefault()
	}

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
			Score:   0,
			Numbers: []int{0, 0, 0, 0, 0, 0, 0},
			Closed:  []bool{false, false, false, false, false, false, false},
		}
		g.Base.Player[i].Score = score
		g.Base.Player[i].ThrowRounds = make([]throw.Round, 0)
		g.Base.Player[i].LastThrows = make([]throw.Throw, 3)
		g.Base.Player[i].TotalThrowCount = 0
	}

	sequence := g.Base.UndoLog.CreateSequence()
	sequence.AddActionToSequence(undo.Action{
		Action: "CREATEGAME",
	})

	// Update scoreboard
	utils.WSSendUpdate(g.Base.UID, h)

	return nil
}

func (controller *CricketGameController) initGhost() *CricketGameController {
	controller.Numbers = utils.GetRandomCricketNumbers(false)
	controller.NumberClosed = []bool{false, false, false, false, false, false, false}
	controller.NumberRevealed = []bool{false, false, false, false, false, false, false}
	return controller
}

func (controller *CricketGameController) initDefault() *CricketGameController {
	controller.Numbers = []int{15, 16, 17, 18, 19, 20, 25}
	controller.NumberClosed = []bool{false, false, false, false, false, false, false}
	controller.NumberRevealed = []bool{true, true, true, true, true, true, true}
	return controller
}

func (controller *CricketGameController) initRandom() *CricketGameController {
	controller.Numbers = utils.GetRandomCricketNumbers(true)
	controller.NumberClosed = []bool{false, false, false, false, false, false, false}
	controller.NumberRevealed = []bool{true, true, true, true, true, true, true}
	return controller
}

func scoreCut(g *CricketGame, points, index int, activeID string, allPlayer []player.Player, sequence *undo.Sequence) {
	for i := range allPlayer {
		if g.Base.Player[i].UID != activeID {
			if !checkClosed(g.Base.Player[i].Score, index) {
				g.Base.Player[i].Score.Score += points
				sequence.AddActionToSequence(undo.Action{
					Action: "GAINPOINTS",
					Player: &g.Base.Player[i],
					Points: points,
				})
			}
		}
	}

}

func scoreNormal(g *CricketGame, points, index int, player player.Player, sequence *undo.Sequence) {
	if checkClosed(player.Score, index) {
		player.Score.Score += points
		sequence.AddActionToSequence(undo.Action{
			Action: "GAINPOINTS",
			Player: &player,
			Points: points,
		})
	}
}

func checkWin(g *CricketGame) bool {
	var activePlayer = g.Base.Player[g.Base.ActivePlayer]
	var activePlayerScore = activePlayer.Score.Score

	// Switch over cut or normal (default is noscore)
	switch g.Base.Variant {
	case "cut":
		// It does not matter if the active player has everything closed
		// as he could make someone other win
		// Player wins if all closed and least score
		for _, p := range g.Base.Player {
			allClosed := checkAllClosed(p.Score)
			leastScore := true

			for _, pl := range g.Base.Player {
				if p.UID != pl.UID {
					if p.Score.Score > pl.Score.Score {
						leastScore = false
					}
				}
			}

			if allClosed && leastScore {
				return true
			}

		}
		return false
	case "normal":
		// If not all numbers of active player are closed there is no win possible
		if checkAllClosed(activePlayer.Score) {
			// If player has biggest score he wins
			for _, p := range g.Base.Player {
				if activePlayer.UID != p.UID {
					if activePlayerScore < p.Score.Score {
						return true
					}
				}
			}
		}

	default:
		// This applies to noscore
		if checkAllClosed(activePlayer.Score) {
			return true
		}
		break
	}
	return false
}

func cricketWin(g *CricketGame, activePlayer *player.Player, sequence *undo.Sequence, throwRound *throw.Round) {
	previousMessage := g.Base.Message
	previousState := g.Base.GameState

	if g.Base.Settings.Podium {
		// Do podium and continue game
		doPodium(&g.Base, activePlayer, sequence)
		return
	}

	g.Base.GameState = "WON"
	g.Base.Message = "Game shot!"

	throwRound.Done = true
	sequence.AddActionToSequence(undo.Action{
		Action:            "DOWIN",
		GameID:            g.Base.UID,
		PreviousGameState: previousState,
		PreviousMessage:   previousMessage,
		Player:            activePlayer,
	})
}

// checkClosed will satisfy Score interface for Cricket
func checkClosed(score score.BaseScore, index int) bool {
	return score.Closed[index]
}

// checkToClose will satisfy Score interdace for Cricket
func checkToClose(score score.BaseScore, index int) bool {
	return score.Numbers[index] >= 3
}

// checkAllClosed will satisfy Score interface for Cricket
func checkAllClosed(score score.BaseScore) bool {
	for _, n := range score.Closed {
		if !n {
			return false
		}
	}
	return true
}

// increaseHitCount will increase hit count and return
// relevant modifier for applying points
func increaseHitCount(score score.BaseScore, index, modifier int) int {
	// Calculate relevant modifier against count
	alreadyHit := 3 - score.Numbers[index]

	// No negative
	if alreadyHit < 0 {
		alreadyHit = 0
	}

	relevantMod := modifier - alreadyHit

	// Increse hitcount
	score.Numbers[index] += modifier

	return relevantMod
}
