package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameState int

const (
	Playing GameState = iota
	GameOver
)

var (
	currentGameState GameState = Playing
	winner           int       = 0
)

var player1 Creature
var player2 Creature
var player1HealthBar HealthBar
var player2HealthBar HealthBar
var platformRect rl.Rectangle
var platformY float32

func ResetGame() {
	fmt.Println("Resetting Game...")
	//reset player 1
	player1.Pos = rl.NewVector2(float32(rl.GetScreenWidth())*0.2, platformY-player1.Size) // Start slightly above platform
	player1.Vel = rl.Vector2Zero()
	player1.Health = player1.MaxHealth
	player1.Direction = 1
	player1.IsGrounded = false
	player1.IsAttacking = false
	player1.IsBlocking = false
	player1.CanDamage = false
	player1.AttackTimer = 0
	player1.AttackCooldownTimer = 0
	player1.AnimationFSM.ChangeAnimationState("idle") // Reset animation state

	//reset player 2
	player2.Pos = rl.NewVector2(float32(rl.GetScreenWidth())*0.8-player2.Size, platformY-player2.Size) // Start slightly above platform
	player2.Vel = rl.Vector2Zero()
	player2.Health = player2.MaxHealth
	player2.Direction = -1
	player2.IsGrounded = false
	player2.IsAttacking = false
	player2.IsBlocking = false
	player2.CanDamage = false
	player2.AttackTimer = 0
	player2.AttackCooldownTimer = 0
	player2.AnimationFSM.ChangeAnimationState("idle") // Reset animation state

	//reset game state
	currentGameState = Playing
	winner = 0
}

