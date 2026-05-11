//go:build windows

package input

import (
	"fmt"
	"unsafe"
)

const (
	inputKeyboard   = 1
	keyEventKeyUp   = 0x0002
	keyEventUnicode = 0x0004
	keyEventExtended = 0x0001 // arrows, Delete, Home, End, PgUp, PgDn
)

// keyboardInput mirrors the Win32 INPUT struct for keyboard events on 64-bit Windows.
// Layout: type(4) + pad(4) + wVk(2) + wScan(2) + dwFlags(4) + time(4) + pad+dwExtraInfo+trailing(20) = 40 bytes
type keyboardInput struct {
	inputType uint32
	_         [4]byte
	wVk       uint16
	wScan     uint16
	dwFlags   uint32
	time      uint32
	_         [20]byte // alignment(4) + dwExtraInfo(8) + trailing padding(8)
}

// vkMap maps Web KeyboardEvent.key names to Win32 Virtual Key codes.
var vkMap = map[string]uint16{
	"Backspace":  0x08,
	"Tab":        0x09,
	"Enter":      0x0D,
	"Escape":     0x1B,
	"PageUp":     0x21,
	"PageDown":   0x22,
	"End":        0x23,
	"Home":       0x24,
	"ArrowLeft":  0x25,
	"ArrowUp":    0x26,
	"ArrowRight": 0x27,
	"ArrowDown":  0x28,
	"Delete":     0x2E,
	"F1":  0x70, "F2":  0x71, "F3":  0x72, "F4":  0x73,
	"F5":  0x74, "F6":  0x75, "F7":  0x76, "F8":  0x77,
	"F9":  0x78, "F10": 0x79, "F11": 0x7A, "F12": 0x7B,
}

// extendedVK holds VK codes that require KEYEVENTF_EXTENDEDKEY to work correctly
// outside the numpad (arrows, navigation cluster).
var extendedVK = map[uint16]bool{
	0x21: true, 0x22: true, 0x23: true, 0x24: true,
	0x25: true, 0x26: true, 0x27: true, 0x28: true,
	0x2E: true,
}

// sendKeyPress sends a key-down + key-up pair for the given key name.
// Named keys (e.g. "Enter", "ArrowLeft") use Win32 virtual key codes.
// Single characters are sent as Unicode events and work with any keyboard layout.
func sendKeyPress(key string) error {
	if vk, ok := vkMap[key]; ok {
		return sendVirtualKey(vk)
	}
	for _, r := range key {
		if err := sendUnicodeChar(r); err != nil {
			return err
		}
	}
	return nil
}

func sendVirtualKey(vk uint16) error {
	flags := uint32(0)
	if extendedVK[vk] {
		flags = keyEventExtended
	}
	if err := callSendInputKb(&keyboardInput{inputType: inputKeyboard, wVk: vk, dwFlags: flags}); err != nil {
		return err
	}
	return callSendInputKb(&keyboardInput{inputType: inputKeyboard, wVk: vk, dwFlags: flags | keyEventKeyUp})
}

func sendUnicodeChar(r rune) error {
	scan := uint16(r)
	if err := callSendInputKb(&keyboardInput{inputType: inputKeyboard, wScan: scan, dwFlags: keyEventUnicode}); err != nil {
		return err
	}
	return callSendInputKb(&keyboardInput{inputType: inputKeyboard, wScan: scan, dwFlags: keyEventUnicode | keyEventKeyUp})
}

func callSendInputKb(in *keyboardInput) error {
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
