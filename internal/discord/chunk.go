package discord

import (
	"strings"
)

type Chunks struct {
	Chunks []*strings.Builder
}

func NewChunks() *Chunks {
	return &Chunks{Chunks: make([]*strings.Builder, 0, 10)}
}

func (c *Chunks) Write(s string) {
	for len(s) > 0 {
		current := c.currentChunk()
		currentLength := current.Len()
		remainingLength := 1500 - currentLength

		if len(s) <= remainingLength {
			current.WriteString(s)
			return
		}

		splitIndex := findSplitIndex(s, remainingLength)
		current.WriteString(s[:splitIndex])
		s = s[splitIndex:]
		c.Chunks = append(c.Chunks, &strings.Builder{})
	}
}

func findSplitIndex(s string, maxLength int) int {
	if maxLength >= len(s) {
		return len(s)
	}

	// Using LastIndexFunc to find the last newline before maxLength
	newlineIndex := strings.LastIndexFunc(s[:maxLength], func(r rune) bool {
		return r == '\n'
	})

	if newlineIndex != -1 {
		return newlineIndex + 1 // Include the newline
	}

	// If no newline, find the last space before maxLength
	spaceIndex := strings.LastIndexFunc(s[:maxLength], func(r rune) bool {
		return r == ' '
	})

	if spaceIndex != -1 {
		return spaceIndex + 1 // Include the space
	}

	return maxLength // No newline or space found
}

func (c *Chunks) currentChunk() *strings.Builder {
	if len(c.Chunks) == 0 {
		c.Chunks = append(c.Chunks, &strings.Builder{})
	}
	return c.Chunks[len(c.Chunks)-1]
}
