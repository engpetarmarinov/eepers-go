package game

import (
	"time"

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
