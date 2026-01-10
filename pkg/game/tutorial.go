package game

import (
	"time"

	"github.com/engpetarmarinov/eepers-go/pkg/input"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// PopupState represents the state of a pop-up message.
type PopupState struct {
	Label     string
	Visible   bool
	Animation float32
}

// TutorialPhase represents the different phases of the tutorial.
type TutorialPhase int

const (
	TutorialMove TutorialPhase = iota
	TutorialWaitingForBombPick
	TutorialPlaceBombs
	TutorialWaitingForSprint
	TutorialSprint
	TutorialDone
)

// TutorialState represents the state of the tutorial.
type TutorialState struct {
	Phase                TutorialPhase
	Waiting              float32
	PrevStepTimestamp    time.Time
	HurryCount           int
	Popup                PopupState
	KnowsHowToMove       bool
	KnowsHowToPlaceBombs bool
	KnowsHowToSprint     bool
}

// ShowPopup displays a pop-up message.
func (gs *State) ShowPopup(text string) {
	gs.Tutorial.Popup.Visible = true
	gs.Tutorial.Popup.Label = text
}

// HidePopup hides the pop-up message.
func (gs *State) HidePopup() {
	gs.Tutorial.Popup.Visible = false
}

// DrawPopup draws the pop-up message.
func (gs *State) DrawPopup(screenWidth, screenHeight int32) {
	popup := &gs.Tutorial.Popup
	if popup.Visible {
		if popup.Animation < 1.0 {
			popup.Animation += rl.GetFrameTime() * 5.0
			if popup.Animation > 1.0 {
				popup.Animation = 1.0
			}
		}
	} else {
		if popup.Animation > 0.0 {
			popup.Animation -= rl.GetFrameTime() * 5.0
			if popup.Animation < 0.0 {
				popup.Animation = 0.0
			}
		}
	}

	if popup.Animation > 0.0 {
		fontSize := float32(42) * popup.Animation
		textSize := rl.MeasureText(popup.Label, int32(fontSize))
		textPos := rl.NewVector2(float32(screenWidth)/2-float32(textSize)/2, float32(screenHeight)-100)

		rl.DrawText(popup.Label, int32(textPos.X), int32(textPos.Y), int32(fontSize), rl.White)
	}
}

// UpdateTutorial handles all tutorial phase logic and progression
func (gs *State) UpdateTutorial(inputState input.InputState) {
	switch gs.Tutorial.Phase {
	case TutorialMove:
		gs.ShowPopup("Use arrow keys or left stick to move.")
		if inputState.MoveRight || inputState.MoveLeft || inputState.MoveUp || inputState.MoveDown {
			gs.Tutorial.KnowsHowToMove = true
			gs.HidePopup()
			gs.Tutorial.Phase = TutorialWaitingForSprint
		}
	case TutorialWaitingForSprint:
		// Silent phase - waiting for player to move quickly (hurry)
		// When HurryCount >= 10, advance to TutorialSprint
		if gs.Tutorial.HurryCount >= 10 {
			gs.Tutorial.Phase = TutorialSprint
		}
	case TutorialSprint:
		gs.ShowPopup("Hold SHIFT or trigger to sprint.")
		if inputState.IsRunning {
			gs.Tutorial.KnowsHowToSprint = true
			gs.HidePopup()
			gs.Tutorial.Phase = TutorialWaitingForBombPick
		}
	case TutorialWaitingForBombPick:
		// Silent phase - waiting for player to pick up first bomb
		// Phase advances to TutorialPlaceBombs in player_turn.go when bomb is picked up
	case TutorialPlaceBombs:
		gs.ShowPopup("Press space or A button to plant a bomb.")
		if inputState.PlaceBomb {
			gs.Tutorial.KnowsHowToPlaceBombs = true
			gs.HidePopup()
			gs.Tutorial.Phase = TutorialDone
		}
	case TutorialDone:
		// Tutorial complete - do nothing
	}
}

// TutorialTrackMovementSpeed tracks how quickly the player is moving for the sprint tutorial
func (gs *State) TutorialTrackMovementSpeed(isRunning bool) {
	// Only track during waiting for sprint phase
	if gs.Tutorial.Phase != TutorialWaitingForSprint {
		return
	}

	currentTime := rl.GetTime()
	deltaTime := currentTime - float64(gs.Tutorial.PrevStepTimestamp.Unix())

	if deltaTime < 0.2 {
		gs.Tutorial.HurryCount++
	} else if gs.Tutorial.HurryCount > 0 {
		gs.Tutorial.HurryCount--
	}

	gs.Tutorial.PrevStepTimestamp = time.Unix(int64(currentTime), 0)

	// Mark as knowing how to sprint if holding shift
	if isRunning {
		gs.Tutorial.KnowsHowToSprint = true
	}
}
