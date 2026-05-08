# RemoteMouse — Mobile

Cliente web em React/TypeScript que funciona como controle remoto no navegador do celular.

---

## Pré-requisitos

- [Node.js 18+](https://nodejs.org/)
- Celular e computador na mesma rede Wi-Fi

---

## Iniciando

```bash
npm install
npm run dev -- --host
```

O flag `--host` expõe o servidor na rede local. Abra o endereço exibido no terminal no navegador do celular.

---

## Outros scripts

| Comando         | Descrição                        |
|-----------------|----------------------------------|
| `npm run dev`   | Servidor de desenvolvimento      |
| `npm run build` | Build de produção (pasta `dist/`) |
| `npm run lint`  | Verificação de código com ESLint |

---

## Estrutura

```
mobile/
├── src/
│   ├── components/
│   │   ├── ConnectScreen.tsx  # Tela de conexão (IP + porta)
│   │   ├── RemoteScreen.tsx   # Tela principal (alterna mouse/teclado)
│   │   ├── Trackpad.tsx       # Área de toque para mover o mouse
│   │   ├── ClickButtons.tsx   # Botões de clique esquerdo/direito
│   │   └── Keyboard.tsx       # Teclado remoto com teclas especiais
│   ├── hooks/
│   │   └── useWebSocket.ts    # Gerencia a conexão WebSocket
│   ├── App.tsx                # Raiz: ConnectScreen → RemoteScreen
│   ├── main.tsx               # Entry point
│   └── index.css              # Estilos globais (dark theme)
└── index.html
```

---

## Gestos no trackpad

| Gesto         | Ação           |
|---------------|----------------|
| 1 dedo        | Mover o mouse  |
| 2 dedos       | Scroll         |

---

## Modo teclado

Toque no ícone `⌨` na barra superior para alternar para o teclado remoto. O teclado nativo do celular aparece automaticamente para digitar texto. Teclas especiais (Esc, Tab, setas, Del) ficam disponíveis como botões na tela.
