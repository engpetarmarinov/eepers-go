package game

import "github.com/engpetarmarinov/eepers-go/pkg/world"

// UpdateExplosions updates the state of all explosions.
func (gs *State) UpdateExplosions() {
	for i := len(gs.Explosions) - 1; i >= 0; i-- {
		explosion := &gs.Explosions[i]
		explosion.Timer--

		if explosion.Timer <= 0 {
			gs.Map[explosion.Position.Y][explosion.Position.X] = world.CellFloor
			gs.Explosions = append(gs.Explosions[:i], gs.Explosions[i+1:]...)
		}
	}
}
