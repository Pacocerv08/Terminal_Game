package red

import (
	"context"

	"github.com/Pacocerv08/Terminal_Game/internal/ciclo"
	"github.com/Pacocerv08/Terminal_Game/internal/entrada"
	"github.com/Pacocerv08/Terminal_Game/internal/graficos"
	"github.com/Pacocerv08/Terminal_Game/internal/juego"
)

// Sesion es la conexión de un jugador. Guarda el ID que le dio el ciclo y lo
// pone en cada mensaje: así la identidad la asigna el servidor, no el cliente.
type Sesion struct {
	ID        juego.IDJugador
	mensajes  chan<- ciclo.Mensaje
	dibujante graficos.Dibujante
	traductor entrada.Traductor
}

// Atender une al jugador al ciclo, traduce y valida sus teclas, dibuja cada
// Cuadro que llega y, al terminar, avisa con ciclo.Salir.
func (s *Sesion) Atender(ctx context.Context, nombre string) error {
	// TODO: ciclo.Unirse, leer teclas, ciclo.Actuar, dibujar cuadros y
	// defer ciclo.Salir (con un ctx propio con plazo, porque ctx ya puede
	// estar cancelado).
	return nil
}
