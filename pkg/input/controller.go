package input

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	// Gamepad constants for Xbox controller
	GamepadPlayer1 = int32(0)

	// Dead zone for analog stick to avoid drift
	AnalogDeadZone = 0.25

	// Right trigger threshold for running
	TriggerThreshold = 0.1
)

// InputState represents the current input state from keyboard or controller
type InputState struct {
	MoveRight        bool
	MoveLeft         bool
	MoveUp           bool
	MoveDown         bool
	PlaceBomb        bool
	IsRunning        bool
	IsPressed        bool // For turn-based movement (false means held down for continuous)
	MenuToggle       bool // ESC key or Menu/Start button
	MenuConfirm      bool // Enter key or A button (when in menu)
	MenuNavigateUp   bool // Up arrow/stick (for menu navigation)
	MenuNavigateDown bool // Down arrow/stick (for menu navigation)
}

// AnalogState tracks previous analog stick state for detecting new presses
type AnalogState struct {
	PrevStickX float32
	PrevStickY float32
}

var analogState = AnalogState{}

// GetInput gets unified input from both keyboard and gamepad, combining inputs
func GetInput() InputState {
	var input InputState

	// Get keyboard input
	keyboardInput := getKeyboardInput()

	// Check if gamepad is available and get its input
	gamepadAvailable := rl.IsGamepadAvailable(GamepadPlayer1)
	var gamepadInput InputState
	if gamepadAvailable {
		gamepadInput = getGamepadInput()
	} else {
		// When no gamepad, set IsPressed to true (default for turn-based)
		gamepadInput.IsPressed = true
	}

	// Combine inputs - either keyboard OR gamepad can trigger movement
	input.MoveRight = keyboardInput.MoveRight || gamepadInput.MoveRight
	input.MoveLeft = keyboardInput.MoveLeft || gamepadInput.MoveLeft
	input.MoveUp = keyboardInput.MoveUp || gamepadInput.MoveUp
	input.MoveDown = keyboardInput.MoveDown || gamepadInput.MoveDown
	input.PlaceBomb = keyboardInput.PlaceBomb || gamepadInput.PlaceBomb

	// Menu inputs
	input.MenuToggle = keyboardInput.MenuToggle || gamepadInput.MenuToggle
	input.MenuConfirm = keyboardInput.MenuConfirm || gamepadInput.MenuConfirm
	input.MenuNavigateUp = keyboardInput.MenuNavigateUp || gamepadInput.MenuNavigateUp
	input.MenuNavigateDown = keyboardInput.MenuNavigateDown || gamepadInput.MenuNavigateDown

	// Running mode is active if either keyboard shift OR gamepad trigger is held
	input.IsRunning = keyboardInput.IsRunning || gamepadInput.IsRunning

	// IsPressed is true only if BOTH inputs are in pressed mode (turn-based)
	// If either is in continuous mode (IsPressed=false), we use continuous mode
	input.IsPressed = keyboardInput.IsPressed && gamepadInput.IsPressed

	return input
}

