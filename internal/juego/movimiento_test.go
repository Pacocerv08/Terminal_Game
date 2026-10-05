package juego

import "testing"

// estadoConJugadores crea un estado sobre el mapa de texto y le pone los
// jugadores dados, sin pasar por AgregarJugador.
func estadoConJugadores(t *testing.T, texto string, jugadores ...Jugador) *Estado {
	t.Helper()
	e := NuevoEstado(mapaDeTexto(t, texto))
	for i := range jugadores {
		j := jugadores[i]
		e.Jugadores[j.ID] = &j
	}
	return e
}

func TestJugadorMover(t *testing.T) {
	const mapa = `
		.....
		.#...
		.....
	`
	casos := []struct {
		nombre      string
		inicio      Posicion
		dir         Direccion
		quiereMover bool
		quierePos   Posicion
	}{
		{"arriba", Posicion{2, 1}, Arriba, true, Posicion{2, 0}},
		{"abajo", Posicion{2, 1}, Abajo, true, Posicion{2, 2}},
		{"izquierda", Posicion{2, 0}, Izquierda, true, Posicion{1, 0}},
		{"derecha", Posicion{2, 1}, Derecha, true, Posicion{3, 1}},
		{"pared", Posicion{0, 1}, Derecha, false, Posicion{0, 1}},
		{"borde arriba", Posicion{0, 0}, Arriba, false, Posicion{0, 0}},
		{"borde abajo", Posicion{4, 2}, Abajo, false, Posicion{4, 2}},
		{"borde izquierda", Posicion{0, 0}, Izquierda, false, Posicion{0, 0}},
		{"borde derecha", Posicion{4, 2}, Derecha, false, Posicion{4, 2}},
		{"dirección inválida", Posicion{2, 1}, Direccion(99), false, Posicion{2, 1}},
		{"dirección negativa", Posicion{2, 1}, Direccion(-1), false, Posicion{2, 1}},
	}
	m := mapaDeTexto(t, mapa)
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			j := &Jugador{ID: 1, Pos: c.inicio, Vida: VidaInicial}
			if got := j.Mover(c.dir, m); got != c.quiereMover {
				t.Errorf("Mover = %v, quiere %v", got, c.quiereMover)
			}
			if j.Pos != c.quierePos {
				t.Errorf("Pos = %v, quiere %v", j.Pos, c.quierePos)
			}
		})
	}
}

func TestMoverAccionAplicar(t *testing.T) {
	// Mapa con una pared en (1,1) y sin borde, para probar pared y salida del
	// mapa por separado.
	const mapa = `
		.....
		.#...
		.....
	`
	vivo := func(id IDJugador, x, y int) Jugador {
		return Jugador{ID: id, Nombre: "j", Pos: Posicion{x, y}, Vida: VidaInicial}
	}
	muerto := func(id IDJugador, x, y int) Jugador {
		return Jugador{ID: id, Nombre: "j", Pos: Posicion{x, y}, Vida: 0}
	}

	casos := []struct {
		nombre     string
		jugadores  []Jugador
		autor      IDJugador
		dir        Direccion
		quierenPos map[IDJugador]Posicion
	}{
		{
			nombre:     "casilla libre",
			jugadores:  []Jugador{vivo(1, 2, 0)},
			autor:      1,
			dir:        Abajo,
			quierenPos: map[IDJugador]Posicion{1: {2, 1}},
		},
		{
			nombre:     "pared",
			jugadores:  []Jugador{vivo(1, 0, 1)},
			autor:      1,
			dir:        Derecha,
			quierenPos: map[IDJugador]Posicion{1: {0, 1}},
		},
		{
			nombre:     "borde del mapa",
			jugadores:  []Jugador{vivo(1, 0, 0)},
			autor:      1,
			dir:        Arriba,
			quierenPos: map[IDJugador]Posicion{1: {0, 0}},
		},
		{
			nombre:     "casilla ocupada por jugador vivo",
			jugadores:  []Jugador{vivo(1, 2, 0), vivo(2, 3, 0)},
			autor:      1,
			dir:        Derecha,
			quierenPos: map[IDJugador]Posicion{1: {2, 0}, 2: {3, 0}},
		},
		{
			nombre:     "casilla de jugador muerto se puede ocupar",
			jugadores:  []Jugador{vivo(1, 2, 0), muerto(2, 3, 0)},
			autor:      1,
			dir:        Derecha,
			quierenPos: map[IDJugador]Posicion{1: {3, 0}, 2: {3, 0}},
		},
		{
			nombre:     "autor inexistente",
			jugadores:  []Jugador{vivo(1, 2, 0)},
			autor:      9,
			dir:        Abajo,
			quierenPos: map[IDJugador]Posicion{1: {2, 0}},
		},
		{
			nombre:     "autor muerto",
			jugadores:  []Jugador{muerto(1, 2, 0)},
			autor:      1,
			dir:        Abajo,
			quierenPos: map[IDJugador]Posicion{1: {2, 0}},
		},
		{
			nombre:     "dirección inválida",
			jugadores:  []Jugador{vivo(1, 2, 0)},
			autor:      1,
			dir:        Direccion(99),
			quierenPos: map[IDJugador]Posicion{1: {2, 0}},
		},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			e := estadoConJugadores(t, mapa, c.jugadores...)
			MoverAccion{Dir: c.dir}.Aplicar(e, c.autor)
			if len(e.Jugadores) != len(c.quierenPos) {
				t.Errorf("hay %d jugadores, quiere %d", len(e.Jugadores), len(c.quierenPos))
			}
			for id, quiere := range c.quierenPos {
				if got := e.Jugadores[id].Pos; got != quiere {
					t.Errorf("jugador %d en %v, quiere %v", id, got, quiere)
				}
			}
		})
	}
}

