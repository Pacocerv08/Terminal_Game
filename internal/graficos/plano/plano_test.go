package plano

import (
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/Pacocerv08/Terminal_Game/internal/juego"
	"github.com/Pacocerv08/Terminal_Game/internal/partida"
)

const mapaGrande = `
##########
#........#
#........#
#........#
#........#
#........#
#........#
##########
`

const mapaChico = `
####
#..#
####
`

func nuevoMapa(t *testing.T, texto string) *juego.Mapa {
	t.Helper()
	m, err := juego.MapaDesdeTexto(texto)
	if err != nil {
		t.Fatalf("mapa de prueba no válido: %v", err)
	}
	return m
}

// vistaCon arma una vista donde el jugador 1 (yo) está en pos con vida 100.
func vistaCon(t *testing.T, texto string, pos juego.Posicion) juego.Vista {
	t.Helper()
	return juego.Vista{
		Yo:        1,
		Mapa:      nuevoMapa(t, texto),
		Jugadores: []juego.JugadorVisible{{ID: 1, Nombre: "yo", Pos: pos, Vida: 100}},
	}
}

// lineasDeMapa dibuja, comprueba que hay exactamente alto líneas y devuelve
// todas menos el indicador.
func lineasDeMapa(t *testing.T, v juego.Vista, ancho, alto int) []string {
	t.Helper()
	lineas := strings.Split(NuevoPlano().Dibujar(v, partida.EnCurso, ancho, alto), "\n")
	if len(lineas) != alto {
		t.Fatalf("hay %d líneas y se esperaban %d", len(lineas), alto)
	}
	return lineas[:alto-1]
}

func comparar(t *testing.T, got, want []string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("salida distinta\n got: %q\nwant: %q", got, want)
	}
}

// Pantalla de 10x4: ventana de 5 casillas de ancho por 3 de alto.
const (
	anchoPantalla = 10
	altoPantalla  = 4
)

func TestJugadorCentrado(t *testing.T) {
	v := vistaCon(t, mapaGrande, juego.Posicion{X: 5, Y: 4})
	comparar(t, lineasDeMapa(t, v, anchoPantalla, altoPantalla), []string{
		". . . . . ",
		". . @@. . ",
		". . . . . ",
	})
}

func TestCamaraSeDetieneEnLosBordes(t *testing.T) {
	casos := []struct {
		nombre string
		pos    juego.Posicion
		want   []string
	}{
		{"arriba", juego.Posicion{X: 5, Y: 1}, []string{"##########", ". . @@. . ", ". . . . . "}},
		{"abajo", juego.Posicion{X: 5, Y: 6}, []string{". . . . . ", ". . @@. . ", "##########"}},
		{"izquierda", juego.Posicion{X: 1, Y: 4}, []string{"##. . . . ", "##@@. . . ", "##. . . . "}},
		{"derecha", juego.Posicion{X: 8, Y: 4}, []string{". . . . ##", ". . . @@##", ". . . . ##"}},
		{"esquina", juego.Posicion{X: 1, Y: 1}, []string{"##########", "##@@. . . ", "##. . . . "}},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			v := vistaCon(t, mapaGrande, c.pos)
			comparar(t, lineasDeMapa(t, v, anchoPantalla, altoPantalla), c.want)
		})
	}
}

func TestMapaMasChicoQuePantalla(t *testing.T) {
	v := vistaCon(t, mapaChico, juego.Posicion{X: 1, Y: 1})
	t.Run("medidas pares", func(t *testing.T) {
		comparar(t, lineasDeMapa(t, v, 12, 6), []string{
			"",
			"  ########",
			"  ##@@. ##",
			"  ########",
			"",
		})
	})
	t.Run("medidas impares", func(t *testing.T) {
		comparar(t, lineasDeMapa(t, v, 11, 5), []string{
			"########",
			"##@@. ##",
			"########",
			"",
		})
	})
}

func TestOtrosJugadores(t *testing.T) {
	casos := []struct {
		nombre string
		otro   juego.JugadorVisible
		want   []string
	}{
		{"dentro de la pantalla", juego.JugadorVisible{ID: 2, Pos: juego.Posicion{X: 6, Y: 4}, Vida: 50},
			[]string{". . . . . ", ". . @@&&. ", ". . . . . "}},
		{"fuera de la pantalla", juego.JugadorVisible{ID: 2, Pos: juego.Posicion{X: 1, Y: 1}, Vida: 50},
			[]string{". . . . . ", ". . @@. . ", ". . . . . "}},
		{"muerto no se dibuja", juego.JugadorVisible{ID: 2, Pos: juego.Posicion{X: 6, Y: 4}, Vida: 0},
			[]string{". . . . . ", ". . @@. . ", ". . . . . "}},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			v := vistaCon(t, mapaGrande, juego.Posicion{X: 5, Y: 4})
			v.Jugadores = append(v.Jugadores, c.otro)
			comparar(t, lineasDeMapa(t, v, anchoPantalla, altoPantalla), c.want)
		})
	}
}

func TestYoMuertoNoSeDibujaPeroLaCamaraLoSigue(t *testing.T) {
	v := vistaCon(t, mapaGrande, juego.Posicion{X: 1, Y: 1})
	v.Jugadores[0].Vida = 0
	comparar(t, lineasDeMapa(t, v, anchoPantalla, altoPantalla), []string{
		"##########",
		"##. . . . ",
		"##. . . . ",
	})
}

