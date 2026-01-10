package game

import (
	"github.com/engpetarmarinov/eepers-go/pkg/audio"
	"github.com/engpetarmarinov/eepers-go/pkg/entities"
	"github.com/engpetarmarinov/eepers-go/pkg/world"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	bombCountdown   = 3
	explosionDamage = 0.45
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
	// Reset damaged flag for all eepers at the start of the turn
	for i := range gs.Eepers {
		gs.Eepers[i].Damaged = false
	}

	// Update bomb countdowns and explode if needed
	for i := len(gs.Bombs) - 1; i >= 0; i-- {
		bomb := &gs.Bombs[i]
		bomb.Countdown--

		if bomb.Countdown == 0 {
			gs.Explode(bomb.Position)
			gs.Bombs = append(gs.Bombs[:i], gs.Bombs[i+1:]...)
		}
	}

	// Process damage to eepers after all explosions
	for i := range gs.Eepers {
		eeper := &gs.Eepers[i]
		if !eeper.Dead && eeper.Damaged {
			switch eeper.Kind {
			case entities.EeperGuard:
				eeper.Eyes = entities.EyesCringe
				eeper.Health -= explosionDamage
				if eeper.Health <= 0 {
					eeper.Dead = true
				}
			case entities.EeperMother:
				// Mother spawns 4 guards when killed
				position := eeper.Position
				eeper.Dead = true
				gs.SpawnGuard(world.IVector2{X: position.X, Y: position.Y})
				gs.SpawnGuard(world.IVector2{X: position.X + 4, Y: position.Y})
				gs.SpawnGuard(world.IVector2{X: position.X, Y: position.Y + 4})
				gs.SpawnGuard(world.IVector2{X: position.X + 4, Y: position.Y + 4})
			case entities.EeperGnome:
				// Gnome drops a key when killed
				eeper.Dead = true
				gs.AllocateItem(eeper.Position, entities.ItemKey)
			case entities.EeperFather:
				// Father is immune to explosions
			}
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

	// Damage eepers and player at explosion position
	gs.damageAtPosition(position)

	// And in all four directions
	for _, dir := range Directions {
		for i := 1; i <= 4; i++ {
			pos := position.Add(dir.Mul(i))
			mapWidth := len(gs.Map[0])
			mapHeight := len(gs.Map)
			if pos.X < 0 || pos.X >= mapWidth || pos.Y < 0 || pos.Y >= mapHeight {
				break // Stop if we go out of bounds
			}

			if gs.Map[pos.Y][pos.X] == world.CellWall || gs.Map[pos.Y][pos.X] == world.CellDoor {
				break // Stop if we hit a wall or a door
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

			// Damage eepers and player at this position
			gs.damageAtPosition(pos)
		}
	}

	rl.PlaySound(audio.BlastSound)
}

// damageAtPosition damages player and eepers at the given position
func (gs *State) damageAtPosition(pos world.IVector2) {
	// Damage player if at this position
	if gs.Player.Position.X == pos.X && gs.Player.Position.Y == pos.Y {
		gs.KillPlayer()
	}

	// Damage eepers that overlap with this position
	for i := range gs.Eepers {
		eeper := &gs.Eepers[i]
		if !eeper.Dead && gs.isInsideRect(eeper.Position, eeper.Size, pos) {
			eeper.Damaged = true
		}
	}
}

// isInsideRect checks if a point is inside a rectangle
func (gs *State) isInsideRect(rectPos world.IVector2, rectSize world.IVector2, point world.IVector2) bool {
	return point.X >= rectPos.X && point.X < rectPos.X+rectSize.X &&
		point.Y >= rectPos.Y && point.Y < rectPos.Y+rectSize.Y
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
