package ui

import (
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

// DrawEeperEyes draws an eeper's eyes with direction offset and interpolation.
func DrawEeperEyes(eeper entities.EeperState, interpPos rl.Vector2, turnAnimation float32) {
	eyeColor := palette.Colors["COLOR_EYES"]

	// Calculate rendered size (accounts for gnome 70% scaling)
	var eeperRenderSize rl.Vector2
	if eeper.Kind == entities.EeperGnome {
		// Gnomes are rendered at 70% size
		gnomeRatio := float32(0.7)
		eeperRenderSize = rl.NewVector2(float32(eeper.Size.X)*50.0*gnomeRatio, float32(eeper.Size.Y)*50.0*gnomeRatio)
	} else {
		// All other eepers rendered at full size
		eeperRenderSize = rl.NewVector2(float32(eeper.Size.X)*50.0, float32(eeper.Size.Y)*50.0)
	}

	// Eyes size as ratio of eeper size
	eyeRatioWidth := float32(13.0 / 64.0)
	eyeRatioHeight := float32(23.0 / 64.0)
	eyeSize := rl.NewVector2(eeperRenderSize.X*eyeRatioWidth, eeperRenderSize.Y*eyeRatioHeight)

	// Calculate eye offset based on looking direction
	lookDir := rl.NewVector2(
		float32(eeper.EyesTarget.X-eeper.Position.X),
		float32(eeper.EyesTarget.Y-eeper.Position.Y),
	)

	// Normalize and scale the look direction for eye offset
	length := float32(math.Sqrt(float64(lookDir.X*lookDir.X + lookDir.Y*lookDir.Y)))
	eyeOffset := rl.NewVector2(0, 0)
	if length > 0.01 {
		eyeOffsetScale := eyeSize.X * 0.6
		eyeOffset = rl.NewVector2(
			(lookDir.X/length)*eyeOffsetScale,
			(lookDir.Y/length)*eyeOffsetScale,
		)
	}

	// Position eyes at center of eeper
	center := rl.NewVector2(interpPos.X+eeperRenderSize.X*0.5, interpPos.Y+eeperRenderSize.Y*0.5)
	eyePosition := rl.Vector2Add(center, eyeOffset)

	// Calculate left and right eye positions
	leftEyePos := rl.NewVector2(
		eyePosition.X-eyeSize.X*1.0,
		eyePosition.Y+eyeSize.Y*0.2,
	)
	rightEyePos := rl.NewVector2(
		eyePosition.X+eyeSize.X*1.0,
		eyePosition.Y+eyeSize.Y*0.2,
	)

	drawEye(leftEyePos, eyeSize, entities.EyesMeshes[eeper.Eyes][0], eyeColor)
	drawEye(rightEyePos, eyeSize, entities.EyesMeshes[eeper.Eyes][1], eyeColor)
}

// drawEye draws a single eye using the provided mesh
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
