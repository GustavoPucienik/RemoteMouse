// Ponto de entrada do servidor RemoteMouse para desktop.
// Inicializa a configuração padrão e sobe o servidor WebSocket que recebe
// comandos do aplicativo mobile e os traduz em eventos de mouse do Windows.
package main

import (
	"log"

	"remotemouse/desktop/internal/config"
	"remotemouse/desktop/internal/server"
)

func main() {
	cfg := config.Default()
	srv := server.New(cfg)
	log.Printf("RemoteMouse listening on ws://%s/ws", cfg.Addr())
	// log.Fatal termina o processo com código de saída não-zero se o servidor falhar,
	// garantindo que erros de bind/listen (ex: porta em uso) não passem silenciosamente.
	log.Fatal(srv.Start())
}
