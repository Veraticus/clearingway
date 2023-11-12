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
		remainingLength := 1500 - currentLength

		// If the current chunk can accommodate the entire string `s`
		if len(s) <= remainingLength {
			current.WriteString(s)
			break
		}

		// Find the best place to split the string
		splitIndex := findSplitIndex(s, remainingLength, currentLength)

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
func findSplitIndex(s string, maxLength, currentLength int) int {
	minIndex := 1000 - currentLength // Adjusted based on the current length of the chunk
	if minIndex < 0 {
		minIndex = 0
	}

	for i := maxLength; i > minIndex; i-- {
		if i > 1 && s[i-1] == '\n' && s[i-2] == '\n' {
			return i // Include the double newline in the first part
		}
	}

	for i := maxLength; i > minIndex; i-- {
		if s[i-1] == '\n' {
			return i // Include the newline in the first part
		}
	}

	for i := maxLength; i > minIndex; i-- {
		if s[i-1] == ' ' {
			return i // Split at the space
		}
	}

	return maxLength // No suitable split character found
}
func (c *Chunks) currentChunk() *strings.Builder {
	return c.Chunks[len(c.Chunks)-1]
}
