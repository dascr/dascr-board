package podium

import (
	"github.com/dascr/dascr-board/player"
)

// Placement will hold the actual place on Podium
type Placement struct {
	Number int
	Player *player.Player
}

// Podium will hold Podium information
type Podium []*Placement

// AddPlayerToPodium will add a player to Podium
func (p *Podium) AddPlayerToPodium(player *player.Player) {
	lastPlacement := 1
	if len(*p) > 0 {
		lastPlacement = (*p)[len(*p)-1].Number + 1
	}

	*p = append(*p, &Placement{
		Number: lastPlacement,
		Player: player,
	})
}

// GetPodiumLength will return count of player on podium as int
func (p *Podium) GetPodiumLength() int {
	return len(*p)
}

// GetPodium will return the Podium itself
func (p *Podium) GetPodium() *Podium {
	return p
}

// ResetPodium will clear the podium for a new game
func (p *Podium) ResetPodium() {
	*p = Podium{}
}

// RemoveLastPlacement will remove last placement from podium
func (p *Podium) RemoveLastPlacement() {
	if len(*p) > 0 {
		*p = (*p)[:len(*p)-1]
	}
}
