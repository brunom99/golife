package life

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBubbles_Add(t *testing.T) {
	var bubbles Bubbles
	bubble := &Bubble{ID: "1", Rarity: RarityCommon}
	bubbles.Add(bubble)

	assert.Equal(t, 1, len(bubbles))
	assert.Equal(t, bubble, bubbles[0])

	// Test adding a nil bubble
	bubbles.Add(nil)
	assert.Equal(t, 1, len(bubbles))
}

func TestBubbles_Remove(t *testing.T) {
	bubble1 := &Bubble{ID: "1", Rarity: RarityCommon}
	bubble2 := &Bubble{ID: "2", Rarity: RarityLight}
	bubbles := Bubbles{bubble1, bubble2}

	bubbles.Remove(bubble1)
	assert.Equal(t, 1, len(bubbles))
	assert.Equal(t, bubble2, bubbles[0])

	bubbles.Remove(nil) // Should not panic
	assert.Equal(t, 1, len(bubbles))

	bubbles.Remove(&Bubble{ID: "3"}) // Removing non-existent bubble
	assert.Equal(t, 1, len(bubbles))
}

func TestBubbles_Contains(t *testing.T) {
	bubble1 := &Bubble{ID: "1", Rarity: RarityCommon}
	bubble2 := &Bubble{ID: "2", Rarity: RarityLight}
	bubbles := Bubbles{bubble1}

	assert.True(t, bubbles.Contains(bubble1))
	assert.False(t, bubbles.Contains(bubble2))
	assert.False(t, bubbles.Contains(nil))
}

func TestBubbles_Count(t *testing.T) {
	bubble1 := &Bubble{ID: "1", Rarity: RarityCommon}
	bubble2 := &Bubble{ID: "2", Rarity: RarityLight}
	bubble3 := &Bubble{ID: "3", Rarity: RarityCommon}
	bubbles := Bubbles{bubble1, bubble2, bubble3}

	counts := bubbles.Count()

	assert.Equal(t, 2, len(counts[RarityCommon]))
	assert.Equal(t, 1, len(counts[RarityLight]))

	// Test with an empty list
	var emptyBubbles Bubbles
	assert.Equal(t, 0, len(emptyBubbles.Count()))
}

func TestBubbles_IsEmpty(t *testing.T) {
	var bubbles Bubbles
	assert.True(t, bubbles.IsEmpty())

	bubble := &Bubble{ID: "1", Rarity: RarityCommon}
	bubbles.Add(bubble)
	assert.False(t, bubbles.IsEmpty())
}

func TestBubbles_HasBubblesRarity(t *testing.T) {
	bubble := &Bubble{ID: "1", Rarity: RarityLight}
	bubbles := Bubbles{bubble}

	assert.True(t, bubbles.HasBubblesRarity(RarityLight))
	assert.False(t, bubbles.HasBubblesRarity(RarityDark))
	assert.False(t, bubbles.HasBubblesRarity(RarityCommon))

	// Test with an empty list
	var emptyBubbles Bubbles
	assert.False(t, emptyBubbles.HasBubblesRarity(RarityLight))
}

func TestBubble_IsDark(t *testing.T) {
	bubble := &Bubble{ID: "1", Rarity: RarityDark}
	assert.True(t, bubble.IsDark())

	bubble.Rarity = RarityLight
	assert.False(t, bubble.IsDark())

	var nilBubble *Bubble
	assert.False(t, nilBubble.IsDark())
}

func TestBubble_IsLight(t *testing.T) {
	bubble := &Bubble{ID: "1", Rarity: RarityLight}
	assert.True(t, bubble.IsLight())

	bubble.Rarity = RarityDark
	assert.False(t, bubble.IsLight())

	var nilBubble *Bubble
	assert.False(t, nilBubble.IsLight())
}