func TestBotin(t *testing.T) {
	v := vistaCon(t, mapaGrande, juego.Posicion{X: 5, Y: 4})
	v.Botin = []juego.ObjetoEnSuelo{
		{Objeto: juego.Objeto{Tipo: juego.ObjetoMunicion, Cantidad: 5}, Pos: juego.Posicion{X: 4, Y: 4}},
		{Objeto: juego.Objeto{Tipo: juego.ObjetoBotiquin, Cantidad: 1}, Pos: juego.Posicion{X: 5, Y: 4}}, // bajo mis pies
	}
	comparar(t, lineasDeMapa(t, v, anchoPantalla, altoPantalla), []string{
		". . . . . ",
		". **@@. . ",
		". . . . . ",
	})
}

func TestIndicador(t *testing.T) {
	v := vistaCon(t, mapaGrande, juego.Posicion{X: 12, Y: 7})
	casos := []struct {
		fase  partida.Fase
		ancho int
		want  string
	}{
		{partida.SalaDeEspera, 60, "Vida: 100 | Pos: 12,7 | Sala de espera"},
		{partida.EnCurso, 60, "Vida: 100 | Pos: 12,7 | En curso"},
		{partida.Final, 60, "Vida: 100 | Pos: 12,7 | Final"},
		{partida.Fase(99), 60, "Vida: 100 | Pos: 12,7 | ?"},
		{partida.EnCurso, 10, "Vida: 100 "},
	}
	for _, c := range casos {
		if got := NuevoPlano().Dibujar(v, c.fase, c.ancho, 1); got != c.want {
			t.Errorf("fase %v ancho %d: got %q, want %q", c.fase, c.ancho, got, c.want)
		}
	}
}

func TestTruncarCuentaRunas(t *testing.T) {
	casos := []struct {
		s    string
		n    int
		want string
	}{
		{"ñandú", 3, "ñan"},
		{"ñandú", 5, "ñandú"},
		{"ñandú", 9, "ñandú"},
		{"日本語", 2, "日本"},
		{"abc", 0, ""},
	}
	for _, c := range casos {
		got := truncar(c.s, c.n)
		if got != c.want || !utf8.ValidString(got) {
			t.Errorf("truncar(%q, %d) = %q, se esperaba %q", c.s, c.n, got, c.want)
		}
	}
}

func TestSimbolosMidenDosCaracteresASCII(t *testing.T) {
	simbolos := map[string]string{
		"pared": simboloPared, "suelo": simboloSuelo, "agua": simboloAgua,
		"yo": simboloYo, "otro": simboloOtro, "botin": simboloBotin, "vacio": simboloVacio,
	}
	for nombre, s := range simbolos {
		if utf8.RuneCountInString(s) != 2 || len(s) != 2 {
			t.Errorf("%s = %q no mide exactamente dos caracteres", nombre, s)
		}
		for _, r := range s {
			if r < 0x20 || r > 0x7E {
				t.Errorf("%s = %q tiene un carácter que no es ASCII imprimible: %q", nombre, s, r)
			}
		}
	}
}

func TestTamanosInvalidosYLimites(t *testing.T) {
	v := vistaCon(t, mapaGrande, juego.Posicion{X: 5, Y: 4})
	v.Jugadores = append(v.Jugadores, juego.JugadorVisible{ID: 2, Pos: juego.Posicion{X: 6, Y: 4}, Vida: 10})
	medidas := [][2]int{
		{0, 0}, {-5, -5}, {0, 10}, {10, 0}, {-1, 5}, {5, -1},
		{1, 1}, {1, 5}, {5, 1}, {2, 2}, {3, 3}, {2, 1},
		{199, 99}, {200, 100}, {201, 101}, {100000, 100000}, {80, 24},
	}
	for _, m := range medidas {
		ancho, alto := m[0], m[1]
		salida := NuevoPlano().Dibujar(v, partida.EnCurso, ancho, alto)

		altoEsperado := min(max(alto, 0), altoMaximo)
		anchoMax := min(max(ancho, 0), anchoMaximo)
		if altoEsperado == 0 {
			if salida != "" {
				t.Errorf("%dx%d: se esperaba salida vacía y salió %q", ancho, alto, salida)
			}
			continue
		}
		lineas := strings.Split(salida, "\n")
		if len(lineas) != altoEsperado {
			t.Errorf("%dx%d: %d líneas, se esperaban %d", ancho, alto, len(lineas), altoEsperado)
		}
		for i, l := range lineas {
			if n := utf8.RuneCountInString(l); n > anchoMax {
				t.Errorf("%dx%d: la línea %d mide %d y el máximo es %d", ancho, alto, i, n, anchoMax)
			}
		}
	}
}

func TestVistasIncompletasNoFallan(t *testing.T) {
	t.Run("mapa nil", func(t *testing.T) {
		v := juego.Vista{Yo: 1, Jugadores: []juego.JugadorVisible{{ID: 1, Vida: 100}}}
		comparar(t, lineasDeMapa(t, v, 10, 4), []string{"", "", ""})
	})
	t.Run("yo inexistente", func(t *testing.T) {
		v := vistaCon(t, mapaChico, juego.Posicion{X: 1, Y: 1})
		v.Yo = 99
		salida := NuevoPlano().Dibujar(v, partida.EnCurso, 60, 5)
		lineas := strings.Split(salida, "\n")
		if got, want := lineas[len(lineas)-1], "Vida: 0 | Pos: -,- | En curso"; got != want {
			t.Errorf("indicador %q, se esperaba %q", got, want)
		}
	})
}
