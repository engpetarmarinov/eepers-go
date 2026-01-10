package main

import (
	"fmt"

	"github.com/engpetarmarinov/eepers-go/pkg/audio"
	"github.com/engpetarmarinov/eepers-go/pkg/entities"
	"github.com/engpetarmarinov/eepers-go/pkg/game"
	"github.com/engpetarmarinov/eepers-go/pkg/input"
	"github.com/engpetarmarinov/eepers-go/pkg/palette"
	"github.com/engpetarmarinov/eepers-go/pkg/ui"
	"github.com/engpetarmarinov/eepers-go/pkg/world"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// restartGame resets all game state and reloads the level
func restartGame(gs *game.State) error {
	// Clear all dynamic game state
	gs.Bombs = nil
	gs.Explosions = nil
	gs.Eepers = nil
	gs.Items = nil
	gs.TurnAnimation = 0
	gs.ShouldQuit = false

	// Reload the map
	err := game.LoadGameFromImage("assets/map.png", gs, true)
	if err != nil {
		return err
	}

	// Reset player state
	gs.Player.Health = 1.0
	gs.Player.BombSlots = 1
	gs.Player.Bombs = 0
	gs.Player.Keys = 0
	gs.Player.Dead = false

	// Reset camera
	gs.Camera.Zoom = 1.0

	// Reset menu and tutorial
	gs.Menu = game.NewMenuState()
	gs.Tutorial = game.TutorialState{
		Phase: game.TutorialMove,
	}

	return nil
}

func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	screenWidth := int32(rl.GetMonitorWidth(0))
	screenHeight := int32(rl.GetMonitorHeight(0))

	rl.InitWindow(screenWidth, screenHeight, "Eepers - Go Edition")
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
	err = game.LoadGameFromImage("assets/map.png", gs, true)
	if err != nil {
		panic(err)
	}
	gs.Player.Health = 1.0
	gs.Player.BombSlots = 1
	gs.Camera.Zoom = 1.0
	gs.Menu = game.NewMenuState()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() && !gs.ShouldQuit {
		screenWidth := int32(rl.GetScreenWidth())
		screenHeight := int32(rl.GetScreenHeight())
		gs.Camera.Offset = rl.NewVector2(float32(screenWidth/2), float32(screenHeight/2))
		inputState := input.GetInput()

		// Handle menu toggle
		if inputState.MenuToggle {
			gs.Menu.ToggleMenu()
		}

		// Handle menu input when menu is open
		if gs.Menu.IsOpen {
			if inputState.MenuNavigateUp {
				gs.Menu.MoveUp()
			}
			if inputState.MenuNavigateDown {
				gs.Menu.MoveDown()
			}
			if inputState.MenuConfirm {
				switch gs.Menu.SelectedOption {
				case game.MenuContinue:
					gs.Menu.CloseMenu()
				case game.MenuRestart:
					err = restartGame(gs)
					if err != nil {
						panic(err)
					}
				case game.MenuQuit:
					// Set quit flag to exit gracefully
					gs.ShouldQuit = true
				}
			}
		}

		// Only process game input when menu is closed
		if !gs.Menu.IsOpen {
			switch gs.Tutorial.Phase {
			case game.TutorialMove:
				gs.ShowPopup("Use arrow keys or left stick to move")
				if inputState.MoveRight || inputState.MoveLeft || inputState.MoveUp || inputState.MoveDown {
					gs.Tutorial.KnowsHowToMove = true
					gs.HidePopup()
					gs.Tutorial.Phase = game.TutorialPlaceBombs
				}
			case game.TutorialPlaceBombs:
				gs.ShowPopup("Press space or A button to plant a bomb")
				if inputState.PlaceBomb {
					gs.Tutorial.KnowsHowToPlaceBombs = true
					gs.HidePopup()
					gs.Tutorial.Phase = game.TutorialDone
				}
			}

			if !gs.Player.Dead {
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
			ui.UpdatePlayerEyes(&gs.Player)

			// Update turn animation with faster speed when running (shift/trigger held)
			if gs.TurnAnimation > 0 {
				animSpeed := float32(10.0)
				if inputState.IsRunning {
					animSpeed = 12.5 // 25% faster when sprinting (1.0 / 0.8 = 1.25)
				}
				gs.TurnAnimation -= rl.GetFrameTime() * animSpeed
			}

			if gs.Player.Dead && rl.GetTime() > gs.Player.DeathTime+2.0 {
				err = restartGame(gs)
				if err != nil {
					panic(err)
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
			rl.DrawRectangle(int32(eeper.Position.X*50), int32(eeper.Position.Y*50), int32(eeper.Size.X*50), int32(eeper.Size.Y*50), color)
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

		gs.DrawPopup(screenWidth, screenHeight)

		rl.EndMode2D()

		// Draw UI in screen space (outside of Mode2D)
		ui.DrawUI(gs, screenWidth)

		rl.DrawText("Eepers in Go!", 10, 10, 20, rl.LightGray)

		// Draw menu on top of everything
		gs.Menu.DrawMenu(screenWidth, screenHeight)

		rl.EndDrawing()
	}
}
