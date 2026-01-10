package game

import (
	"math/rand"

	"github.com/engpetarmarinov/eepers-go/pkg/audio"
	"github.com/engpetarmarinov/eepers-go/pkg/entities"
	"github.com/engpetarmarinov/eepers-go/pkg/pathfinding"
	"github.com/engpetarmarinov/eepers-go/pkg/world"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	guardAttackCooldown   = 10
	guardStepsLimit       = 100 // How many pathfinding steps to search
	guardStepLengthLimit  = 100 // How far to look in each direction during pathfinding
	guardTurnRegeneration = 0.01
	gnomeStepsLimit       = 9 // Gnomes can detect player up to 9 steps away
	gnomeStepLengthLimit  = 1 // Gnomes move only 1 cell at a time
	fatherWakeUpRadius    = 3 // Father wakes up when player is within 3 cells
)

// UpdateEepers updates the state of all eepers.
func (gs *State) UpdateEepers() {
	for i := range gs.Eepers {
		eeper := &gs.Eepers[i]
		if eeper.Dead {
			continue
		}

		switch eeper.Kind {
		case entities.EeperGuard:
			gs.updateGuard(eeper)
		case entities.EeperMother:
			gs.updateMother(eeper)
		case entities.EeperGnome:
			gs.updateGnome(eeper)
		case entities.EeperFather:
			gs.updateFather(eeper)
		}
	}
}

func (gs *State) updateGuard(eeper *entities.EeperState) {
	if eeper.Health <= 0 {
		eeper.Dead = true
		return
	}

	// Store the current position before any movement
	oldPosition := eeper.Position
	oldEyes := eeper.Eyes

	// Recompute the distance map for this guard
	gs.recomputePathForEeper(eeper)

	// Check the distance to player at guard's current position
	currentDist := eeper.Path[eeper.Position.Y][eeper.Position.X]

	// If guard is at player position (distance 0), kill player
	if currentDist == 0 {
		gs.KillPlayer()
		eeper.Eyes = entities.EyesSurprised
		eeper.PrevPosition = oldPosition
		eeper.PrevEyes = oldEyes
		return
	}

	// If player is reachable (distance > 0)
	if currentDist > 0 {
		// Guard is awake and tracking player
		if eeper.AttackCooldown <= 0 {
			// Try to move closer to player
			moved := gs.moveGuardTowardPlayer(eeper)
			if moved {
				rl.PlaySound(audio.GuardStepSound)
			}
			eeper.AttackCooldown = guardAttackCooldown
		} else {
			// Decrement cooldown while waiting
			eeper.AttackCooldown--
		}

		// Set eye state based on distance
		if eeper.Path[eeper.Position.Y][eeper.Position.X] == 1 {
			eeper.Eyes = entities.EyesAngry
		} else {
			eeper.Eyes = entities.EyesOpen
		}
		eeper.EyesTarget = gs.Player.Position

		// Check if guard caught player by moving onto them
		if gs.isPlayerInAttackRange(eeper) {
			gs.KillPlayer()
		}
	} else {
		// Player is not reachable - guard is sleeping
		eeper.Eyes = entities.EyesClosed
		eeper.EyesTarget = world.IVector2{
			X: eeper.Position.X + eeper.Size.X/2,
			Y: eeper.Position.Y + eeper.Size.Y,
		}
		eeper.AttackCooldown = guardAttackCooldown + 1
	}

	// Health regeneration
	if eeper.Health < 1.0 {
		eeper.Health += guardTurnRegeneration
		if eeper.Health > 1.0 {
			eeper.Health = 1.0
		}
	}

	// Set previous position AFTER all movement and state changes
	// This ensures interpolation works correctly
	eeper.PrevPosition = oldPosition
	eeper.PrevEyes = oldEyes
}

// updateMother updates a Mother eeper - behaves exactly like guards but larger
func (gs *State) updateMother(eeper *entities.EeperState) {
	// Mother eepers behave identically to guards, just with different size
	// Reuse the guard update logic
	gs.updateGuard(eeper)
}

