package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(1920, 1080, "Fighting Game")

	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	//animations
	idleAnimation := NewAnimation("idle", rl.LoadTexture("Sprites/IdleAnimation.png"), 4, .2)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Blue)

		rl.EndDrawing()
	}
}
