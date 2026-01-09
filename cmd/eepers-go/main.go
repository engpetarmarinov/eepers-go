package main

import (
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
	err = game.LoadGameFromImage("assets/map.png", gs, true, true)
	if err != nil {
		panic(err)
	}
	gs.Player.Health = 1.0
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
			if rl.IsKeyPressed(rl.KeyRight) {
				gs.PlayerTurn(game.Right)
				gs.TurnAnimation = 1.0
				gs.UpdateEepers()
			}
			if rl.IsKeyPressed(rl.KeyLeft) {
				gs.PlayerTurn(game.Left)
				gs.TurnAnimation = 1.0
				gs.UpdateEepers()
			}
			if rl.IsKeyPressed(rl.KeyUp) {
				gs.PlayerTurn(game.Up)
				gs.TurnAnimation = 1.0
				gs.UpdateEepers()
			}
			if rl.IsKeyPressed(rl.KeyDown) {
				gs.PlayerTurn(game.Down)
				gs.TurnAnimation = 1.0
				gs.UpdateEepers()
			}
			if rl.IsKeyPressed(rl.KeySpace) {
				gs.PlantBomb()
			}
		}

		ui.UpdatePlayerEyes(&gs.Player)

		if gs.TurnAnimation > 0 {
			gs.TurnAnimation -= rl.GetFrameTime() * 10
		} else {
			gs.UpdateBombs()
		}

		if gs.Player.Dead && rl.GetTime() > gs.Player.DeathTime+2.0 {
			// Restart game logic
			err = game.LoadGameFromImage("assets/map.png", gs, true, true)
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

		for y, row := range gs.Map {
			for x, cell := range row {
				color := world.CellColor(cell)
				rl.DrawRectangle(int32(x*50), int32(y*50), 50, 50, color)
			}
		}

		for _, item := range gs.Items {
			if item.Kind != entities.ItemNone {
				var color rl.Color
				switch item.Kind {
				case entities.ItemKey:
					color = palette.Colors["COLOR_DOORKEY"]
				case entities.ItemBombRefill:
					color = palette.Colors["COLOR_BOMB"]
				case entities.ItemBombSlot:
					color = palette.Colors["COLOR_DOORKEY"]
				case entities.ItemCheckpoint:
					color = palette.Colors["COLOR_CHECKPOINT"]
				}
				rl.DrawCircle(int32(item.Position.X*50+25), int32(item.Position.Y*50+25), 20, color)
			}
		}

		for _, bomb := range gs.Bombs {
			rl.DrawCircle(int32(bomb.Position.X*50+25), int32(bomb.Position.Y*50+25), 20, palette.Colors["COLOR_BOMB"])
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

		if gs.Player.Dead {
			rl.DrawText("YOU DIED", screenWidth/2-100, screenHeight/2-50, 50, rl.Red)
		}

		gs.DrawPopup(screenWidth, screenHeight)

		ui.DrawUI(gs, screenWidth)

		rl.EndMode2D()

		rl.DrawText("Eepers in Go!", 10, 10, 20, rl.LightGray)

		rl.EndDrawing()
	}
}
