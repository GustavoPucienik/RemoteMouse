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
  // composingRef: true while Android IME (GBoard) is in the middle of composition.
  // Clearing target.value during composition causes GBoard to fire a ghost
  // deleteContentBackward that cancels every character typed — so we must
  // wait for compositionEnd before touching the textarea.
  const composingRef = useRef(false)
  // After compositionEnd we send the committed text ourselves.  The browser
  // then fires a redundant insertText with the same chars (plus any suffix
  // like a trailing space added by the commit key).  pendingCompositionRef
  // holds what we already sent so onInput can skip the duplicate and only
  // forward the extra suffix.
  const pendingCompositionRef = useRef<string | null>(null)

  useEffect(() => {
    inputRef.current?.focus()
  }, [])

  const sendKey = useCallback(
    (key: string) => send({ type: 'key_press', key }),
    [send],
  )

  function onCompositionStart() {
    composingRef.current = true
    pendingCompositionRef.current = null
  }

  function onCompositionEnd(e: React.CompositionEvent<HTMLTextAreaElement>) {
    composingRef.current = false
    const text = e.data ?? ''
    pendingCompositionRef.current = text
    for (const char of text) sendKey(char)
    e.currentTarget.value = ''
  }

  function onInput(e: React.FormEvent<HTMLTextAreaElement>) {
    const ev = e.nativeEvent as InputEvent
    const target = e.target as HTMLTextAreaElement

    // During composition do nothing — compositionEnd will handle it.
    if (composingRef.current || ev.inputType === 'insertCompositionText') return

    const pending = pendingCompositionRef.current
    pendingCompositionRef.current = null

    if (pending !== null && ev.inputType === 'insertText') {
      // This is the insertText that fires right after compositionEnd.
      // We already sent `pending`; only send any suffix (e.g. the space
      // that the user pressed to commit the composition word).
      const suffix = (ev.data ?? '').startsWith(pending)
        ? (ev.data ?? '').slice(pending.length)
        : ''
      for (const char of suffix) sendKey(char)
      target.value = ''
      return
    }

    if (ev.inputType === 'deleteContentBackward') {
      sendKey('Backspace')
    } else if (ev.inputType === 'insertLineBreak' || ev.inputType === 'insertParagraph') {
      sendKey('Enter')
    } else if (ev.data) {
      for (const char of ev.data) sendKey(char)
    }

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
        onCompositionStart={onCompositionStart}
        onCompositionEnd={onCompositionEnd}
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
