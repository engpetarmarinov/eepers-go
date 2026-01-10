package audio

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	BlastSound      rl.Sound
	KeyPickupSound  rl.Sound
	BombPickupSound rl.Sound
	OpenDoorSound   rl.Sound
	CheckpointSound rl.Sound
	PlantBombSound  rl.Sound
	GuardStepSound  rl.Sound
	FootstepsSounds []rl.Sound
	AmbientMusic    rl.Music
)

// LoadAudio loads all the audio files for the game.
func LoadAudio() {
	rl.InitAudioDevice()

	BlastSound = rl.LoadSound("assets/sounds/blast.ogg")
	KeyPickupSound = rl.LoadSound("assets/sounds/key-pickup.wav")
	BombPickupSound = rl.LoadSound("assets/sounds/bomb-pickup.ogg")
	OpenDoorSound = rl.LoadSound("assets/sounds/open-door.wav")
	CheckpointSound = rl.LoadSound("assets/sounds/checkpoint.ogg")
	PlantBombSound = rl.LoadSound("assets/sounds/plant-bomb.wav")
	GuardStepSound = rl.LoadSound("assets/sounds/guard-step.ogg")

	FootstepsSounds = make([]rl.Sound, 4)
	FootstepsSounds[0] = rl.LoadSound("assets/sounds/footsteps.mp3")
	FootstepsSounds[1] = rl.LoadSound("assets/sounds/footsteps.mp3")
	FootstepsSounds[2] = rl.LoadSound("assets/sounds/footsteps.mp3")
	FootstepsSounds[3] = rl.LoadSound("assets/sounds/footsteps.mp3")

	AmbientMusic = rl.LoadMusicStream("assets/sounds/ambient.wav")
}

// UnloadAudio unloads all the audio files.
func UnloadAudio() {
	rl.UnloadSound(BlastSound)
	rl.UnloadSound(KeyPickupSound)
	rl.UnloadSound(BombPickupSound)
	rl.UnloadSound(OpenDoorSound)
	rl.UnloadSound(CheckpointSound)
	rl.UnloadSound(PlantBombSound)
	rl.UnloadSound(GuardStepSound)
	for _, s := range FootstepsSounds {
		rl.UnloadSound(s)
	}
	rl.UnloadMusicStream(AmbientMusic)
	rl.CloseAudioDevice()
}
