package ciclo

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Pacocerv08/Terminal_Game/internal/bots"
	"github.com/Pacocerv08/Terminal_Game/internal/juego"
	"github.com/Pacocerv08/Terminal_Game/internal/partida"
)

// Estas pruebas llaman a tick() directamente: no hay temporizadores ni
// esperas de tiempo real. Para sincronizar con otras goroutines solo se usa
// esperar, con un límite de intentos.

// maxIntentos limita las esperas con runtime.Gosched().
const maxIntentos = 1_000_000

// esperar cede el procesador hasta que cumple() sea verdadero. Si tras
// maxIntentos intentos no se cumple, la prueba falla con descripcion.
func esperar(t *testing.T, descripcion string, cumple func() bool) {
	t.Helper()
	for i := 0; i < maxIntentos; i++ {
		if cumple() {
			return
		}
		runtime.Gosched()
	}
	t.Fatalf("tras %d intentos no se cumplió: %s", maxIntentos, descripcion)
}

// recibir espera un valor de canal con el mismo límite de intentos. Un canal
// cerrado devuelve el valor cero enseguida.
func recibir[T any](t *testing.T, canal <-chan T, descripcion string) T {
	t.Helper()
	var valor T
	esperar(t, descripcion, func() bool {
		select {
		case valor = <-canal:
			return true
		default:
			return false
		}
	})
	return valor
}

// esperarCierre vacía el canal hasta que esté cerrado.
func esperarCierre(t *testing.T, cuadros <-chan Cuadro, descripcion string) {
	t.Helper()
	esperar(t, descripcion, func() bool {
		select {
		case _, abierto := <-cuadros:
			return !abierto
		default:
			return false
		}
	})
}

const mapaLibre = `
	.....
	.....
	.....
`

// banco agrupa lo que necesita una prueba del ciclo.
type banco struct {
	c       *Ciclo
	estado  *juego.Estado
	entrada chan Mensaje
}

func nuevoBanco(t *testing.T, mapa string, cola int) *banco {
	t.Helper()
	m, err := juego.MapaDesdeTexto(mapa)
	if err != nil {
		t.Fatalf("MapaDesdeTexto: %v", err)
	}
	estado := juego.NuevoEstado(m)
	entrada := make(chan Mensaje, cola)
	return &banco{
		c:       Nuevo(time.Hour, estado, partida.NuevaMaquina(), entrada),
		estado:  estado,
		entrada: entrada,
	}
}

type resultadoUnirse struct {
	s   Suscripcion
	err error
}

// iniciarUnirse lanza Unirse en otra goroutine y espera a que su mensaje esté
// en la cola, sin ejecutar ningún tick.
func (b *banco) iniciarUnirse(t *testing.T, ctx context.Context, nombre string) <-chan resultadoUnirse {
	t.Helper()
	antes := len(b.entrada)
	res := make(chan resultadoUnirse, 1)
	go func() {
		s, err := Unirse(ctx, b.entrada, nombre)
		res <- resultadoUnirse{s, err}
	}()
	esperar(t, "el mensaje de Unirse en la cola", func() bool { return len(b.entrada) == antes+1 })
	return res
}

// unir hace un Unirse completo: lo lanza, ejecuta un tick y devuelve la
// respuesta.
func (b *banco) unir(t *testing.T, nombre string) (Suscripcion, error) {
	t.Helper()
	res := b.iniciarUnirse(t, context.Background(), nombre)
	b.c.tick()
	r := recibir(t, res, "la respuesta de Unirse")
	return r.s, r.err
}

func (b *banco) unirOk(t *testing.T, nombre string) Suscripcion {
	t.Helper()
	s, err := b.unir(t, nombre)
	if err != nil {
		t.Fatalf("Unirse(%q): %v", nombre, err)
	}
	return s
}

func (b *banco) actuar(t *testing.T, id juego.IDJugador, a juego.Accion) {
	t.Helper()
	if err := Actuar(context.Background(), b.entrada, id, a); err != nil {
		t.Fatalf("Actuar: %v", err)
	}
}

func (b *banco) pos(id juego.IDJugador) juego.Posicion {
	return b.estado.Jugadores[id].Pos
}

