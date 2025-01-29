package life

type Bubbles []*Bubble

func (b *Bubbles) Add(bubble *Bubble) {
	if b == nil {
		*b = []*Bubble{}
	}
	if bubble == nil {
		return
	}
	*b = append(*b, bubble)
}

func (b *Bubbles) Remove(other *Bubble) {
	if b == nil || other == nil {
		return
	}
	var a []*Bubble
	for _, bubble := range *b {
		if bubble.ID != other.ID {
			a = append(a, bubble)
		}
	}
	*b = a
}

func (b *Bubbles) Contains(other *Bubble) bool {
	if b == nil || other == nil {
		return false
	}
	for _, bubble := range *b {
		if bubble.ID == other.ID {
			return true
		}
	}
	return false
}

func (b *Bubbles) Count() map[Rarity]Bubbles {
	if b == nil {
		return nil
	}
	m := make(map[Rarity]Bubbles)
	for _, bubble := range *b {
		m[bubble.Rarity] = append(m[bubble.Rarity], bubble)
	}
	return m
}

func (b *Bubbles) IsEmpty() bool {
	return b == nil || len(*b) == 0
}

func (b *Bubbles) HasBubblesRarity(rarity Rarity) bool {
	if b == nil {
		return false
	}
	for _, bubble := range *b {
		if bubble.Rarity == rarity {
			return true
		}
	}
	return false
}
