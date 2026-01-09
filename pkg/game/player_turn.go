package game

import (
	"math/rand"

	"github.com/engpetarmarinov/eepers-go/pkg/audio"
	"github.com/engpetarmarinov/eepers-go/pkg/entities"
	"github.com/engpetarmarinov/eepers-go/pkg/world"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Direction represents a direction of movement.
type Direction int

const (
	Left Direction = iota
	Right
	Up
	Down
)

// DirectionVector maps directions to their corresponding vectors.
var DirectionVector = map[Direction]world.IVector2{
	Left:  {X: -1, Y: 0},
	Right: {X: 1, Y: 0},
	Up:    {X: 0, Y: -1},
	Down:  {X: 0, Y: 1},
}

// PlayerTurn handles the player's turn.
func (gs *State) PlayerTurn(dir Direction) {
	gs.Player.PrevPosition = gs.Player.Position
	newPos := gs.Player.Position.Add(DirectionVector[dir])

	if !gs.WithinMap(newPos) {
		return
	}

	switch gs.Map[newPos.Y][newPos.X] {
	case world.CellFloor:
		gs.Player.Position = newPos
		rl.PlaySound(audio.FootstepsSounds[rand.Intn(len(audio.FootstepsSounds))])
		for i := range gs.Items {
			item := &gs.Items[i]
			if item.Position == newPos {
				switch item.Kind {
				case entities.ItemKey:
					gs.Player.Keys++
					item.Kind = entities.ItemNone // Mark as collected
					rl.PlaySound(audio.KeyPickupSound)
				case entities.ItemBombRefill:
					if gs.Player.Bombs < gs.Player.BombSlots {
						gs.Player.Bombs++
						rl.PlaySound(audio.BombPickupSound)
						item.Kind = entities.ItemNone // Mark as collected
					}
				case entities.ItemBombSlot:
					gs.Player.BombSlots++
					item.Kind = entities.ItemNone // Mark as collected
				case entities.ItemCheckpoint:
					// Handle checkpoint logic here
					item.Kind = entities.ItemNone // Mark as collected
					rl.PlaySound(audio.CheckpointSound)
				}
			}
		}
	case world.CellDoor:
		if gs.Player.Keys > 0 {
			gs.Player.Keys--
			gs.Map[newPos.Y][newPos.X] = world.CellFloor
			gs.Player.Position = newPos
			rl.PlaySound(audio.OpenDoorSound)
		}
	}
}

// WithinMap checks if a position is within the map boundaries.
func (gs *State) WithinMap(pos world.IVector2) bool {
	return pos.Y >= 0 && pos.Y < len(gs.Map) && pos.X >= 0 && pos.X < len(gs.Map[0])
}
