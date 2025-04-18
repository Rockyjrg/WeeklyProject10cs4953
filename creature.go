package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Creature struct {
	Pos        rl.Vector2
	Vel        rl.Vector2
	Size       float32
	Color      rl.Color
	Speed      float32
	Direction  float32
	IsGrounded bool
	JumpPower  float32
	Health     float32
	MaxHealth  float32
	AnimationFSM

	//attack properties?
	IsAttacking        bool
	IsBlocking         bool
	AttackTimer        float32
	AttackDuration     float32
	AttackActiveStart  float32
	AttackActiveEnd    float32
	AttackDamage       float32
	AttackHitboxOffset rl.Vector2
	AttackHitboxSize   rl.Vector2
	CanDamage          bool

	AttackCooldownTimer float32
	AttackCooldown      float32
}

const defaultAttackDuration = 0.5
const defaultAttackActiveStart = 0.1
const defaultAttackActiveEnd = 0.35
const defaultAttackCooldown = 0.6

func (c *Creature) GetHitbox() rl.Rectangle {
	return rl.NewRectangle(c.Pos.X, c.Pos.Y, c.Size, c.Size)
}

func (c *Creature) Attack() {
	if !c.IsAttacking && !c.IsBlocking && c.AttackCooldownTimer <= 0 {
		c.IsAttacking = true
		c.AttackTimer = 0
		c.AnimationFSM.ChangeAnimationState("punch")
		c.AttackCooldownTimer = c.AttackCooldown
		c.CanDamage = true
	}
}

func (c *Creature) StartBlocking() {
	if !c.IsAttacking && c.IsGrounded {
		if !c.IsBlocking {
			c.IsBlocking = true
			c.AnimationFSM.ChangeAnimationState("block")
			c.Vel.X = 0
		}
	}
}

func (c *Creature) StopBlocking() {
	if c.IsBlocking {
		c.IsBlocking = false
	}
}

func (c *Creature) GetAttackHitbox() (rl.Rectangle, bool) {
	isActive := c.IsAttacking &&
		c.AttackTimer >= c.AttackActiveStart &&
		c.AttackTimer < c.AttackActiveEnd

	offsetX := c.AttackHitboxOffset.X
	if c.Direction < 0 {
		offsetX = -(c.AttackHitboxOffset.X + c.AttackHitboxSize.X)
	}

	hitboxX := c.Pos.X + offsetX
	hitboxY := c.Pos.Y + c.AttackHitboxOffset.Y

	attackBox := rl.NewRectangle(hitboxX, hitboxY, c.AttackHitboxSize.X, c.AttackHitboxSize.Y)

	return attackBox, isActive
}

func (c *Creature) TakeDamage(damage float32) {
	if c.Health > 0 {
		c.Health -= damage
		if c.Health < 0 {
			c.Health = 0
		}
	}
}

func (c *Creature) ApplyGravity(g rl.Vector2) {
	c.Vel = rl.Vector2Add(c.Vel, rl.Vector2Scale(g, rl.GetFrameTime()))
}

func (c *Creature) Move(x float32) {
	if x != 0 {
		c.Direction = x
	}
	c.Vel.X = x * c.Speed
}

func (c *Creature) DrawCreature() {
	c.AnimationFSM.DrawWithFSM(c.Pos, c.Size, c.Direction)
}

func (c *Creature) UpdateCreature(groundY float32, blockKeyDown bool) {

	dt := rl.GetFrameTime()
	if c.AttackCooldownTimer > 0 {
		c.AttackCooldownTimer -= dt
	}
	if c.IsAttacking {
		c.AttackTimer += dt
		if c.AttackTimer >= c.AttackDuration {
			c.IsAttacking = false
			c.CanDamage = false
			c.AttackTimer = 0
		}
	}

	if c.IsBlocking && (!blockKeyDown || c.IsAttacking) {
		c.StopBlocking()
	}

	if !(c.IsGrounded && c.Vel.Y >= 0) {
		c.Pos.Y += c.Vel.Y * dt
	}
	c.Pos.X += c.Vel.X * dt

	creatureBottom := c.Pos.Y + c.Size
	if creatureBottom >= groundY && c.Vel.Y >= 0 {
		c.Vel.Y = 0
		c.Pos.Y = groundY - c.Size
		if !c.IsGrounded { // Just landed
			c.IsGrounded = true
		}
	} else if creatureBottom < groundY { // Check if above ground
		c.IsGrounded = false
	}
	if c.IsAttacking {
		c.AnimationFSM.ChangeAnimationState("punch")
	} else if c.IsBlocking {
		c.AnimationFSM.ChangeAnimationState("block")
	} else if !c.IsGrounded {
		c.AnimationFSM.ChangeAnimationState("jump")
	} else { // Grounded and not attacking/blocking
		if c.Vel.X == 0 {
			c.AnimationFSM.ChangeAnimationState("idle")
		} else {
			c.AnimationFSM.ChangeAnimationState("walk")
		}
	}

	// if c.IsAttacking {
	// 	c.AnimationFSM.ChangeAnimationState("punch")
	// } else if
}

func (c *Creature) Jump() {
	if c.IsGrounded {
		c.Vel.Y = -c.JumpPower
		c.IsGrounded = false
	}
}
