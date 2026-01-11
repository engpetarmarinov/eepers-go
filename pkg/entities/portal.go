package entities

import "github.com/engpetarmarinov/eepers-go/pkg/world"

// PortalState represents the state of a portal
type PortalState struct {
	ID           int              // Portal ID (1-4)
	CenterPos    world.IVector2   // Center position of the 3x3 portal
	Cells        []world.IVector2 // All 9 cells that make up the portal
	OpenProgress float32          // 0.0 = closed, 1.0 = fully open
	IsActivated  bool             // Whether the portal has been entered
}

// NewPortal creates a new portal state
func NewPortal(id int, centerPos world.IVector2) PortalState {
	cells := make([]world.IVector2, 0, 9)
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			cells = append(cells, world.IVector2{
				X: centerPos.X + dx,
				Y: centerPos.Y + dy,
			})
		}
	}

	return PortalState{
		ID:           id,
		CenterPos:    centerPos,
		Cells:        cells,
		OpenProgress: 0.0,
		IsActivated:  false,
	}
}

// ContainsPosition checks if a position is within this portal's cells
func (p *PortalState) ContainsPosition(pos world.IVector2) bool {
	for _, cell := range p.Cells {
		if cell.X == pos.X && cell.Y == pos.Y {
			return true
		}
	}
	return false
}

// DistanceToPlayer calculates the distance from the portal center to a position
func (p *PortalState) DistanceToPlayer(playerPos world.IVector2) float32 {
	dx := float32(p.CenterPos.X - playerPos.X)
	dy := float32(p.CenterPos.Y - playerPos.Y)
	return dx*dx + dy*dy // Squared distance (no need for sqrt for comparison)
}
