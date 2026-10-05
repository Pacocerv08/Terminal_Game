package juego

import "testing"

// mapaDeTexto construye un mapa para una prueba y falla si el texto no vale.
func mapaDeTexto(t *testing.T, texto string) *Mapa {
	t.Helper()
	m, err := MapaDesdeTexto(texto)
	if err != nil {
		t.Fatalf("MapaDesdeTexto: %v", err)
	}
	return m
}

func TestNuevoMapa(t *testing.T) {
	m := NuevoMapa(4, 3)
	if m.Ancho != 4 || m.Alto != 3 {
		t.Fatalf("tamaño = %dx%d, quiere 4x3", m.Ancho, m.Alto)
	}
	for y := 0; y < 3; y++ {
		for x := 0; x < 4; x++ {
			quiere := Suelo
			if x == 0 || y == 0 || x == 3 || y == 2 {
				quiere = Pared
			}
			if got, dentro := m.CasillaEn(Posicion{x, y}); !dentro || got != quiere {
				t.Errorf("CasillaEn(%d,%d) = %v, %v; quiere %v, true", x, y, got, dentro, quiere)
			}
		}
	}
}

func TestNuevoMapaSinInterior(t *testing.T) {
	casos := []struct {
		nombre      string
		ancho, alto int
		quiereAncho int
		quiereAlto  int
	}{
		{"2x2", 2, 2, 2, 2},
		{"3x2", 3, 2, 3, 2},
		{"medida cero", 0, 5, 0, 5},
		{"medida negativa", -4, -1, 0, 0},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			m := NuevoMapa(c.ancho, c.alto)
			if m.Ancho != c.quiereAncho || m.Alto != c.quiereAlto {
				t.Fatalf("tamaño = %dx%d, quiere %dx%d", m.Ancho, m.Alto, c.quiereAncho, c.quiereAlto)
			}
			for y := 0; y < m.Alto; y++ {
				for x := 0; x < m.Ancho; x++ {
					if m.EsTransitable(Posicion{x, y}) {
						t.Errorf("(%d,%d) es transitable, quiere pared", x, y)
					}
				}
			}
		})
	}
}

func TestMapaDesdeTexto(t *testing.T) {
	casos := []struct {
		nombre      string
		texto       string
		quiereError bool
		ancho, alto int
	}{
		{"válido", "#.#\n...", false, 3, 2},
		{"con sangría y líneas vacías", `
			#..
			.#.
		`, false, 3, 2},
		{"una sola casilla", ".", false, 1, 1},
		{"vacío", "", true, 0, 0},
		{"solo espacios", " \n\t\n", true, 0, 0},
		{"filas desiguales", "###\n..", true, 0, 0},
		{"línea vacía en medio", "..\n\n..", true, 0, 0},
		{"carácter no válido", "#.x", true, 0, 0},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			m, err := MapaDesdeTexto(c.texto)
			if c.quiereError {
				if err == nil {
					t.Fatalf("quiere error, obtuvo mapa %dx%d", m.Ancho, m.Alto)
				}
				return
			}
			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if m.Ancho != c.ancho || m.Alto != c.alto {
				t.Errorf("tamaño = %dx%d, quiere %dx%d", m.Ancho, m.Alto, c.ancho, c.alto)
			}
		})
	}
}

func TestMapaDesdeTextoCasillas(t *testing.T) {
	m := mapaDeTexto(t, `
		#.
		.#
	`)
	casos := []struct {
		p      Posicion
		quiere Casilla
	}{
		{Posicion{0, 0}, Pared},
		{Posicion{1, 0}, Suelo},
		{Posicion{0, 1}, Suelo},
		{Posicion{1, 1}, Pared},
	}
	for _, c := range casos {
		if got, dentro := m.CasillaEn(c.p); !dentro || got != c.quiere {
			t.Errorf("CasillaEn(%v) = %v, %v; quiere %v, true", c.p, got, dentro, c.quiere)
		}
	}
}

func TestCasillaEnFueraDelMapa(t *testing.T) {
	m := mapaDeTexto(t, "..\n..")
	casos := []Posicion{{-1, 0}, {0, -1}, {2, 0}, {0, 2}, {100, 100}}
	for _, p := range casos {
		got, dentro := m.CasillaEn(p)
		if dentro || got != Pared {
			t.Errorf("CasillaEn(%v) = %v, %v; quiere Pared, false", p, got, dentro)
		}
		if m.EsTransitable(p) {
			t.Errorf("EsTransitable(%v) = true, quiere false", p)
		}
	}
}

func TestEsTransitable(t *testing.T) {
	m := mapaDeTexto(t, "#.")
	if m.EsTransitable(Posicion{0, 0}) {
		t.Error("la pared no debe ser transitable")
	}
	if !m.EsTransitable(Posicion{1, 0}) {
		t.Error("el suelo debe ser transitable")
	}
	var nulo *Mapa
	if nulo.EsTransitable(Posicion{0, 0}) {
		t.Error("un mapa nulo no tiene casillas transitables")
	}
}

func TestDesplazar(t *testing.T) {
	casos := []struct {
		d      Direccion
		quiere Posicion
	}{
		{Arriba, Posicion{5, 4}},
		{Abajo, Posicion{5, 6}},
		{Izquierda, Posicion{4, 5}},
		{Derecha, Posicion{6, 5}},
		{Direccion(99), Posicion{5, 5}},
	}
	for _, c := range casos {
		if got := (Posicion{5, 5}).Desplazar(c.d); got != c.quiere {
			t.Errorf("Desplazar(%d) = %v, quiere %v", c.d, got, c.quiere)
		}
	}
}
