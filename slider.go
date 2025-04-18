package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type HealthBar struct {
	Rect         rl.Rectangle
	FillColor    rl.Color
	BackColor    rl.Color
	OutlineColor rl.Color
	BorderSize   float32
}

func NewHealthBar(x, y, width, height float32, fill, back, outline rl.Color, border float32) HealthBar {
	return HealthBar{
		Rect:         rl.NewRectangle(x, y, width, height),
		FillColor:    fill,
		BackColor:    back,
		OutlineColor: outline,
		BorderSize:   border,
	}
}

func (hb *HealthBar) Draw(currentHealth, maxHealth float32) {
	currentHealth = float32(math.Max(0, float64(currentHealth)))
	maxHealth = float32(math.Max(1, float64(maxHealth)))

	//calculate ratio of current health
	healthRatio := currentHealth / maxHealth
	fillWidth := hb.Rect.Width * healthRatio

	//draw background
	rl.DrawRectangleRec(hb.Rect, hb.BackColor)

	//draw current health fill
	fillRec := rl.NewRectangle(hb.Rect.X, hb.Rect.Y, fillWidth, hb.Rect.Height)
	rl.DrawRectangleRec(fillRec, hb.FillColor)

	//outline
	if hb.BorderSize > 0 {
		rl.DrawRectangleLinesEx(hb.Rect, hb.BorderSize, hb.OutlineColor)
	}
}
