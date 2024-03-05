package engine

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Framebuffer struct {
	id            uint32
	Texture       uint32
	Width, Height int
}

// NewFramebuffer creates a new framebuffer with the specified width and height.
func NewFramebuffer(width, height int) (*Framebuffer, error) {
	framebuffer := &Framebuffer{
		Width:  width,
		Height: height,
	}

	// Create framebuffer
	gl.GenFramebuffers(1, &framebuffer.id)
	gl.BindFramebuffer(gl.FRAMEBUFFER, framebuffer.id)

	// Create texture attachment for framebuffer
	gl.GenTextures(1, &framebuffer.Texture)
	gl.BindTexture(gl.TEXTURE_2D, framebuffer.Texture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, int32(width), int32(height), 0, gl.RGB, gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, framebuffer.Texture, 0)

	// Check framebuffer completeness
	if status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); status != gl.FRAMEBUFFER_COMPLETE {
		return nil, fmt.Errorf("Framebuffer incomplete: %x", status)
	}

	// Unbind framebuffer
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	return framebuffer, nil
}

// Bind sets the framebuffer as the current render target.
func (fb *Framebuffer) Bind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, fb.id)
}

// Unbind resets the render target to the default framebuffer.
func (fb *Framebuffer) Unbind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

// Cleanup releases the resources associated with the framebuffer.
func (fb *Framebuffer) Cleanup() {
	gl.DeleteFramebuffers(1, &fb.id)
	gl.DeleteTextures(1, &fb.Texture)
}
