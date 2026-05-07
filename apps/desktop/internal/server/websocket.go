package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"remotemouse/desktop/internal/config"
	"remotemouse/desktop/internal/input"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Server struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Server {
	return &Server{cfg: cfg}
}

func (s *Server) Start() error {
	http.HandleFunc("/ws", s.handleWS)
	return http.ListenAndServe(s.cfg.Addr(), nil)
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer conn.Close()
	log.Printf("client connected: %s", r.RemoteAddr)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("client disconnected: %s", r.RemoteAddr)
			break
		}
		if err := input.Handle(msg); err != nil {
			log.Println("input error:", err)
		}
	}
}
