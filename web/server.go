package web

import (
	"fmt"
	"github.com/gorilla/websocket"
	"golife/client"
	"golife/config"
	"golife/utils"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	Config  config.Config
	Clients client.Clients
}

func (s *Server) Start() error {
	// create http server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.Config.Server.Port),
		Handler:      s.getRouter(),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}
	// listen & serve
	log.Printf("Starting server (%s)", server.Addr)
	return server.ListenAndServe()
}

func (s *Server) getRouter() *chi.Mux {
	router := chi.NewRouter()
	// default route for all static files
	router.Handle("/*", http.FileServer(http.Dir("./static")))
	// websocket handler
	router.HandleFunc("/ws", s.wsHandler)
	// api last activity
	router.HandleFunc("/api/activity", s.activityHandler)
	//
	return router
}

func (s *Server) wsHandler(w http.ResponseWriter, r *http.Request) {
	// upgrader
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	// upgrade connection to a webSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		return
	}
	// new client
	var idClient string
	// defer function when connection with the client is over
	defer func(ws *websocket.Conn) {
		// close websocket (error not wanted)
		_ = ws.Close()
		// destoy client (with all bubbles)
		go s.Clients.Delete(idClient)
	}(ws)
	// load config dynamically (to avoid rebuild project when file config changes)
	if err = s.Config.LoadFile("./config.toml"); err != nil {
		log.Fatal(err)
	}
	// client connected
	if idClient, err = s.Clients.Add(ws, s.Config); err != nil {
		// error when starting new client
		return
	}
	// waiting message from client
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			// close websocket (error not wanted)
			_ = ws.Close()
			break
		}
		// send msg to all bubbles ?
		_ = msg
	}
	log.Printf("Client %s: websocket is closed", idClient)
}

func (s *Server) activityHandler(w http.ResponseWriter, _ *http.Request) {
	utils.HttpAccept(w, struct {
		LastActivity int64 `json:"last_activity"`
	}{
		LastActivity: s.Clients.GetLastActivity(),
	})
}
