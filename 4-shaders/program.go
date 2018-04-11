package main

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"strings"
)

type Program struct {
	shaders   *Shaders
	glProgram uint32
}

func NewProgram(shaders *Shaders) *Program {
	if err := gl.Init(); err != nil {
		check("initializing gl", err)
	}
	glP := gl.CreateProgram()
	program := &Program{
		shaders:   shaders,
		glProgram: glP,
	}
	program.LoadShaders()
	return program
}

func (p *Program) LoadShaders() {
	vert, frag := p.shaders.GetSource()
	p.attachShaders(vert, frag)
	gl.UseProgram(p.glProgram)
}

func (p *Program) UpdateShaders() {
	updated, vert, frag := p.shaders.GetUpdatedSource()
	if updated {
		glP := gl.CreateProgram()
		p.glProgram = glP
		p.attachShaders(vert, frag)
	}
}

func (p *Program) attachShaders(vert string, frag string) {
	vertexShader, err := compileShader(vert, gl.VERTEX_SHADER)
	check("attaching vertex shader", err)
	fragmentShader, err := compileShader(frag, gl.FRAGMENT_SHADER)
	check("attaching fragment shader", err)

	gl.AttachShader(p.glProgram, vertexShader)
	gl.AttachShader(p.glProgram, fragmentShader)
	gl.LinkProgram(p.glProgram)

	var success int32
	gl.GetProgramiv(p.glProgram, gl.LINK_STATUS, &success)
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
