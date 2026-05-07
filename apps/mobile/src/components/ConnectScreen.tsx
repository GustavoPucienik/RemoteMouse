import { useState } from 'react'

interface Props {
  onConnect: (ip: string, port: string) => void
}

export default function ConnectScreen({ onConnect }: Props) {
  const [ip, setIp] = useState('')
  const [port, setPort] = useState('8080')

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (ip.trim()) onConnect(ip.trim(), port.trim())
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
