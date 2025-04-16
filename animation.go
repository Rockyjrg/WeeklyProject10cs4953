package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Animation struct {
	Name         string
	SpriteSheet  rl.Texture2D
	MaxIndex     int
	CurrentIndex int
	Timer        float32
	SwitchTime   float32
	Loop         bool
}

func (a *Animation) Reset() {
	a.Timer = 0
	a.CurrentIndex = 0
}

func NewAnimation(newName string, newSheet rl.Texture2D, spriteNum int, newTime float32) Animation {
	newAnimation := Animation{
		Name:         newName,
		SpriteSheet:  newSheet,
		MaxIndex:     spriteNum - 1,
		CurrentIndex: 0,
		Timer:        0,
		SwitchTime:   newTime,
		Loop:         true,
	}
	return newAnimation
}

func (a *Animation) UpdateTime() {

	a.Timer += rl.GetFrameTime()
	if a.Timer > a.SwitchTime {
		a.Timer = 0
		a.CurrentIndex++
	}

	if a.CurrentIndex > a.MaxIndex {
		if a.Loop {
			a.CurrentIndex = 0
		} else {
			a.CurrentIndex = a.MaxIndex
		}
	}
}

func (a Animation) DrawAnimation(pos rl.Vector2, size float32, direction float32) {
	sourceRect := rl.NewRectangle(float32(32*a.CurrentIndex), 0, 32*direction, 32)
	destRect := rl.NewRectangle(pos.X, pos.Y, size, size)
	rl.DrawTexturePro(a.SpriteSheet, sourceRect, destRect, rl.Vector2Zero(), 0, rl.White)
}
