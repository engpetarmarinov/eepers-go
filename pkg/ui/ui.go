package ui

import (
	"github.com/engpetarmarinov/eepers-go/pkg/game"
	"github.com/engpetarmarinov/eepers-go/pkg/palette"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const cellSize float32 = 50.0

// DrawUI draws the game's UI with visual inventory display.
func DrawUI(gs *game.State, screenWidth int32) {
	// Draw keys as circles (like in original ADA version)
	for i := 0; i < gs.Player.Keys; i++ {
		position := rl.NewVector2(100.0+float32(i)*cellSize, 100.0)
		rl.DrawCircleV(position, cellSize*0.25, palette.Colors["COLOR_DOORKEY"])
	}

	// Draw bombs as circles - only show available bombs
	var padding float32 = cellSize * 0.5
	for i := 0; i < gs.Player.Bombs; i++ {
		position := rl.NewVector2(100.0+float32(i)*(cellSize+padding), 200.0)
		rl.DrawCircleV(position, cellSize*0.5, palette.Colors["COLOR_BOMB"])
	}

	// Draw health bar
	healthBarWidth := int32(200)
	healthBarHeight := int32(20)
	healthBarX := screenWidth - healthBarWidth - 10
	healthBarY := int32(10)
	rl.DrawRectangle(healthBarX, healthBarY, healthBarWidth, healthBarHeight, rl.Gray)
	currentHealthWidth := int32(float32(healthBarWidth) * gs.Player.Health)
	rl.DrawRectangle(healthBarX, healthBarY, currentHealthWidth, healthBarHeight, palette.Colors["COLOR_HEALTHBAR"])
}
