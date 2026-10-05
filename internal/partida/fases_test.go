package partida

import (
	"testing"

	"github.com/Pacocerv08/Terminal_Game/internal/juego"
)

func TestNuevaMaquinaEmpiezaEnSalaDeEspera(t *testing.T) {
	if got := NuevaMaquina().Fase(); got != SalaDeEspera {
		t.Errorf("Fase() = %v, quiere SalaDeEspera", got)
	}
}

func TestPermiteUnirse(t *testing.T) {
	casos := []struct {
		fase   Fase
		quiere bool
	}{
		{SalaDeEspera, true},
		{EnCurso, false},
		{Final, false},
	}
	for _, c := range casos {
		m := &Maquina{fase: c.fase}
		if got := m.PermiteUnirse(); got != c.quiere {
			t.Errorf("fase %v: PermiteUnirse() = %v, quiere %v", c.fase, got, c.quiere)
		}
	}
}

func TestAccionPermitida(t *testing.T) {
	casos := []struct {
		nombre string
		fase   Fase
		accion juego.Accion
		quiere bool
	}{
		{"mover en sala de espera", SalaDeEspera, juego.MoverAccion{Dir: juego.Arriba}, true},
		{"disparar en sala de espera", SalaDeEspera, juego.DispararAccion{Dir: juego.Arriba}, false},
		{"recoger en sala de espera", SalaDeEspera, juego.RecogerAccion{}, false},
		{"acción nula en sala de espera", SalaDeEspera, nil, false},
		{"mover en curso", EnCurso, juego.MoverAccion{Dir: juego.Arriba}, false},
		{"mover en final", Final, juego.MoverAccion{Dir: juego.Arriba}, false},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			m := &Maquina{fase: c.fase}
			if got := m.AccionPermitida(c.accion); got != c.quiere {
				t.Errorf("AccionPermitida = %v, quiere %v", got, c.quiere)
			}
		})
	}
}
