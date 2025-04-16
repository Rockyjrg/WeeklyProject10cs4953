package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(1920, 1080, "Fighting Game")

	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	//stage variables
	platformHeight := float32(50)
	platformY := float32(rl.GetScreenHeight() - int(platformHeight) - 100)
	platformRect := rl.NewRectangle(0, platformY, float32(rl.GetScreenWidth()), platformY)
	platformColor := rl.DarkGray

	//animations
	idleAnimation := NewAnimation("idle", rl.LoadTexture("Sprites/IdleAnimation(2).png"), 5, .2)
	walkAnimation := NewAnimation("walk", rl.LoadTexture("Sprites/WalkingAnimation(1).png"), 3, .075)
	jumpAnimation := NewAnimation("jump", rl.LoadTexture("Sprites/JumpAnimation(2).png"), 5, .1)
	jumpAnimation.Loop = false

	//fighting animations
	//blockingAnimation := NewAnimation("block", rl.LoadTexture("Sprites/BlockingAnimation(1).png"), 2, 0.7)

	animationFSM := NewAnimationFSM()
	animationFSM.AddAnimation(walkAnimation)
	animationFSM.AddAnimation(jumpAnimation)
	animationFSM.AddAnimation(idleAnimation)
	animationFSM.ChangeAnimationState("idle")

	player1 := Creature{
		Pos:          rl.NewVector2(380, 100),
		Vel:          rl.NewVector2(0, 0),
		Size:         200,
		Color:        rl.Red,
		Speed:        300,
		Direction:    1,
		AnimationFSM: animationFSM,
		IsGrounded:   false,
	}

	gravity := rl.NewVector2(0, 1000)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Blue)

		player1VelX := float32(0.0)
		if rl.IsKeyDown(rl.KeyD) {
			player1VelX = 1
		}
		if rl.IsKeyDown(rl.KeyA) {
			player1VelX = -1
		}

		player1.Move(player1VelX)

		player1.ApplyGravity(gravity)
		player1.UpdateCreature(platformY)
		rl.DrawRectangleRec(platformRect, platformColor)
		player1.DrawCreature()
		rl.EndDrawing()
	}
}
