package ui

import (
	"fmt"

	"github.com/engpetarmarinov/eepers-go/pkg/entities"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// DrawPortal renders a portal with opening animation
func DrawPortal(portal entities.PortalState) {
	// Portal door colors based on ID
	doorColors := map[int]rl.Color{
		1: rl.NewColor(40, 80, 40, 255), // Dark green door
		2: rl.NewColor(40, 40, 80, 255), // Dark blue door
		3: rl.NewColor(80, 80, 40, 255), // Dark yellow door
		4: rl.NewColor(80, 40, 40, 255), // Dark red door
	}

	doorColor, exists := doorColors[portal.ID]
	if !exists {
		doorColor = rl.NewColor(60, 60, 60, 255)
	}

	// Calculate the full 3x3 portal area (150x150 pixels)
	// Portal center is at portal.CenterPos, so the 3x3 grid spans from -1 to +1 in each direction
	portalX := int32((portal.CenterPos.X - 1) * 50)
	portalY := int32((portal.CenterPos.Y - 1) * 50)
	portalWidth := int32(150)  // 3 cells * 50 pixels
	portalHeight := int32(150) // 3 cells * 50 pixels

	// Draw black hole background for entire 3x3 area
	rl.DrawRectangle(portalX, portalY, portalWidth, portalHeight, rl.Black)

	// Calculate door height based on open progress
	// Door slides UP, revealing black hole from BOTTOM to TOP
	// Progress 0.0 = closed (door covers all), 1.0 = open (door gone, lifted up)
	doorHeight := int32(float32(portalHeight) * (1.0 - portal.OpenProgress))

	if doorHeight > 0 {
		// Keep the door's TOP edge fixed at portalY
		// The door shrinks upward from the bottom, revealing the black hole from bottom to top
		doorY := portalY

		// Draw the entire door as one solid piece
		rl.DrawRectangle(portalX, doorY, portalWidth, doorHeight, doorColor)

		// Add horizontal line details across the entire door
		if doorHeight > 20 {
			lineColor := rl.ColorBrightness(doorColor, 0.3)
			// Top line
			rl.DrawLine(portalX+10, doorY+10, portalX+portalWidth-10, doorY+10, lineColor)
			// Bottom line
			if doorHeight > 40 {
				rl.DrawLine(portalX+10, doorY+doorHeight-10, portalX+portalWidth-10, doorY+doorHeight-10, lineColor)
			}
			// Middle line
			if doorHeight > 80 {
				rl.DrawLine(portalX+10, doorY+doorHeight/2, portalX+portalWidth-10, doorY+doorHeight/2, lineColor)
			}
		}

		// Draw door border
		rl.DrawRectangleLines(portalX, doorY, portalWidth, doorHeight, rl.ColorBrightness(doorColor, -0.3))
	}

	// Draw outer frame for entire 3x3 portal
	frameColor := rl.NewColor(100, 100, 100, 255)
	rl.DrawRectangleLines(portalX, portalY, portalWidth, portalHeight, frameColor)

	// Draw portal number in the center of the black hole (visible when door opens)
	if portal.OpenProgress > 0.1 {
		centerX := portalX + portalWidth/2
		centerY := portalY + portalHeight/2

		portalText := fmt.Sprintf("%d", portal.ID)
		fontSize := int32(80)
		textWidth := rl.MeasureText(portalText, fontSize)

		// Text fades in as door opens
		textAlpha := uint8(float32(255) * portal.OpenProgress)

		// Draw text glow effect
		glowColor := rl.NewColor(255, 255, 255, textAlpha/4)
		for dx := int32(-3); dx <= 3; dx += 3 {
			for dy := int32(-3); dy <= 3; dy += 3 {
				if dx != 0 || dy != 0 {
					rl.DrawText(portalText, centerX-textWidth/2+dx, centerY-fontSize/2+dy, fontSize, glowColor)
				}
			}
		}

		// Draw main text in white
		textColor := rl.NewColor(255, 255, 255, textAlpha)
		rl.DrawText(portalText, centerX-textWidth/2, centerY-fontSize/2, fontSize, textColor)
	}
}
