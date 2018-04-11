package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// Shaders represents the shaders used by the gl program
type Shaders struct {
	vert *Shader
	frag *Shader
}

// NewShaders returns an object containing shaders from the paths listed
func NewShaders(vertPath string, fragPath string) *Shaders {
	return &Shaders{
		vert: NewShader(vertPath),
		frag: NewShader(fragPath),
	}
}

// Shader contains a path to the shader source and a timestamp of its last
// modification so that it may be updated
type Shader struct {
	ModTime time.Time
	Path    string
}

// NewShader returns a Shader
func NewShader(path string) *Shader {
	fileinfo, err := os.Stat(path)
	check("new shader; getting file info", err)
	return &Shader{
		Path:    path,
		ModTime: fileinfo.ModTime(),
	}
}

// GetSource returns the source string for each shader
func (ss *Shaders) GetSource() (vert string, frag string) {
	vert = readShaderFile(ss.vert.Path)
	frag = readShaderFile(ss.frag.Path)
	return
}

// GetUpdatedSource returns the source string for each shader
// that has been modified.
func (ss *Shaders) GetUpdatedSource() (updated bool, vert string, frag string) {
	if ss.vert.Update() || ss.frag.Update() {
		updated = true
		vert = readShaderFile(ss.vert.Path)
		frag = readShaderFile(ss.frag.Path)
	}
	return
}

// Update checks if the shader can be updated and sets the latest ModTime
// then returns a bool representing whether or not it was updated
func (s *Shader) Update() bool {
	fileinfo, err := os.Stat(s.Path)
	check("stat on shader file", err)
	if fileinfo.ModTime().After(s.ModTime) {
		return true
	}
	return false
}

// PRIVATE UTILS

func check(msg string, err error) {
	if err != nil {
		panic(fmt.Sprintf("%s; error:%v", msg, err))
	}
}

// reads a shader file and returns the source
func readShaderFile(path string) string {
	shaderBuf, err := ioutil.ReadFile(path)
	check("reading shader file", err)
	code := string(shaderBuf)
	code += "\x00"
	return code
}
