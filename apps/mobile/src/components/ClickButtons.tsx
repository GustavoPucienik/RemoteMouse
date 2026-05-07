interface Props {
  send: (data: object) => void
}

export default function ClickButtons({ send }: Props) {
  function click(button: 'left' | 'right') {
    send({ type: 'mouse_click', button, double: false })
  }

  return (
    <div className="click-buttons">
      <button onPointerDown={() => click('left')}>Esquerdo</button>
      <button onPointerDown={() => click('right')}>Direito</button>
    </div>
  )
}
