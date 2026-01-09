package game

import (
	"github.com/engpetarmarinov/eepers-go/pkg/audio"
	"github.com/engpetarmarinov/eepers-go/pkg/entities"
	"github.com/engpetarmarinov/eepers-go/pkg/world"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	bombCountdown = 5
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
func (gs *State) Explode(position world.IVector2) {
	// Create an explosion at the bomb's location
	gs.Explosions = append(gs.Explosions, entities.ExplosionState{
		Position:     position,
		Timer:        20, // Explosion lasts for 20 frames
		InitialTimer: 20,
	})
	gs.Map[position.Y][position.X] = world.CellExplosion

	// And in all four directions
	for _, dir := range Directions {
		for i := 1; i <= 4; i++ {
			pos := position.Add(dir.Mul(i))

			mapWidth := len(gs.Map[0])
			mapHeight := len(gs.Map)
			if pos.X < 0 || pos.X >= mapWidth || pos.Y < 0 || pos.Y >= mapHeight {
				break // Stop if we go out of bounds
			}

			if gs.Map[pos.Y][pos.X] == world.CellWall {
				break // Stop if we hit a wall
			}

			// If we hit a barricade, flood fill it with explosions and stop
			if gs.Map[pos.Y][pos.X] == world.CellBarricade {
				gs.FloodFill(pos, world.CellBarricade, world.CellExplosion)
				break
			}

			gs.Explosions = append(gs.Explosions, entities.ExplosionState{
				Position:     pos,
				Timer:        20,
				InitialTimer: 20,
			})
			gs.Map[pos.Y][pos.X] = world.CellExplosion
		}
	}

	rl.PlaySound(audio.BlastSound)
}

// FloodFill fills all connected cells of the same type with a new cell type.
// This is used to destroy entire barricades when an explosion hits them.
func (gs *State) FloodFill(start world.IVector2, background world.Cell, fill world.Cell) {
	mapWidth := len(gs.Map[0])
	mapHeight := len(gs.Map)

	// Check if start position is valid
	if start.X < 0 || start.X >= mapWidth || start.Y < 0 || start.Y >= mapHeight {
		return
	}

	// Initialize queue with start position
	queue := []world.IVector2{start}
	gs.Map[start.Y][start.X] = fill

	// Add explosion at the start position
	gs.Explosions = append(gs.Explosions, entities.ExplosionState{
		Position:     start,
		Timer:        20,
		InitialTimer: 20,
	})

	// BFS flood fill
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// Check all four directions
		for _, dir := range Directions {
			newPos := current.Add(dir)

			// Check if position is valid and has the background cell type
			if newPos.X >= 0 && newPos.X < mapWidth && newPos.Y >= 0 && newPos.Y < mapHeight {
				if gs.Map[newPos.Y][newPos.X] == background {
					gs.Map[newPos.Y][newPos.X] = fill
					queue = append(queue, newPos)

					// Add explosion at this position
					gs.Explosions = append(gs.Explosions, entities.ExplosionState{
						Position:     newPos,
						Timer:        20,
						InitialTimer: 20,
					})
				}
			}
		}
	}
}
