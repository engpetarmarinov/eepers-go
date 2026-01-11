package game

// World represents a collection of levels
type World struct {
	Name     string   // Display name of the world
	HubLevel string   // Path to the hub/gallery level for this world
	Levels   []string // Paths to all levels in this world (portal 1 -> Levels[0], etc.)
}

// WorldConfig holds all worlds in the game
type WorldConfig struct {
	Worlds       []World
	CurrentWorld int // Index of the current world (0-based)
}

// GetCurrentHub returns the hub level path for the current world
func (wc *WorldConfig) GetCurrentHub() string {
	if wc.CurrentWorld < 0 || wc.CurrentWorld >= len(wc.Worlds) {
		return ""
	}
	return wc.Worlds[wc.CurrentWorld].HubLevel
}

// GetLevel returns the level path for a given portal number in the current world
func (wc *WorldConfig) GetLevel(portalNumber int) string {
	if wc.CurrentWorld < 0 || wc.CurrentWorld >= len(wc.Worlds) {
		return ""
	}

	world := wc.Worlds[wc.CurrentWorld]
	levelIndex := portalNumber - 1 // Portal 1 -> index 0, Portal 2 -> index 1, etc.

	if levelIndex < 0 || levelIndex >= len(world.Levels) {
		return ""
	}

	return world.Levels[levelIndex]
}

// HasLevel checks if a portal number exists in the current world
func (wc *WorldConfig) HasLevel(portalNumber int) bool {
	if wc.CurrentWorld < 0 || wc.CurrentWorld >= len(wc.Worlds) {
		return false
	}

	world := wc.Worlds[wc.CurrentWorld]
	levelIndex := portalNumber - 1

	return levelIndex >= 0 && levelIndex < len(world.Levels)
}

// GetTotalLevelsInCurrentWorld returns the number of levels in the current world
func (wc *WorldConfig) GetTotalLevelsInCurrentWorld() int {
	if wc.CurrentWorld < 0 || wc.CurrentWorld >= len(wc.Worlds) {
		return 0
	}
	return len(wc.Worlds[wc.CurrentWorld].Levels)
}

// GetCurrentWorldName returns the name of the current world
func (wc *WorldConfig) GetCurrentWorldName() string {
	if wc.CurrentWorld < 0 || wc.CurrentWorld >= len(wc.Worlds) {
		return ""
	}
	return wc.Worlds[wc.CurrentWorld].Name
}

// NextWorld moves to the next world, returns true if successful
func (wc *WorldConfig) NextWorld() bool {
	if wc.CurrentWorld+1 < len(wc.Worlds) {
		wc.CurrentWorld++
		return true
	}
	return false
}
