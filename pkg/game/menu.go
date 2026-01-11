package game

import (
	"github.com/engpetarmarinov/eepers-go/pkg/palette"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// MenuOption represents a menu option
type MenuOption int

const (
	MenuContinue MenuOption = iota
	MenuExitLevel
	MenuRestart
	MenuQuit
)

// MenuState represents the state of the pause menu
type MenuState struct {
	IsOpen         bool
	SelectedOption MenuOption
	TotalOptions   int
}

// NewMenuState creates a new menu state
func NewMenuState() MenuState {
	return MenuState{
		IsOpen:         false,
		SelectedOption: MenuContinue,
		TotalOptions:   4, // Continue, Exit Level, Restart, Quit
	}
}

// ToggleMenu toggles the menu open/closed
func (ms *MenuState) ToggleMenu() {
	ms.IsOpen = !ms.IsOpen
	if ms.IsOpen {
		// Reset to first option when opening
		ms.SelectedOption = MenuContinue
	}
}

// OpenMenu opens the menu
func (ms *MenuState) OpenMenu() {
	ms.IsOpen = true
	ms.SelectedOption = MenuContinue
}

// CloseMenu closes the menu
func (ms *MenuState) CloseMenu() {
	ms.IsOpen = false
}

// MoveUp moves selection up
func (ms *MenuState) MoveUp() {
	if ms.SelectedOption > MenuContinue {
		ms.SelectedOption--
	} else {
		// Wrap to bottom
		ms.SelectedOption = MenuQuit
	}
}

// MoveUpInHub moves selection up, skipping Exit Level
func (ms *MenuState) MoveUpInHub() {
	if ms.SelectedOption > MenuContinue {
		ms.SelectedOption--
		// Skip Exit Level in hub
		if ms.SelectedOption == MenuExitLevel {
			ms.SelectedOption--
		}
	} else {
		// Wrap to bottom (Quit)
		ms.SelectedOption = MenuQuit
	}
}

// MoveDown moves selection down
func (ms *MenuState) MoveDown() {
	if ms.SelectedOption < MenuQuit {
		ms.SelectedOption++
	} else {
		// Wrap to top
		ms.SelectedOption = MenuContinue
	}
}

// MoveDownInHub moves selection down, skipping Exit Level
func (ms *MenuState) MoveDownInHub() {
	if ms.SelectedOption < MenuQuit {
		ms.SelectedOption++
		// Skip Exit Level in hub
		if ms.SelectedOption == MenuExitLevel {
			ms.SelectedOption++
		}
	} else {
		// Wrap to top
		ms.SelectedOption = MenuContinue
	}
}

// GetOptionText returns the text for a menu option
func GetOptionText(option MenuOption) string {
	switch option {
	case MenuContinue:
		return "Continue"
	case MenuRestart:
		return "Restart"
	case MenuExitLevel:
		return "Exit Level"
	case MenuQuit:
		return "Quit"
	default:
		return ""
	}
}

// DrawMenu draws the pause menu in the center of the screen
func (ms *MenuState) DrawMenu(inHub bool) {
	if !ms.IsOpen {
		return
	}

	currentWidth := int32(rl.GetScreenWidth())
	currentHeight := int32(rl.GetScreenHeight())

	// Get render dimensions (which might differ from screen dimensions on HiDPI/Retina)
	renderWidth := int32(rl.GetRenderWidth())
	renderHeight := int32(rl.GetRenderHeight())

	// Use the larger of the two to ensure full coverage
	overlayWidth := currentWidth
	overlayHeight := currentHeight
	if renderWidth > currentWidth {
		overlayWidth = renderWidth
	}
	if renderHeight > currentHeight {
		overlayHeight = renderHeight
	}

	// Semi-transparent overlay covering the entire screen
	// Use DrawRectangleRec for precise full-screen coverage
	fullScreenRect := rl.NewRectangle(0, 0, float32(overlayWidth), float32(overlayHeight))
	rl.DrawRectangleRec(fullScreenRect, rl.Fade(rl.Black, 0.7))

	// Menu box dimensions - use current screen dimensions for centering
	menuWidth := int32(400)
	menuHeight := int32(300)
	menuX := (currentWidth - menuWidth) / 2
	menuY := (currentHeight - menuHeight) / 2

	// Draw menu background
	rl.DrawRectangle(menuX, menuY, menuWidth, menuHeight, palette.Colors["COLOR_BACKGROUND"])

	// Draw menu border
	borderColor := palette.Colors["COLOR_LABEL"]
	rl.DrawRectangleLines(menuX, menuY, menuWidth, menuHeight, borderColor)
	rl.DrawRectangleLines(menuX+1, menuY+1, menuWidth-2, menuHeight-2, borderColor)

	// Draw title
	title := "PAUSED"
	titleSize := int32(40)
	titleWidth := rl.MeasureText(title, titleSize)
	titleX := menuX + (menuWidth-titleWidth)/2
	titleY := menuY + 30
	rl.DrawText(title, titleX, titleY, titleSize, palette.Colors["COLOR_LABEL"])

	// Draw menu options
	optionSize := int32(30)
	optionY := menuY + 100
	optionSpacing := int32(50)

	optionIndex := int32(0)
	for i := MenuOption(0); i <= MenuQuit; i++ {
		// Skip "Exit Level" if we're already in the hub
		if i == MenuExitLevel && inHub {
			continue
		}

		optionText := GetOptionText(i)
		textWidth := rl.MeasureText(optionText, optionSize)
		textX := menuX + (menuWidth-textWidth)/2
		textY := optionY + optionIndex*optionSpacing

		// Highlight selected option
		if i == ms.SelectedOption {
			// Draw selection background
			highlightPadding := int32(10)
			highlightWidth := textWidth + highlightPadding*2
			highlightHeight := optionSize + highlightPadding
			highlightX := textX - highlightPadding
			highlightY := textY - highlightPadding/2

			rl.DrawRectangle(highlightX, highlightY, highlightWidth, highlightHeight, palette.Colors["COLOR_PLAYER"])
			rl.DrawText(optionText, textX, textY, optionSize, palette.Colors["COLOR_BACKGROUND"])
		} else {
			rl.DrawText(optionText, textX, textY, optionSize, palette.Colors["COLOR_LABEL"])
		}

		optionIndex++
	}
}
