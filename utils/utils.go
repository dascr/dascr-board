package utils

import (
	"math/rand"
	"sort"
	"time"

	"github.com/dascr/dascr-board/player"
	"github.com/dascr/dascr-board/podium"
	"github.com/dascr/dascr-board/throw"
	"github.com/dascr/dascr-board/ws"
)

// WSSendUpdate will trigger an update for a game
// using the websocket connection
func WSSendUpdate(id string, h *ws.Hub) {
	message := ws.Message{
		Room: id,
		Data: []byte("update"),
	}

	h.Broadcast <- message
}

// CheckOngoingRound will see if there is already an ongoing round
func CheckOngoingRound(rnds []throw.Round, currentRound int) bool {
	for _, rnd := range rnds {
		if rnd.Round == currentRound {
			return true
		}
	}
	return false
}

// CheckRoundDone will see if everyone has thrown to the ongoing game.ThrowRound
func CheckRoundDone(players []player.Player, currentRound int, podium *podium.Podium) bool {
	var rnds = 0
	var playersInGame []player.Player

	// Fill slice with players still in game
	for _, pl := range players {
		var contained = false
		for _, p := range *podium {
			if pl.UID == p.Player.UID {
				contained = true
			}
		}
		if !contained {
			playersInGame = append(playersInGame, pl)
		}
	}

	for _, pl := range playersInGame {
		for _, rnd := range pl.ThrowRounds {
			if rnd.Round == currentRound && rnd.Done {
				rnds++
			}
		}
	}

	return rnds == len(playersInGame)
}

// GetSingleRandomCricketNumber will return 1 random cricket number
func GetSingleRandomCricketNumber() int {
	n := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 25}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(n), func(i, j int) { n[i], n[j] = n[j], n[i] })
	return n[0]
}

// GetRandomCricketNumbers will return 7 random cricket numbers
func GetRandomCricketNumbers(order bool) []int {
	n := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 25}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(n), func(i, j int) { n[i], n[j] = n[j], n[i] })

	r := n[:7]

	// Order for mode random
	// Skip order for ghost
	if order {
		sort.Ints(r)
	}

	return r
}

// ChooseNextPlayer will return next index for player slice
// Currently used in X01, Cricket
func ChooseNextPlayer(playerSlice []player.Player, activeIndex int, podium *podium.Podium) int {
	playerCount := len(playerSlice)
	playerIndex := activeIndex

	for {
		playerIndex++
		if playerIndex >= playerCount {
			playerIndex -= playerCount
		}

		var contained = false
		for _, p := range *podium {
			if playerSlice[playerIndex].UID == p.Player.UID {
				contained = true
			}
		}
		if !contained {
			return playerIndex
		}
	}
}
