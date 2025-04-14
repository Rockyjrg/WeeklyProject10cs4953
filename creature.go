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
