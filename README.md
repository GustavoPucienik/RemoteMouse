# RemoteMouse

Controle o computador pelo celular via rede local. Sem conta, sem internet, sem anúncios.

---

## O problema

Muitas vezes precisamos fazer pequenas ações no computador — pular um anúncio, pausar uma música, controlar uma apresentação — sem querer levantar ou sem ter o teclado/mouse por perto.

As soluções existentes costumam ter anúncios, limitações na versão gratuita, dependência de servidores externos ou aplicativos pesados e desatualizados.

---

## A solução

O RemoteMouse transforma qualquer celular em um controle remoto para o computador usando apenas a rede Wi-Fi local. Open-source, leve e focado em privacidade.

---

## Funcionalidades

| Funcionalidade   | Status |
|------------------|--------|
| Mover o mouse    | ✅ |
| Clique esquerdo/direito | ✅ |
| Scroll           | ✅ |
| Teclado remoto   | ✅ |
| Controle de mídia | 🔜 |
| Controle de volume | 🔜 |

---

## Como usar

**1. No computador** — inicie o servidor:

```bash
cd apps/desktop
go run .\cmd\main.go
```

O terminal exibirá o IP da máquina na rede local.

**2. No celular** — inicie o app web:

```bash
cd apps/mobile
npm install
npm run dev -- --host
```

Abra o endereço exibido no terminal no navegador do celular, digite o IP do computador e conecte.

---

## Tecnologias

| Camada | Tecnologia |
|---|---|
| Servidor (desktop) | Go + WebSocket + Win32 API (`user32.dll`) |
| Cliente (mobile/web) | React + TypeScript + Vite |

---

## Estrutura do projeto

```
RemoteMouse/
├── apps/
│   ├── desktop/    # Servidor Go para Windows
│   └── mobile/     # Cliente web (React)
└── shared/
    └── protocol/   # Especificação do protocolo WebSocket
```

---

## Futuras funcionalidades

- Conexão por QR Code
- Controle de mídia e volume
- Atalhos personalizados
- Comunicação criptografada
- PWA (instalar como app no celular)
- Suporte multi-monitor
- Wake-on-LAN

---

## Filosofia

Minimalista, rápido, privado e fácil de contribuir. Nenhum dado sai da sua rede local.

---

## Status

🚧 Em desenvolvimento — MVP funcional.
