package game

import (
	"github.com/engpetarmarinov/eepers-go/pkg/entities"
	"github.com/engpetarmarinov/eepers-go/pkg/pathfinding"
	"github.com/engpetarmarinov/eepers-go/pkg/world"
)

const (
	guardAttackCooldown = 5
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
			// To be implemented
		case entities.EeperGnome:
			// To be implemented
		case entities.EeperFather:
			// To be implemented
		}
	}
}

func (gs *State) updateGuard(eeper *entities.EeperState) {
	if eeper.Health <= 0 {
		eeper.Dead = true
		return
	}

	if eeper.AttackCooldown > 0 {
		eeper.AttackCooldown--
	}

	path := pathfinding.BFS(gs.Map, pathfinding.Point{X: eeper.Position.X, Y: eeper.Position.Y}, pathfinding.Point{X: gs.Player.Position.X, Y: gs.Player.Position.Y})

	if len(path) > 1 {
		nextPos := path[1]
		if gs.eeperCanStandHere(world.IVector2{X: nextPos.X, Y: nextPos.Y}, eeper) {
			eeper.PrevPosition = eeper.Position
			eeper.Position.X = nextPos.X
			eeper.Position.Y = nextPos.Y
		}
	} else if eeper.AttackCooldown == 0 {
		// Check if player is within attack range
		if gs.isPlayerInAttackRange(eeper) {
			gs.Player.Health -= 0.35
			if gs.Player.Health <= 0 {
				gs.KillPlayer()
			}
			eeper.AttackCooldown = guardAttackCooldown
		}
	}
}

func (gs *State) eeperCanStandHere(pos world.IVector2, currentEeper *entities.EeperState) bool {
	if !gs.WithinMap(pos) || gs.Map[pos.Y][pos.X] != world.CellFloor {
		return false
	}

	for i := range gs.Eepers {
		eeper := &gs.Eepers[i]
		if eeper != currentEeper && !eeper.Dead {
			if pos.X >= eeper.Position.X && pos.X < eeper.Position.X+eeper.Size.X &&
				pos.Y >= eeper.Position.Y && pos.Y < eeper.Position.Y+eeper.Size.Y {
				return false
			}
		}
	}

	return true
}

func (gs *State) isPlayerInAttackRange(eeper *entities.EeperState) bool {
	return gs.Player.Position.X >= eeper.Position.X && gs.Player.Position.X < eeper.Position.X+eeper.Size.X &&
		gs.Player.Position.Y >= eeper.Position.Y && gs.Player.Position.Y < eeper.Position.Y+eeper.Size.Y
}