func TestNuevoConIntervaloInvalido(t *testing.T) {
	for _, intervalo := range []time.Duration{0, -time.Second} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("Nuevo(%v) debía entrar en pánico", intervalo)
				}
			}()
			Nuevo(intervalo, nil, nil, nil)
		}()
	}
}

func TestUnirse(t *testing.T) {
	b := nuevoBanco(t, mapaLibre, 16)
	s := b.unirOk(t, "ana")

	if s.ID != 1 {
		t.Errorf("ID = %d, quiere 1", s.ID)
	}
	j := b.estado.Jugadores[s.ID]
	if j == nil || j.Nombre != "ana" || j.Pos != (juego.Posicion{X: 0, Y: 0}) {
		t.Fatalf("jugador = %+v, quiere ana en (0,0)", j)
	}
	cuadro := recibir(t, s.Cuadros, "el cuadro del tick de unirse")
	if cuadro.Vista.Yo != s.ID || cuadro.Fase != partida.SalaDeEspera || cuadro.Vista.Tick != 1 {
		t.Errorf("cuadro = Yo %d, Fase %v, Tick %d; quiere Yo %d, SalaDeEspera, Tick 1",
			cuadro.Vista.Yo, cuadro.Fase, cuadro.Vista.Tick, s.ID)
	}
}

func TestUnirseDosJugadores(t *testing.T) {
	b := nuevoBanco(t, mapaLibre, 16)
	s1 := b.unirOk(t, "ana")
	s2 := b.unirOk(t, "beto")
	if s1.ID != 1 || s2.ID != 2 {
		t.Errorf("IDs = %d y %d, quiere 1 y 2", s1.ID, s2.ID)
	}
	if b.pos(1) == b.pos(2) {
		t.Errorf("los dos jugadores están en %v", b.pos(1))
	}
}

func TestUnirseSinLugar(t *testing.T) {
	b := nuevoBanco(t, ".", 16)
	b.unirOk(t, "ana")

	_, err := b.unir(t, "beto")

	if !errors.Is(err, ErrSinLugar) {
		t.Errorf("err = %v, quiere ErrSinLugar", err)
	}
	if len(b.estado.Jugadores) != 1 || len(b.c.suscripciones) != 1 {
		t.Errorf("hay %d jugadores y %d suscripciones; quiere 1 y 1",
			len(b.estado.Jugadores), len(b.c.suscripciones))
	}
}

// Si quien se une se cancela antes de recibir la respuesta, el ciclo retira
// al jugador que acaba de agregar.
func TestUnirseCanceladoRetiraAlJugador(t *testing.T) {
	b := nuevoBanco(t, mapaLibre, 16)
	ctx, cancelar := context.WithCancel(context.Background())
	res := b.iniciarUnirse(t, ctx, "ana")

	cancelar()
	r := recibir(t, res, "que Unirse devuelva el error de ctx")
	if !errors.Is(r.err, context.Canceled) {
		t.Fatalf("err = %v, quiere context.Canceled", r.err)
	}

	b.c.tick() // procesa el mensaje: agrega a ana, no puede responder y la retira

	if len(b.estado.Jugadores) != 0 || len(b.c.suscripciones) != 0 {
		t.Errorf("quedaron %d jugadores y %d suscripciones; quiere 0 y 0",
			len(b.estado.Jugadores), len(b.c.suscripciones))
	}
}

func TestActuarMueveAlJugador(t *testing.T) {
	b := nuevoBanco(t, mapaLibre, 16)
	s := b.unirOk(t, "ana")

	b.actuar(t, s.ID, juego.MoverAccion{Dir: juego.Derecha})
	b.c.tick()

	if got := b.pos(s.ID); got != (juego.Posicion{X: 1, Y: 0}) {
		t.Errorf("pos = %v, quiere (1,0)", got)
	}
}

