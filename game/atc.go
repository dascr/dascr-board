package game

import (
	"math/rand"
	"time"

	"github.com/dascr/dascr-board/player"
	"github.com/dascr/dascr-board/score"
	"github.com/dascr/dascr-board/throw"
	"github.com/dascr/dascr-board/undo"
	"github.com/dascr/dascr-board/utils"
	"github.com/dascr/dascr-board/ws"
)

// ATCGame will hold the ATC game information
type ATCGame struct {
	Base BaseGame
}

// StartGame will satisfy interface Game for ATC
func (g *ATCGame) StartGame() error {
	// CreateScore for each player
	// and init empty thtrow splice
	for i := range g.Base.Player {
		score := score.BaseScore{
			CurrentNumber: 1,
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

// GetStatus will satisfy interface Game for ATC
func (g *ATCGame) GetStatus() BaseGame {
	return g.Base
}

// GetStatusDisplay will satisfy interface Game for game ATC
func (g *ATCGame) GetStatusDisplay() BaseGame {
	return stripDisplay(&g.Base)
}

// NextPlayer will satisfy interface Game for game ATC
func (g *ATCGame) NextPlayer(h *ws.Hub) {
	switchToNextPlayer(&g.Base, h)
}

// RequestThrow will satisfy inteface Game for ATC
func (g *ATCGame) RequestThrow(number, modifier int, h *ws.Hub) error {
	sequence := g.Base.UndoLog[len(g.Base.UndoLog)-1].Sequence + 1
	activePlayer := &g.Base.Player[g.Base.ActivePlayer]
	previousMessage := g.Base.Message
	previousState := g.Base.GameState

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

		// Score logic
		if activePlayer.Score.CurrentNumber == number {
			switch g.Base.Variant {
			case "normal":
				increaseNumberByOne(g, throwRound, activePlayer, sequence)
				break
			case "fast":
				for i := 0; i < modifier; i++ {
					increaseNumberByOne(g, throwRound, activePlayer, sequence)
				}
				break
			default:
				break
			}
		}
		// Check if 3 throws in round and close round
		// Also set gameState and perhaps increase game.Base.ThrowRound
		// if everyone has already thrown to this round
		if len(throwRound.Throws) == 3 {
			throwRound.Done = true
			g.Base.GameState = "NEXTPLAYER"
			g.Base.Message = "Remove Darts!"
			g.Base.UndoLog = append(g.Base.UndoLog, undo.Undo{Sequence: sequence, Action: "CLOSEPLAYERTHROWROUND", Player: activePlayer, RoundNumber: throwRound.Round, GameID: g.Base.UID, PreviousGameState: previousState, PreviousMessage: previousMessage})
		}
	}

	// Set assets for Frontend
	setFrontendAssets(activePlayer, &g.Base)

	// Update scoreboard
	utils.WSSendUpdate(g.Base.UID, h)

	return nil
}

// Undo will satisfy interface Game for game ATC
func (g *ATCGame) Undo(h *ws.Hub) error {
	if err := triggerUndo(&g.Base, h); err != nil {
		return err
	}

	return nil
}

// Rematch will satisfy interface Game for game ATC
func (g *ATCGame) Rematch(h *ws.Hub) error {
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

	for i := range g.Base.Player {
		score := score.BaseScore{
			CurrentNumber: 1,
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

// increases active players number by one
// will go from 20 to 25
// will go from 25 to win
func increaseNumberByOne(game *ATCGame, throwRound *throw.Round, p *player.Player, sequence int) {
	previousNumberToHit := p.Score.CurrentNumber
	if p.Score.CurrentNumber == 25 {
		atcwin(game, throwRound, p, sequence)
	} else if p.Score.CurrentNumber != 20 {
		p.Score.CurrentNumber++
	} else {
		p.Score.CurrentNumber = 25
	}
	game.Base.UndoLog = append(game.Base.UndoLog, undo.Undo{Sequence: sequence, Action: "ATCINCREASENUMBER", GameID: game.Base.UID, Player: p, PreviousNumberToHit: previousNumberToHit})
}

// atcwin is for winning the atc game
func atcwin(game *ATCGame, throwRound *throw.Round, activePlayer *player.Player, sequence int) {
	previousMessage := game.Base.Message
	previousState := game.Base.GameState

	if game.Base.Settings.Podium {
		// Do podium and continue game
		doPodium(&game.Base, activePlayer, sequence)
		return
	}

	doWin(&game.Base)
	throwRound.Done = true
	game.Base.UndoLog = append(game.Base.UndoLog, undo.Undo{Sequence: sequence, Action: "DOWIN", GameID: game.Base.UID, PreviousGameState: previousState, PreviousMessage: previousMessage, Player: activePlayer})
}
