package main

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/loov/zombies-on-ice/g"
)

type Controller struct {
	ID int

	Connected bool
	Sticky    bool
	Updater   ControllerUpdater

	DPad        DPad
	Start, Back bool
	A, B, X, Y  bool

	Left  Analog
	Right Analog
}

type Analog struct {
	Direction g.V2
	Hold      bool
	Trigger   bool
}

type ControllerUpdater interface {
	Update(input *Controller, window *glfw.Window)
}

type DPad struct {
	Up    bool
	Down  bool
	Left  bool
	Right bool
}

func (dpad DPad) Direction() (r g.V2) {
	if dpad.Down {
		r.Y -= 1
	}
	if dpad.Up {
		r.Y += 1
	}
	if dpad.Left {
		r.X -= 1
	}
	if dpad.Right {
		r.X += 1
	}
	return
}

func (a *Controller) Merge(b *Controller) {
	if !b.Connected {
		return
	}

	a.Connected = a.Connected || b.Connected
	a.Sticky = a.Sticky || b.Sticky

	a.DPad.Up = a.DPad.Up || b.DPad.Up
	a.DPad.Down = a.DPad.Down || b.DPad.Down
	a.DPad.Left = a.DPad.Left || b.DPad.Left
	a.DPad.Right = a.DPad.Right || b.DPad.Right

	a.Start = a.Start || b.Start
	a.Back = a.Back || b.Back
	a.A = a.A || b.A
	a.B = a.B || b.B
	a.X = a.X || b.X
	a.Y = a.Y || b.Y
}

func (a *Controller) Active() bool {
	return a.DPad.Up || a.DPad.Down || a.DPad.Left || a.DPad.Right ||
		a.Start || a.A || a.B || a.X || a.Y
}

type Keyboard struct {
	Connected bool

	Up, Down, Left, Right glfw.Key
	Start, Back           glfw.Key
	A, B, X, Y            glfw.Key
}

func (key *Keyboard) Update(input *Controller, window *glfw.Window) {
	getkey := func(button glfw.Key) bool {
		if button == glfw.KeyUnknown || button == 0 {
			return false
		}

		return window.GetKey(button) == glfw.Press
	}

	input.DPad.Up = getkey(key.Up)
	input.DPad.Down = getkey(key.Down)
	input.DPad.Left = getkey(key.Left)
	input.DPad.Right = getkey(key.Right)

	input.Start = getkey(key.Start)
	input.Back = getkey(key.Back)

	input.A = getkey(key.A)
	input.B = getkey(key.B)
	input.X = getkey(key.X)
	input.Y = getkey(key.Y)

	input.Left = Analog{
		Direction: input.DPad.Direction(),
		Hold:      input.A,
		Trigger:   input.B,
	}

	if input.Active() {
		key.Connected = true
	}
	if input.Back {
		key.Connected = false
	}

	input.Connected = key.Connected
}

type Gamepad struct {
	Id       glfw.Joystick
	DeadZone float32
}

func (gamepad Gamepad) Update(input *Controller, window *glfw.Window) {
	// clear state
	*input = Controller{ID: input.ID, Updater: input.Updater}

	axes := glfw.GetJoystickAxes(gamepad.Id)
	buttons := glfw.GetJoystickButtons(gamepad.Id)

	input.Connected = len(axes) > 0 && len(buttons) > 0
	if !input.Connected {
		return
	}

	button := func(i int) bool {
		if i < len(buttons) {
			return buttons[i] == 1
		}
		return false
	}

	input.DPad.Up = button(10)
	input.DPad.Right = button(11)
	input.DPad.Down = button(12)
	input.DPad.Left = button(13)

	input.A = button(0)
	input.B = button(1)
	input.X = button(2)
	input.Y = button(3)

	input.Back = button(6)
	input.Start = button(7)

	input.Left.Direction = g.V2{ // left thumb
		X: g.ApplyDeadZone(axes[0], gamepad.DeadZone),
		Y: -g.ApplyDeadZone(axes[1], gamepad.DeadZone),
	}
	input.Left.Hold = button(8)    // left thumb pressed
	input.Left.Trigger = button(4) // left trigger

	input.Right.Direction = g.V2{ // right thumb
		X: g.ApplyDeadZone(axes[4], gamepad.DeadZone),
		Y: -g.ApplyDeadZone(axes[3], gamepad.DeadZone),
	}
	input.Right.Hold = button(9)    // right thumb pressed
	input.Right.Trigger = button(5) // right trigger
}