func TestLaUltimaAccionGana(t *testing.T) {
	b := nuevoBanco(t, mapaLibre, 16)
	s := b.unirOk(t, "ana")

	b.actuar(t, s.ID, juego.MoverAccion{Dir: juego.Derecha})
	b.actuar(t, s.ID, juego.MoverAccion{Dir: juego.Abajo})
	b.c.tick()

	if got := b.pos(s.ID); got != (juego.Posicion{X: 0, Y: 1}) {
		t.Errorf("pos = %v, quiere (0,1): solo debía contar la última acción", got)
	}
}

func TestUnaAccionPorTick(t *testing.T) {
	b := nuevoBanco(t, mapaLibre, 16)
	s := b.unirOk(t, "ana")

	b.actuar(t, s.ID, juego.MoverAccion{Dir: juego.Derecha})
	b.actuar(t, s.ID, juego.MoverAccion{Dir: juego.Derecha})
	b.c.tick()

	if got := b.pos(s.ID); got != (juego.Posicion{X: 1, Y: 0}) {
		t.Errorf("pos = %v, quiere (1,0): un solo paso por tick", got)
	}
}

func TestDosJugadoresPorLaMismaCasillaGanaElPrimeroDelOrden(t *testing.T) {
	b := nuevoBanco(t, mapaLibre, 16)
	b.unirOk(t, "ana")
	b.unirOk(t, "beto")
	inicio := map[juego.IDJugador]juego.Posicion{1: {X: 0, Y: 1}, 2: {X: 2, Y: 1}}
	for id, p := range inicio {
		b.estado.Jugadores[id].Pos = p
	}
	b.actuar(t, 1, juego.MoverAccion{Dir: juego.Derecha})
	b.actuar(t, 2, juego.MoverAccion{Dir: juego.Izquierda})
	orden := ordenDeAplicacion([]juego.IDJugador{1, 2}, 0, b.estado.Tick)

	b.c.tick()

	ganador, perdedor := orden[0], orden[1]
	if got := b.pos(ganador); got != (juego.Posicion{X: 1, Y: 1}) {
		t.Errorf("el ganador (%d) está en %v, quiere (1,1)", ganador, got)
	}
	if got := b.pos(perdedor); got != inicio[perdedor] {
		t.Errorf("el perdedor (%d) está en %v, quiere %v", perdedor, got, inicio[perdedor])
	}
}

func TestOrdenDeAplicacion(t *testing.T) {
	t.Run("es una permutación y no modifica la entrada", func(t *testing.T) {
		ids := []juego.IDJugador{4, 2, 5, 1, 3}
		copia := slices.Clone(ids)
		orden := ordenDeAplicacion(ids, 7, 11)
		if !slices.Equal(ids, copia) {
			t.Errorf("la entrada cambió: %v", ids)
		}
		ordenado := slices.Clone(orden)
		slices.Sort(ordenado)
		if !slices.Equal(ordenado, []juego.IDJugador{1, 2, 3, 4, 5}) {
			t.Errorf("orden = %v, no es una permutación de 1..5", orden)
		}
	})

	t.Run("no depende del orden de entrada", func(t *testing.T) {
		a := ordenDeAplicacion([]juego.IDJugador{3, 1, 2, 5, 4}, 9, 42)
		b := ordenDeAplicacion([]juego.IDJugador{1, 2, 3, 4, 5}, 9, 42)
		if !slices.Equal(a, b) {
			t.Errorf("%v != %v con los mismos ids en otro orden", a, b)
		}
	})

	t.Run("es repetible", func(t *testing.T) {
		ids := []juego.IDJugador{1, 2, 3, 4, 5, 6}
		if a, b := ordenDeAplicacion(ids, 3, 8), ordenDeAplicacion(ids, 3, 8); !slices.Equal(a, b) {
			t.Errorf("%v != %v con los mismos datos", a, b)
		}
	})

	t.Run("no favorece siempre al mismo", func(t *testing.T) {
		const ticks = 64
		primeroUno := 0
		for tick := uint64(0); tick < ticks; tick++ {
			if ordenDeAplicacion([]juego.IDJugador{1, 2}, 0, tick)[0] == 1 {
				primeroUno++
			}
		}
		if primeroUno < ticks/4 || primeroUno > ticks*3/4 {
			t.Errorf("el ID 1 fue primero %d de %d veces; quiere un reparto parejo", primeroUno, ticks)
		}
	})
}

