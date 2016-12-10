package main

import (
	"flag"
	"log"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

func init() { runtime.LockOSThread() }

func main() {
	flag.Parse()

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.Visible, glfw.False) // do not steal focus

	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	glfw.WindowHint(glfw.Samples, 4)

	window, err := glfw.CreateWindow(800, 600, "Zombie Room", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	window.Restore() // do not steal focus

	window.SetPos(32, 64)

	if err := gl.Init(); err != nil {
		panic(err)
	}

	game := NewGame()

	last := float32(glfw.GetTime())
	for !window.ShouldClose() {
		if window.GetKey(glfw.KeyF10) == glfw.Press {
			game.Unload()
			game = NewGame()
		}

		now := float32(glfw.GetTime())
		game.Update(window, now-last)
		last = now

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
