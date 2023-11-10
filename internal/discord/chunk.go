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
	for len(s) > 0 {
		current := c.currentChunk()
		currentLength := len(current.String())
		remainingLength := 1900 - currentLength

		// If the current chunk can accommodate the entire string `s`
		if len(s) <= remainingLength {
			current.WriteString(s)
			break
		}

		// Find the best place to split the string
		splitIndex := findSplitIndex(s, remainingLength)

		// Append the first part of the string to the current chunk
		current.WriteString(s[:splitIndex])

		// Create a new chunk for the remaining part of the string
		s = s[splitIndex:]
		c.Chunks = append(c.Chunks, &strings.Builder{})
	}
}

// findSplitIndex finds the best index to split the string.
// It looks for a newline close to the specified max length without exceeding it.
// If no newline is found, it returns the max length.
func findSplitIndex(s string, maxLength int) int {
	if maxLength >= len(s) {
		return len(s)
	}

	for i := maxLength; i > maxLength-1000 && i > 0; i-- {
		if s[i] == '\n' {
			return i + 1 // Include the newline in the first part
		}
	}
	return maxLength // No suitable newline found, split at maxLength
}

func (c *Chunks) currentChunk() *strings.Builder {
	return c.Chunks[len(c.Chunks)-1]
}
