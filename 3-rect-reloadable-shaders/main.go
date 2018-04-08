package main

import (
	"fmt"
	"go/build"
	_ "image/png"
	"log"
	"os"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var (
	texture uint32
)

const (
	sizeof_float32 = 4
)

var (
	VAO uint32
	VBO uint32
	EBO uint32
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {

	window := initGlfwWindow()
	shaders := NewShaders("shader.vert", "shader.frag")
	program := NewProgram(shaders)
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)
	gl.UseProgram(program.glProgram)
	defer func() { destroyScene() }()

	initBuffers()

	for !(window.ShouldClose()) {
		program.UpdateShaders()
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.UseProgram(program.glProgram)
		gl.BindVertexArray(VAO)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
		gl.BindVertexArray(0)
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

var vertices = []float32{
	// indexed to be an EBO
	0.5, 0.5, 0.0, // top right
	0.5, -0.5, 0.0, // bottom right
	-0.5, -0.5, 0.0, // bottom left
	-0.5, 0.5, 0.0, // top left
}

var indices = []uint32{
	0, 1, 3,
	1, 2, 3,
}

func destroyScene() {
	defer glfw.Terminate()
	gl.Flush()
}

func initGlfwWindow() *glfw.Window {
	// initialize glfw window
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(800, 600, "hello-rectangle", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	window.SetFramebufferSizeCallback(glfw.FramebufferSizeCallback(fbcallback))
	return window
}

func initBuffers() {
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)
	gl.GenBuffers(1, &EBO)

	// bind Vertex Array first
	gl.BindVertexArray(VAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*sizeof_float32, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*sizeof_float32, gl.Ptr(indices), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*sizeof_float32, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	gl.BindVertexArray(0)
	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
}

func fbcallback(w *glfw.Window, width int, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
}

func init() {
	dir, err := importPathToDir("github.com/mrbeskin/shader-learning/3-rect-reloadable-shaders")
	if err != nil {
		log.Fatalln("could not locate assets on GOPATH:", err)
	}
	err = os.Chdir(dir)
	if err != nil {
		log.Panicln("os.Chdir:", err)
	}
}

func importPathToDir(importPath string) (string, error) {
	p, err := build.Import(importPath, "", build.FindOnly)
	if err != nil {
		return "", err
	}
	return p.Dir, nil
}
