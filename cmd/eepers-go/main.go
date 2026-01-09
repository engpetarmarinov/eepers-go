package main

import (
	"fmt"

	"github.com/engpetarmarinov/eepers-go/pkg/audio"
	"github.com/engpetarmarinov/eepers-go/pkg/entities"
	"github.com/engpetarmarinov/eepers-go/pkg/game"
	"github.com/engpetarmarinov/eepers-go/pkg/palette"
	"github.com/engpetarmarinov/eepers-go/pkg/ui"
	"github.com/engpetarmarinov/eepers-go/pkg/world"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	screenWidth := int32(rl.GetMonitorWidth(0))
	screenHeight := int32(rl.GetMonitorHeight(0))

	rl.InitWindow(screenWidth, screenHeight, "Eepers - Go Edition")
	rl.MaximizeWindow()
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

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		screenWidth := int32(rl.GetScreenWidth())
		screenHeight := int32(rl.GetScreenHeight())
		gs.Camera.Offset = rl.NewVector2(float32(screenWidth/2), float32(screenHeight/2))
		// Update
		switch gs.Tutorial.Phase {
		case game.TutorialMove:
			gs.ShowPopup("Use arrow keys to move")
			if rl.IsKeyPressed(rl.KeyRight) || rl.IsKeyPressed(rl.KeyLeft) || rl.IsKeyPressed(rl.KeyUp) || rl.IsKeyPressed(rl.KeyDown) {
				gs.Tutorial.KnowsHowToMove = true
				gs.HidePopup()
				gs.Tutorial.Phase = game.TutorialPlaceBombs
			}
		case game.TutorialPlaceBombs:
			gs.ShowPopup("Press space to plant a bomb")
			if rl.IsKeyPressed(rl.KeySpace) {
				gs.Tutorial.KnowsHowToPlaceBombs = true
				gs.HidePopup()
				gs.Tutorial.Phase = game.TutorialDone
			}
		}

		if !gs.Player.Dead {
			// Check if shift is held for running/sprinting
			holdingShift := rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift)

			// When shift is held, use IsKeyDown for continuous movement
			// When shift is not held, use IsKeyPressed for turn-based movement
			var rightKey, leftKey, upKey, downKey bool
			if holdingShift && gs.TurnAnimation <= 0 {
				rightKey = rl.IsKeyDown(rl.KeyRight) || rl.IsKeyDown(rl.KeyD)
				leftKey = rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyA)
				upKey = rl.IsKeyDown(rl.KeyUp) || rl.IsKeyDown(rl.KeyW)
				downKey = rl.IsKeyDown(rl.KeyDown) || rl.IsKeyDown(rl.KeyS)
			} else {
				rightKey = rl.IsKeyPressed(rl.KeyRight) || rl.IsKeyPressed(rl.KeyD)
				leftKey = rl.IsKeyPressed(rl.KeyLeft) || rl.IsKeyPressed(rl.KeyA)
				upKey = rl.IsKeyPressed(rl.KeyUp) || rl.IsKeyPressed(rl.KeyW)
				downKey = rl.IsKeyPressed(rl.KeyDown) || rl.IsKeyPressed(rl.KeyS)
			}

			// Right: Arrow Right or D
			if rightKey {
				gs.PlayerTurn(game.Right)
				gs.TurnAnimation = 1.0
				gs.ItemsTurn()
				gs.UpdateEepers()
				gs.UpdateBombs()
			}
			// Left: Arrow Left or A
			if leftKey {
				gs.PlayerTurn(game.Left)
				gs.TurnAnimation = 1.0
				gs.ItemsTurn()
				gs.UpdateEepers()
				gs.UpdateBombs()
			}
			// Up: Arrow Up or W
			if upKey {
				gs.PlayerTurn(game.Up)
				gs.TurnAnimation = 1.0
				gs.ItemsTurn()
				gs.UpdateEepers()
				gs.UpdateBombs()
			}
			// Down: Arrow Down or S
			if downKey {
				gs.PlayerTurn(game.Down)
				gs.TurnAnimation = 1.0
				gs.ItemsTurn()
				gs.UpdateEepers()
				gs.UpdateBombs()
			}
			if rl.IsKeyPressed(rl.KeySpace) {
				gs.PlantBomb()
			}
		}

		gs.UpdateExplosions()
		ui.UpdatePlayerEyes(&gs.Player)

		// Update turn animation with faster speed when shift is held
		if gs.TurnAnimation > 0 {
			holdingShift := rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift)
			animSpeed := float32(10.0)
			if holdingShift {
				animSpeed = 12.5 // 25% faster when sprinting (1.0 / 0.8 = 1.25)
			}
			gs.TurnAnimation -= rl.GetFrameTime() * animSpeed
		}

		if gs.Player.Dead && rl.GetTime() > gs.Player.DeathTime+2.0 {
			// Restart game logic
			err = game.LoadGameFromImage("assets/map.png", gs, true)
			if err != nil {
				panic(err)
			}
			gs.Player.Health = 1.0
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

		rl.EndDrawing()
	}
}
