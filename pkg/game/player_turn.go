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

func (gs *State) RemoveDoor(startPos world.IVector2) {
	q := []world.IVector2{startPos}
	visited := make(map[world.IVector2]bool)
	visited[startPos] = true

	for len(q) > 0 {
		curr := q[0]
		q = q[1:]

		if gs.WithinMap(curr) && (gs.Map[curr.Y][curr.X] == world.CellDoor || gs.Map[curr.Y][curr.X] == world.CellBarricade) {
			isBarricade := gs.Map[curr.Y][curr.X] == world.CellBarricade
			gs.Map[curr.Y][curr.X] = world.CellFloor

			// Cardinal directions
			for _, dir := range DirectionVector {
				next := curr.Add(dir)
				if _, found := visited[next]; !found {
					if gs.WithinMap(next) && ((isBarricade && gs.Map[next.Y][next.X] == world.CellBarricade) || (!isBarricade && gs.Map[next.Y][next.X] == world.CellDoor)) {
						q = append(q, next)
						visited[next] = true
					}
				}
			}
			// Diagonal directions
			for _, diag := range []world.IVector2{
				{X: -1, Y: -1},
				{X: -1, Y: 1},
				{X: 1, Y: -1},
				{X: 1, Y: 1},
			} {
				next := curr.Add(diag)
				if _, found := visited[next]; !found {
					if gs.WithinMap(next) && ((isBarricade && gs.Map[next.Y][next.X] == world.CellBarricade) || (!isBarricade && gs.Map[next.Y][next.X] == world.CellDoor)) {
						q = append(q, next)
						visited[next] = true
					}
				}
			}
		}
	}
}

// PlayerTurn handles the player's turn.
func (gs *State) PlayerTurn(dir Direction) {
	gs.Player.PrevPosition = gs.Player.Position
	newPos := gs.Player.Position.Add(DirectionVector[dir])

	// Set eyes target to look in the direction of movement
	gs.Player.EyesTarget = newPos.Add(DirectionVector[dir])

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
					// Only pick up if we have space and the item is not on cooldown
					if gs.Player.Bombs < gs.Player.BombSlots && item.Cooldown <= 0 {
						gs.Player.Bombs++
						item.Cooldown = 10 // BOMB_GENERATOR_COOLDOWN
						rl.PlaySound(audio.BombPickupSound)
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
			gs.RemoveDoor(newPos)
			gs.Player.Position = newPos
			rl.PlaySound(audio.OpenDoorSound)
		}
	case world.CellBarricade:
		// Player cannot move through barricades
	}
}

// WithinMap checks if a position is within the map boundaries.
func (gs *State) WithinMap(pos world.IVector2) bool {
	return pos.Y >= 0 && pos.Y < len(gs.Map) && pos.X >= 0 && pos.X < len(gs.Map[0])
}
