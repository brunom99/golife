package life

import (
	"golife/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBubble_Create(t *testing.T) {
	// channel in/out to communicate with bubble
	chOut := make(chan MessageFromBubble)
	chIn := make(chan MessageToBubble)
	// create bubble
	pos := Position{Row: 1, Column: 2}
	bubble := CreateBubble(config.Config{}, 0, chOut, chIn, pos)
	assert.Equal(t, 1, bubble.Position.Row)
	assert.Equal(t, 2, bubble.Position.Column)
}
