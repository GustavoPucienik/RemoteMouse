import { useState } from 'react'
import ConnectScreen from './components/ConnectScreen'
import RemoteScreen from './components/RemoteScreen'

export default function App() {
  const [wsUrl, setWsUrl] = useState<string | null>(null)

  if (!wsUrl) {
    return (
      <ConnectScreen
        onConnect={(ip, port) => setWsUrl(`ws://${ip}:${port}/ws`)}
      />
    )
  }

  return (
    <RemoteScreen
      wsUrl={wsUrl}
      onDisconnect={() => setWsUrl(null)}
    />
  )
}
