package game

import "github.com/engpetarmarinov/eepers-go/pkg/world"

// Directions defines the four cardinal directions for movement and explosions.
var Directions = [4]world.IVector2{
	{X: 0, Y: -1}, // Up
	{X: 0, Y: 1},  // Down
	{X: -1, Y: 0}, // Left
	{X: 1, Y: 0},  // Right
}
