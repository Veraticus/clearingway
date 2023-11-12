package discord

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	tests := []struct {
		name                string
		input               string
		expectedChunks      int
		expectedLastLen     int
		expectedLastContent string
	}{
		{
			name:           "Short string",
			input:          "Hello, world!",
			expectedChunks: 1,
		},
		{
			name:           "Exact 2000 characters",
			input:          strings.Repeat("a", 1500),
			expectedChunks: 1,
		},
		{
			name:           "Just over 2000 characters",
			input:          strings.Repeat("a", 1501),
			expectedChunks: 2,
		},
		{
			name:                "String with newlines",
			input:               strings.Repeat("this is a string that should not split\n", 39),
			expectedChunks:      2,
			expectedLastContent: "this is a string that should not split\n",
		},
		{
			name:                "String with double newlines",
			input:               strings.Repeat("this is a string that should not split\nthis is a string that should split\n\n", 21),
			expectedChunks:      2,
			expectedLastContent: "this is a string that should not split\nthis is a string that should split\n\n",
		},
		{
			name:            "Long string without newlines",
			input:           strings.Repeat("a", 3000),
			expectedChunks:  2,
			expectedLastLen: 1500, // Second chunk will be 1000 characters
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chunks := NewChunks()
			chunks.Write(tt.input)

			assert.Equal(t, tt.expectedChunks, len(chunks.Chunks), "Number of chunks should match")

			for i, chunk := range chunks.Chunks {
				length := len(chunk.String())
				if i < len(chunks.Chunks)-1 {
					// Check all but the last chunk
					assert.GreaterOrEqual(t, length, 1000, "Chunk length should be at least 1500 characters")
					assert.LessOrEqual(t, length, 1500, "Chunk length should be less than 2000 characters")
				} else {
					if tt.expectedLastLen > 0 {
						// Check last chunk if specific length is expected
						assert.Equal(t, tt.expectedLastLen, length, "Last chunk length should match expected length. Last chunk is: %+v", chunks.Chunks[i].String())
						// Check last chunk if specific length is expected
					}
					if tt.expectedLastContent != "" {
						assert.Equal(t, tt.expectedLastContent, chunks.Chunks[i].String(), "Last chunk content should match expected content. Last chunk is: %+v", chunks.Chunks[i].String())
					}
				}
			}
		})
	}
}
