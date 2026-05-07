//go:build windows

package input

import (
	"encoding/json"
	"fmt"
	"syscall"
	"unsafe"
)

var (
	user32      = syscall.NewLazyDLL("user32.dll")
	fnSendInput = user32.NewProc("SendInput")
)

const (
	inputMouse     = 0
	mouseFlagMove  = 0x0001
	mouseFlagLDown = 0x0002
	mouseFlagLUp   = 0x0004
	mouseFlagRDown = 0x0008
	mouseFlagRUp   = 0x0010
	mouseFlagMDown = 0x0020
	mouseFlagMUp   = 0x0040
	mouseFlagWheel = 0x0800
	wheelDelta     = 120
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

type baseCmd struct {
	Type string `json:"type"`
}

type mouseMoveCmd struct {
	DX int32 `json:"dx"`
	DY int32 `json:"dy"`
}

type mouseClickCmd struct {
	Button string `json:"button"` // "left", "right", "middle"
	Double bool   `json:"double"`
}

type scrollCmd struct {
	DY int32 `json:"dy"` // positive = scroll up, negative = scroll down
}

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

func sendMouseMove(dx, dy int32) error {
	return callSendInput(&mouseInput{
		inputType: inputMouse,
		dx:        dx,
		dy:        dy,
		dwFlags:   mouseFlagMove,
	})
}

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
	delta := dy * wheelDelta
	return callSendInput(&mouseInput{
		inputType: inputMouse,
		mouseData: uint32(delta),
		dwFlags:   mouseFlagWheel,
	})
}

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
