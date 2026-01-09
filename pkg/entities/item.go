package entities

import "github.com/engpetarmarinov/eepers-go/pkg/world"

// ItemKind represents the type of item.
type ItemKind int

const (
	ItemNone ItemKind = iota
	ItemKey
	ItemBombRefill
	ItemCheckpoint
	ItemBombSlot
)

// Item represents an item in the game.
type Item struct {
	Kind     ItemKind
	Position world.IVector2
	Cooldown int
}
