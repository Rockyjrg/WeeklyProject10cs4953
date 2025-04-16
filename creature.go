package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Creature struct {
	Pos       rl.Vector2
	Vel       rl.Vector2
	Size      float32
	Color     rl.Color
	Speed     float32
	Direction float32
	AnimationFSM
	IsGrounded bool
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

func (c *Creature) UpdateCreature(groundY float32) {
	c.Pos = rl.Vector2Add(c.Pos, rl.Vector2Scale(c.Vel, rl.GetFrameTime()))

	creatureBottom := c.Pos.Y + c.Size
	if creatureBottom >= groundY && c.Vel.Y >= 0 {
		c.Vel.Y = 0
		c.Pos.Y = groundY - c.Size
		c.IsGrounded = true
	} else {
		c.IsGrounded = false
	}

	if !c.IsGrounded {
		c.AnimationFSM.ChangeAnimationState("jump")
	}

	if c.Vel.X == 0 && c.IsGrounded {
		c.AnimationFSM.ChangeAnimationState("idle")
	} else if c.Vel.X != 0 && c.IsGrounded {
		c.AnimationFSM.ChangeAnimationState("walk")
	}
}
