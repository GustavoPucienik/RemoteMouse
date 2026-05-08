import { useRef, useEffect, useCallback } from 'react'

interface Props {
  send: (data: object) => void
}

const SPECIAL_KEYS = [
  { label: 'Esc', key: 'Escape' },
  { label: 'Tab', key: 'Tab' },
  { label: '←', key: 'ArrowLeft' },
  { label: '↑', key: 'ArrowUp' },
  { label: '↓', key: 'ArrowDown' },
  { label: '→', key: 'ArrowRight' },
  { label: 'Del', key: 'Delete' },
]

export default function Keyboard({ send }: Props) {
  const inputRef = useRef<HTMLTextAreaElement>(null)

  useEffect(() => {
    inputRef.current?.focus()
  }, [])

  const sendKey = useCallback(
    (key: string) => send({ type: 'key_press', key }),
    [send],
  )

  function onInput(e: React.FormEvent<HTMLTextAreaElement>) {
    const ev = e.nativeEvent as InputEvent
    const target = e.target as HTMLTextAreaElement

    if (ev.inputType === 'deleteContentBackward') {
      sendKey('Backspace')
    } else if (ev.inputType === 'insertLineBreak' || ev.inputType === 'insertParagraph') {
      sendKey('Enter')
    } else if (ev.data) {
      for (const char of ev.data) sendKey(char)
    }

    // Keep textarea empty so future inputs are always detected
    target.value = ''
  }

  function onKeyDown(e: React.KeyboardEvent<HTMLTextAreaElement>) {
    // On desktop, special keys don't always fire an input event (e.g. Backspace
    // on an empty textarea). Catch them here; preventDefault stops the input event
    // from also firing, preventing double-send.
    const special = new Set([
      'Backspace', 'Tab', 'Enter', 'Escape', 'Delete',
      'ArrowLeft', 'ArrowRight', 'ArrowUp', 'ArrowDown',
      'Home', 'End', 'PageUp', 'PageDown',
      'F1', 'F2', 'F3', 'F4', 'F5', 'F6',
      'F7', 'F8', 'F9', 'F10', 'F11', 'F12',
    ])
    if (special.has(e.key)) {
      e.preventDefault()
      sendKey(e.key)
    }
  }

  return (
    <div className="keyboard-panel">
      <div className="keyboard-special-keys">
        {SPECIAL_KEYS.map(({ label, key }) => (
          <button
            key={key}
            className="key-btn"
            onPointerDown={(e) => { e.preventDefault(); sendKey(key) }}
          >
            {label}
          </button>
        ))}
      </div>
      <textarea
        ref={inputRef}
        className="keyboard-input"
        onInput={onInput}
        onKeyDown={onKeyDown}
        autoComplete="off"
        autoCorrect="off"
        autoCapitalize="off"
        spellCheck={false}
        rows={3}
        placeholder="Toque aqui e comece a digitar..."
      />
    </div>
  )
}
