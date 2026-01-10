package entities

import "github.com/engpetarmarinov/eepers-go/pkg/world"

// PlayerState represents the state of the player.
type PlayerState struct {
	PrevPosition  world.IVector2
	Position      world.IVector2
	PrevEyes      EyesKind
	Eyes          EyesKind
	EyesTarget    world.IVector2
	Keys          int
	Bombs         int
	BombSlots     int
	Health        float32
	Dead          bool
	DeathTime     float64
	ReachedFather bool // Victory condition - player reached Father
}
