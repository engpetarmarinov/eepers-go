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
	Portals            []entities.PortalState
	TurnAnimation      float32
	Camera             rl.Camera2D
	Tutorial           TutorialState
	Menu               MenuState
	ShouldQuit         bool
	DurationOfLastTurn float64
	Checkpoint         CheckpointState
	WorldConfig        WorldConfig // Configuration for all worlds and levels
	InHub              bool        // Whether player is currently in a hub level
	CurrentLevelPath   string      // Path to the currently loaded level
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

// LoadLevel loads a specific level by path
func (gs *State) LoadLevel(levelPath string, isHub bool) error {
	if levelPath == "" {
		return nil // Invalid level path
	}

	// Clear all dynamic game state
	gs.Bombs = nil
	gs.Explosions = nil
	gs.Eepers = nil
	gs.Items = nil
	gs.Portals = nil
	gs.TurnAnimation = 0

	// Load the level
	err := LoadGameFromImage(levelPath, gs, true)
	if err != nil {
		return err
	}

	// Set current level info
	gs.CurrentLevelPath = levelPath
	gs.InHub = isHub

	// Reset player state
	gs.Player.Health = 1.0
	gs.Player.Dead = false
	gs.Player.ReachedFather = false
	gs.Player.VictoryTime = 0
	gs.Player.BombSlots = 1 // Player starts with 1 bomb slot
	gs.Player.Bombs = 0     // Player starts with no bombs
	gs.Player.Keys = 0      // Player starts with no keys

	// Reset camera
	gs.Camera.Zoom = 1.0

	// Save checkpoint for this level
	gs.SaveCheckpoint()

	return nil
}

// LoadHub loads the hub level for the current world
func (gs *State) LoadHub() error {
	hubPath := gs.WorldConfig.GetCurrentHub()
	return gs.LoadLevel(hubPath, true)
}

// LoadNextLevel loads the next level, returns true if there is a next level
func (gs *State) LoadNextLevel() (bool, error) {
	// If we're in a regular level (not hub), return to hub
	if !gs.InHub {
		err := gs.LoadHub()
		if err != nil {
			return false, err
		}
		return true, nil
	}

	// If we're in the hub and reached Father, we've completed all levels in this world
	// Try to advance to the next world
	if gs.WorldConfig.NextWorld() {
		err := gs.LoadHub()
		if err != nil {
			return false, err
		}
		return true, nil
	}

	// No more worlds
	return false, nil
}

// RestartFromFirstLevel restarts the game from the first world's hub
func (gs *State) RestartFromFirstLevel() error {
	// Reset tutorial
	gs.Tutorial = TutorialState{
		Phase: TutorialMove,
	}

	// Reset to first world
	gs.WorldConfig.CurrentWorld = 0

	// Load first world's hub
	return gs.LoadHub()
}

// LoadLevelFromPortal loads a level based on portal number in the current world
func (gs *State) LoadLevelFromPortal(portalNumber int) error {
	levelPath := gs.WorldConfig.GetLevel(portalNumber)
	if levelPath == "" {
		return nil // Invalid portal number
	}

	return gs.LoadLevel(levelPath, false)
}
