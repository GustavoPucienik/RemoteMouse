// Package server implementa o servidor WebSocket que recebe comandos do app mobile.
package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"remotemouse/desktop/internal/config"
	"remotemouse/desktop/internal/input"
)

// upgrader converte uma conexão HTTP comum em WebSocket.
// CheckOrigin retorna sempre true para aceitar qualquer origem — adequado para
// uso em rede local onde não há risco de CSRF via browser.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Server encapsula a configuração necessária para rodar o servidor HTTP/WebSocket.
type Server struct {
	cfg *config.Config
}

// New cria um novo Server com a configuração fornecida.
func New(cfg *config.Config) *Server {
	return &Server{cfg: cfg}
}

// Start registra a rota WebSocket e bloqueia escutando por conexões.
// Retorna erro se o bind na porta falhar (ex: porta já em uso).
func (s *Server) Start() error {
	http.HandleFunc("/ws", s.handleWS)
	return http.ListenAndServe(s.cfg.Addr(), nil)
}

// handleWS gerencia uma conexão WebSocket individual.
// Cada cliente recebe sua própria goroutine implícita (o http.Server faz isso).
// O loop lê mensagens JSON em sequência e as repassa ao handler de input;
// quando o cliente desconecta, ReadMessage retorna erro e o loop encerra.
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
			// Erro esperado ao fechar o app ou perder conexão Wi-Fi.
			log.Printf("client disconnected: %s", r.RemoteAddr)
			break
		}
		if err := input.Handle(msg); err != nil {
			log.Println("input error:", err)
		}
	}
}