// Dos jugadores quieren la misma casilla: gana quien se aplica primero.
func TestMoverAccionEmpate(t *testing.T) {
	const mapa = "..."
	crear := func() *Estado {
		return estadoConJugadores(t, mapa,
			Jugador{ID: 1, Pos: Posicion{0, 0}, Vida: VidaInicial},
			Jugador{ID: 2, Pos: Posicion{2, 0}, Vida: VidaInicial},
		)
	}

	e := crear()
	MoverAccion{Dir: Derecha}.Aplicar(e, 1)
	MoverAccion{Dir: Izquierda}.Aplicar(e, 2)
	if e.Jugadores[1].Pos != (Posicion{1, 0}) || e.Jugadores[2].Pos != (Posicion{2, 0}) {
		t.Errorf("aplicando 1 primero: pos 1 = %v, pos 2 = %v; quiere (1,0) y (2,0)",
			e.Jugadores[1].Pos, e.Jugadores[2].Pos)
	}

	e = crear()
	MoverAccion{Dir: Izquierda}.Aplicar(e, 2)
	MoverAccion{Dir: Derecha}.Aplicar(e, 1)
	if e.Jugadores[2].Pos != (Posicion{1, 0}) || e.Jugadores[1].Pos != (Posicion{0, 0}) {
		t.Errorf("aplicando 2 primero: pos 1 = %v, pos 2 = %v; quiere (0,0) y (1,0)",
			e.Jugadores[1].Pos, e.Jugadores[2].Pos)
	}
}

func TestAgregarJugador(t *testing.T) {
	e := NuevoEstado(mapaDeTexto(t, `
		#..
		.#.
	`))
	casos := []struct {
		nombre string
		id     IDJugador
		pos    Posicion
	}{
		{"primero", 1, Posicion{1, 0}},
		{"segundo", 2, Posicion{2, 0}},
		{"tercero", 3, Posicion{0, 1}},
		{"cuarto", 4, Posicion{2, 1}},
	}
	for _, c := range casos {
		j := e.AgregarJugador(c.nombre)
		if j == nil {
			t.Fatalf("%s: AgregarJugador devolvió nil", c.nombre)
		}
		if j.ID != c.id || j.Pos != c.pos {
			t.Errorf("%s: ID %d en %v, quiere ID %d en %v", c.nombre, j.ID, j.Pos, c.id, c.pos)
		}
		if j.Nombre != c.nombre || j.Vida != VidaInicial || !j.EstaVivo() {
			t.Errorf("%s: nombre %q, vida %d", c.nombre, j.Nombre, j.Vida)
		}
		if e.Jugadores[j.ID] != j {
			t.Errorf("%s: no quedó registrado en Estado.Jugadores", c.nombre)
		}
	}

	// Mapa lleno: devuelve nil y no consume ningún ID.
	if j := e.AgregarJugador("sobra"); j != nil {
		t.Errorf("con el mapa lleno devolvió %+v, quiere nil", j)
	}
	if len(e.Jugadores) != 4 || e.siguienteID != 4 {
		t.Errorf("tras el fallo hay %d jugadores y siguienteID %d; quiere 4 y 4", len(e.Jugadores), e.siguienteID)
	}
}

func TestAgregarJugadorReutilizaCasillaDeMuerto(t *testing.T) {
	e := NuevoEstado(mapaDeTexto(t, "#."))
	primero := e.AgregarJugador("a")
	primero.Vida = 0
	segundo := e.AgregarJugador("b")
	if segundo == nil || segundo.Pos != primero.Pos {
		t.Errorf("segundo = %+v, quiere un jugador en %v", segundo, primero.Pos)
	}
}

func TestAgregarJugadorSinMapaTransitable(t *testing.T) {
	casos := []struct {
		nombre string
		estado *Estado
	}{
		{"mapa de pared", NuevoEstado(mapaDeTexto(t, "##\n##"))},
		{"sin mapa", NuevoEstado(nil)},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if j := c.estado.AgregarJugador("x"); j != nil {
				t.Errorf("devolvió %+v, quiere nil", j)
			}
		})
	}
}
