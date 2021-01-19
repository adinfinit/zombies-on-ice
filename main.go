package main

import (
	"flag"
	"log"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func init() { runtime.LockOSThread() }

func main() {
	flag.Parse()

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.True)

	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	glfw.WindowHint(glfw.Samples, 4)

	window, err := glfw.CreateWindow(800, 600, "Zombies on Ice", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if runtime.GOOS == "windows" {
		window.SetPos(32, 64)
	}

	if err := gl.Init(); err != nil {
		panic(err)
	}

	game := NewGame()
	game.Clock = glfw.GetTime()
	for !window.ShouldClose() {
		if window.GetKey(glfw.KeyF10) == glfw.Press {
			game.Unload()
			game = NewGame()
		}

		game.Update(window, glfw.GetTime())
		game.Render(window)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
