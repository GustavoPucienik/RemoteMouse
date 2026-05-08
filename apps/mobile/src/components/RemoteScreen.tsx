import { useState } from 'react'
import { useWebSocket } from '../hooks/useWebSocket'
import Trackpad from './Trackpad'
import ClickButtons from './ClickButtons'
import Keyboard from './Keyboard'

interface Props {
  wsUrl: string
  onDisconnect: () => void
}

type Mode = 'mouse' | 'keyboard'

export default function RemoteScreen({ wsUrl, onDisconnect }: Props) {
  const { connected, send } = useWebSocket(wsUrl)
  const [mode, setMode] = useState<Mode>('mouse')

  if (!connected) {
    return (
      <div className="connecting">
        <span>Conectando...</span>
        <button onClick={onDisconnect}>Cancelar</button>
      </div>
    )
  }

  return (
    <div className="remote">
      <div className="status-bar">
        <span>
          <span className="status-dot" />
          Conectado
        </span>
        <div className="status-actions">
          <button
            className={`mode-btn${mode === 'keyboard' ? ' active' : ''}`}
            onClick={() => setMode(m => m === 'mouse' ? 'keyboard' : 'mouse')}
          >
            ⌨
          </button>
          <button onClick={onDisconnect}>Desconectar</button>
        </div>
      </div>

      {mode === 'mouse' ? (
        <>
          <Trackpad send={send} />
          <ClickButtons send={send} />
        </>
      ) : (
        <Keyboard send={send} />
      )}
    </div>
  )
}
