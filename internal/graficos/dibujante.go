package graficos

import (
	"github.com/Pacocerv08/Terminal_Game/internal/juego"
	"github.com/Pacocerv08/Terminal_Game/internal/partida"
)

// Dibujante convierte lo que ve un jugador en el texto que se muestra en su
// terminal.
type Dibujante interface {
	// Dibujar devuelve la pantalla completa para una terminal de
	// ancho x alto caracteres. La fase permite mostrar la sala de espera o
	// el resultado final.
	Dibujar(v juego.Vista, fase partida.Fase, ancho, alto int) string
}