// updateFather updates the Father eeper - the goal of the game!
func (gs *State) updateFather(eeper *entities.EeperState) {
	// Set previous position at START of turn
	eeper.PrevPosition = eeper.Position
	eeper.PrevEyes = eeper.Eyes

	// Check if player is touching Father (victory condition!)
	if gs.isPlayerInAttackRange(eeper) {
		// Player reached Father - trigger victory!
		if !gs.Player.ReachedFather {
			// First time reaching - record the time
			gs.Player.ReachedFather = true
			gs.Player.VictoryTime = rl.GetTime()
		}
		return
	}

	// Check if player is within wake-up radius (3 cells around Father)
	wakeUpRect := struct {
		X, Y, W, H int
	}{
		X: eeper.Position.X - fatherWakeUpRadius,
		Y: eeper.Position.Y - fatherWakeUpRadius,
		W: eeper.Size.X + fatherWakeUpRadius*2,
		H: eeper.Size.Y + fatherWakeUpRadius*2,
	}

	playerInWakeRadius := gs.Player.Position.X >= wakeUpRect.X &&
		gs.Player.Position.X < wakeUpRect.X+wakeUpRect.W &&
		gs.Player.Position.Y >= wakeUpRect.Y &&
		gs.Player.Position.Y < wakeUpRect.Y+wakeUpRect.H

	if playerInWakeRadius {
		// Player is nearby - Father wakes up and tracks with eyes
		eeper.Eyes = entities.EyesOpen
		eeper.EyesTarget = gs.Player.Position
	} else {
		// Player is far - Father sleeps
		eeper.Eyes = entities.EyesClosed
		eeper.EyesTarget = world.IVector2{
			X: eeper.Position.X + eeper.Size.X/2,
			Y: eeper.Position.Y + eeper.Size.Y,
		}
	}
}

func (gs *State) moveGuardTowardPlayer(eeper *entities.EeperState) bool {
	currentDist := eeper.Path[eeper.Position.Y][eeper.Position.X]
	if currentDist <= 0 {
		return false
	}

	// Find all positions that are one step closer to the player
	// The guard will JUMP to these positions
	directions := []world.IVector2{
		{X: 0, Y: 1},
		{X: 0, Y: -1},
		{X: 1, Y: 0},
		{X: -1, Y: 0},
	}
	var availablePositions []world.IVector2

	for _, dir := range directions {
		pos := eeper.Position

		// Keep moving in this direction, checking each step
		for {
			// Try to move one step in this direction
			nextPos := world.IVector2{X: pos.X + dir.X, Y: pos.Y + dir.Y}

			// Check if we're out of bounds
			if !gs.WithinMap(nextPos) {
				break
			}

			// Check if we can stand at this new position
			if !gs.eeperCanStandHere(nextPos, eeper) {
				break
			}

			// Move to the new position
			pos = nextPos

			// Check if this position has the right distance (one step closer to player)
			if eeper.Path[pos.Y][pos.X] == currentDist-1 {
				// We found a valid position that's closer to the player
				availablePositions = append(availablePositions, pos)
				break
			}

			// If we reached the player position (distance 0), stop
			if eeper.Path[pos.Y][pos.X] == 0 {
				break
			}
		}
	}

	// If we found valid moves, pick one randomly
	if len(availablePositions) > 0 {
		newPos := availablePositions[rand.Intn(len(availablePositions))]
		eeper.Position = newPos
		return true
	}

	return false
}

func (gs *State) recomputePathForEeper(eeper *entities.EeperState) {
	// Create a function to check if the eeper can stand at a position
	canStand := func(p pathfinding.Point) bool {
		pos := world.IVector2{X: p.X, Y: p.Y}
		return gs.eeperCanStandHere(pos, eeper)
	}

	eeper.Path = pathfinding.ComputeDistanceMap(
		gs.Map,
		pathfinding.Point{X: gs.Player.Position.X, Y: gs.Player.Position.Y},
		pathfinding.Point{X: eeper.Size.X, Y: eeper.Size.Y},
		guardStepsLimit,
		guardStepLengthLimit,
		canStand,
	)
}

