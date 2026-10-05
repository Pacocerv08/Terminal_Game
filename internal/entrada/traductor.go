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

// Traducir implementa Traductor. Mover: w/a/s/d, también en mayúscula por si
// el jugador tiene activado el bloqueo de mayúsculas, y las flechas "up",
// "down", "left" y "right". Cualquier otra tecla devuelve false.
//
// TODO: disparar y recoger.
func (TraductorTeclado) Traducir(tecla string) (juego.Accion, bool) {
	switch tecla {
	case "w", "W", "up":
		return juego.MoverAccion{Dir: juego.Arriba}, true
	case "s", "S", "down":
		return juego.MoverAccion{Dir: juego.Abajo}, true
	case "a", "A", "left":
		return juego.MoverAccion{Dir: juego.Izquierda}, true
	case "d", "D", "right":
		return juego.MoverAccion{Dir: juego.Derecha}, true
	}
	return nil, false
}
