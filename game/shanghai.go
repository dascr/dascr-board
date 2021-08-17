package game

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/dascr/dascr-board/logger"
	"github.com/dascr/dascr-board/player"
	"github.com/dascr/dascr-board/podium"
	"github.com/dascr/dascr-board/score"
	"github.com/dascr/dascr-board/throw"
	"github.com/dascr/dascr-board/undo"
	"github.com/dascr/dascr-board/utils"
	"github.com/dascr/dascr-board/ws"
)

// ShanghaiGame will hold the Shanghai game information
type ShanghaiGame struct {
	Base BaseGame
}

// StartGame will satisfy interface Game for Shanghai
func (g *ShanghaiGame) StartGame() error {
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

	g.Base.Podium = &podium.Podium{}

	g.Base.UndoLog = &undo.Log{}
	sequence := g.Base.UndoLog.CreateSequence()
	sequence.AddActionToSequence(undo.Action{
		Action: "CREATEGAME",
	})

	return nil
}

// GetStatus will satisfy interface Game for Shanghai
func (g *ShanghaiGame) GetStatus() BaseGame {
	return g.Base
}

// GetStatusDisplay will satisfy interface Game for game Shanghai
func (g *ShanghaiGame) GetStatusDisplay() BaseGame {
	return stripDisplay(&g.Base)
}

// NextPlayer will satisfy interface Game for game Shanghai
func (g *ShanghaiGame) NextPlayer(h *ws.Hub) {
	activePlayer := &g.Base.Player[g.Base.ActivePlayer]
	sequence, err := g.Base.UndoLog.GetLastSequence()
	if err != nil {
		logger.Errorf("Error getting last sequence in nextPlayer Split-Score: %+v", err)
	}
	checkIncrease(&g.Base, activePlayer, sequence)
	switchToNextPlayer(&g.Base, h)
	checkShanghaiEnd(&g.Base, sequence)
}

// RequestThrow will satisfy inteface Game for Shanghai
func (g *ShanghaiGame) RequestThrow(number, modifier int, h *ws.Hub) error {
	sequence := g.Base.UndoLog.CreateSequence()
	activePlayer := &g.Base.Player[g.Base.ActivePlayer]
	previousMessage := g.Base.Message
	previousState := g.Base.GameState

	if g.Base.GameState == "THROW" {
		// check if ongoing round else create
		checkOngoingElseCreate(activePlayer, &g.Base, sequence)

		// Add Throw to last round
		throwRound := addThrowToCurrentRound(activePlayer, &g.Base, sequence, number, modifier)

		// Score logic
		scoreIfHit(g, number, modifier, throwRound, activePlayer, sequence)

		// Check if 3 throws in round and close round
		// Also set gameState and perhaps increase game.Base.ThrowRound
		// if everyone has already thrown to this round
		if len(throwRound.Throws) == 3 {
			if !winIfShanghai(g, throwRound, activePlayer, sequence) {
				closePlayerRound(&g.Base, activePlayer, throwRound, sequence, []undo.Action{}, previousState, previousMessage)
			}
		}
	}

	// Set assets for Frontend
	setFrontendAssets(activePlayer, &g.Base)

	// Update scoreboard
	utils.WSSendUpdate(g.Base.UID, h)

	return nil
}

// Undo will satisfy interface Game for game Shanghai
func (g *ShanghaiGame) Undo(h *ws.Hub) error {
	if err := triggerUndo(&g.Base, h); err != nil {
		return err
	}

	return nil
}

// Rematch will satisfy interface Game for game Shanghai
func (g *ShanghaiGame) Rematch(h *ws.Hub) error {
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

	for i := range g.Base.Player {
		score := score.BaseScore{
			CurrentNumber: 1,
		}
		g.Base.Player[i].Score = score
		g.Base.Player[i].ThrowRounds = make([]throw.Round, 0)
		g.Base.Player[i].LastThrows = make([]throw.Throw, 3)
	}

	sequence := g.Base.UndoLog.CreateSequence()
	sequence.AddActionToSequence(undo.Action{
		Action: "CREATEGAME",
	})

	// Update scoreboard
	utils.WSSendUpdate(g.Base.UID, h)

	return nil
}

// scoreIfHit will just increase score if the current number was hit
func scoreIfHit(game *ShanghaiGame, number int, modifier int, throwRound *throw.Round, p *player.Player, sequence *undo.Sequence) {
	previousScore := p.Score.Score
	previousAverage := p.Average
	previousThrowSum := p.ThrowSum
	previousLastThree := p.LastThrows
	previousMessage := game.Base.Message
	previousState := game.Base.GameState
	previousNumberToHit := p.Score.CurrentNumber
	// Check if hit is relevant
	if p.Score.CurrentNumber == number {
		p.Score.Score += number * modifier
		sequence.AddActionToSequence(undo.Action{
			Action:            "UPDATESCORE",
			GameID:            game.Base.UID,
			PreviousGameState: previousState,
			PreviousMessage:   previousMessage,
			Player:            p,
			PreviousScore:     previousScore,
			PreviousAverage:   previousAverage,
			PreviousThrowSum:  previousThrowSum,
			PreviousLastThree: previousLastThree,
		})
	}

	if len(throwRound.Throws) == 3 {
		p.Score.CurrentNumber += 1
		sequence.AddActionToSequence(undo.Action{
			Action:              "ATCINCREASENUMBER",
			GameID:              game.Base.UID,
			Player:              p,
			PreviousNumberToHit: previousNumberToHit,
		})
	}
}

// winIfShanghai will check if throw round was shanghai and win/end the game
func winIfShanghai(game *ShanghaiGame, throwRound *throw.Round, p *player.Player, sequence *undo.Sequence) bool {
	// Before actually closing the round check if shanghai
	single := false
	double := false
	triple := false

	for _, thr := range throwRound.Throws {
		switch thr.Modifier {
		case 1:
			single = true
		case 2:
			double = true
		case 3:
			triple = true
		}
	}

	if single && double && triple {
		logger.Info("Shanghai was shot")
		previousState := game.Base.GameState
		previousMessage := game.Base.Message

		doWin(&game.Base)
		sequence.AddActionToSequence(undo.Action{
			Action:            "DOWIN",
			GameID:            game.Base.UID,
			PreviousGameState: previousState,
			PreviousMessage:   previousMessage,
		})

		winnerName := p.Name
		if p.Nickname != "" {
			winnerName += " - " + p.Nickname
		}
		game.Base.Message = fmt.Sprintf("Shanghai Shot! Winner is %+v", winnerName)
		return true
	}
	return false
}

// checkShanghaiEnd will handle game end after 20 throw rounds
func checkShanghaiEnd(base *BaseGame, sequence *undo.Sequence) {
	// Game ends if every player has 20 throw rounds which are done
	end := true
	for _, p := range base.Player {
		if len(p.ThrowRounds) != 20 {
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
		sequence.AddActionToSequence(undo.Action{
			Action:            "DOWIN",
			GameID:            base.UID,
			PreviousGameState: previousState,
			PreviousMessage:   previousMessage,
		})

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

// checkIncrease will increase players current number if applicable
func checkIncrease(base *BaseGame, p *player.Player, sequence *undo.Sequence) {
	previousNumberToHit := p.Score.CurrentNumber

	if base.ThrowRound != p.Score.CurrentNumber {
		p.Score.CurrentNumber = base.ThrowRound
		sequence.AddActionToSequence(undo.Action{
			Action:              "ATCINCREASENUMBER",
			GameID:              base.UID,
			Player:              p,
			PreviousNumberToHit: previousNumberToHit,
		})

	}
}
