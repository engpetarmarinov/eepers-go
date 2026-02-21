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
		rl.TraceLog(rl.LogInfo, "MENU OPENED: Screen=%dx%d, Render=%dx%d",
			rl.GetScreenWidth(), rl.GetScreenHeight(),
			rl.GetRenderWidth(), rl.GetRenderHeight())
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

func (ms *MenuState) DrawMenu(inHub bool) {
	if !ms.IsOpen {
		return
	}

	// screenWidth := int32(rl.GetScreenWidth())
	// screenHeight := int32(rl.GetScreenHeight())

	renderWidth := int32(rl.GetRenderWidth())
	renderHeight := int32(rl.GetRenderHeight())
	// Menu box dimensions - use screen dimensions for scaling and centering
	menuWidth := renderWidth / 3
	if menuWidth < 400 {
		menuWidth = 400
	}
	menuHeight := renderHeight / 3
	if menuHeight < 300 {
		menuHeight = 300
	}
	menuX := (renderWidth - menuWidth) / 2
	menuY := (renderHeight - menuHeight) / 2

	// Semi-transparent overlay covering the entire screen
	// Use DrawRectangleRec for precise full-screen coverage
	fullScreenRect := rl.NewRectangle(0, 0, float32(renderWidth), float32(renderHeight))
	rl.DrawRectangleRec(fullScreenRect, rl.Fade(rl.Black, 0.7))

	// Draw menu background (DarkGray provides contrast without needing Alpha Blending)
	rl.DrawRectangle(menuX, menuY, menuWidth, menuHeight, palette.Colors["COLOR_BACKGROUND"])

	// Draw menu border using simple basic lines (avoids DrawRectangleLinesEx GL batch bugs)
	borderColor := palette.Colors["COLOR_LABEL"]
	rl.DrawRectangleLines(menuX, menuY, menuWidth, menuHeight, borderColor)
	rl.DrawRectangleLines(menuX+1, menuY+1, menuWidth-2, menuHeight-2, borderColor)

	// Scale fonts according to menu height
	titleSize := menuHeight / 6
	if titleSize < 40 {
		titleSize = 40
	}
	optionSize := menuHeight / 10
	if optionSize < 30 {
		optionSize = 30
	}
	optionSpacing := menuHeight / 6
	if optionSpacing < 50 {
		optionSpacing = 50
	}

	// Draw title
	title := "PAUSED"
	titleWidth := rl.MeasureText(title, titleSize)
	titleX := menuX + (menuWidth-titleWidth)/2
	titleY := menuY + menuHeight/10
	rl.DrawText(title, titleX, titleY, titleSize, palette.Colors["COLOR_LABEL"])

	// Draw menu options
	optionY := menuY + menuHeight/3

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
			highlightPadding := optionSize / 3
			highlightWidth := textWidth + highlightPadding*2
			highlightHeight := optionSize + highlightPadding
			highlightX := textX - highlightPadding
			highlightY := textY - highlightPadding/2

			rl.DrawRectangle(highlightX, highlightY, highlightWidth, highlightHeight, palette.Colors["COLOR_PLAYER"])
			rl.DrawText(optionText, textX, textY, optionSize, palette.Colors["COLOR_BACKGROUND"])
		} else {
			rl.DrawText(optionText, textX, textY, optionSize, rl.LightGray)
		}

		optionIndex++
	}
}
