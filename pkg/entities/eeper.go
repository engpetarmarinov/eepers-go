package entities

import "github.com/engpetarmarinov/eepers-go/pkg/world"

// EeperKind represents the type of eeper.
type EeperKind int

const (
	EeperGuard EeperKind = iota
	EeperMother
	EeperGnome
	EeperFather
)

// EeperState represents the state of an eeper.
type EeperState struct {
	Kind           EeperKind
	Dead           bool
	Position       world.IVector2
	PrevPosition   world.IVector2
	EyesAngle      float32
	EyesTarget     world.IVector2
	PrevEyes       EyesKind
	Eyes           EyesKind
	Size           world.IVector2
	Path           [][]int
	Damaged        bool
	Health         float32
	AttackCooldown int
}
