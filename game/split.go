package game

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/dascr/dascr-board/player"
	"github.com/dascr/dascr-board/score"
	"github.com/dascr/dascr-board/throw"
	"github.com/dascr/dascr-board/undo"
	"github.com/dascr/dascr-board/utils"
	"github.com/dascr/dascr-board/ws"
)

// SplitGame will hold the Split game information
type SplitGame struct {
	Base BaseGame
}

// StartGame will satisfy interface Game for game Split
func (g *SplitGame) StartGame() error {
	// Init start score vor edart
	startscore := 40
	// set score to 0 if steel
	if g.Base.Variant == "steel" {
		startscore = 0
	}

	// CreateScore for each player
	// and init empty throw splice
	for i := range g.Base.Player {
		score := score.BaseScore{
			Score: startscore,
			Split: true,
		}
		g.Base.Player[i].Score = score
		g.Base.Player[i].ThrowRounds = make([]throw.Round, 0)
		g.Base.Player[i].LastThrows = make([]throw.Throw, 3)
	}

	g.Base.Podium = make([]string, 0)

	g.Base.UndoLog = make([]undo.Undo, 0)
	g.Base.UndoLog = append(g.Base.UndoLog, undo.Undo{
		Sequence: 0,
		Action:   "CREATEGAME",
	})

	return nil
}

// GetStatus will satisfy interface Game for game Split
func (g *SplitGame) GetStatus() BaseGame {
	return g.Base
}

// GetStatusDisplay will satisfy interface Game for game Split
func (g *SplitGame) GetStatusDisplay() BaseGame {
	return stripDisplay(&g.Base)
}

// NextPlayer will satisfy interface Game for game Split
func (g *SplitGame) NextPlayer(h *ws.Hub) {
	switchToNextPlayer(&g.Base, h)
}

// RequestThrow will satisfy interface Game for game Split
func (g *SplitGame) RequestThrow(number, modifier int, h *ws.Hub) error {
	sequence := g.Base.UndoLog[len(g.Base.UndoLog)-1].Sequence + 1

	points := number * modifier
	activePlayer := &g.Base.Player[g.Base.ActivePlayer]

	// Check game state
	if g.Base.GameState == "THROW" {
		// Check if ongoing round
		ongoing := utils.CheckOngoingRound(activePlayer.ThrowRounds, g.Base.ThrowRound)
		if !ongoing {
			// If there is no round associated with the current throw round of the game
			// create one
			newRound := &throw.Round{
				Round:  g.Base.ThrowRound,
				Throws: []throw.Throw{},
			}
			activePlayer.ThrowRounds = append(activePlayer.ThrowRounds, *newRound)
			g.Base.UndoLog = append(g.Base.UndoLog, undo.Undo{Sequence: sequence, Action: "CREATETHROWROUND", RoundNumber: newRound.Round, Player: activePlayer})
		}
		// Now add throw to that round or to the existing
		throwRound := &activePlayer.ThrowRounds[g.Base.ThrowRound-1]
		newThrow := &throw.Throw{
			Number:   number,
			Modifier: modifier,
		}

		throwRound.Throws = append(throwRound.Throws, *newThrow)
		g.Base.UndoLog = append(g.Base.UndoLog, undo.Undo{Sequence: sequence, Action: "CREATETHROW", Player: activePlayer, RoundNumber: throwRound.Round})

		// Check if variant steel and if we need to build up score
		if g.Base.Variant == "steel" && g.Base.ThrowRound == 1 {
			chargeupScore(g, activePlayer, points, sequence)
		} else {
			splitLogic(g, activePlayer, number, modifier, sequence)
		}

		// Check if 3 throws in round and close round
		// Also set gameState and perhaps increase game.Base.ThrowRound
		// if everyone has already thrown to this round
		if len(throwRound.Throws) == 3 {
			previousMessage := g.Base.Message
			previousState := g.Base.GameState
			throwRound.Done = true
			g.Base.GameState = "NEXTPLAYER"
			g.Base.Message = "Remove Darts!"
			g.Base.UndoLog = append(g.Base.UndoLog, undo.Undo{Sequence: sequence, Action: "CLOSEPLAYERTHROWROUND", Player: activePlayer, RoundNumber: throwRound.Round, GameID: g.Base.UID, PreviousGameState: previousState, PreviousMessage: previousMessage})

			checkEndGame(&g.Base, sequence)
		}

		// Set assets for Frontend
		setFrontendAssets(activePlayer, &g.Base)

		// Update scoreboard
		utils.WSSendUpdate(g.Base.UID, h)
	}

	return nil
}

// Undo will satisfy interface Game for game Split
func (g *SplitGame) Undo(h *ws.Hub) error {
	if err := triggerUndo(&g.Base, h); err != nil {
		return err
	}

	return nil
}