// getGamepadInput processes Xbox controller input
func getGamepadInput() InputState {
	var input InputState

	// Check if right trigger is held for running
	rightTrigger := rl.GetGamepadAxisMovement(GamepadPlayer1, rl.GamepadAxisRightTrigger)
	input.IsRunning = rightTrigger > TriggerThreshold

	// Get left analog stick values
	leftStickX := rl.GetGamepadAxisMovement(GamepadPlayer1, rl.GamepadAxisLeftX)
	leftStickY := rl.GetGamepadAxisMovement(GamepadPlayer1, rl.GamepadAxisLeftY)

	// Apply dead zone
	if abs(leftStickX) < AnalogDeadZone {
		leftStickX = 0
	}
	if abs(leftStickY) < AnalogDeadZone {
		leftStickY = 0
	}

	// Detect analog stick threshold crossings BEFORE updating previous state
	// This is used for both movement and menu navigation
	stickUpPressed := leftStickY < -AnalogDeadZone && analogState.PrevStickY >= -AnalogDeadZone
	stickDownPressed := leftStickY > AnalogDeadZone && analogState.PrevStickY <= AnalogDeadZone
	stickLeftPressed := leftStickX < -AnalogDeadZone && analogState.PrevStickX >= -AnalogDeadZone
	stickRightPressed := leftStickX > AnalogDeadZone && analogState.PrevStickX <= AnalogDeadZone

	// Determine movement direction based on analog stick
	// When running (trigger held), use continuous movement
	// When not running, detect threshold crossings for turn-based movement
	if input.IsRunning {
		input.IsPressed = false
		input.MoveRight = leftStickX > AnalogDeadZone
		input.MoveLeft = leftStickX < -AnalogDeadZone
		input.MoveDown = leftStickY > AnalogDeadZone
		input.MoveUp = leftStickY < -AnalogDeadZone
	} else {
		// For turn-based movement, detect new threshold crossings
		input.IsPressed = true
		input.MoveRight = stickRightPressed
		input.MoveLeft = stickLeftPressed
		input.MoveDown = stickDownPressed
		input.MoveUp = stickUpPressed
	}

	// Update previous state
	analogState.PrevStickX = leftStickX
	analogState.PrevStickY = leftStickY

	// Also check D-pad for digital movement
	if !input.IsRunning {
		if rl.IsGamepadButtonPressed(GamepadPlayer1, rl.GamepadButtonLeftFaceRight) {
			input.MoveRight = true
			input.IsPressed = true
		}
		if rl.IsGamepadButtonPressed(GamepadPlayer1, rl.GamepadButtonLeftFaceLeft) {
			input.MoveLeft = true
			input.IsPressed = true
		}
		if rl.IsGamepadButtonPressed(GamepadPlayer1, rl.GamepadButtonLeftFaceDown) {
			input.MoveDown = true
			input.IsPressed = true
		}
		if rl.IsGamepadButtonPressed(GamepadPlayer1, rl.GamepadButtonLeftFaceUp) {
			input.MoveUp = true
			input.IsPressed = true
		}
	} else {
		if rl.IsGamepadButtonDown(GamepadPlayer1, rl.GamepadButtonLeftFaceRight) {
			input.MoveRight = true
			input.IsPressed = false
		}
		if rl.IsGamepadButtonDown(GamepadPlayer1, rl.GamepadButtonLeftFaceLeft) {
			input.MoveLeft = true
			input.IsPressed = false
		}
		if rl.IsGamepadButtonDown(GamepadPlayer1, rl.GamepadButtonLeftFaceDown) {
			input.MoveDown = true
			input.IsPressed = false
		}
		if rl.IsGamepadButtonDown(GamepadPlayer1, rl.GamepadButtonLeftFaceUp) {
			input.MoveUp = true
			input.IsPressed = false
		}
	}

	// A button for placing bombs (like Space key)
	input.PlaceBomb = rl.IsGamepadButtonPressed(GamepadPlayer1, rl.GamepadButtonRightFaceDown)

	// Menu controls
	// Start/Menu button to toggle menu (button 7 on most controllers)
	input.MenuToggle = rl.IsGamepadButtonPressed(GamepadPlayer1, rl.GamepadButtonMiddleRight)

	// A button also confirms menu selection
	input.MenuConfirm = rl.IsGamepadButtonPressed(GamepadPlayer1, rl.GamepadButtonRightFaceDown)

	// Menu navigation with D-pad or analog stick (using the pre-calculated press states)
	input.MenuNavigateUp = rl.IsGamepadButtonPressed(GamepadPlayer1, rl.GamepadButtonLeftFaceUp) || stickUpPressed
	input.MenuNavigateDown = rl.IsGamepadButtonPressed(GamepadPlayer1, rl.GamepadButtonLeftFaceDown) || stickDownPressed

	return input
}

// getKeyboardInput processes keyboard input (existing behavior)
func getKeyboardInput() InputState {
	var input InputState

	// Check if shift is held for running/sprinting
	input.IsRunning = rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift)

	// When shift is held, use IsKeyDown for continuous movement
	// When shift is not held, use IsKeyPressed for turn-based movement
	if input.IsRunning {
		input.IsPressed = false
		input.MoveRight = rl.IsKeyDown(rl.KeyRight) || rl.IsKeyDown(rl.KeyD)
		input.MoveLeft = rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyA)
		input.MoveUp = rl.IsKeyDown(rl.KeyUp) || rl.IsKeyDown(rl.KeyW)
		input.MoveDown = rl.IsKeyDown(rl.KeyDown) || rl.IsKeyDown(rl.KeyS)
	} else {
		input.IsPressed = true
		input.MoveRight = rl.IsKeyPressed(rl.KeyRight) || rl.IsKeyPressed(rl.KeyD)
		input.MoveLeft = rl.IsKeyPressed(rl.KeyLeft) || rl.IsKeyPressed(rl.KeyA)
		input.MoveUp = rl.IsKeyPressed(rl.KeyUp) || rl.IsKeyPressed(rl.KeyW)
		input.MoveDown = rl.IsKeyPressed(rl.KeyDown) || rl.IsKeyPressed(rl.KeyS)
	}

	// Space for placing bombs
	input.PlaceBomb = rl.IsKeyPressed(rl.KeySpace)

	// Menu controls
	input.MenuToggle = rl.IsKeyPressed(rl.KeyEscape)
	input.MenuConfirm = rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeyKpEnter)
	input.MenuNavigateUp = rl.IsKeyPressed(rl.KeyUp) || rl.IsKeyPressed(rl.KeyW)
	input.MenuNavigateDown = rl.IsKeyPressed(rl.KeyDown) || rl.IsKeyPressed(rl.KeyS)

	return input
}

// abs returns the absolute value of a float32
func abs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}
