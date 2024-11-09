package life

import "fmt"

type Position struct {
	Row    int `json:"row"`
	Column int `json:"column"`
}

func (p *Position) IsSame(other Position) bool {
	if p == nil {
		return false
	}
	return p.Row == other.Row && p.Column == other.Column
}

func (p *Position) ToString() string {
	if p == nil {
		return ""
	}
	return fmt.Sprintf("%d_%d", p.Row, p.Column)
}

func (p *Position) Move(target Position, allowDiagonalMove bool) {
	if p == nil {
		return
	}
	// define increment
	incRow := p.determineIncrement(p.Row, target.Row)
	incCol := p.determineIncrement(p.Column, target.Column)
	// move to target pos
	switch {
	case allowDiagonalMove && incRow != 0 && incCol != 0:
		p.Column += incCol
		p.Row += incRow
	case incRow != 0:
		p.Row += incRow
	case incCol != 0:
		p.Column += incCol
	}
}

func (p *Position) determineIncrement(current, target int) int {
	switch {
	case current < target:
		return 1
	case current > target:
		return -1
	default:
		return 0
	}
}
