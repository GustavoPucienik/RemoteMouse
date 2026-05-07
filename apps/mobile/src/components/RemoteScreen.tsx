import { useWebSocket } from '../hooks/useWebSocket'
import Trackpad from './Trackpad'
import ClickButtons from './ClickButtons'

interface Props {
  wsUrl: string
  onDisconnect: () => void
}

export default function RemoteScreen({ wsUrl, onDisconnect }: Props) {
  const { connected, send } = useWebSocket(wsUrl)

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
        <button onClick={onDisconnect}>Desconectar</button>
      </div>
      <Trackpad send={send} />
      <ClickButtons send={send} />
    </div>
  )
}
