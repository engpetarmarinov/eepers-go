package entities

import "github.com/engpetarmarinov/eepers-go/pkg/world"

// BombState represents the state of a bomb.
type BombState struct {
	Position  world.IVector2
	Countdown int
}
