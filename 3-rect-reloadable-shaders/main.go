package main

import (
	"fmt"
	"go/build"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var (
	texture uint32
)

const (
	sizeof_float32 = 4
)

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

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

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
	// init Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	vertexShaderSource, err := readShaderFile("shader.vert")
	fragmentShaderSource, err := readShaderFile("shader.frag")

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	program := setupProgram(vertexShaderSource, fragmentShaderSource)
	gl.UseProgram(program)

	var VAO, EBO, VBO uint32

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

	// render loop
	for !window.ShouldClose() {
		reload := false

		newVertexShaderSource, err := readShaderFile("shader.vert")
		if err != nil {
			log.Printf("vertex shader invalid: %v", err)
		} else {
			vertexShaderSource = newVertexShaderSource
			reload = true
		}

		newFragmentShaderSource, err := readShaderFile("shader.frag")
		if err != nil {
			log.Printf("fragment shader invalid: %v", err)
		} else {
			fragmentShaderSource = newFragmentShaderSource
			reload = true
		}
		if reload {
			program = setupProgram(vertexShaderSource, fragmentShaderSource)
		}
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.UseProgram(program)
		gl.BindVertexArray(VAO)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
		gl.BindVertexArray(0)
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func destroyScene() {
	gl.Flush()
}

func setupProgram(vertexShaderSource string, fragmentShaderSource string) uint32 {
	// vertex shader
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	// fragment shader
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	shaderProgram := gl.CreateProgram()

	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)

	var success int32
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		panic("could not link shader program")
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return shaderProgram
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	// compile shader
	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	gl.CompileShader(shader)
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)

	// check failure and log if necessary
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}
	free()
	return shader, nil
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

// reads a shader from a file and returns a string representation
// that is usable in opengl programs
func readShaderFile(shaderPath string) (string, error) {
	shaderBuf, err := ioutil.ReadFile(shaderPath)
	if err != nil {
		return "", fmt.Errorf("shader %q unable to be read: %v", shaderPath, err)
	}
	return string(shaderBuf), nil
}