func (gs *State) eeperCanStandHere(pos world.IVector2, currentEeper *entities.EeperState) bool {
	// Check ALL cells that the eeper occupies (e.g., 3x3 for guards)
	// This is critical - guards can't move through walls!
	for x := pos.X; x < pos.X+currentEeper.Size.X; x++ {
		for y := pos.Y; y < pos.Y+currentEeper.Size.Y; y++ {
			cellPos := world.IVector2{X: x, Y: y}

			// Check map bounds
			if !gs.WithinMap(cellPos) {
				return false
			}

			// Check if cell is floor or explosion (can step into explosions)
			cell := gs.Map[cellPos.Y][cellPos.X]
			if cell != world.CellFloor && cell != world.CellExplosion {
				return false
			}

			// Check collision with other eepers
			for i := range gs.Eepers {
				eeper := &gs.Eepers[i]
				if eeper != currentEeper && !eeper.Dead {
					// Check if this cell overlaps with another eeper
					if cellPos.X >= eeper.Position.X && cellPos.X < eeper.Position.X+eeper.Size.X &&
						cellPos.Y >= eeper.Position.Y && cellPos.Y < eeper.Position.Y+eeper.Size.Y {
						return false
					}
				}
			}
		}
	}

	return true
}

func (gs *State) isPlayerInAttackRange(eeper *entities.EeperState) bool {
	return gs.Player.Position.X >= eeper.Position.X && gs.Player.Position.X < eeper.Position.X+eeper.Size.X &&
		gs.Player.Position.Y >= eeper.Position.Y && gs.Player.Position.Y < eeper.Position.Y+eeper.Size.Y
}

// updateGnome updates a gnome eeper that flees from the player
func (gs *State) updateGnome(eeper *entities.EeperState) {
	if eeper.Health <= 0 {
		eeper.Dead = true
		return
	}

	// Store the current position before any movement
	oldPosition := eeper.Position
	oldEyes := eeper.Eyes

	// Recompute path for gnome (they use different pathfinding params)
	gs.recomputePathForGnome(eeper)

	// Check if player is reachable
	currentDist := eeper.Path[eeper.Position.Y][eeper.Position.X]

	if currentDist >= 0 {
		// Player is reachable - gnome flees (moves to higher distance)
		gs.moveGnomeAwayFromPlayer(eeper)
		eeper.Eyes = entities.EyesOpen
		eeper.EyesTarget = gs.Player.Position
	} else {
		// Player not reachable - gnome sleeps
		eeper.Eyes = entities.EyesClosed
		eeper.EyesTarget = world.IVector2{
			X: eeper.Position.X + eeper.Size.X/2,
			Y: eeper.Position.Y + eeper.Size.Y,
		}
	}

	// Set previous position AFTER all movement and state changes
	eeper.PrevPosition = oldPosition
	eeper.PrevEyes = oldEyes
}

// moveGnomeAwayFromPlayer moves the gnome to a position farther from the player
func (gs *State) moveGnomeAwayFromPlayer(eeper *entities.EeperState) {
	currentDist := eeper.Path[eeper.Position.Y][eeper.Position.X]
	if currentDist < 0 {
		return
	}

	// Find all adjacent positions with HIGHER distance (fleeing)
	directions := []world.IVector2{
		{X: 0, Y: 1},
		{X: 0, Y: -1},
		{X: 1, Y: 0},
		{X: -1, Y: 0},
	}
	var availablePositions []world.IVector2

	for _, dir := range directions {
		newPos := eeper.Position.Add(dir)

		// Check if position is valid
		if !gs.WithinMap(newPos) {
			continue
		}

		// Check if it's a floor cell
		if gs.Map[newPos.Y][newPos.X] != world.CellFloor {
			continue
		}

		// Gnomes flee - move to positions with HIGHER distance
		if eeper.Path[newPos.Y][newPos.X] > currentDist {
			availablePositions = append(availablePositions, newPos)
		}
	}

	// If found positions to flee to, pick one randomly
	if len(availablePositions) > 0 {
		eeper.Position = availablePositions[rand.Intn(len(availablePositions))]
	}
}

// recomputePathForGnome computes distance map for gnome (different params than guards)
func (gs *State) recomputePathForGnome(eeper *entities.EeperState) {
	canStand := func(p pathfinding.Point) bool {
		pos := world.IVector2{X: p.X, Y: p.Y}
		if !gs.WithinMap(pos) {
			return false
		}
		return gs.Map[pos.Y][pos.X] == world.CellFloor
	}

	eeper.Path = pathfinding.ComputeDistanceMap(
		gs.Map,
		pathfinding.Point{X: gs.Player.Position.X, Y: gs.Player.Position.Y},
		pathfinding.Point{X: eeper.Size.X, Y: eeper.Size.Y},
		gnomeStepsLimit,      // 9 steps
		gnomeStepLengthLimit, // 1 cell per step
		canStand,
	)
}

