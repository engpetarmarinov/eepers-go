package ui

import (
	"fmt"
	"math"

	"github.com/engpetarmarinov/eepers-go/pkg/entities"
	"github.com/engpetarmarinov/eepers-go/pkg/palette"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// DrawPlayerEyes draws the player's eyes with direction offset.
func DrawPlayerEyes(player entities.PlayerState, interpPos rl.Vector2) {
	eyeColor := palette.Colors["COLOR_EYES"]
	size := rl.NewVector2(12, 28)
	eyeSize := rl.NewVector2(size.X, size.Y)

	// Calculate eye offset based on looking direction
	// Eyes target is where the player is looking
	lookDir := rl.NewVector2(
		float32(player.EyesTarget.X-player.Position.X),
		float32(player.EyesTarget.Y-player.Position.Y),
	)

	// Normalize and scale the look direction for eye offset
	length := float32(math.Sqrt(float64(lookDir.X*lookDir.X + lookDir.Y*lookDir.Y)))
	eyeOffset := rl.NewVector2(0, 0)
	if length > 0.01 {
		eyeOffset = rl.NewVector2(
			(lookDir.X/length)*3.0, // Small offset in X direction
			(lookDir.Y/length)*3.0, // Small offset in Y direction
		)
	}

	leftEyeOffset := rl.NewVector2(size.X+4, size.Y*0.8)
	rightEyeOffset := rl.NewVector2(size.X*3, size.Y*0.8)

	leftEyePos := rl.Vector2Add(rl.Vector2Add(interpPos, leftEyeOffset), eyeOffset)
	rightEyePos := rl.Vector2Add(rl.Vector2Add(interpPos, rightEyeOffset), eyeOffset)

	drawEye(leftEyePos, eyeSize, entities.EyesMeshes[player.Eyes][0], eyeColor)
	drawEye(rightEyePos, eyeSize, entities.EyesMeshes[player.Eyes][1], eyeColor)
}

func drawEye(position, size rl.Vector2, mesh entities.EyeMesh, color rl.Color) {
	transform := func(v rl.Vector2) rl.Vector2 {
		vCentered := rl.Vector2Subtract(v, rl.NewVector2(0.5, 0.5))
		vScaled := rl.NewVector2(vCentered.X*size.X, vCentered.Y*size.Y)
		return rl.Vector2Add(position, vScaled)
	}

	points := []rl.Vector2{
		transform(mesh[0]),
		transform(mesh[1]),
		transform(mesh[3]),
		transform(mesh[2]),
	}

	rl.DrawTriangleFan(points, color)
}

// UpdatePlayerEyes updates the player's eye state.
func UpdatePlayerEyes(player *entities.PlayerState) {
	player.PrevEyes = player.Eyes

	if player.Dead {
		player.Eyes = entities.EyesCringe
		return
	}

	// Simple logic for now, can be expanded
	if math.Abs(float64(player.Position.X-player.PrevPosition.X)) > 0 || math.Abs(float64(player.Position.Y-player.PrevPosition.Y)) > 0 {
		player.Eyes = entities.EyesSurprised
	} else {
		player.Eyes = entities.EyesOpen
	}
}

// DrawEeperEyes draws an eeper's eyes with direction offset and interpolation.
func DrawEeperEyes(eeper entities.EeperState, interpPos rl.Vector2, turnAnimation float32) {
	eyeColor := palette.Colors["COLOR_EYES"]

	// Eeper eyes are larger and positioned differently than player eyes
	size := rl.NewVector2(18, 40) // Larger eyes for eepers
	eyeSize := rl.NewVector2(size.X, size.Y)

	// Calculate eye offset based on looking direction
	lookDir := rl.NewVector2(
		float32(eeper.EyesTarget.X-eeper.Position.X),
		float32(eeper.EyesTarget.Y-eeper.Position.Y),
	)

	// Normalize and scale the look direction for eye offset
	length := float32(math.Sqrt(float64(lookDir.X*lookDir.X + lookDir.Y*lookDir.Y)))
	eyeOffset := rl.NewVector2(0, 0)
	if length > 0.01 {
		eyeOffset = rl.NewVector2(
			(lookDir.X/length)*4.0, // Eye movement offset
			(lookDir.Y/length)*4.0,
		)
	}

	// Position eyes based on eeper size - centered horizontally, upper portion vertically
	eeperCenterX := float32(eeper.Size.X) * 25.0 // Center of eeper
	eeperEyeY := float32(eeper.Size.Y) * 15.0    // Upper portion for eyes

	eyeSpacing := float32(20.0) // Space between eyes
	leftEyeOffset := rl.NewVector2(eeperCenterX-eyeSpacing, eeperEyeY)
	rightEyeOffset := rl.NewVector2(eeperCenterX+eyeSpacing, eeperEyeY)

	leftEyePos := rl.Vector2Add(rl.Vector2Add(interpPos, leftEyeOffset), eyeOffset)
	rightEyePos := rl.Vector2Add(rl.Vector2Add(interpPos, rightEyeOffset), eyeOffset)

	drawEye(leftEyePos, eyeSize, entities.EyesMeshes[eeper.Eyes][0], eyeColor)
	drawEye(rightEyePos, eyeSize, entities.EyesMeshes[eeper.Eyes][1], eyeColor)
}

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
