package main

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"io/ioutil"
	"strings"
)

type Shader struct {
	ID uint32
}

func NewShader(fragPath string, vertPath string) *Shader {
	vert := readShaderFile(vertPath)
	frag := readShaderFile(fragPath)
	id := gl.CreateProgram()
	shader := &Shader{
		ID: id,
	}
	shader.attachShaders(vert, frag)
	gl.UseProgram(shader.ID)
	return shader
}

func (s *Shader) attachShaders(vert string, frag string) {
	vertexShader, err := compileShader(vert, gl.VERTEX_SHADER)
	check("attaching vertex shader", err)
	fragmentShader, err := compileShader(frag, gl.FRAGMENT_SHADER)
	check("attaching fragment shader", err)

	gl.AttachShader(s.ID, vertexShader)
	gl.AttachShader(s.ID, fragmentShader)
	gl.LinkProgram(s.ID)

	var success int32
	gl.GetProgramiv(s.ID, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		panic("could not link shader program")
	}
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	// compile shader
	csources, free := gl.Strs(source)
	defer free()
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
	return shader, nil
}

func readShaderFile(path string) string {
	shaderBuf, err := ioutil.ReadFile(path)
	check("reading shader file", err)
	return string(shaderBuf) + "\x00"
}

func check(msg string, err error) {
	if err != nil {
		panic(fmt.Sprintf("%s; error:%v", msg, err))
	}
}
