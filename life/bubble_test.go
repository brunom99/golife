package life

import (
	"golife/config"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBubble_Create(t *testing.T) {
	chOut := make(chan MessageFromBubble)
	chIn := make(chan MessageToBubble)
	pos := Position{Row: 1, Column: 2}
	bubble := CreateBubble(config.Config{}, 0, chOut, chIn, pos)

	assert.Equal(t, 1, bubble.Position.Row)
	assert.Equal(t, 2, bubble.Position.Column)
}

func TestBubble_Move(t *testing.T) {
	bubble := &Bubble{
		Position:          Position{Row: 0, Column: 0},
		targetPos:         Position{Row: 2, Column: 2},
		allowDiagonalMove: true,
	}
	bubble.move()
	assert.NotEqual(t, Position{Row: 0, Column: 0}, bubble.Position)
}

func TestBubble_Terminate(t *testing.T) {
	bubble := &Bubble{}
	bubble.Terminate()
	assert.True(t, bubble.IsFinish)
}

func TestBubble_MessageHandling(t *testing.T) {
	chOut := make(chan MessageFromBubble, 1)
	chIn := make(chan MessageToBubble, 1)
	bubble := CreateBubble(config.Config{}, 0, chOut, chIn, Position{})

	msg := MessageToBubble{}
	bubble.Message(msg)
	select {
	case receivedMsg := <-chIn:
		assert.Equal(t, msg, receivedMsg)
	case <-time.After(time.Millisecond * 100):
		t.Fatal("Message not received")
	}
}

func TestBubble_WakeUp(t *testing.T) {
	chOut := make(chan MessageFromBubble, 1)
	chIn := make(chan MessageToBubble, 1)
	bubble := CreateBubble(config.Config{}, 0, chOut, chIn, Position{})

	go bubble.WakeUp()
	time.Sleep(time.Millisecond * 50)
	bubble.Terminate()
	assert.True(t, bubble.IsFinish)
	assert.False(t, bubble.IsInvisible)
}

func TestBubble_SetNewTargetPosition(t *testing.T) {
	bubble := &Bubble{
		config:   config.Config{Grid: config.GridConfig{Size: 10}},
		randSeed: rand.New(rand.NewSource(42)),
	}
	prevPos := bubble.targetPos
	bubble.setNewTargetPosition()
	assert.Less(t, bubble.targetPos.Row, 10)
	assert.Less(t, bubble.targetPos.Column, 10)
	assert.GreaterOrEqual(t, bubble.targetPos.Row, 0)
	assert.GreaterOrEqual(t, bubble.targetPos.Column, 0)
	assert.NotEqual(t, prevPos, bubble.targetPos)
}
