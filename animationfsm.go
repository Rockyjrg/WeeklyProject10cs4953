package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type AnimationFSM struct {
	Animations  map[string]Animation
	CurrentAnim Animation
}

func NewAnimationFSM() AnimationFSM {
	return AnimationFSM{Animations: make(map[string]Animation)}
}

func (a *AnimationFSM) AddAnimation(anim Animation) {
	a.Animations[anim.Name] = anim
}

func (a *AnimationFSM) ChangeAnimationState(newAnim string) {
	if newAnim == a.CurrentAnim.Name {
		return
	}
	_, ok := a.Animations[newAnim]
	if !ok {
		return
	}

	a.CurrentAnim = a.Animations[newAnim]
	a.CurrentAnim.Reset()
}

func (a *AnimationFSM) DrawWithFSM(pos rl.Vector2, size float32, direction float32) {
	a.CurrentAnim.UpdateTime()
	a.CurrentAnim.DrawAnimation(pos, size, direction)
}
