package ciclo

import (
	"context"
	"time"

	"github.com/Pacocerv08/Terminal_Game/internal/bots"
	"github.com/Pacocerv08/Terminal_Game/internal/juego"
	"github.com/Pacocerv08/Terminal_Game/internal/partida"
)

// Ciclo avanza el juego un tick cada intervalo. Es el único que modifica el
// estado. El orden de cada tick está descrito en la documentación del paquete.
type Ciclo struct {
	intervalo     time.Duration
	estado        *juego.Estado
	maquina       *partida.Maquina
	entrada       <-chan Mensaje
	suscripciones map[juego.IDJugador]chan Cuadro
	controlados   []bots.Bot
}

// Nuevo crea un ciclo. Lee mensajes de entrada y no escribe en él.
func Nuevo(intervalo time.Duration, e *juego.Estado, m *partida.Maquina, entrada <-chan Mensaje) *Ciclo {
	return &Ciclo{
		intervalo:     intervalo,
		estado:        e,
		maquina:       m,
		entrada:       entrada,
		suscripciones: make(map[juego.IDJugador]chan Cuadro),
	}
}

// RegistrarBot agrega un bot. Debe llamarse antes de Ejecutar.
func (c *Ciclo) RegistrarBot(b bots.Bot) {
	c.controlados = append(c.controlados, b)
}

// Ejecutar corre el ciclo hasta que ctx se cancele. Al terminar cierra todas
// las suscripciones.
func (c *Ciclo) Ejecutar(ctx context.Context) {
	// TODO: ticker a intervalo fijo y los seis pasos de cada tick, con una
	// sola acción por jugador por tick (la última que llegue).
}
