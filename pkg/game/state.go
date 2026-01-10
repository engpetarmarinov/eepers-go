package game

import (
	"github.com/engpetarmarinov/eepers-go/pkg/entities"
	"github.com/engpetarmarinov/eepers-go/pkg/world"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// State represents the entire state of the game.
type State struct {
	Map                [][]world.Cell
	Player             entities.PlayerState
	Items              []entities.Item
	Eepers             []entities.EeperState
	Bombs              []entities.BombState
	Explosions         []entities.ExplosionState
	TurnAnimation      float32
	Camera             rl.Camera2D
	Tutorial           TutorialState
	Menu               MenuState
	ShouldQuit         bool
	DurationOfLastTurn float64
}

// AllocateItem adds a new item to the game state at the specified position.
func (gs *State) AllocateItem(position world.IVector2, kind entities.ItemKind) {
	gs.Items = append(gs.Items, entities.Item{
		Position: position,
		Kind:     kind,
		Cooldown: 0,
	})
}
