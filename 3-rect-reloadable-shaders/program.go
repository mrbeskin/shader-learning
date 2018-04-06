package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Program struct {
	shaders   *Shaders
	glProgram uint32
}

func NewProgram(shaders *Shaders) *Program {
	program := &Program{
		shaders: shaders,
	}
	program.LoadShaders()
}

func (p *Program) LoadShaders() {
	vert, frag := p.Shaders.GetSource()
	p.glProgram = p.attachShaders(vert, frag)
}

func (p *Program) UpdateShaders() {
	updated, vert, frag := p.shaders.GetUpdatedSource()
	reload := false
	if updated {
		gl.DeleteProgram(p.glProgram)
		gl.Flush()
		p.glProgram = p.attachShaders(vert, frag)
	}
}

func (p *Program) attachShaders(vert string, frag string) {
	vertexShader, err := compileShader(vert, gl.VERTEX_SHADER)
	fragmentShader, err := compileShader(frag, gl.FRAGMENT_SHADER)

	p.glProgram = gl.CreateProgram()

	gl.AttachShader(p.glProgram, vertexShader)
	gl.AttachShader(p.glProgram, fragmentShader)
	gl.LinkProgram(p.glProgram)

	var success int32
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &success)
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
