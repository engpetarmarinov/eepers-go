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
