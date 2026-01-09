package entities

import "github.com/engpetarmarinov/eepers-go/pkg/world"

// ExplosionState represents the state of an explosion.
type ExplosionState struct {
	Position     world.IVector2
	Timer        int
	InitialTimer int
}
