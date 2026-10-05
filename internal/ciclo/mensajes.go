package ciclo

import (
	"context"

	"github.com/Pacocerv08/Terminal_Game/internal/juego"
	"github.com/Pacocerv08/Terminal_Game/internal/partida"
)

// TipoMensaje distingue los mensajes que una sesión envía al ciclo.
type TipoMensaje int

const (
	MensajeUnirse TipoMensaje = iota
	MensajeAccion
	MensajeSalir
)

// Mensaje es lo único que una sesión envía al ciclo.
type Mensaje struct {
	Tipo      TipoMensaje
	Jugador   juego.IDJugador        // lo pone la sesión, nunca viene de las teclas
	Accion    juego.Accion           // solo en MensajeAccion
	Nombre    string                 // solo en MensajeUnirse
	Respuesta chan<- RespuestaUnirse // solo en MensajeUnirse (con buffer de 1)
}

// Cuadro es lo que el ciclo entrega a cada sesión después de un tick.
type Cuadro struct {
	Vista juego.Vista
	Fase  partida.Fase
}

// Suscripcion es el canal por el que una sesión recibe sus cuadros. El ciclo
// es quien lo escribe y lo cierra.
type Suscripcion struct {
	ID      juego.IDJugador
	Cuadros <-chan Cuadro
}

// RespuestaUnirse es la respuesta del ciclo a un MensajeUnirse. Err no es nil
// si no se pudo unir, por ejemplo porque la partida ya empezó.
type RespuestaUnirse struct {
	Suscripcion Suscripcion
	Err         error
}

// Unirse pide al ciclo agregar un jugador y espera su respuesta. No se queda
// bloqueada si ctx se cancela.
//
// TODO: si ctx se cancela después de enviar el mensaje, el ciclo puede crear
// un jugador que nadie va a retirar. El ciclo debe retirar a los jugadores
// cuya respuesta no se pudo entregar.
func Unirse(ctx context.Context, entrada chan<- Mensaje, nombre string) (Suscripcion, error) {
	respuesta := make(chan RespuestaUnirse, 1)
	m := Mensaje{Tipo: MensajeUnirse, Nombre: nombre, Respuesta: respuesta}
	if err := enviar(ctx, entrada, m); err != nil {
		return Suscripcion{}, err
	}
	select {
	case r := <-respuesta:
		return r.Suscripcion, r.Err
	case <-ctx.Done():
		return Suscripcion{}, ctx.Err()
	}
}

// Actuar pide al ciclo aplicar la acción a en nombre del jugador id.
func Actuar(ctx context.Context, entrada chan<- Mensaje, id juego.IDJugador, a juego.Accion) error {
	return enviar(ctx, entrada, Mensaje{Tipo: MensajeAccion, Jugador: id, Accion: a})
}

// Salir avisa al ciclo que el jugador id se fue. Devuelve el error de ctx si
// el ciclo ya terminó y nadie lee el canal, para no quedarse bloqueada.
func Salir(ctx context.Context, entrada chan<- Mensaje, id juego.IDJugador) error {
	return enviar(ctx, entrada, Mensaje{Tipo: MensajeSalir, Jugador: id})
}

// enviar entrega m al ciclo o falla si ctx se cancela antes.
func enviar(ctx context.Context, entrada chan<- Mensaje, m Mensaje) error {
	select {
	case entrada <- m:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
