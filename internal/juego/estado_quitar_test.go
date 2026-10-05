package juego

import "testing"

func TestQuitarJugador(t *testing.T) {
	e := NuevoEstado(mapaDeTexto(t, "..."))
	a := e.AgregarJugador("a")
	b := e.AgregarJugador("b")
	posA := a.Pos

	e.QuitarJugador(a.ID)

	if _, existe := e.Jugadores[a.ID]; existe {
		t.Error("el jugador quitado sigue en Estado.Jugadores")
	}
	if e.Jugadores[b.ID] != b || len(e.Jugadores) != 1 {
		t.Errorf("el otro jugador debía quedar intacto; hay %d jugadores", len(e.Jugadores))
	}

	// Su casilla queda libre y su ID no se reutiliza.
	c := e.AgregarJugador("c")
	if c == nil {
		t.Fatal("AgregarJugador devolvió nil con una casilla libre")
	}
	if c.Pos != posA {
		t.Errorf("c en %v, quiere la casilla liberada %v", c.Pos, posA)
	}
	if c.ID != 3 {
		t.Errorf("c tiene ID %d, quiere 3 (los IDs no se reutilizan)", c.ID)
	}
}

func TestQuitarJugadorInexistente(t *testing.T) {
	e := NuevoEstado(mapaDeTexto(t, "..."))
	a := e.AgregarJugador("a")

	e.QuitarJugador(99)

	if len(e.Jugadores) != 1 || e.Jugadores[a.ID] != a {
		t.Errorf("quitar un ID inexistente cambió el estado: %d jugadores", len(e.Jugadores))
	}
}

func TestAvanzarCuentaTicks(t *testing.T) {
	e := NuevoEstado(mapaDeTexto(t, "."))
	for quiere := uint64(1); quiere <= 3; quiere++ {
		e.Avanzar()
		if e.Tick != quiere {
			t.Errorf("tras %d avances Tick = %d", quiere, e.Tick)
		}
	}
}
