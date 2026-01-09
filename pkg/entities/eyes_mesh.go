package entities

import rl "github.com/gen2brain/raylib-go/raylib"

type EyeMesh [4]rl.Vector2
type EyesMesh [2]EyeMesh

var EyesMeshes = map[EyesKind]EyesMesh{
	EyesOpen: {
		// Left Eye
		EyeMesh{{X: 0.0, Y: 0.0}, {X: 0.0, Y: 1.0}, {X: 1.0, Y: 0.0}, {X: 1.0, Y: 1.0}},
		// Right Eye
		EyeMesh{{X: 0.0, Y: 0.0}, {X: 0.0, Y: 1.0}, {X: 1.0, Y: 0.0}, {X: 1.0, Y: 1.0}},
	},
	EyesClosed: {
		// Left Eye
		EyeMesh{{X: 0.0, Y: 0.8}, {X: 0.0, Y: 1.0}, {X: 1.0, Y: 0.8}, {X: 1.0, Y: 1.0}},
		// Right Eye
		EyeMesh{{X: 0.0, Y: 0.8}, {X: 0.0, Y: 1.0}, {X: 1.0, Y: 0.8}, {X: 1.0, Y: 1.0}},
	},
	EyesAngry: {
		// Left Eye
		EyeMesh{{X: 0.0, Y: 0.0}, {X: 0.0, Y: 1.0}, {X: 1.0, Y: 0.3}, {X: 1.0, Y: 1.0}},
		// Right Eye
		EyeMesh{{X: 0.0, Y: 0.3}, {X: 0.0, Y: 1.0}, {X: 1.0, Y: 0.0}, {X: 1.0, Y: 1.0}},
	},
	EyesCringe: {
		// Left Eye
		EyeMesh{{X: 0.0, Y: 0.5}, {X: 0.25, Y: 0.75}, {X: 1.3, Y: 0.75}, {X: 0.0, Y: 1.0}},
		// Right Eye
		EyeMesh{{X: 1.0, Y: 0.5}, {X: 0.75, Y: 0.75}, {X: -0.3, Y: 0.75}, {X: 1.0, Y: 1.0}},
	},
	EyesSurprised: {
		// Left Eye
		EyeMesh{{X: 0.0, Y: 0.3}, {X: 0.0, Y: 1.0}, {X: 1.0, Y: 0.3}, {X: 1.0, Y: 1.0}},
		// Right Eye
		EyeMesh{{X: 0.0, Y: 0.3}, {X: 0.0, Y: 1.0}, {X: 1.0, Y: 0.3}, {X: 1.0, Y: 1.0}},
	},
}