// SpawnGuard creates a new guard at the specified position
func (gs *State) SpawnGuard(position world.IVector2) {
	size := world.IVector2{X: 3, Y: 3}

	// Initialize path map with correct dimensions
	height := len(gs.Map)
	width := len(gs.Map[0])
	path := make([][]int, height)
	for i := range path {
		path[i] = make([]int, width)
		for j := range path[i] {
			path[i][j] = -1
		}
	}

	guard := entities.EeperState{
		Kind:           entities.EeperGuard,
		Dead:           false,
		Position:       position,
		PrevPosition:   position,
		EyesAngle:      0,
		EyesTarget:     world.IVector2{X: position.X + size.X/2, Y: position.Y + size.Y},
		PrevEyes:       entities.EyesClosed,
		Eyes:           entities.EyesClosed,
		Size:           size,
		Path:           path,
		Damaged:        false,
		Health:         1.0,
		AttackCooldown: guardAttackCooldown,
	}

	gs.Eepers = append(gs.Eepers, guard)
}

// SpawnGnome creates a new gnome at the specified position
func (gs *State) SpawnGnome(position world.IVector2) {
	size := world.IVector2{X: 1, Y: 1} // Gnomes are 1x1

	// Initialize path map with correct dimensions
	height := len(gs.Map)
	width := len(gs.Map[0])
	path := make([][]int, height)
	for i := range path {
		path[i] = make([]int, width)
		for j := range path[i] {
			path[i][j] = -1
		}
	}

	gnome := entities.EeperState{
		Kind:         entities.EeperGnome,
		Dead:         false,
		Position:     position,
		PrevPosition: position,
		EyesAngle:    0,
		EyesTarget:   world.IVector2{X: position.X + size.X/2, Y: position.Y + size.Y},
		PrevEyes:     entities.EyesClosed,
		Eyes:         entities.EyesClosed,
		Size:         size,
		Path:         path,
		Damaged:      false,
		Health:       1.0, // Gnomes start alive, die in one explosion hit
	}

	gs.Eepers = append(gs.Eepers, gnome)
}

// SpawnMother creates a new Mother eeper at the specified position
func (gs *State) SpawnMother(position world.IVector2) {
	size := world.IVector2{X: 7, Y: 7} // Mothers are 7x7 (large)

	// Initialize path map with correct dimensions
	height := len(gs.Map)
	width := len(gs.Map[0])
	path := make([][]int, height)
	for i := range path {
		path[i] = make([]int, width)
		for j := range path[i] {
			path[i][j] = -1
		}
	}

	mother := entities.EeperState{
		Kind:           entities.EeperMother,
		Dead:           false,
		Position:       position,
		PrevPosition:   position,
		EyesAngle:      0,
		EyesTarget:     world.IVector2{X: position.X + size.X/2, Y: position.Y + size.Y},
		PrevEyes:       entities.EyesClosed,
		Eyes:           entities.EyesClosed,
		Size:           size,
		Path:           path,
		Damaged:        false,
		Health:         1.0,
		AttackCooldown: guardAttackCooldown,
	}

	gs.Eepers = append(gs.Eepers, mother)
}

// SpawnFather creates a new Father eeper at the specified position - the goal!
func (gs *State) SpawnFather(position world.IVector2) {
	size := world.IVector2{X: 7, Y: 7} // Father is 7x7 (large, same as Mother)

	// Father doesn't need path map (doesn't chase)
	father := entities.EeperState{
		Kind:         entities.EeperFather,
		Dead:         false,
		Position:     position,
		PrevPosition: position,
		EyesAngle:    0,
		EyesTarget:   world.IVector2{X: position.X + size.X/2, Y: position.Y + size.Y},
		PrevEyes:     entities.EyesClosed,
		Eyes:         entities.EyesClosed,
		Size:         size,
		Path:         nil, // Father doesn't use pathfinding
		Damaged:      false,
		Health:       1.0, // Immune to damage anyway
	}

	gs.Eepers = append(gs.Eepers, father)
}