func TestSalir(t *testing.T) {
	b := nuevoBanco(t, mapaLibre, 16)
	s1 := b.unirOk(t, "ana")
	b.unirOk(t, "beto")
	posBeto := b.pos(2)

	b.actuar(t, s1.ID, juego.MoverAccion{Dir: juego.Derecha}) // se descarta al salir
	if err := Salir(context.Background(), b.entrada, s1.ID); err != nil {
		t.Fatalf("Salir: %v", err)
	}
	b.c.tick()

	if _, existe := b.estado.Jugadores[s1.ID]; existe {
		t.Error("ana sigue en el estado")
	}
	if len(b.estado.Jugadores) != 1 || b.pos(2) != posBeto {
		t.Errorf("beto debía quedar intacto: %d jugadores, pos %v", len(b.estado.Jugadores), b.pos(2))
	}
	if len(b.c.suscripciones) != 1 {
		t.Errorf("hay %d suscripciones, quiere 1", len(b.c.suscripciones))
	}
	esperarCierre(t, s1.Cuadros, "que el canal de ana se cierre")
}

// Si Salir no intentara el envío antes de mirar ctx, el select de enviar
// elegiría al azar entre los dos casos listos y fallaría solo a veces. Por eso
// el escenario completo se repite muchas veces: con 200 repeticiones una
// implementación con esa elección al azar pasaría todas con probabilidad 2⁻²⁰⁰.
func TestSalirConContextoYaCancelado(t *testing.T) {
	const repeticiones = 200
	for i := 0; i < repeticiones; i++ {
		b := nuevoBanco(t, mapaLibre, 16)
		s := b.unirOk(t, "ana")
		ctx, cancelar := context.WithCancel(context.Background())
		cancelar()

		if err := Salir(ctx, b.entrada, s.ID); err != nil {
			t.Fatalf("repetición %d: Salir con cola con espacio y ctx cancelado: %v, quiere nil", i, err)
		}
		b.c.tick()

		if len(b.estado.Jugadores) != 0 || len(b.c.suscripciones) != 0 {
			t.Fatalf("repetición %d: ana no se retiró: %d jugadores, %d suscripciones",
				i, len(b.estado.Jugadores), len(b.c.suscripciones))
		}
	}
}

func TestSalirConColaLlenaYContextoCancelado(t *testing.T) {
	b := nuevoBanco(t, mapaLibre, 1)
	s := b.unirOk(t, "ana")
	b.entrada <- Mensaje{} // la cola, de capacidad 1, queda llena
	ctx, cancelar := context.WithCancel(context.Background())
	cancelar()

	if err := Salir(ctx, b.entrada, s.ID); !errors.Is(err, context.Canceled) {
		t.Errorf("err = %v, quiere context.Canceled", err)
	}
}

func TestMensajesQueSeIgnoran(t *testing.T) {
	b := nuevoBanco(t, mapaLibre, 16)
	s := b.unirOk(t, "ana")
	antes := b.pos(s.ID)

	b.entrada <- Mensaje{}                                                              // valor cero
	b.entrada <- Mensaje{tipo: mensajeAccion, jugador: s.ID}                            // acción nula
	b.entrada <- Mensaje{tipo: mensajeAccion, jugador: 99, accion: juego.MoverAccion{}} // sin suscripción
	b.entrada <- Mensaje{tipo: mensajeUnirse, nombre: "sin canal de respuesta"}         // mal formado
	b.entrada <- Mensaje{tipo: mensajeSalir, jugador: 99}                               // desconocido
	b.entrada <- Mensaje{tipo: tipoMensaje(99), jugador: s.ID}                          // tipo inventado
	b.c.tick()

	if len(b.estado.Jugadores) != 1 || b.pos(s.ID) != antes || len(b.c.suscripciones) != 1 {
		t.Errorf("el estado cambió: %d jugadores, pos %v, %d suscripciones",
			len(b.estado.Jugadores), b.pos(s.ID), len(b.c.suscripciones))
	}
}

