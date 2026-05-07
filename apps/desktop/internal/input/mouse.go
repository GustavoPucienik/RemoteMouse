// Compilado apenas em Windows — usa a API Win32 diretamente via syscall.
//go:build windows

// Package input traduz comandos JSON recebidos via WebSocket em eventos
// reais de mouse, simulados através da função SendInput da user32.dll.
package input

import (
	"encoding/json"
	"fmt"
	"syscall"
	"unsafe"
)

var (
	// user32.dll é a biblioteca do Windows que expõe APIs de interface gráfica e input.
	user32 = syscall.NewLazyDLL("user32.dll")
	// SendInput injeta eventos sintéticos de teclado/mouse na fila de input do sistema.
	fnSendInput = user32.NewProc("SendInput")
)

// Flags de tipo de input e de evento de mouse conforme a documentação Win32.
// Referência: https://learn.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-mouseinput
const (
	inputMouse     = 0      // valor do campo dwType para eventos de mouse
	mouseFlagMove  = 0x0001 // movimento relativo do cursor (dx/dy)
	mouseFlagLDown = 0x0002 // botão esquerdo pressionado
	mouseFlagLUp   = 0x0004 // botão esquerdo liberado
	mouseFlagRDown = 0x0008 // botão direito pressionado
	mouseFlagRUp   = 0x0010 // botão direito liberado
	mouseFlagMDown = 0x0020 // botão do meio pressionado
	mouseFlagMUp   = 0x0040 // botão do meio liberado
	mouseFlagWheel = 0x0800 // evento de roda de scroll
	wheelDelta     = 120    // unidade padrão Win32 para um "clique" de scroll
)

// mouseInput mirrors the Win32 INPUT struct for mouse events on 64-bit Windows.
// Layout: type(4) + pad(4) + MOUSEINPUT(dx=4, dy=4, mouseData=4, dwFlags=4, time=4, [auto-pad=4], extra=8)
// Total: 40 bytes — matches sizeof(INPUT) on 64-bit.
type mouseInput struct {
	inputType uint32
	_         [4]byte // padding: Win32 aligns the union to offset 8
	dx        int32
	dy        int32
	mouseData uint32
	dwFlags   uint32
	time      uint32
	extra     uintptr // Go auto-aligns to 8 bytes (adds 4 bytes padding before this)
}

// baseCmd é desserializado primeiro para descobrir o tipo do comando
// antes de fazer o parse completo nos structs específicos abaixo.
type baseCmd struct {
	Type string `json:"type"`
}

// mouseMoveCmd representa um deslocamento relativo do cursor em pixels.
type mouseMoveCmd struct {
	DX int32 `json:"dx"` // deslocamento horizontal: positivo = direita
	DY int32 `json:"dy"` // deslocamento vertical: positivo = baixo
}

// mouseClickCmd representa um clique de mouse, simples ou duplo.
type mouseClickCmd struct {
	Button string `json:"button"` // "left", "right", "middle"
	Double bool   `json:"double"` // true para duplo clique
}

// scrollCmd representa rolagem da roda do mouse.
type scrollCmd struct {
	DY int32 `json:"dy"` // positive = scroll up, negative = scroll down
}

// Handle decodifica uma mensagem JSON recebida pelo WebSocket e executa o
// comando de input correspondente. Usa double-unmarshal: primeiro lê só o
// campo "type" para evitar alocar todos os structs de uma vez, depois faz
// o parse completo apenas no struct relevante.
func Handle(data []byte) error {
	var base baseCmd
	if err := json.Unmarshal(data, &base); err != nil {
		return err
	}
	switch base.Type {
	case "mouse_move":
		var cmd mouseMoveCmd
		json.Unmarshal(data, &cmd)
		return sendMouseMove(cmd.DX, cmd.DY)
	case "mouse_click":
		var cmd mouseClickCmd
		json.Unmarshal(data, &cmd)
		return sendMouseClick(cmd.Button, cmd.Double)
	case "scroll":
		var cmd scrollCmd
		json.Unmarshal(data, &cmd)
		return sendScroll(cmd.DY)
	default:
		return fmt.Errorf("unknown command type: %s", base.Type)
	}
}

// sendMouseMove desloca o cursor de forma relativa à posição atual.
// Sem MOUSEEVENTF_ABSOLUTE, dx/dy são interpretados como delta, não coordenadas absolutas.
func sendMouseMove(dx, dy int32) error {
	return callSendInput(&mouseInput{
		inputType: inputMouse,
		dx:        dx,
		dy:        dy,
		dwFlags:   mouseFlagMove,
	})
}

// sendMouseClick simula um clique completo (press + release) em um dos três botões.
// Para duplo clique envia o par press/release duas vezes seguidas, que é o
// comportamento esperado pelo Windows para disparar o evento WM_LBUTTONDBLCLK.
func sendMouseClick(button string, double bool) error {
	var downFlag, upFlag uint32
	switch button {
	case "left":
		downFlag, upFlag = mouseFlagLDown, mouseFlagLUp
	case "right":
		downFlag, upFlag = mouseFlagRDown, mouseFlagRUp
	case "middle":
		downFlag, upFlag = mouseFlagMDown, mouseFlagMUp
	default:
		return fmt.Errorf("unknown button: %s", button)
	}
	clicks := 1
	if double {
		clicks = 2
	}
	for range clicks {
		if err := callSendInput(&mouseInput{inputType: inputMouse, dwFlags: downFlag}); err != nil {
			return err
		}
		if err := callSendInput(&mouseInput{inputType: inputMouse, dwFlags: upFlag}); err != nil {
			return err
		}
	}
	return nil
}

func sendScroll(dy int32) error {
	// Win32: positive delta = scroll up, negative = scroll down.
	// Multiplica por wheelDelta (120) porque o Windows acumula deltas fracionados
	// internamente; 120 equivale a exatamente um "notch" da roda física.
	delta := dy * wheelDelta
	return callSendInput(&mouseInput{
		inputType: inputMouse,
		mouseData: uint32(delta),
		dwFlags:   mouseFlagWheel,
	})
}

// callSendInput chama a API Win32 SendInput com um único evento de mouse.
// Assinatura Win32: UINT SendInput(UINT cInputs, LPINPUT pInputs, int cbSize)
//   - cInputs: número de structs no array (sempre 1 aqui)
//   - pInputs: ponteiro para o array de INPUT
//   - cbSize: tamanho em bytes de cada INPUT (necessário para compatibilidade 32/64-bit)
//
// Retorna 0 em caso de falha; o erro do sistema fica em err via GetLastError.
func callSendInput(in *mouseInput) error {
	ret, _, err := fnSendInput.Call(
		1,
		uintptr(unsafe.Pointer(in)),
		uintptr(unsafe.Sizeof(*in)),
	)
	if ret == 0 {
		return fmt.Errorf("SendInput: %w", err)
	}
	return nil
}
