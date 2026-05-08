# RemoteMouse — Protocolo WebSocket

Todas as mensagens são JSON enviados pelo cliente (mobile) para o servidor (desktop) via WebSocket em `ws://<ip>:8080/ws`.

---

## mouse_move

Move o cursor relativamente à posição atual.

```json
{ "type": "mouse_move", "dx": 10, "dy": -5 }
```

| Campo | Tipo  | Descrição                        |
|-------|-------|----------------------------------|
| `dx`  | int   | Deslocamento horizontal (pixels) |
| `dy`  | int   | Deslocamento vertical (pixels)   |

---

## mouse_click

Realiza um clique simples ou duplo.

```json
{ "type": "mouse_click", "button": "left", "double": false }
```

| Campo    | Tipo    | Valores                      |
|----------|---------|------------------------------|
| `button` | string  | `"left"`, `"right"`, `"middle"` |
| `double` | boolean | `true` para clique duplo     |

---

## key_press

Pressiona e solta uma tecla no desktop.

```json
{ "type": "key_press", "key": "a" }
{ "type": "key_press", "key": "Enter" }
```

| Campo | Tipo   | Descrição |
|-------|--------|-----------|
| `key` | string | Caractere único (ex: `"a"`, `"@"`) ou nome de tecla especial (Web `KeyboardEvent.key`) |

**Teclas especiais suportadas:**

| `key`        | Ação           |
|--------------|----------------|
| `Backspace`  | Apagar         |
| `Tab`        | Tabulação      |
| `Enter`      | Confirmar      |
| `Escape`     | Cancelar       |
| `Delete`     | Delete         |
| `Home`       | Início da linha |
| `End`        | Fim da linha   |
| `PageUp`     | Página acima   |
| `PageDown`   | Página abaixo  |
| `ArrowLeft`  | Seta esquerda  |
| `ArrowRight` | Seta direita   |
| `ArrowUp`    | Seta acima     |
| `ArrowDown`  | Seta abaixo    |
| `F1`–`F12`   | Teclas de função |

---

## scroll

Rola a tela verticalmente.

```json
{ "type": "scroll", "dy": 3 }
```

| Campo | Tipo | Descrição                                   |
|-------|------|---------------------------------------------|
| `dy`  | int  | Notches de rolagem: positivo = cima, negativo = baixo |
