import { useRef } from 'react'

interface Props {
  send: (data: object) => void
}

interface Point {
  x: number
  y: number
}

export default function Trackpad({ send }: Props) {
  const lastSingle = useRef<Point | null>(null)
  const lastTwoY = useRef<number | null>(null)
  const scrollAccum = useRef(0)

  function onTouchStart(e: React.TouchEvent) {
    if (e.touches.length === 1) {
      lastSingle.current = { x: e.touches[0].clientX, y: e.touches[0].clientY }
      lastTwoY.current = null
      scrollAccum.current = 0
    } else if (e.touches.length === 2) {
      lastTwoY.current = (e.touches[0].clientY + e.touches[1].clientY) / 2
      lastSingle.current = null
    }
  }

  function onTouchMove(e: React.TouchEvent) {
    if (e.touches.length === 1 && lastSingle.current) {
      const dx = Math.round(e.touches[0].clientX - lastSingle.current.x)
      const dy = Math.round(e.touches[0].clientY - lastSingle.current.y)
      if (dx !== 0 || dy !== 0) {
        send({ type: 'mouse_move', dx, dy })
      }
      lastSingle.current = { x: e.touches[0].clientX, y: e.touches[0].clientY }
    } else if (e.touches.length === 2 && lastTwoY.current !== null) {
      const avgY = (e.touches[0].clientY + e.touches[1].clientY) / 2
      scrollAccum.current += avgY - lastTwoY.current
      lastTwoY.current = avgY

      // Send one notch per ~30px of finger movement
      const notches = Math.trunc(scrollAccum.current / 30)
      if (notches !== 0) {
        // fingers move down (positive) → scroll down → negative protocol dy
        send({ type: 'scroll', dy: -notches })
        scrollAccum.current -= notches * 30
      }
    }
  }

  function onTouchEnd() {
    lastSingle.current = null
    lastTwoY.current = null
    scrollAccum.current = 0
  }

  return (
    <div
      className="trackpad"
      onTouchStart={onTouchStart}
      onTouchMove={onTouchMove}
      onTouchEnd={onTouchEnd}
    >
      <span className="trackpad-hint">1 dedo · mover &nbsp;&nbsp; 2 dedos · scroll</span>
    </div>
  )
}
