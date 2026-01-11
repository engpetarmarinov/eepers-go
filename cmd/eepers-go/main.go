package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/engpetarmarinov/eepers-go/pkg/audio"
	"github.com/engpetarmarinov/eepers-go/pkg/entities"
	"github.com/engpetarmarinov/eepers-go/pkg/game"
	"github.com/engpetarmarinov/eepers-go/pkg/input"
	"github.com/engpetarmarinov/eepers-go/pkg/palette"
	"github.com/engpetarmarinov/eepers-go/pkg/ui"
	"github.com/engpetarmarinov/eepers-go/pkg/world"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	setupResourcePath()
	rl.SetConfigFlags(rl.FlagWindowMaximized | rl.FlagMsaa4xHint | rl.FlagVsyncHint | rl.FlagWindowHighdpi)
	screenWidth := int32(rl.GetMonitorWidth(0))
	screenHeight := int32(rl.GetMonitorHeight(0))
	rl.InitWindow(screenWidth, screenHeight, "Eepers Go")
	rl.ToggleFullscreen()
	// Disable ESC key as default exit key so we can use it for menu
	rl.SetExitKey(0)
	defer rl.CloseWindow()
	audio.LoadAudio()
	defer audio.UnloadAudio()
	err := game.LoadColors("assets/colors.txt")
	if err != nil {
		panic(err)
	}

	gs := &game.State{}

	// Configure worlds and their levels
	gs.WorldConfig = game.WorldConfig{
		CurrentWorld: 0,
		Worlds: []game.World{
			{
				Name:     "World 1",
				HubLevel: "assets/worlds/1/hub.png",
				Levels: []string{
					//"assets/worlds/1/levels/1.png",
					//"assets/worlds/1/levels/2.png",
					//"assets/worlds/1/levels/3.png",
					//"assets/worlds/1/levels/4.png",
					"assets/worlds/1/levels/1-debug.png",
					"assets/worlds/1/levels/2-debug.png",
					"assets/worlds/1/levels/3-debug.png",
					"assets/worlds/1/levels/4-debug.png",
				},
			},
		},
	}

	// Load the first world's hub level
	err = gs.LoadHub()
	if err != nil {
		panic(err)
	}
	gs.Camera.Zoom = 1.0
	gs.Menu = game.NewMenuState()

	// Start playing ambient music
	rl.PlayMusicStream(audio.AmbientMusic)

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() && !gs.ShouldQuit {
		screenWidth := int32(rl.GetScreenWidth())
		screenHeight := int32(rl.GetScreenHeight())
		gs.Camera.Offset = rl.NewVector2(float32(screenWidth/2), float32(screenHeight/2))
		inputState := input.GetInput()

		// Update music stream
		rl.UpdateMusicStream(audio.AmbientMusic)

		// Handle menu toggle
		if inputState.MenuToggle {
			gs.Menu.ToggleMenu()
			// Pause/resume music based on menu state
			if gs.Menu.IsOpen {
				rl.PauseMusicStream(audio.AmbientMusic)
			} else {
				rl.ResumeMusicStream(audio.AmbientMusic)
			}
		}

		// Handle menu input when menu is open
		if gs.Menu.IsOpen {
			if inputState.MenuNavigateUp {
				if gs.InHub {
					gs.Menu.MoveUpInHub()
				} else {
					gs.Menu.MoveUp()
				}
			}
			if inputState.MenuNavigateDown {
				if gs.InHub {
					gs.Menu.MoveDownInHub()
				} else {
					gs.Menu.MoveDown()
				}
			}
			if inputState.MenuConfirm {
				switch gs.Menu.SelectedOption {
				case game.MenuContinue:
					gs.Menu.CloseMenu()
					rl.ResumeMusicStream(audio.AmbientMusic)
				case game.MenuRestart:
					err = restartGame(gs)
					if err != nil {
						panic(err)
					}
					gs.Menu.CloseMenu()
					rl.ResumeMusicStream(audio.AmbientMusic)
				case game.MenuExitLevel:
					// Return to hub level
					err = gs.LoadHub()
					if err != nil {
						panic(err)
					}
					gs.Menu.CloseMenu()
					rl.ResumeMusicStream(audio.AmbientMusic)
				case game.MenuQuit:
					// Set quit flag to exit gracefully
					gs.ShouldQuit = true
				}
			}
		}

		// Only process game input when menu is closed
		if !gs.Menu.IsOpen {
			gs.UpdateTutorial(inputState)

			// Prevent input during portal entry animation
			if !gs.Player.Dead && !gs.Player.EnteringPortal {
				// Handle movement based on input state
				// When running (shift/trigger held), allow continuous movement
				// When not running, use turn-based movement
				var shouldMove bool
				if inputState.IsRunning && gs.TurnAnimation <= 0 {
					shouldMove = true
				} else if inputState.IsPressed {
					shouldMove = true
				}

				if shouldMove {
					// Track movement speed for sprint tutorial
					gs.TutorialTrackMovementSpeed(inputState.IsRunning)

					// Right
					if inputState.MoveRight {
						gs.PlayerTurn(game.Right)
						gs.TurnAnimation = 1.0
						gs.ItemsTurn()
						gs.UpdateEepers()
						gs.UpdateBombs()
					}
					// Left
					if inputState.MoveLeft {
						gs.PlayerTurn(game.Left)
						gs.TurnAnimation = 1.0
						gs.ItemsTurn()
						gs.UpdateEepers()
						gs.UpdateBombs()
					}
					// Up
					if inputState.MoveUp {
						gs.PlayerTurn(game.Up)
						gs.TurnAnimation = 1.0
						gs.ItemsTurn()
						gs.UpdateEepers()
						gs.UpdateBombs()
					}
					// Down
					if inputState.MoveDown {
						gs.PlayerTurn(game.Down)
						gs.TurnAnimation = 1.0
						gs.ItemsTurn()
						gs.UpdateEepers()
						gs.UpdateBombs()
					}
				}

				if inputState.PlaceBomb {
					gs.PlantBomb()
				}
			}

			gs.UpdateExplosions()
			gs.UpdatePortals()
			ui.UpdatePlayerEyes(&gs.Player)

			// Handle portal entry animation
			if gs.Player.EnteringPortal {
				portalDuration := rl.GetTime() - gs.Player.PortalEntryTime
				const portalAnimationTime = 0.8 // Short 0.8 second animation

				if portalDuration < portalAnimationTime {
					// Quick zoom in - from 1.0 to 2.5 in 0.8 seconds
					zoomProgress := float32(portalDuration / portalAnimationTime)
					// Use ease-in curve for acceleration effect (falling feeling)
					easeProgress := zoomProgress * zoomProgress
					gs.Camera.Zoom = 1.0 + easeProgress*1.5 // Zoom from 1.0 to 2.5
				} else {
					// Animation complete - activate the portal
					err = gs.LoadLevelFromPortal(gs.Player.PortalToActivate)
					if err != nil {
						panic(err)
					}
					// Reset portal entry state
					gs.Player.EnteringPortal = false
					gs.Player.PortalToActivate = 0
					gs.Camera.Zoom = 1.0 // Reset zoom
				}
			}

			if gs.Player.Dead && rl.GetTime() > gs.Player.DeathTime+2.0 {
				// Restore from last checkpoint
				gs.RestoreCheckpoint()
			}

			// Check victory condition - player reached Father!
			if gs.Player.ReachedFather {
				victoryDuration := rl.GetTime() - gs.Player.VictoryTime

				// Play victory sound once at the start
				if victoryDuration < 0.1 {
					rl.PlaySound(audio.VictorySound)
				}

				// Victory animation: zoom in camera over 11 seconds
				const victoryAnimationTime = 11.0
				if victoryDuration < victoryAnimationTime {
					// Smoothly zoom from 1.0 to 3.0 over 3 seconds
					zoomProgress := float32(victoryDuration / victoryAnimationTime)
					// Use ease-in-out curve for smooth animation
					easeProgress := zoomProgress * zoomProgress * (3.0 - 2.0*zoomProgress)
					gs.Camera.Zoom = 1.0 + easeProgress*2.0 // Zoom from 1.0 to 3.0
				} else {
					// Animation complete - try to load next level
					hasNextLevel, err := gs.LoadNextLevel()
					if err != nil {
						panic(err)
					}

					if !hasNextLevel {
						// No more levels - game complete!
						gs.ShowPopup("Congratulations! You completed all levels!")
						gs.Player.ReachedFather = false
						gs.Player.VictoryTime = 0
						gs.Camera.Zoom = 1.0
						gs.HidePopup()
						// Could add a "game complete" state here
						// For now, we'll just restart from the first level
						err = gs.RestartFromFirstLevel()
						if err != nil {
							panic(err)
						}
					} else {
						// Next level loaded successfully
						gs.HidePopup()
					}
				}
			}
		}

		// Update camera
		cameraTarget := rl.NewVector2(float32(gs.Player.Position.X*50), float32(gs.Player.Position.Y*50))
		gs.Camera.Target = rl.Vector2Lerp(gs.Camera.Target, cameraTarget, rl.GetFrameTime()*5.0)

		// Draw
		rl.BeginDrawing()
		rl.ClearBackground(palette.Colors["COLOR_BACKGROUND"])

		rl.BeginMode2D(gs.Camera)

		// Draw map cells first
		for y, row := range gs.Map {
			for x, cell := range row {
				color := world.CellColor(cell)
				rl.DrawRectangle(int32(x*50), int32(y*50), 50, 50, color)
			}
		}

		// Draw portals on top of floor
		for _, portal := range gs.Portals {
			ui.DrawPortal(portal)
		}

		// Then draw explosions on top
		for _, explosion := range gs.Explosions {
			alpha := float32(explosion.Timer) / float32(explosion.InitialTimer)
			color := rl.Fade(palette.Colors["COLOR_BOMB_FLASH"], alpha)
			rl.DrawRectangle(int32(explosion.Position.X*50), int32(explosion.Position.Y*50), 50, 50, color)
		}

		for _, item := range gs.Items {
			if item.Kind != entities.ItemNone {
				var color rl.Color
				switch item.Kind {
				case entities.ItemKey:
					color = palette.Colors["COLOR_DOORKEY"]
					rl.DrawCircle(int32(item.Position.X*50+25), int32(item.Position.Y*50+25), 20, color)
				case entities.ItemBombRefill:
					// Show dimmed bomb and cooldown timer if on cooldown
					if item.Cooldown > 0 {
						color = rl.ColorBrightness(palette.Colors["COLOR_BOMB"], -0.5)
						rl.DrawCircle(int32(item.Position.X*50+25), int32(item.Position.Y*50+25), 20, color)
						// Draw cooldown timer
						countdownText := fmt.Sprintf("%d", item.Cooldown)
						textWidth := rl.MeasureText(countdownText, 20)
						rl.DrawText(countdownText, int32(item.Position.X*50+25)-textWidth/2, int32(item.Position.Y*50+15), 20, palette.Colors["COLOR_LABEL"])
					} else {
						color = palette.Colors["COLOR_BOMB"]
						rl.DrawCircle(int32(item.Position.X*50+25), int32(item.Position.Y*50+25), 20, color)
					}
				case entities.ItemBombSlot:
					color = palette.Colors["COLOR_DOORKEY"]
					rl.DrawCircle(int32(item.Position.X*50+25), int32(item.Position.Y*50+25), 20, color)
				case entities.ItemCheckpoint:
					color = palette.Colors["COLOR_CHECKPOINT"]
					rl.DrawCircle(int32(item.Position.X*50+25), int32(item.Position.Y*50+25), 20, color)
				}
			}
		}

		for _, eeper := range gs.Eepers {
			if eeper.Dead {
				continue
			}

			var color rl.Color
			switch eeper.Kind {
			case entities.EeperGuard:
				color = palette.Colors["COLOR_GUARD"]
			case entities.EeperMother:
				color = palette.Colors["COLOR_MOTHER"]
			case entities.EeperGnome:
				color = palette.Colors["COLOR_DOORKEY"]
			case entities.EeperFather:
				color = palette.Colors["COLOR_FATHER"]
			}

			// Interpolate eeper position for smooth movement
			eeperPrevPos := rl.NewVector2(float32(eeper.PrevPosition.X*50), float32(eeper.PrevPosition.Y*50))
			eeperPos := rl.NewVector2(float32(eeper.Position.X*50), float32(eeper.Position.Y*50))
			eeperInterpPos := rl.Vector2Lerp(eeperPos, eeperPrevPos, gs.TurnAnimation)
			eeperSize := rl.NewVector2(float32(eeper.Size.X*50), float32(eeper.Size.Y*50))

			// Gnomes are rendered smaller (70% size) and centered
			renderPos := eeperInterpPos
			renderSize := eeperSize
			if eeper.Kind == entities.EeperGnome {
				gnomeRatio := float32(0.7)
				renderSize = rl.NewVector2(eeperSize.X*gnomeRatio, eeperSize.Y*gnomeRatio)
				offset := rl.NewVector2((eeperSize.X-renderSize.X)*0.5, (eeperSize.Y-renderSize.Y)*0.5)
				renderPos = rl.Vector2Add(eeperInterpPos, offset)
			}

			// Draw eeper body
			rl.DrawRectangleV(renderPos, renderSize, color)

			// Draw health bar for guards and mothers
			if eeper.Kind == entities.EeperGuard || eeper.Kind == entities.EeperMother {
				ui.DrawEeperHealthBar(eeper, eeperInterpPos, eeperSize)

				// Draw cooldown bubble only when the eeper can see the player (path >= 0)
				if eeper.Path != nil && eeper.Position.Y >= 0 && eeper.Position.Y < len(eeper.Path) &&
					eeper.Position.X >= 0 && eeper.Position.X < len(eeper.Path[0]) &&
					eeper.Path[eeper.Position.Y][eeper.Position.X] >= 0 {
					ui.DrawEeperCooldownBubble(eeper, eeperInterpPos, eeperSize, color)
				}
			}

			// Draw eeper eyes (use renderPos and renderSize for gnomes)
			ui.DrawEeperEyes(eeper, renderPos, gs.TurnAnimation)
		}

		playerPrevPos := rl.NewVector2(float32(gs.Player.PrevPosition.X*50), float32(gs.Player.PrevPosition.Y*50))
		playerPos := rl.NewVector2(float32(gs.Player.Position.X*50), float32(gs.Player.Position.Y*50))
		interpPos := rl.Vector2Lerp(playerPos, playerPrevPos, gs.TurnAnimation)
		rl.DrawRectangleV(interpPos, rl.NewVector2(50, 50), palette.Colors["COLOR_PLAYER"])
		ui.DrawPlayerEyes(gs.Player, interpPos)

		// Draw bombs AFTER player so they appear on top
		for _, bomb := range gs.Bombs {
			rl.DrawCircle(int32(bomb.Position.X*50+25), int32(bomb.Position.Y*50+25), 20, palette.Colors["COLOR_BOMB"])
			countdownText := fmt.Sprintf("%d", bomb.Countdown)
			textWidth := rl.MeasureText(countdownText, 20)
			rl.DrawText(countdownText, int32(bomb.Position.X*50+25)-textWidth/2, int32(bomb.Position.Y*50+15), 20, rl.White)
		}

		if gs.Player.Dead {
			rl.DrawText("YOU DIED", screenWidth/2-100, screenHeight/2-50, 50, rl.Red)
		}

		rl.EndMode2D()

		// Draw popup in screen space (not affected by camera zoom)
		gs.DrawPopup(screenWidth, screenHeight)

		// Draw UI in screen space (outside of Mode2D)
		ui.DrawUI(gs, screenWidth)

		// Draw menu on top of everything
		gs.Menu.DrawMenu(gs.InHub)

		rl.EndDrawing()

		// Update turn animation AFTER rendering to ensure first frame shows correct positions
		// Update with faster speed when running (shift/trigger held)
		if gs.TurnAnimation > 0 {
			animSpeed := float32(10.0)
			if inputState.IsRunning {
				animSpeed = 12 // 20% faster when sprinting
			}
			gs.TurnAnimation -= rl.GetFrameTime() * animSpeed
			// Clamp to 0 to prevent negative values that cause extrapolation
			if gs.TurnAnimation < 0 {
				gs.TurnAnimation = 0
			}
		}
	}
}

// restartGame resets all game state and reloads the first level
func restartGame(gs *game.State) error {
	return gs.RestartFromFirstLevel()
}

// setupResourcePath changes the working directory to the Resources folder
// when running from a macOS .app bundle
func setupResourcePath() {
	if runtime.GOOS == "darwin" {
		exePath, err := os.Executable()
		if err != nil {
			return
		}

		// Check if we're running from a .app bundle
		// Executable path will be: Eepers.app/Contents/MacOS/eepers
		exeDir := filepath.Dir(exePath)
		if filepath.Base(exeDir) == "MacOS" {
			resourcesPath := filepath.Join(filepath.Dir(exeDir), "Resources")
			if _, err := os.Stat(resourcesPath); err == nil {
				_ = os.Chdir(resourcesPath)
			}
		}
	}
}
