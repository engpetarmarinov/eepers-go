package entities

// EyesKind represents the different states of the eyes.
type EyesKind int

const (
	EyesOpen EyesKind = iota
	EyesClosed
	EyesAngry
	EyesCringe
	EyesSurprised
)
