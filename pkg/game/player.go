package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// KillPlayer marks the player as dead and records the time of death.
func (gs *State) KillPlayer() {
	if !gs.Player.Dead {
		gs.Player.Health = 0
		gs.Player.Dead = true
		gs.Player.DeathTime = rl.GetTime()
	}
}
