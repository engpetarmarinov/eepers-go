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
	Checkpoint         CheckpointState
}

// CheckpointState stores a snapshot of the game state for respawning
type CheckpointState struct {
	Map             [][]world.Cell
	PlayerPosition  world.IVector2
	PlayerKeys      int
	PlayerBombs     int
	PlayerBombSlots int
	Eepers          []entities.EeperState
	Items           []entities.Item
	Bombs           []entities.BombState
}

// AllocateItem adds a new item to the game state at the specified position.
func (gs *State) AllocateItem(position world.IVector2, kind entities.ItemKind) {
	gs.Items = append(gs.Items, entities.Item{
		Position: position,
		Kind:     kind,
		Cooldown: 0,
	})
}

// SaveCheckpoint saves the current game state to checkpoint
func (gs *State) SaveCheckpoint() {
	// Clone the map
	gs.Checkpoint.Map = make([][]world.Cell, len(gs.Map))
	for i := range gs.Map {
		gs.Checkpoint.Map[i] = make([]world.Cell, len(gs.Map[i]))
		copy(gs.Checkpoint.Map[i], gs.Map[i])
	}

	// Save player state
	gs.Checkpoint.PlayerPosition = gs.Player.Position
	gs.Checkpoint.PlayerKeys = gs.Player.Keys
	gs.Checkpoint.PlayerBombs = gs.Player.Bombs
	gs.Checkpoint.PlayerBombSlots = gs.Player.BombSlots

	// Clone eepers
	gs.Checkpoint.Eepers = make([]entities.EeperState, len(gs.Eepers))
	for i := range gs.Eepers {
		gs.Checkpoint.Eepers[i] = gs.Eepers[i]
		// Deep copy the path map
		if gs.Eepers[i].Path != nil {
			gs.Checkpoint.Eepers[i].Path = make([][]int, len(gs.Eepers[i].Path))
			for j := range gs.Eepers[i].Path {
				gs.Checkpoint.Eepers[i].Path[j] = make([]int, len(gs.Eepers[i].Path[j]))
				copy(gs.Checkpoint.Eepers[i].Path[j], gs.Eepers[i].Path[j])
			}
		}
	}

	// Clone items
	gs.Checkpoint.Items = make([]entities.Item, len(gs.Items))
	copy(gs.Checkpoint.Items, gs.Items)

	// Clone bombs
	gs.Checkpoint.Bombs = make([]entities.BombState, len(gs.Bombs))
	copy(gs.Checkpoint.Bombs, gs.Bombs)
}

// RestoreCheckpoint restores the game state from checkpoint
func (gs *State) RestoreCheckpoint() {
	// Restore the map
	gs.Map = make([][]world.Cell, len(gs.Checkpoint.Map))
	for i := range gs.Checkpoint.Map {
		gs.Map[i] = make([]world.Cell, len(gs.Checkpoint.Map[i]))
		copy(gs.Map[i], gs.Checkpoint.Map[i])
	}

	// Restore player state
	gs.Player.Position = gs.Checkpoint.PlayerPosition
	gs.Player.PrevPosition = gs.Checkpoint.PlayerPosition // Set PrevPosition to avoid interpolation issues
	gs.Player.Keys = gs.Checkpoint.PlayerKeys
	gs.Player.Bombs = gs.Checkpoint.PlayerBombs
	gs.Player.BombSlots = gs.Checkpoint.PlayerBombSlots
	gs.Player.Dead = false
	gs.Player.Health = 1.0

	// Restore eepers
	gs.Eepers = make([]entities.EeperState, len(gs.Checkpoint.Eepers))
	for i := range gs.Checkpoint.Eepers {
		gs.Eepers[i] = gs.Checkpoint.Eepers[i]
		// Deep copy the path map
		if gs.Checkpoint.Eepers[i].Path != nil {
			gs.Eepers[i].Path = make([][]int, len(gs.Checkpoint.Eepers[i].Path))
			for j := range gs.Checkpoint.Eepers[i].Path {
				gs.Eepers[i].Path[j] = make([]int, len(gs.Checkpoint.Eepers[i].Path[j]))
				copy(gs.Eepers[i].Path[j], gs.Checkpoint.Eepers[i].Path[j])
			}
		}
	}

	// Restore items
	gs.Items = make([]entities.Item, len(gs.Checkpoint.Items))
	copy(gs.Items, gs.Checkpoint.Items)

	// Restore bombs
	gs.Bombs = make([]entities.BombState, len(gs.Checkpoint.Bombs))
	copy(gs.Bombs, gs.Checkpoint.Bombs)

	// Clear explosions
	gs.Explosions = nil

	// Reset turn animation to prevent visual glitches
	gs.TurnAnimation = 0
}
