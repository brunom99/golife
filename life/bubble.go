package life

import (
	"golife/config"
	"golife/utils"
	"math/rand"
	"time"
)

type Rarity string

const (
	RarityCommon Rarity = "common"
	RarityLight  Rarity = "light"
	RarityDark   Rarity = "dark"
)

type Bubble struct {
	ID                string   `json:"id"`
	Position          Position `json:"pos"`
	PositionPrev      Position `json:"-"`
	Rarity            Rarity   `json:"rarity"`
	IsFinish          bool     `json:"is_finish"`
	IsInvisible       bool     `json:"is_invisible"`
	randSeed          *rand.Rand
	chOut             chan MessageFromBubble
	chIn              chan MessageToBubble
	config            config.Config
	speed             int
	targetPos         Position
	allowDiagonalMove bool
}

func CreateBubble(config config.Config, seed int64,
	chOut chan MessageFromBubble, chIn chan MessageToBubble, posStart Position, rarity ...Rarity) *Bubble {
	bubble := Bubble{
		ID:           utils.Uuid(),
		Position:     posStart,
		PositionPrev: posStart,
		IsInvisible:  true,
		randSeed:     rand.New(rand.NewSource(seed)),
		chOut:        chOut,
		chIn:         chIn,
		config:       config,
		targetPos:    posStart,
	}
	bubble.init(rarity...)
	return &bubble
}

func (b *Bubble) WakeUp() {
	// check channel: message to bubble
	go b.readMessage()
	// while bubble is alive
	for !b.IsFinish { // avoid goroutine leaks
		// send message to the client
		b.sendMessage()
		// random waiting
		time.Sleep(time.Duration(b.speed) * time.Millisecond)
		// move bubble
		b.move()
		// after first turn, bubble looses invincibility
		b.IsInvisible = false
	}
	// bubble is finish: close read channel
	close(b.chIn)
	// bubble is finish: send a last message
	b.sendMessage()
}

func (b *Bubble) Message(msg MessageToBubble) {
	b.chIn <- msg
}

func (b *Bubble) Terminate() {
	b.IsFinish = true
}

func (b *Bubble) IsDark() bool {
	return b != nil && b.Rarity == RarityDark
}

func (b *Bubble) IsLight() bool {
	return b != nil && b.Rarity == RarityLight
}

func (b *Bubble) sendMessage() {
	b.chOut <- MessageFromBubble{
		Bubble: b,
	}
}

func (b *Bubble) readMessage() {
	// while bubble is alive
	for !b.IsFinish {
		// wait for channel message
		message := <-b.chIn
		_ = message
	}
}

func (b *Bubble) init(rarity ...Rarity) {
	// bubble rarity
	if len(rarity) > 0 {
		// force rarity
		b.Rarity = rarity[0]
	} else {
		// random rarity
		b.Rarity = RarityCommon
		randPoolPosition := utils.RandInt(1, b.config.Bubble.Pool, b.randSeed)
		incrementPool := 0
		for _, rar := range []Rarity{RarityDark, RarityLight} {
			rarityPool := b.config.Bubbles[string(rar)].Pool
			incrementPool += rarityPool
			if randPoolPosition <= incrementPool {
				b.Rarity = rar
				break
			}
		}
	}
	// config bubble by rarity
	confBubble := b.config.Bubbles[string(b.Rarity)]
	// bubble speed
	b.speed = utils.RandInt(confBubble.MinSpeed, confBubble.MaxSpeed, b.randSeed)
	// allow diagonal movement
	b.allowDiagonalMove = confBubble.Diagonal
}

func (b *Bubble) move() {
	// define new targetPos ?
	if b.Position.IsSame(b.targetPos) {
		b.setNewTargetPosition()
	}
	// save old position
	b.PositionPrev = b.Position
	// move
	b.Position.Move(b.targetPos, b.allowDiagonalMove)
}

func (b *Bubble) setNewTargetPosition() {
	gridSize := b.config.Grid.Size
	b.targetPos = Position{
		Row:    utils.RandInt(0, gridSize, b.randSeed),
		Column: utils.RandInt(0, gridSize, b.randSeed),
	}
}
