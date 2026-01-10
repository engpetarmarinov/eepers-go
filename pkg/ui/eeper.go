package ui

import (
	"fmt"

	"github.com/engpetarmarinov/eepers-go/pkg/entities"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// DrawEeperHealthBar draws a health bar for an eeper
func DrawEeperHealthBar(eeper entities.EeperState, interpPos rl.Vector2, size rl.Vector2) {
	if eeper.Health >= 1.0 {
		return // Don't draw health bar when at full health
	}

	barWidth := size.X
	barHeight := float32(5.0)
	barPos := rl.NewVector2(interpPos.X, interpPos.Y-barHeight-2)

	// Background (dark)
	rl.DrawRectangleV(barPos, rl.NewVector2(barWidth, barHeight), rl.NewColor(40, 40, 40, 200))

	// Health fill (green to red based on health)
	healthColor := rl.NewColor(
		uint8(255*(1.0-eeper.Health)),
		uint8(255*eeper.Health),
		0,
		255,
	)
	rl.DrawRectangleV(barPos, rl.NewVector2(barWidth*eeper.Health, barHeight), healthColor)
}

// DrawEeperCooldownBubble draws a countdown bubble above an eeper
func DrawEeperCooldownBubble(eeper entities.EeperState, interpPos rl.Vector2, size rl.Vector2, backgroundColor rl.Color) {
	bubbleRadius := float32(30.0)
	bubbleCenter := rl.NewVector2(
		interpPos.X+size.X*0.5,
		interpPos.Y-bubbleRadius*2.0,
	)

	// Draw bubble circle
	rl.DrawCircleV(bubbleCenter, bubbleRadius, backgroundColor)

	// Draw countdown number
	countdownText := fmt.Sprintf("%d", eeper.AttackCooldown)
	fontSize := int32(40)
	textWidth := rl.MeasureText(countdownText, fontSize)
	textPos := rl.NewVector2(
		bubbleCenter.X-float32(textWidth)/2,
		bubbleCenter.Y-float32(fontSize)/2,
	)

	rl.DrawText(countdownText, int32(textPos.X), int32(textPos.Y), fontSize, rl.Black)
}
