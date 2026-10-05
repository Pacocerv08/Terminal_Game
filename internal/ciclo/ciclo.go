package ciclo

import (
	"context"
	"fmt"
	"math/rand/v2"
	"slices"
	"time"

	"github.com/Pacocerv08/Terminal_Game/internal/bots"
	"github.com/Pacocerv08/Terminal_Game/internal/juego"
	"github.com/Pacocerv08/Terminal_Game/internal/partida"
)

// botControlado es un bot junto con el jugador que maneja.
type botControlado struct {
	id  juego.IDJugador
	bot bots.Bot
}

// Ciclo avanza el juego un tick cada intervalo. Es el único que modifica el
// estado. El orden de cada tick está descrito en la documentación del paquete.
type Ciclo struct {
	intervalo     time.Duration
	estado        *juego.Estado
	maquina       *partida.Maquina
	entrada       <-chan Mensaje
	suscripciones map[juego.IDJugador]chan Cuadro
	controlados   []botControlado
	semilla       uint64
}

// Nuevo crea un ciclo. Lee mensajes de entrada y no escribe en él. Entra en
// pánico si intervalo no es positivo, porque es un error de programación.
func Nuevo(intervalo time.Duration, e *juego.Estado, m *partida.Maquina, entrada <-chan Mensaje) *Ciclo {
	if intervalo <= 0 {
		panic("ciclo: el intervalo debe ser positivo")
	}
	return &Ciclo{
		intervalo:     intervalo,
		estado:        e,
		maquina:       m,
		entrada:       entrada,
		suscripciones: make(map[juego.IDJugador]chan Cuadro),
	}
}

// FijarSemilla fija la semilla del orden de desempate (ver ordenDeAplicacion).
// Por defecto vale 0. Debe llamarse antes de Ejecutar.
func (c *Ciclo) FijarSemilla(semilla uint64) {
	c.semilla = semilla
}

// RegistrarBot agrega un bot como jugador y devuelve su ID. Devuelve
// ErrNoSePuedeUnir o ErrSinLugar si no se puede agregar. Debe llamarse antes
// de Ejecutar.
func (c *Ciclo) RegistrarBot(b bots.Bot) (juego.IDJugador, error) {
	if !c.maquina.PermiteUnirse() {
		return 0, ErrNoSePuedeUnir
	}
	j := c.estado.AgregarJugador(fmt.Sprintf("bot-%d", len(c.controlados)+1))
	if j == nil {
		return 0, ErrSinLugar
	}
	c.controlados = append(c.controlados, botControlado{id: j.ID, bot: b})
	return j.ID, nil
}

// Ejecutar llama a tick cada intervalo hasta que ctx se cancele. Al terminar
// cierra todas las suscripciones.
func (c *Ciclo) Ejecutar(ctx context.Context) {
	defer c.cerrarSuscripciones()
	reloj := time.NewTicker(c.intervalo)
	defer reloj.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-reloj.C:
			c.tick()
		}
	}
}

// tick ejecuta un solo paso del juego, sin depender del reloj. El orden está
// descrito en la documentación del paquete.
func (c *Ciclo) tick() {
	pendientes := make(map[juego.IDJugador]juego.Accion)
	c.procesarMensajes(pendientes)
	c.recolectarBots(pendientes)
	c.aplicar(pendientes)
	c.estado.Avanzar()
	c.maquina.Avanzar(c.estado)
	c.publicar()
}

// procesarMensajes atiende los mensajes que había en la cola al empezar el
// tick, en orden de llegada. Los que lleguen mientras tanto esperan al
// siguiente tick, así un productor rápido no impide que el tick termine. De
// las acciones se recuerda la última de cada jugador en pendientes.
func (c *Ciclo) procesarMensajes(pendientes map[juego.IDJugador]juego.Accion) {
	for n := len(c.entrada); n > 0; n-- {
		select {
		case m := <-c.entrada:
			c.procesar(m, pendientes)
		default:
			return
		}
	}
}

func (c *Ciclo) procesar(m Mensaje, pendientes map[juego.IDJugador]juego.Accion) {
	switch m.tipo {
	case mensajeUnirse:
		c.unir(m)
	case mensajeAccion:
		if m.accion == nil {
			return
		}
		// Solo actúan los jugadores con sesión; los bots no pasan por aquí.
		if _, ok := c.suscripciones[m.jugador]; !ok {
			return
		}
		pendientes[m.jugador] = m.accion
	case mensajeSalir:
		c.retirar(m.jugador)
		delete(pendientes, m.jugador)
	}
	// mensajeNulo (el valor cero) y cualquier otro valor se ignoran.
}

