# RemoteMouse — Desktop

Servidor WebSocket escrito em Go que recebe comandos do app mobile e os traduz em eventos reais de mouse e teclado no Windows.

---

## Pré-requisitos

- [Go 1.22+](https://go.dev/dl/)
- Windows (o servidor usa a API Win32 diretamente via `user32.dll`)

---

## Iniciando

```bash
go run .\cmd\main.go
```

Por padrão o servidor escuta em `0.0.0.0:8080`. O endereço IP da máquina na rede local é exibido no terminal ao iniciar.

Para compilar um executável:

```bash
go build -o remotemouse.exe .\cmd\main.go
```

---

## Estrutura

```
desktop/
├── cmd/
│   └── main.go          # Ponto de entrada
└── internal/
    ├── config/
    │   └── config.go    # Host e porta (padrão: 0.0.0.0:8080)
    ├── input/
    │   ├── mouse.go     # Mouse: mover, clicar, scroll + roteador de comandos
    │   └── keyboard.go  # Teclado: teclas Unicode e virtuais Win32
    └── server/
        └── websocket.go # Servidor WebSocket (gorilla/websocket)
```

---

## Protocolo

Todos os comandos recebidos são JSON. Veja a especificação completa em [`/shared/protocol/protocol.md`](../../shared/protocol/protocol.md).

| Comando       | Descrição                        |
|---------------|----------------------------------|
| `mouse_move`  | Move o cursor (delta x/y)        |
| `mouse_click` | Clique simples ou duplo          |
| `scroll`      | Rola a tela (notches)            |
| `key_press`   | Pressiona uma tecla ou caractere |

---

## Dependências

| Pacote                  | Uso                        |
|-------------------------|----------------------------|
| `gorilla/websocket`     | Servidor WebSocket         |
| `syscall` (stdlib)      | Chamadas diretas à Win32   |
