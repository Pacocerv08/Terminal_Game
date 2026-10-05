package entrada

import (
	"testing"

	"github.com/Pacocerv08/Terminal_Game/internal/juego"
)

func TestTraducir(t *testing.T) {
	casos := []struct {
		tecla string
		want  juego.Accion
		ok    bool
	}{
		{"w", juego.MoverAccion{Dir: juego.Arriba}, true},
		{"W", juego.MoverAccion{Dir: juego.Arriba}, true},
		{"up", juego.MoverAccion{Dir: juego.Arriba}, true},
		{"s", juego.MoverAccion{Dir: juego.Abajo}, true},
		{"S", juego.MoverAccion{Dir: juego.Abajo}, true},
		{"down", juego.MoverAccion{Dir: juego.Abajo}, true},
		{"a", juego.MoverAccion{Dir: juego.Izquierda}, true},
		{"A", juego.MoverAccion{Dir: juego.Izquierda}, true},
		{"left", juego.MoverAccion{Dir: juego.Izquierda}, true},
		{"d", juego.MoverAccion{Dir: juego.Derecha}, true},
		{"D", juego.MoverAccion{Dir: juego.Derecha}, true},
		{"right", juego.MoverAccion{Dir: juego.Derecha}, true},
		{"", nil, false},
		{" ", nil, false},
		{"q", nil, false},
		{"x", nil, false},
		{"ww", nil, false},
		{"up ", nil, false},
		{"UP", nil, false},
		{"Up", nil, false},
		{"enter", nil, false},
		{"ctrl+c", nil, false},
		{"ñ", nil, false},
	}
	var tr TraductorTeclado
	for _, c := range casos {
		t.Run(c.tecla, func(t *testing.T) {
			got, ok := tr.Traducir(c.tecla)
			if ok != c.ok || got != c.want {
				t.Errorf("Traducir(%q) = (%v, %v), se esperaba (%v, %v)", c.tecla, got, ok, c.want, c.ok)
			}
		})
	}
}
