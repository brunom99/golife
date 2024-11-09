package client

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"golife/config"
	"sync"
	"time"
)

type Clients struct {
	mapClients      map[string]*Client
	mutexMapClients sync.Mutex
	lastActivity    int64 // milliseconds, to check if a goroutine is still alive
}

func (cs *Clients) Add(wsClient *websocket.Conn, config config.Config) (string, error) {
	// lock mutex
	cs.mutexMapClients.Lock()
	defer cs.mutexMapClients.Unlock()
	// create map clients
	if cs.mapClients == nil {
		cs.mapClients = make(map[string]*Client)
	}
	// new client
	client := cs.createClient(wsClient, config)
	cs.mapClients[client.id] = client
	// start client
	return client.id, client.start()
}

func (cs *Clients) Delete(idClient string) {
	// lock mutex
	cs.mutexMapClients.Lock()
	defer cs.mutexMapClients.Unlock()
	// send message to client
	if client, ok := cs.mapClients[idClient]; ok {
		// delete client's bubbles
		client.terminateBubbles()
		// delete pointer
		cs.mapClients[idClient] = nil
		// remove client from map
		delete(cs.mapClients, idClient)
	}
}

func (cs *Clients) UpdateLastActivity() {
	cs.lastActivity = time.Now().UnixMilli()
}

func (cs *Clients) GetLastActivity() int64 {
	return cs.lastActivity
}

func (cs *Clients) createClient(wsClient *websocket.Conn, config config.Config) *Client {
	// procedural random source
	seed := config.Grid.Seed
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	// generate client id
	idClient := uuid.New().String()
	// new client
	return &Client{
		id:                   idClient,
		seed:                 seed,
		ws:                   wsClient,
		config:               config,
		fnUpdateLastActivity: cs.UpdateLastActivity,
	}
}
