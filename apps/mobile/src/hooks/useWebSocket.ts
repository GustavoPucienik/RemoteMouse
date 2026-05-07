import { useRef, useCallback, useEffect, useState } from 'react'

type SendFn = (data: object) => void

export function useWebSocket(url: string): { connected: boolean; send: SendFn } {
  const ws = useRef<WebSocket | null>(null)
  const [connected, setConnected] = useState(false)

  useEffect(() => {
    const socket = new WebSocket(url)
    socket.onopen = () => setConnected(true)
    socket.onclose = () => setConnected(false)
    socket.onerror = () => setConnected(false)
    ws.current = socket
    return () => {
      socket.close()
      ws.current = null
    }
  }, [url])

  const send = useCallback((data: object) => {
    if (ws.current?.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify(data))
    }
  }, [])

  return { connected, send }
}
