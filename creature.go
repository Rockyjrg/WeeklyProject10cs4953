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
	Texture   rl.Texture2D
}

func (c *Creature) NewCreature() Creature {

}
