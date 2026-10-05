package juego

import (
	"reflect"
	"testing"
)

// estadoParaVista crea un estado con dos jugadores (el 1 con inventario), un
// objeto en el suelo y el mapa de una sola casilla.
func estadoParaVista(t *testing.T) *Estado {
	t.Helper()
	e := estadoConJugadores(t, ".",
		Jugador{ID: 2, Nombre: "dos", Pos: Posicion{0, 0}, Vida: 50,
			Inventario: Inventario{Capacidad: 4, objetos: []Objeto{{ObjetoBotiquin, 1}}}},
		Jugador{ID: 1, Nombre: "uno", Pos: Posicion{0, 0}, Vida: VidaInicial,
			Inventario: Inventario{Capacidad: 4, objetos: []Objeto{{ObjetoMunicion, 5}}}},
	)
	e.Botin = []ObjetoEnSuelo{{Objeto: Objeto{ObjetoPistola, 1}, Pos: Posicion{0, 0}}}
	return e
}

func TestVistaParaContenido(t *testing.T) {
	e := estadoParaVista(t)
	e.Tick = 7
	v := e.VistaPara(1)

	if v.Yo != 1 || v.Tick != 7 {
		t.Errorf("Yo = %d, Tick = %d; quiere 1 y 7", v.Yo, v.Tick)
	}
	quieren := []JugadorVisible{
		{ID: 1, Nombre: "uno", Pos: Posicion{0, 0}, Vida: VidaInicial},
		{ID: 2, Nombre: "dos", Pos: Posicion{0, 0}, Vida: 50},
	}
	if !reflect.DeepEqual(v.Jugadores, quieren) {
		t.Errorf("Jugadores = %+v, quiere %+v (ordenados por ID)", v.Jugadores, quieren)
	}
	if !reflect.DeepEqual(v.Botin, e.Botin) {
		t.Errorf("Botin = %+v, quiere %+v", v.Botin, e.Botin)
	}
	if !reflect.DeepEqual(v.Mi, e.Jugadores[1].Inventario) {
		t.Errorf("Mi = %+v, quiere %+v", v.Mi, e.Jugadores[1].Inventario)
	}
	if v.Mapa != e.Mapa {
		t.Error("el mapa debe compartirse por puntero porque es inmutable")
	}
}

func TestVistaParaJugadorInexistente(t *testing.T) {
	e := estadoParaVista(t)
	v := e.VistaPara(99)
	if v.Mi.Capacidad != 0 || len(v.Mi.objetos) != 0 {
		t.Errorf("Mi = %+v, quiere un inventario vacío", v.Mi)
	}
}

// La vista no comparte memoria con el estado, en ninguno de los dos sentidos.
// El mapa queda fuera porque se comparte a propósito (es inmutable).
func TestVistaParaNoComparteMemoria(t *testing.T) {
	t.Run("cambiar la vista no cambia el estado", func(t *testing.T) {
		e := estadoParaVista(t)
		antes := estadoParaVista(t)
		v := e.VistaPara(1)

		v.Jugadores[0].Vida = 999
		v.Jugadores[0].Pos = Posicion{9, 9}
		v.Botin[0].Cantidad = 999
		v.Botin[0].Pos = Posicion{9, 9}
		v.Mi.objetos[0].Cantidad = 999

		if !reflect.DeepEqual(e.Jugadores, antes.Jugadores) {
			t.Error("cambiar la vista cambió Estado.Jugadores")
		}
		if !reflect.DeepEqual(e.Botin, antes.Botin) {
			t.Error("cambiar la vista cambió Estado.Botin")
		}
	})

	t.Run("cambiar el estado no cambia la vista", func(t *testing.T) {
		e := estadoParaVista(t)
		v := e.VistaPara(1)

		e.Jugadores[1].Vida = 1
		e.Jugadores[1].Pos = Posicion{9, 9}
		e.Jugadores[1].Inventario.objetos[0].Cantidad = 777
		e.Jugadores[2].Nombre = "otro"
		e.Botin[0].Cantidad = 777

		// Valores esperados escritos a mano, independientes del estado: son
		// los que tenía estadoParaVista cuando se pidió la vista.
		quierenJugadores := []JugadorVisible{
			{ID: 1, Nombre: "uno", Pos: Posicion{0, 0}, Vida: VidaInicial},
			{ID: 2, Nombre: "dos", Pos: Posicion{0, 0}, Vida: 50},
		}
		quiereBotin := []ObjetoEnSuelo{
			{Objeto: Objeto{ObjetoPistola, 1}, Pos: Posicion{0, 0}},
		}
		quiereMi := Inventario{Capacidad: 4, objetos: []Objeto{{ObjetoMunicion, 5}}}

		if !reflect.DeepEqual(v.Jugadores, quierenJugadores) {
			t.Errorf("Jugadores = %+v, quiere %+v", v.Jugadores, quierenJugadores)
		}
		if !reflect.DeepEqual(v.Botin, quiereBotin) {
			t.Errorf("Botin = %+v, quiere %+v", v.Botin, quiereBotin)
		}
		if !reflect.DeepEqual(v.Mi, quiereMi) {
			t.Errorf("Mi = %+v, quiere %+v", v.Mi, quiereMi)
		}
	})
}
