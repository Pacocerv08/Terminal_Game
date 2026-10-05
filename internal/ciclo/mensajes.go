package ciclo

import (
	"context"
	"errors"

	"github.com/Pacocerv08/Terminal_Game/internal/juego"
	"github.com/Pacocerv08/Terminal_Game/internal/partida"
)

// Errores con los que el ciclo rechaza un Unirse.
var (
	ErrNoSePuedeUnir = errors.New("la partida no admite jugadores nuevos")
	ErrSinLugar      = errors.New("no hay una casilla libre en el mapa")
)

// tipoMensaje distingue los mensajes que una sesión envía al ciclo. El valor
// cero no es ningún mensaje.
type tipoMensaje int

const (
	mensajeNulo tipoMensaje = iota
	mensajeUnirse
	mensajeAccion
	mensajeSalir
)

// Mensaje es lo único que una sesión envía al ciclo. Sus campos son privados:
// solo Unirse, Actuar y Salir construyen mensajes válidos. El valor cero de
// Mensaje no significa nada y el ciclo lo ignora.
type Mensaje struct {
	tipo    tipoMensaje
	jugador juego.IDJugador // lo pone la sesión, nunca viene de las teclas
	accion  juego.Accion    // solo en mensajeAccion
	nombre  string          // solo en mensajeUnirse

	// Solo en mensajeUnirse. respuesta no tiene buffer: el ciclo la entrega
	// con un select contra cancelado, y quien envía el mensaje siempre está
	// esperando en un select entre respuesta y su contexto. Si cancelado se
	// cierra antes de que la respuesta se entregue, el ciclo retira al
	// jugador que acaba de agregar.
	respuesta chan<- RespuestaUnirse
	cancelado <-chan struct{}
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

// RespuestaUnirse es la respuesta del ciclo a un Unirse. Err no es nil si no
// se pudo unir: ErrNoSePuedeUnir o ErrSinLugar.
type RespuestaUnirse struct {
	Suscripcion Suscripcion
	Err         error
}

// Unirse pide al ciclo agregar un jugador y espera su respuesta. No se queda
// bloqueada si ctx se cancela. Si ctx se cancela después de que el ciclo
// recibió el mensaje, el ciclo retira al jugador que había agregado.
//
// ctx debe cancelarse cuando el ciclo termine, por ejemplo derivándolo del
// contexto del servidor: si no, un Unirse que quedó en la cola de un ciclo
// ya terminado esperaría para siempre.
func Unirse(ctx context.Context, entrada chan<- Mensaje, nombre string) (Suscripcion, error) {
	respuesta := make(chan RespuestaUnirse)
	m := Mensaje{tipo: mensajeUnirse, nombre: nombre, respuesta: respuesta, cancelado: ctx.Done()}
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

// Actuar pide al ciclo aplicar la acción a en nombre del jugador id. El ciclo
// aplica como máximo una acción por jugador por tick y conserva la última.
func Actuar(ctx context.Context, entrada chan<- Mensaje, id juego.IDJugador, a juego.Accion) error {
	return enviar(ctx, entrada, Mensaje{tipo: mensajeAccion, jugador: id, accion: a})
}

// Salir avisa al ciclo que el jugador id se fue.
//
// Intenta el envío primero y solo después atiende ctx.Done(): si la cola
// tiene espacio, el aviso se entrega aunque ctx ya esté cancelado. Por eso ctx
// debe ser el del servidor y no el de la sesión: la sesión suele terminar
// justamente porque su contexto se canceló, y aun así debe poder retirar a su
// jugador. Devuelve el error de ctx solo si la cola está llena y ctx se
// cancela, para no quedarse bloqueada con un ciclo que ya terminó.
func Salir(ctx context.Context, entrada chan<- Mensaje, id juego.IDJugador) error {
	m := Mensaje{tipo: mensajeSalir, jugador: id}
	select {
	case entrada <- m:
		return nil
	default:
	}
	return enviar(ctx, entrada, m)
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