func TestSuscriptorLentoSoloTieneElCuadroMasNuevo(t *testing.T) {
	b := nuevoBanco(t, mapaLibre, 16)
	lento := b.unirOk(t, "lento") // tick 1: deja un cuadro sin leer
	b.c.tick()                    // tick 2
	b.c.tick()                    // tick 3: ninguno de los dos se bloquea

	cuadro := recibir(t, lento.Cuadros, "el cuadro más nuevo")
	if cuadro.Vista.Tick != 3 {
		t.Errorf("Vista.Tick = %d, quiere 3", cuadro.Vista.Tick)
	}
	if n := len(lento.Cuadros); n != 0 {
		t.Errorf("quedaron %d cuadros viejos", n)
	}
}

func TestSuscriptorAlDiaRecibeCadaCuadro(t *testing.T) {
	b := nuevoBanco(t, mapaLibre, 16)
	s := b.unirOk(t, "ana")
	for quiere := uint64(1); quiere <= 3; quiere++ {
		if quiere > 1 {
			b.c.tick()
		}
		cuadro := recibir(t, s.Cuadros, fmt.Sprintf("el cuadro del tick %d", quiere))
		if cuadro.Vista.Tick != quiere {
			t.Errorf("Vista.Tick = %d, quiere %d", cuadro.Vista.Tick, quiere)
		}
	}
}

// botFijo siempre pide la misma acción y guarda la última vista que recibió.
type botFijo struct {
	accion juego.Accion
	vista  juego.Vista
}

func (b *botFijo) Decidir(v juego.Vista) juego.Accion {
	b.vista = v
	return b.accion
}

var _ bots.Bot = (*botFijo)(nil)

func TestBotActuaEnElTick(t *testing.T) {
	b := nuevoBanco(t, mapaLibre, 16)
	bot := &botFijo{accion: juego.MoverAccion{Dir: juego.Derecha}}
	id, err := b.c.RegistrarBot(bot)
	if err != nil {
		t.Fatalf("RegistrarBot: %v", err)
	}

	b.c.tick()

	if got := b.pos(id); got != (juego.Posicion{X: 1, Y: 0}) {
		t.Errorf("pos del bot = %v, quiere (1,0)", got)
	}
	if bot.vista.Yo != id {
		t.Errorf("el bot recibió la vista de %d, quiere la de %d", bot.vista.Yo, id)
	}
}

func TestRegistrarBotSinLugar(t *testing.T) {
	b := nuevoBanco(t, ".", 16)
	b.unirOk(t, "ana")
	if _, err := b.c.RegistrarBot(&botFijo{}); !errors.Is(err, ErrSinLugar) {
		t.Errorf("err = %v, quiere ErrSinLugar", err)
	}
}

// Un mensaje de acción con el ID de un bot no cuenta: los bots no tienen
// suscripción y solo actúan con lo que deciden ellos.
func TestMensajeConElIDDeUnBotSeIgnora(t *testing.T) {
	b := nuevoBanco(t, mapaLibre, 16)
	id, err := b.c.RegistrarBot(&botFijo{}) // no decide nada
	if err != nil {
		t.Fatalf("RegistrarBot: %v", err)
	}
	antes := b.pos(id)

	b.entrada <- Mensaje{tipo: mensajeAccion, jugador: id, accion: juego.MoverAccion{Dir: juego.Abajo}}
	b.c.tick()

	if got := b.pos(id); got != antes {
		t.Errorf("el bot se movió de %v a %v por un mensaje ajeno", antes, got)
	}
}

// accionEspia cuenta cuántas veces se aplica.
type accionEspia struct{ llamadas *int }

func (a accionEspia) Aplicar(e *juego.Estado, autor juego.IDJugador) { *a.llamadas++ }

func TestSoloSeAplicanLasAccionesQueLaFasePermite(t *testing.T) {
	b := nuevoBanco(t, mapaLibre, 16)
	s := b.unirOk(t, "ana")
	llamadas := 0

	b.actuar(t, s.ID, accionEspia{&llamadas}) // la sala de espera solo permite mover
	b.c.tick()

	if llamadas != 0 {
		t.Errorf("se aplicó %d veces una acción que la fase no permite", llamadas)
	}
}

