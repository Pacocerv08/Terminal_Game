package entrada

import "github.com/Pacocerv08/Terminal_Game/internal/juego"

// Traductor convierte una tecla en una acción del juego.
type Traductor interface {
	// Traducir devuelve la acción de la tecla. El segundo valor es false si
	// la tecla no corresponde a ninguna acción.
	Traducir(tecla string) (juego.Accion, bool)
}

// TraductorTeclado es el Traductor para teclado: flechas o WASD para mover,
// una tecla para disparar y otra para recoger.
type TraductorTeclado struct{}

// Traducir implementa Traductor.
func (TraductorTeclado) Traducir(tecla string) (juego.Accion, bool) {
	// TODO: mapear teclas a acciones.
	return nil, false
}
