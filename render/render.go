package render

import (
	"math"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/loov/zombies-on-ice/g"
)

type State struct {
	Background g.Color
	Textures   *Textures
}

func NewState() *State {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	return &State{
		Background: g.Black,
		Textures:   NewTextures(),
	}
}

func (state *State) BeginFrame(size g.V2) {
	state.Textures.BeginFrame()

	gl.ClearColor(state.Background.Float())
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.Enable(gl.MULTISAMPLE)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.ALPHA_TEST)

	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	gl.Viewport(0, 0, int32(size.X), int32(size.Y))
}

func (state *State) Ortho(bounds g.Rect, near, far float32) {
	gl.Ortho(
		float64(bounds.Min.X),
		float64(bounds.Max.X),
		float64(bounds.Min.Y),
		float64(bounds.Max.Y),
		1, -1000)
}

func (state *State) EndFrame() {
	state.Textures.EndFrame()
}

func (state *State) PushMatrix() {
	gl.PushMatrix()
}

func (state *State) PopMatrix() {
	gl.PopMatrix()
}

func (state *State) Translate(p g.V2) {
	gl.Translatef(p.X, p.Y, 0.0)
}

func (state *State) Rotate(radians float32) {
	gl.Rotatef(radians*180.0/math.Pi, 0, 0, -1)
}