func TestEjecutarTerminaAlCancelarYCierraSuscripciones(t *testing.T) {
	b := nuevoBanco(t, mapaLibre, 16)
	cuadros := make(chan Cuadro, 1)
	b.c.suscripciones[1] = cuadros
	ctx, cancelar := context.WithCancel(context.Background())
	terminado := make(chan struct{})

	go func() {
		b.c.Ejecutar(ctx)
		close(terminado)
	}()
	cancelar()

	recibir(t, terminado, "que Ejecutar termine tras cancelar ctx")
	if _, abierto := <-cuadros; abierto {
		t.Error("el canal de la suscripción sigue abierto")
	}
}

// Varias goroutines se unen, actúan y salen mientras el hilo de la prueba
// ejecuta ticks. No se comparan valores exactos, solo invariantes.
func TestVariasGoroutinesEnviandoMensajes(t *testing.T) {
	const (
		jugadores = 8
		acciones  = 50
	)
	fila := strings.Repeat(".", 10) + "\n"
	b := nuevoBanco(t, strings.Repeat(fila, 10), 16)

	type destino struct {
		id    juego.IDJugador
		salio bool
	}
	destinos := make(chan destino, jugadores)
	fallos := make(chan error, jugadores)
	var grupo sync.WaitGroup
	for i := 0; i < jugadores; i++ {
		grupo.Add(1)
		go func() {
			defer grupo.Done()
			ctx := context.Background()
			s, err := Unirse(ctx, b.entrada, fmt.Sprintf("j%d", i))
			if err != nil {
				fallos <- err
				return
			}
			for k := 0; k < acciones; k++ {
				dir := juego.Direccion((i + k) % 4)
				if err := Actuar(ctx, b.entrada, s.ID, juego.MoverAccion{Dir: dir}); err != nil {
					fallos <- err
					return
				}
				select {
				case <-s.Cuadros:
				default:
				}
			}
			salio := i%2 == 0
			if salio {
				if err := Salir(ctx, b.entrada, s.ID); err != nil {
					fallos <- err
					return
				}
			}
			destinos <- destino{s.ID, salio}
		}()
	}
	terminaron := make(chan struct{})
	go func() {
		grupo.Wait()
		close(terminaron)
	}()

	hecho := false
	for i := 0; i < maxIntentos && !hecho; i++ {
		b.c.tick()
		select {
		case <-terminaron:
			hecho = true
		default:
			runtime.Gosched()
		}
	}
	if !hecho {
		t.Fatalf("tras %d ticks las goroutines no terminaron", maxIntentos)
	}
	b.c.tick() // atiende lo que haya quedado en la cola
	close(fallos)
	for err := range fallos {
		t.Errorf("una goroutine falló: %v", err)
	}
	close(destinos)

	quedan := map[juego.IDJugador]bool{}
	vistos := map[juego.IDJugador]bool{}
	for d := range destinos {
		if vistos[d.id] {
			t.Errorf("el ID %d se asignó dos veces", d.id)
		}
		vistos[d.id] = true
		if !d.salio {
			quedan[d.id] = true
		}
	}
	if len(vistos) != jugadores {
		t.Errorf("se asignaron %d IDs distintos, quiere %d", len(vistos), jugadores)
	}
	if len(b.estado.Jugadores) != len(quedan) || len(b.c.suscripciones) != len(quedan) {
		t.Errorf("hay %d jugadores y %d suscripciones; quiere %d",
			len(b.estado.Jugadores), len(b.c.suscripciones), len(quedan))
	}
	ocupadas := map[juego.Posicion]juego.IDJugador{}
	for id, j := range b.estado.Jugadores {
		if !quedan[id] {
			t.Errorf("el jugador %d debía haber salido", id)
		}
		if !b.estado.Mapa.EsTransitable(j.Pos) {
			t.Errorf("el jugador %d está en %v, que no es transitable", id, j.Pos)
		}
		if otro, hay := ocupadas[j.Pos]; hay {
			t.Errorf("los jugadores %d y %d comparten la casilla %v", otro, id, j.Pos)
		}
		ocupadas[j.Pos] = id
	}
}
