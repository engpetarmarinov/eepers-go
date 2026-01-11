package game

import (
	"github.com/engpetarmarinov/eepers-go/pkg/audio"
	"github.com/engpetarmarinov/eepers-go/pkg/entities"
	"github.com/engpetarmarinov/eepers-go/pkg/world"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	PortalOpenDistance = 9.0  // Distance squared (3 cells = 9 when squared)
	PortalOpenSpeed    = 0.15 // Fast opening animation
	PortalCloseSpeed   = 0.15 // Fast closing animation
)

// SpawnPortal creates a portal at the specified center position
func (gs *State) SpawnPortal(centerPos world.IVector2, portalID int) {
	portal := entities.NewPortal(portalID, centerPos)
	gs.Portals = append(gs.Portals, portal)
}

// UpdatePortals updates all portal animations based on player proximity
func (gs *State) UpdatePortals() {
	for i := range gs.Portals {
		portal := &gs.Portals[i]

		// Store previous open progress to detect state changes
		prevProgress := portal.OpenProgress

		// Calculate distance to player
		distance := portal.DistanceToPlayer(gs.Player.Position)

		// Open portal if player is close (within 3 cells)
		if distance < PortalOpenDistance {
			portal.OpenProgress += PortalOpenSpeed
			if portal.OpenProgress > 1.0 {
				portal.OpenProgress = 1.0
			}

			// Play sound when starting to open (transition from closed to opening)
			if prevProgress <= 0.0 && portal.OpenProgress > 0.0 {
				rl.PlaySound(audio.OpenPortalSound)
			}
		} else {
			// Close portal if player is far
			portal.OpenProgress -= PortalCloseSpeed
			if portal.OpenProgress < 0.0 {
				portal.OpenProgress = 0.0
			}

			// Play sound when starting to close (transition from open to closing)
			if prevProgress >= 1.0 && portal.OpenProgress < 1.0 {
				rl.PlaySound(audio.OpenPortalSound)
			}
		}
	}
}

// GetPortalAtPosition returns the portal at the given position, or nil if none exists
func (gs *State) GetPortalAtPosition(pos world.IVector2) *entities.PortalState {
	for i := range gs.Portals {
		if gs.Portals[i].ContainsPosition(pos) {
			return &gs.Portals[i]
		}
	}
	return nil
}

// ActivatePortal triggers the portal and loads the corresponding level
func (gs *State) ActivatePortal(portal *entities.PortalState) error {
	if portal == nil || portal.IsActivated {
		return nil
	}

	portal.IsActivated = true
	return gs.LoadLevelFromPortal(portal.ID)
}
