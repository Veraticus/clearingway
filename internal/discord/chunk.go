package discord

import (
	"strings"
)

type Chunks struct {
	Chunks []*strings.Builder
}

func NewChunks() *Chunks {
	return &Chunks{Chunks: []*strings.Builder{{}}}
}

func (c *Chunks) Write(s string) {
	current := c.currentChunk()
	currentLength := len(current.String())
	length := len(s)

	if length+currentLength >= 1900 {
		current = &strings.Builder{}
		c.Chunks = append(c.Chunks, current)
	}

	current.WriteString(s)
}

func (c *Chunks) currentChunk() *strings.Builder {
	return c.Chunks[len(c.Chunks)-1]
}
