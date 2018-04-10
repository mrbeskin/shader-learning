#version 330 core
out vec4 FragColor;

uniform vec4 newColor;

void main() {
    FragColor = newColor;
}
