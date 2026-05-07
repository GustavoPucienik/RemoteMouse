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
	log.Fatal(srv.Start())
}
