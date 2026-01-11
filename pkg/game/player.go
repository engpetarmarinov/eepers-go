package game

import (
	"math/rand"

	"github.com/engpetarmarinov/eepers-go/pkg/audio"
	"github.com/engpetarmarinov/eepers-go/pkg/entities"
	"github.com/engpetarmarinov/eepers-go/pkg/world"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// playerDirection represents a playerDirection of movement.
type playerDirection int

const (
	Left playerDirection = iota
	Right
	Up
	Down
)

// playerDirectionVector maps directions to their corresponding vectors.
var playerDirectionVector = map[playerDirection]world.IVector2{
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
			for _, dir := range playerDirectionVector {
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
func (gs *State) PlayerTurn(dir playerDirection) {
	gs.Player.PrevPosition = gs.Player.Position
	newPos := gs.Player.Position.Add(playerDirectionVector[dir])

	// Set eyes target to look in the playerDirection of movement
	gs.Player.EyesTarget = newPos.Add(playerDirectionVector[dir])

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
						// Advance tutorial when picking up first bomb
						if gs.Player.Bombs == 0 && gs.Tutorial.Phase == TutorialWaitingForBombPick {
							gs.Tutorial.Phase = TutorialPlaceBombs
						}
						gs.Player.Bombs++
						item.Cooldown = 10 // BOMB_GENERATOR_COOLDOWN
						rl.PlaySound(audio.BombPickupSound)
					}
				case entities.ItemBombSlot:
					gs.Player.BombSlots++
					item.Kind = entities.ItemNone // Mark as collected
				case entities.ItemCheckpoint:
					// Mark as collected first, then save state
					item.Kind = entities.ItemNone
					gs.SaveCheckpoint()
					rl.PlaySound(audio.CheckpointSound)
				}
			}
		}

		// Check if player stepped on a portal
		portal := gs.GetPortalAtPosition(newPos)
		if portal != nil && portal.OpenProgress > 0.8 {
			// Start portal entry animation instead of immediately activating
			gs.Player.EnteringPortal = true
			rl.PlaySound(audio.EnterPortalSound)
			gs.Player.PortalEntryTime = rl.GetTime()
			gs.Player.PortalToActivate = portal.ID
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

// KillPlayer marks the player as dead and records the time of death.
func (gs *State) KillPlayer() {
	if !gs.Player.Dead {
		rl.PlaySound(audio.HurtSound)
		gs.Player.Health = 0
		gs.Player.Dead = true
		gs.Player.DeathTime = rl.GetTime()
	}
}
