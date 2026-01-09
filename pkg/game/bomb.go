package game

import (
	"github.com/engpetarmarinov/eepers-go/pkg/audio"
	"github.com/engpetarmarinov/eepers-go/pkg/entities"
	"github.com/engpetarmarinov/eepers-go/pkg/world"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	bombCountdown = 3
	explosionTime = 2
)

// PlantBomb creates a new bomb at the player's position.
func (gs *State) PlantBomb() {
	if gs.Player.Bombs > 0 {
		gs.Player.Bombs--
		gs.Bombs = append(gs.Bombs, entities.BombState{
			Position:  gs.Player.Position,
			Countdown: bombCountdown,
		})
		rl.PlaySound(audio.PlantBombSound)
	}
}

// UpdateBombs updates the state of all bombs.
func (gs *State) UpdateBombs() {
	for i := len(gs.Bombs) - 1; i >= 0; i-- {
		bomb := &gs.Bombs[i]
		bomb.Countdown--

		if bomb.Countdown == 0 {
			gs.Explode(bomb.Position)
			gs.Bombs = append(gs.Bombs[:i], gs.Bombs[i+1:]...)
		}
	}
}

// Explode creates an explosion at a given position.
func (gs *State) Explode(pos world.IVector2) {
	rl.PlaySound(audio.BlastSound)
	for y := pos.Y - 1; y <= pos.Y+1; y++ {
		for x := pos.X - 1; x <= pos.X+1; x++ {
			if gs.WithinMap(world.IVector2{X: x, Y: y}) {
				if gs.Player.Position.X == x && gs.Player.Position.Y == y {
					gs.KillPlayer()
				}

				for i := range gs.Eepers {
					eeper := &gs.Eepers[i]
					if !eeper.Dead && x >= eeper.Position.X && x < eeper.Position.X+eeper.Size.X && y >= eeper.Position.Y && y < eeper.Position.Y+eeper.Size.Y {
						eeper.Health -= 0.5
						if eeper.Kind == entities.EeperMother {
							gs.spawnGnome(eeper.Position)
						}
					}
				}

				gs.Map[y][x] = world.CellExplosion
			}
		}
	}
}
