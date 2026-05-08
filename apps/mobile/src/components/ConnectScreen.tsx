import { useState } from 'react'

const STORAGE_IP_KEY = 'remotemouse_ip'
const STORAGE_PORT_KEY = 'remotemouse_port'

interface Props {
  onConnect: (ip: string, port: string) => void
}

export default function ConnectScreen({ onConnect }: Props) {
  const [ip, setIp] = useState(() => localStorage.getItem(STORAGE_IP_KEY) ?? '')
  const [port, setPort] = useState(() => localStorage.getItem(STORAGE_PORT_KEY) ?? '8080')

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    const trimmedIp = ip.trim()
    const trimmedPort = port.trim()
    if (!trimmedIp) return
    localStorage.setItem(STORAGE_IP_KEY, trimmedIp)
    localStorage.setItem(STORAGE_PORT_KEY, trimmedPort)
    onConnect(trimmedIp, trimmedPort)
  }

  return (
    <div className="connect-screen">
      <h1>RemoteMouse</h1>
      <p>Controle seu PC pelo celular</p>
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          placeholder="IP do computador (ex: 192.168.1.10)"
          value={ip}
          onChange={e => setIp(e.target.value)}
          inputMode="numeric"
          autoComplete="off"
          autoCorrect="off"
          spellCheck={false}
        />
        <input
          type="text"
          placeholder="Porta"
          value={port}
          onChange={e => setPort(e.target.value)}
          inputMode="numeric"
        />
        <button type="submit">Conectar</button>
      </form>
    </div>
  )
}
