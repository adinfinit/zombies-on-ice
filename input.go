package main

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/loov/zombieroom/g"
)

type Input struct {
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

type Controller struct {
	ID int

	Connected bool
	Sticky    bool
	Updater   ControllerUpdater

	DPad        DPad
	Start, Back bool
	A, B, X, Y  bool

	Inputs []Input
}

type Controllers []*Controller

func (controllers Controllers) Merged() *Controller {
	all := &Controller{}
	for _, controller := range controllers {
		all.Merge(controller)
	}
	return all
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

func (controllers Controllers) Update(window *glfw.Window) {
	for _, controller := range controllers {
		wasActive := controller.Active()
		controller.Updater.Update(controller, window)
		controller.Sticky = wasActive
	}
}

func (controllers Controllers) FindInput(id, input int) Input {
	for _, controller := range controllers {
		if controller.ID == id {
			if input < len(controller.Inputs) {
				return controller.Inputs[input]
			}
			break
		}
	}
	return Input{}
}

type Keyboard struct {
	Connected bool

	Up, Down, Left, Right glfw.Key
	Start, Back           glfw.Key
	A, B, X, Y            glfw.Key
}

func (key *Keyboard) Update(input *Controller, window *glfw.Window) {
	input.DPad.Up = window.GetKey(key.Up) == glfw.Press
	input.DPad.Down = window.GetKey(key.Down) == glfw.Press
	input.DPad.Left = window.GetKey(key.Left) == glfw.Press
	input.DPad.Right = window.GetKey(key.Right) == glfw.Press

	input.Start = window.GetKey(key.Start) == glfw.Press
	input.Back = window.GetKey(key.Back) == glfw.Press

	input.A = window.GetKey(key.A) == glfw.Press
	input.B = window.GetKey(key.B) == glfw.Press
	input.X = window.GetKey(key.X) == glfw.Press
	input.Y = window.GetKey(key.Y) == glfw.Press

	input.Inputs = []Input{
		{
			Direction: input.DPad.Direction(),
			Hold:      input.A,
			Trigger:   input.B,
		},
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
	input.Inputs = []Input{{}, {}}

	axes := glfw.GetJoystickAxes(gamepad.Id)
	buttons := glfw.GetJoystickButtons(gamepad.Id)

	input.Connected = len(axes) > 0 && len(buttons) > 0
	if !input.Connected {
		return
	}

	input.DPad.Up = buttons[10] == 1
	input.DPad.Right = buttons[11] == 1
	input.DPad.Down = buttons[12] == 1
	input.DPad.Left = buttons[13] == 1

	input.A = buttons[0] == 1
	input.B = buttons[1] == 1
	input.X = buttons[2] == 1
	input.Y = buttons[3] == 1

	input.Back = buttons[6] == 1
	input.Start = buttons[7] == 1

	// controls for diciple 0
	input.Inputs[0].Direction = g.V2{ // left thumb
		X: g.ApplyDeadZone(axes[0], gamepad.DeadZone),
		Y: -g.ApplyDeadZone(axes[1], gamepad.DeadZone),
	}
	input.Inputs[0].Hold = buttons[8] == 1    // left thumb pressed
	input.Inputs[0].Trigger = buttons[4] == 1 // left trigger

	// controls for diciple 1
	input.Inputs[1].Direction = g.V2{ // right thumb
		X: g.ApplyDeadZone(axes[4], gamepad.DeadZone),
		Y: -g.ApplyDeadZone(axes[3], gamepad.DeadZone),
	}
	input.Inputs[1].Hold = buttons[9] == 1    // right thumb pressed
	input.Inputs[1].Trigger = buttons[5] == 1 // right trigger
}

var (
	Keyboard_0 = Keyboard{
		Up:    glfw.KeyUp,
		Down:  glfw.KeyDown,
		Left:  glfw.KeyLeft,
		Right: glfw.KeyRight,

		Start: glfw.KeyEnter,
		Back:  glfw.KeyI,

		A: glfw.KeyO,
		B: glfw.KeyP,
	}

	Keyboard_1 = Keyboard{
		Up:    glfw.KeyW,
		Down:  glfw.KeyS,
		Left:  glfw.KeyA,
		Right: glfw.KeyD,

		Start: glfw.KeyEnter,
		Back:  glfw.KeyR,

		A: glfw.KeyTab,
		B: glfw.KeyQ,
	}

	Gamepad_0 = Gamepad{glfw.Joystick1, 0.05}
	Gamepad_1 = Gamepad{glfw.Joystick2, 0.05}
)
