package ui

import (
	"fmt"

	"github.com/engpetarmarinov/eepers-go/pkg/game"
	"github.com/engpetarmarinov/eepers-go/pkg/palette"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// DrawUI draws the game's UI.
func DrawUI(gs *game.State, screenWidth int32) {
	// Draw bomb counter
	bombText := fmt.Sprintf("Bombs: %d/%d", gs.Player.Bombs, gs.Player.BombSlots)
	rl.DrawText(bombText, 10, 10, 20, palette.Colors["COLOR_LABEL"])

	// Draw key counter
	keyText := fmt.Sprintf("Keys: %d", gs.Player.Keys)
	rl.DrawText(keyText, 10, 40, 20, palette.Colors["COLOR_LABEL"])

	// Draw health bar
	healthBarWidth := int32(200)
	healthBarHeight := int32(20)
	healthBarX := screenWidth - healthBarWidth - 10
	healthBarY := int32(10)
	rl.DrawRectangle(healthBarX, healthBarY, healthBarWidth, healthBarHeight, rl.Gray)
	currentHealthWidth := int32(float32(healthBarWidth) * gs.Player.Health)
	rl.DrawRectangle(healthBarX, healthBarY, currentHealthWidth, healthBarHeight, palette.Colors["COLOR_HEALTHBAR"])
}
