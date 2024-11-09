package client

import (
	"encoding/json"
	"fmt"
	"golife/config"
	"golife/life"
	"math/rand"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	id                   string
	seed                 int64
	config               config.Config
	ws                   *websocket.Conn
	mutexWsClient        sync.Mutex
	mapBubbles           map[string]life.Bubbles
	mutexMapBubbles      sync.Mutex
	randSeed             *rand.Rand
	disconnected         bool
	fnUpdateLastActivity func()
	totalBubbles         int
}

type InfoClient struct {
	Seed         string `json:"seed"`
	GridSize     int    `json:"grid_size"`
	TotalBubbles int    `json:"total_bubbles"`
}

func (c *Client) start() error {
	// procedural random
	randSource := rand.NewSource(c.seed)
	c.randSeed = rand.New(randSource)
	// send message to the client
	if err := c.sendMessageByWs(nil); err != nil {
		return err
	}
	// mutex lock
	c.mutexMapBubbles.Lock()
	defer c.mutexMapBubbles.Unlock()
	// create n bubbles
	c.mapBubbles = make(map[string]life.Bubbles)
	gridSize := c.config.Grid.Size
	for row := 0; row < gridSize; row++ {
		for col := 0; col < gridSize; col++ {
			var bubbles life.Bubbles
			// bubble position
			pos := life.Position{Row: row, Column: col}
			// check probability to create bubble
			if c.randSeed.Float64() < c.config.Bubble.Proba {
				bubbles = append(bubbles, c.createBubble(pos))
			}
			// add bubble to the map
			c.mapBubbles[pos.ToString()] = bubbles
		}
	}
	return nil
}

func (c *Client) createBubble(pos life.Position, rarity ...life.Rarity) *life.Bubble {
	// channel in/out to communicate with bubble
	chOut := make(chan life.MessageFromBubble)
	chIn := make(chan life.MessageToBubble)
	// create bubble
	bubble := life.CreateBubble(c.config, c.randSeed.Int63(), chOut, chIn, pos, rarity...)
	c.totalBubbles++
	// for each bubble, wait for message to the client
	go c.waitingBubbleMsg(chOut)
	// wake up bubble
	go bubble.WakeUp()
	return bubble
}

func (c *Client) waitingBubbleMsg(chanFromBubble chan life.MessageFromBubble) {
	for !c.disconnected { // avoid goroutine leaks
		// waiting for msg on bubble chan
		msgFromBubble := <-chanFromBubble
		// total bubble --
		if msgFromBubble.Bubble.IsFinish {
			c.totalBubbles--
		}
		// update bubble position
		bubblesToTerminate := c.updateBubblePosition(msgFromBubble)
		// terminate some bubble
		for _, b := range bubblesToTerminate {
			b.Terminate()
		}
		// transmit bubble msg to the browser
		go func() {
			_ = c.sendMessageByWs(msgFromBubble.Bubble)
		}()
		// update last activity
		c.fnUpdateLastActivity()
	}
}

func (c *Client) updateBubblePosition(msg life.MessageFromBubble) (toTerminate life.Bubbles) {
	// if no bubble, bubble is finish or not moving
	if msg.Bubble == nil || msg.Bubble.IsFinish || msg.Bubble.Position.IsSame(msg.Bubble.PositionPrev) {
		return
	}
	// bubble information
	bubble := msg.Bubble
	posPrev, posNew := bubble.PositionPrev.ToString(), bubble.Position.ToString()
	// mutex lock
	c.mutexMapBubbles.Lock()
	defer c.mutexMapBubbles.Unlock()
	// bubbles position prev and new
	bubblesPrevPos, foundPrev := c.mapBubbles[posPrev]
	bubblesNewPos, foundNew := c.mapBubbles[posNew]
	// remove bubble from previous position
	if foundPrev {
		bubblesPrevPos.Remove(bubble)
		c.mapBubbles[posPrev] = bubblesPrevPos
	}
	// do action for the new position
	if foundNew {
		toTerminate = c.doInteraction(bubble, &bubblesNewPos)
		c.mapBubbles[posNew] = bubblesNewPos
	}
	return
}

func (c *Client) doInteraction(bubble *life.Bubble, bubblesNewPos *life.Bubbles) (toTerminate life.Bubbles) {
	// new pos is empty: no interaction
	if bubblesNewPos.IsEmpty() {
		bubblesNewPos.Add(bubble)
		return
	}
	// bubble dark arrives at cell
	// if light bubble exists in the destination: terminate dark bubble
	// else bubble dark terminates all gray bubbles and duplicates itself
	switch {
	case bubble.IsDark():
		// light bubble in destination ?
		if bubblesNewPos.HasBubblesRarity(life.RarityLight) {
			// dark bubble is terminate
			toTerminate.Add(bubble)
			return
		}
		for _, b := range *bubblesNewPos {
			if b.IsInvisible || b.IsFinish {
				// ignore invisible or finish bubble
				continue
			}
			if !b.IsDark() {
				// for each bubble terminate, new dark bubble
				toTerminate.Add(b)
				// new dark bubble
				bubblesNewPos.Add(c.createBubble(bubble.Position, life.RarityDark))
			}
		}
		// add bubble dark
		bubblesNewPos.Add(bubble)
	case bubble.IsLight():
		// bubble light arrives at cell
		// terminate all dark bubbles
		for _, b := range *bubblesNewPos {
			if b.IsDark() && !b.IsInvisible && !b.IsFinish {
				toTerminate.Add(b)
			}
		}
	case bubblesNewPos.HasBubblesRarity(life.RarityDark):
		// common bubble is terminate
		toTerminate.Add(bubble)
	}
	return
}

func (c *Client) terminateBubbles() {
	// mutex lock
	c.mutexMapBubbles.Lock()
	defer c.mutexMapBubbles.Unlock()
	// send message to all bubbles
	for k, bubblesList := range c.mapBubbles {
		for _, bubble := range bubblesList {
			bubble.Terminate()
		}
		c.mapBubbles[k] = nil
	}
	c.mapBubbles = nil
	// client is disconnected
	c.disconnected = true
}

func (c *Client) sendMessageByWs(bubble *life.Bubble) error {
	// lock websocket mutex
	c.mutexWsClient.Lock()
	defer c.mutexWsClient.Unlock()
	// message to the client
	msg := struct {
		Bubble *life.Bubble `json:"bubble"`
		Info   InfoClient   `json:"info"`
	}{
		bubble,
		InfoClient{
			Seed:         fmt.Sprintf("%d", c.seed),
			GridSize:     c.config.Grid.Size,
			TotalBubbles: c.totalBubbles,
		},
	}
	// send message
	if bytes, err := json.Marshal(msg); err == nil {
		return c.ws.WriteMessage(1, bytes)
	}
	return nil
}