// unir atiende un mensajeUnirse. Si la respuesta no se puede entregar porque
// quien la esperaba se canceló, retira al jugador que acaba de agregar.
func (c *Ciclo) unir(m Mensaje) {
	if m.respuesta == nil {
		return
	}
	var r RespuestaUnirse
	switch {
	case !c.maquina.PermiteUnirse():
		r.Err = ErrNoSePuedeUnir
	default:
		j := c.estado.AgregarJugador(m.nombre)
		if j == nil {
			r.Err = ErrSinLugar
			break
		}
		cuadros := make(chan Cuadro, 1)
		c.suscripciones[j.ID] = cuadros
		r.Suscripcion = Suscripcion{ID: j.ID, Cuadros: cuadros}
	}

	select {
	case m.respuesta <- r:
	case <-m.cancelado:
		if r.Err == nil {
			c.retirar(r.Suscripcion.ID)
		}
	}
}

// retirar saca de la partida al jugador id y cierra su suscripción. Si el
// jugador no tiene suscripción, no hace nada.
func (c *Ciclo) retirar(id juego.IDJugador) {
	cuadros, ok := c.suscripciones[id]
	if !ok {
		return
	}
	c.estado.QuitarJugador(id)
	close(cuadros)
	delete(c.suscripciones, id)
}

// recolectarBots pide su acción a cada bot, todos con el estado de inicio de
// tick, y la suma a las pendientes.
func (c *Ciclo) recolectarBots(pendientes map[juego.IDJugador]juego.Accion) {
	for _, b := range c.controlados {
		j, ok := c.estado.Jugadores[b.id]
		if !ok || !j.EstaVivo() {
			continue
		}
		if a := b.bot.Decidir(c.estado.VistaPara(b.id)); a != nil {
			pendientes[b.id] = a
		}
	}
}

// aplicar aplica las acciones pendientes en el orden de ordenDeAplicacion,
// saltando las que la fase no permite.
func (c *Ciclo) aplicar(pendientes map[juego.IDJugador]juego.Accion) {
	ids := make([]juego.IDJugador, 0, len(pendientes))
	for id := range pendientes {
		ids = append(ids, id)
	}
	for _, id := range ordenDeAplicacion(ids, c.semilla, c.estado.Tick) {
		if a := pendientes[id]; c.maquina.AccionPermitida(a) {
			a.Aplicar(c.estado, id)
		}
	}
}

// ordenDeAplicacion devuelve los ids en el orden en que se aplican sus
// acciones. Primero los ordena de menor a mayor, para que el resultado no
// dependa del orden en que llegaron (por ejemplo, del recorrido de un map), y
// luego los baraja con un generador sembrado con (semilla, tick). Así el
// orden es determinista, pero no favorece siempre a los mismos jugadores: en
// un empate cada uno gana cerca de la mitad de las veces a largo plazo.
func ordenDeAplicacion(ids []juego.IDJugador, semilla, tick uint64) []juego.IDJugador {
	orden := slices.Clone(ids)
	slices.Sort(orden)
	azar := rand.New(rand.NewPCG(semilla, tick))
	azar.Shuffle(len(orden), func(i, j int) { orden[i], orden[j] = orden[j], orden[i] })
	return orden
}

// publicar envía un Cuadro a cada suscripción sin bloquear nunca: si el
// canal tiene un cuadro viejo sin leer, lo descarta y pone el nuevo.
func (c *Ciclo) publicar() {
	fase := c.maquina.Fase()
	for id, cuadros := range c.suscripciones {
		cuadro := Cuadro{Vista: c.estado.VistaPara(id), Fase: fase}
		select {
		case cuadros <- cuadro:
			continue
		default:
		}
		select {
		case <-cuadros:
		default:
		}
		// El ciclo es el único que escribe en el canal, así que tras vaciarlo
		// cabe un cuadro; el default solo evita bloquear ante lo imposible.
		select {
		case cuadros <- cuadro:
		default:
		}
	}
}

// cerrarSuscripciones cierra el canal de cada suscripción, para que las
// sesiones sepan que el ciclo terminó.
func (c *Ciclo) cerrarSuscripciones() {
	for id, cuadros := range c.suscripciones {
		close(cuadros)
		delete(c.suscripciones, id)
	}
}
