package bots

import "github.com/Pacocerv08/Terminal_Game/internal/juego"

// Bot es un jugador controlado por el servidor.
type Bot interface {
	// Decidir elige la acción del bot en este tick a partir de lo que ve.
	Decidir(v juego.Vista) juego.Accion
}
