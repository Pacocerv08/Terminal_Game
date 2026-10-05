package juego

import "testing"

func TestEstaVivo(t *testing.T) {
	casos := []struct {
		nombre string
		vida   int
		quiere bool
	}{
		{"con vida", 10, true},
		{"con una vida", 1, true},
		{"sin vida", 0, false},
		{"vida negativa", -5, false},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			j := &Jugador{Vida: c.vida}
			if got := j.EstaVivo(); got != c.quiere {
				t.Errorf("EstaVivo() con vida %d = %v, quiere %v", c.vida, got, c.quiere)
			}
		})
	}
}