func main() {
	rl.InitWindow(1920, 1080, "Fighting Game")

	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	//stage variables
	platformHeight := float32(50)
	platformY := float32(rl.GetScreenHeight() - int(platformHeight) - 100)
	platformRect = rl.NewRectangle(0, platformY, float32(rl.GetScreenWidth()), platformY)
	platformColor := rl.DarkGray

	//animations
	idleAnimation := NewAnimation("idle", rl.LoadTexture("Sprites/IdleAnimation(2).png"), 5, .2)
	walkAnimation := NewAnimation("walk", rl.LoadTexture("Sprites/WalkingAnimation(1).png"), 3, .075)
	jumpAnimation := NewAnimation("jump", rl.LoadTexture("Sprites/JumpAnimation(2).png"), 5, .1)
	punchAnimation := NewAnimation("punch", rl.LoadTexture("Sprites/PunchingAnimation(1).png"), 2, .075)
	blockAnimation := NewAnimation("block", rl.LoadTexture("Sprites/BlockingAnimation(1).png"), 2, 0.2)
	jumpAnimation.Loop = false
	blockAnimation.Loop = false

	//fighting animations
	//blockingAnimation := NewAnimation("block", rl.LoadTexture("Sprites/BlockingAnimation(1).png"), 2, 0.7)

	animationFSM := NewAnimationFSM()
	animationFSM.AddAnimation(walkAnimation)
	animationFSM.AddAnimation(jumpAnimation)
	animationFSM.AddAnimation(idleAnimation)
	animationFSM.AddAnimation(punchAnimation)
	animationFSM.AddAnimation(blockAnimation)
	animationFSM.ChangeAnimationState("idle")

	player1 = Creature{
		Pos:                rl.NewVector2(380, 100),
		Vel:                rl.NewVector2(0, 0),
		Size:               200,
		Color:              rl.Red,
		Speed:              300,
		Direction:          1,
		IsGrounded:         false,
		JumpPower:          600,
		MaxHealth:          100,
		Health:             100,
		AttackDuration:     defaultAttackDuration,
		AttackActiveStart:  defaultAttackActiveStart,
		AttackActiveEnd:    defaultAttackActiveEnd,
		AttackDamage:       5,
		AttackHitboxOffset: rl.NewVector2(150, 50), // Offset from top-left (Pos) when facing right
		AttackHitboxSize:   rl.NewVector2(100, 80), // Width/Height of punch hitbox
		AttackCooldown:     defaultAttackCooldown,
		AnimationFSM:       animationFSM,
	}

	//second player
	animationFSM2 := NewAnimationFSM()
	animationFSM2.AddAnimation(walkAnimation)
	animationFSM2.AddAnimation(jumpAnimation)
	animationFSM2.AddAnimation(idleAnimation)
	animationFSM2.AddAnimation(punchAnimation)
	animationFSM2.AddAnimation(blockAnimation)
	animationFSM2.ChangeAnimationState("idle")

	player2 = Creature{
		Pos:                rl.NewVector2(float32(rl.GetScreenWidth()-400), 100),
		Vel:                rl.NewVector2(0, 0),
		Size:               200,
		Color:              rl.Red,
		Speed:              300,
		Direction:          -1,
		IsGrounded:         false,
		JumpPower:          600,
		MaxHealth:          100,
		Health:             100,
		AttackDuration:     defaultAttackDuration,
		AttackActiveStart:  defaultAttackActiveStart,
		AttackActiveEnd:    defaultAttackActiveEnd,
		AttackDamage:       5,
		AttackHitboxOffset: rl.NewVector2(150, 50), // Offset from top-left (Pos) when facing right
		AttackHitboxSize:   rl.NewVector2(100, 80), // Width/Height of punch hitbox
		AttackCooldown:     defaultAttackCooldown,
		AnimationFSM:       animationFSM2,
	}

	//health bars
	hbWidth := float32(rl.GetScreenWidth() / 4)
	hbHeight := float32(40)
	hbMargin := float32(20)
	player1HealthBar = NewHealthBar(
		hbMargin, hbMargin, hbWidth, hbHeight, rl.Green, rl.DarkGray, rl.Black, 2,
	)
	player2HealthBar = NewHealthBar(
		float32(rl.GetScreenWidth())-hbWidth-hbMargin, hbMargin, hbWidth, hbHeight, rl.Green, rl.DarkGray, rl.Black, 2,
	)

	gravity := rl.NewVector2(0, 1000)

	ResetGame()

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Blue)

		if currentGameState == Playing {
			//player 1 controls
			player1VelX := float32(0.0)
			if rl.IsKeyDown(rl.KeyD) {
				player1VelX = 1
			}
			if rl.IsKeyDown(rl.KeyA) {
				player1VelX = -1
			}
			player1.Move(player1VelX)
			if rl.IsKeyPressed(rl.KeyW) {
				player1.Jump()
			}
			if rl.IsKeyDown(rl.KeyS) {
				player1.StartBlocking()
			}
			if rl.IsKeyPressed(rl.KeyE) {
				player1.Attack()
			}

			//player 2 controls
			player2VelX := float32(0.0)
			if rl.IsKeyDown(rl.KeyRight) {
				player2VelX = 1
			}
			if rl.IsKeyDown(rl.KeyLeft) {
				player2VelX = -1
			}
			player2.Move(player2VelX)
			if rl.IsKeyPressed(rl.KeyUp) {
				player2.Jump()
			}
			if rl.IsKeyDown(rl.KeyDown) {
				player2.StartBlocking()
			}
			if rl.IsKeyPressed(rl.KeyRightShift) {
				player2.Attack()
			}
		} else if currentGameState == GameOver {
			if rl.IsKeyPressed(rl.KeyR) {
				ResetGame()
			}
		}

		// --- Updates (Only when Playing) ---
		if currentGameState == Playing {
			// Physics
			player1.ApplyGravity(gravity)
			player2.ApplyGravity(gravity)

			// Update Creatures (Pass block key status)
			p1Blocking := rl.IsKeyDown(rl.KeyS)
			p2Blocking := rl.IsKeyDown(rl.KeyDown)
			player1.UpdateCreature(platformY, p1Blocking)
			player2.UpdateCreature(platformY, p2Blocking)

			// --- Hit Detection ---
			p1AttackBox, p1IsActive := player1.GetAttackHitbox()
			p2BodyBox := player2.GetHitbox()
			if p1IsActive && rl.CheckCollisionRecs(p1AttackBox, p2BodyBox) {
				if !player2.IsBlocking {
					player2.TakeDamage(player1.AttackDamage)
					player1.CanDamage = false
				} else {
					player1.CanDamage = false
				}
			}

			p2AttackBox, p2IsActive := player2.GetAttackHitbox()
			p1BodyBox := player1.GetHitbox()
			if p2IsActive && rl.CheckCollisionRecs(p2AttackBox, p1BodyBox) {
				if !player1.IsBlocking {
					player1.TakeDamage(player2.AttackDamage)
					player2.CanDamage = false
				} else {
					player2.CanDamage = false
					// Block feedback
				}
			}

			// --- Check for Game Over ---
			if player1.Health <= 0 {
				currentGameState = GameOver
				winner = 2 // Player 2 wins
			} else if player2.Health <= 0 {
				currentGameState = GameOver
				winner = 1 // Player 1 wins
			}
		}

		//player 1
		//player1.Move(player1VelX)
		//player1.ApplyGravity(gravity)
		//player1.UpdateCreature(platformY)

		//player 2
		//player2.Move(player2VelX)
		//player2.ApplyGravity(gravity)
		//player2.UpdateCreature(platformY)

		//stage
		rl.DrawRectangleRec(platformRect, platformColor)

		//draw players
		player1.DrawCreature()
		player2.DrawCreature()

		//draw health bars
		player1HealthBar.Draw(player1.Health, player1.MaxHealth)
		player2HealthBar.Draw(player2.Health, player2.MaxHealth)

		if currentGameState == GameOver {
			rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), rl.Fade(rl.Black, 0.6)) // Darken screen
			winText := fmt.Sprintf("Player %d Wins!", winner)
			restartText := "Press [R] to Restart"

			winTextWidth := rl.MeasureText(winText, 60)
			restartTextWidth := rl.MeasureText(restartText, 30)

			rl.DrawText(winText, (int32(rl.GetScreenWidth())-winTextWidth)/2, int32(rl.GetScreenHeight())/2-40, 60, rl.Yellow)
			rl.DrawText(restartText, (int32(rl.GetScreenWidth())-restartTextWidth)/2, int32(rl.GetScreenHeight())/2+40, 30, rl.White)
		}

		rl.EndDrawing()
	}
}