// Rematch will satisfy interface Game for game Split
func (g *SplitGame) Rematch(h *ws.Hub) error {
	// Init random number generator
	s := rand.NewSource(time.Now().Unix())
	rg := rand.New(s)

	// Reset game state
	g.Base.Message = ""
	g.Base.GameState = "THROW"
	g.Base.Podium = make([]string, 0)
	g.Base.UndoLog = make([]undo.Undo, 0)
	g.Base.ActivePlayer = rg.Intn(len(g.Base.Player))
	g.Base.ThrowRound = 1

	// Init start score vor edart
	startscore := 40
	// set score to 0 if steel
	if g.Base.Variant == "steel" {
		startscore = 0
	}

	// CreateScore for each player
	// and init empty throw splice
	for i := range g.Base.Player {
		score := score.BaseScore{
			Score: startscore,
			Split: true,
		}
		g.Base.Player[i].Score = score
		g.Base.Player[i].ThrowRounds = make([]throw.Round, 0)
		g.Base.Player[i].LastThrows = make([]throw.Throw, 3)
	}

	g.Base.UndoLog = append(g.Base.UndoLog, undo.Undo{
		Sequence: 0,
		Action:   "CREATEGAME",
	})

	// Update scoreboard
	utils.WSSendUpdate(g.Base.UID, h)

	return nil
}

// This is to build score in steel variant
func chargeupScore(g *SplitGame, player *player.Player, points, sequence int) {
	previousScore := player.Score.Score
	player.Score.Score += points
	g.Base.UndoLog = append(g.Base.UndoLog, undo.Undo{Sequence: sequence, Action: "UPDATESPLITSCORE", PreviousScore: previousScore, Player: player})
}

// This is the logic to handle Split Game
func splitLogic(g *SplitGame, player *player.Player, number, modifier, sequence int) {
	rnd := g.Base.ThrowRound
	if g.Base.Variant == "steel" {
		// Decrease by one cause round 1 is to build up score
		rnd--
	}
	// Switch over Throw round as this indicates what needs to be hit
	switch rnd {
	case 1:
		// 15
		if number == 15 {
			chargeupScore(g, player, number*modifier, sequence)
			player.Score.Split = false
		}
	case 2:
		// 16
		if number == 16 {
			chargeupScore(g, player, number*modifier, sequence)
			player.Score.Split = false
		}
	case 3:
		// Double
		if modifier == 2 {
			chargeupScore(g, player, number*modifier, sequence)
			player.Score.Split = false
		}
	case 4:
		// 17
		if number == 17 {
			chargeupScore(g, player, number*modifier, sequence)
			player.Score.Split = false
		}
	case 5:
		// 18
		if number == 18 {
			chargeupScore(g, player, number*modifier, sequence)
			player.Score.Split = false
		}
	case 6:
		// Triple
		if modifier == 3 {
			chargeupScore(g, player, number*modifier, sequence)
			player.Score.Split = false
		}
	case 7:
		// 19
		if number == 19 {
			chargeupScore(g, player, number*modifier, sequence)
			player.Score.Split = false
		}
	case 8:
		// 20
		if number == 20 {
			chargeupScore(g, player, number*modifier, sequence)
			player.Score.Split = false
		}
	case 9:
		// 25
		if number == 25 {
			chargeupScore(g, player, number*modifier, sequence)
			player.Score.Split = false
		}
	default:
		break
	}

	// If throw round is done and not hit yet, split score
	// reset split indicator for next round
	if len(player.ThrowRounds[g.Base.ThrowRound-1].Throws) == 3 {
		checkAndSplit(&g.Base, player, sequence)
	}
}

func checkAndSplit(base *BaseGame, player *player.Player, sequence int) {
	if base.Variant == "steel" && base.ThrowRound == 1 {
		return
	}

	if player.Score.Split {
		previousScore := player.Score.Score
		player.Score.Score = player.Score.Score / 2
		base.UndoLog = append(base.UndoLog, undo.Undo{Sequence: sequence, Action: "UPDATESPLITSCORE", PreviousScore: previousScore, Player: player})
	}
	player.Score.Split = true
}

func checkEndGame(base *BaseGame, sequence int) {
	// Game ends if every player has 9 throw rounds which are done in edart
	// And 10 throw rounds which are done in steel dart
	var threshold int
	end := true
	switch base.Variant {
	case "steel":
		threshold = 10
	case "edart":
		threshold = 9
	default:
		break
	}
	for _, p := range base.Player {
		if len(p.ThrowRounds) != threshold {
			end = false
			break
		}
		for _, rnd := range p.ThrowRounds {
			if !rnd.Done {
				end = false
				break
			}
		}
	}

	if end {
		previousState := base.GameState
		previousMessage := base.Message

		doWin(base)
		base.UndoLog = append(base.UndoLog, undo.Undo{Sequence: sequence, Action: "DOWIN", GameID: base.UID, PreviousGameState: previousState, PreviousMessage: previousMessage})

		// Construct winner output
		highestScorePlayer := base.Player[0]
		for _, p := range base.Player {
			if p.Score.Score > highestScorePlayer.Score.Score {
				highestScorePlayer = p
			}

		}
		winnerName := highestScorePlayer.Name
		if highestScorePlayer.Nickname != "" {
			winnerName += " - " + highestScorePlayer.Nickname
		}
		base.Message = fmt.Sprintf("Game Shot! Winner is %+v", winnerName)
	}
}
