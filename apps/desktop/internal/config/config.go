// Package config centraliza as configurações de rede do servidor.
package config

import "fmt"

// Config armazena o endereço de escuta do servidor WebSocket.
type Config struct {
	Host string // endereço IP de bind; "0.0.0.0" aceita conexões de qualquer interface
	Port int    // porta TCP do servidor
}

// Default retorna a configuração padrão: escuta em todas as interfaces na porta 8080.
// O host "0.0.0.0" é necessário para que o celular na mesma rede Wi-Fi consiga conectar;
// usar "127.0.0.1" restringiria ao loopback, bloqueando conexões externas.
func Default() *Config {
	return &Config{
		Host: "0.0.0.0",
		Port: 8080,
	}
}

// Addr formata o endereço no padrão "host:porta" esperado por http.ListenAndServe.
func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
