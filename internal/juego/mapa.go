package juego

import (
	"errors"
	"fmt"
	"strings"
)

// Posicion es una coordenada entera en el mapa. X crece hacia la derecha e Y
// crece hacia abajo, como en la pantalla.
type Posicion struct{ X, Y int }

// Desplazar devuelve la posición vecina en la dirección d. Si d no es válida
// devuelve la misma posición.
func (p Posicion) Desplazar(d Direccion) Posicion {
	switch d {
	case Arriba:
		p.Y--
	case Abajo:
		p.Y++
	case Izquierda:
		p.X--
	case Derecha:
		p.X++
	}
	return p
}

// Direccion es uno de los cuatro sentidos en los que se puede mover o disparar.
type Direccion int

const (
	Arriba Direccion = iota
	Abajo
	Izquierda
	Derecha
)

// EsValida indica si d es una de las cuatro direcciones.
func (d Direccion) EsValida() bool {
	return d >= Arriba && d <= Derecha
}

// Casilla es el tipo de terreno de una celda del mapa.
type Casilla int

const (
	Suelo Casilla = iota
	Pared
	Agua
)

// Mapa es la cuadrícula de terreno de la partida.
//
// El mapa es inmutable una vez creado. Por eso una Vista lo comparte por
// puntero en lugar de copiarlo. Si algún día el mapa llega a cambiar durante
// la partida (por ejemplo con construcción), habrá que copiarlo en VistaPara.
type Mapa struct {
	Ancho, Alto int
	casillas    []Casilla // por filas: el índice de (x, y) es y*Ancho+x
}

// NuevoMapa crea un mapa del tamaño indicado con pared en todo el borde y
// suelo en el interior. Una medida negativa se trata como 0. Si alguna medida
// es menor que 3 no queda interior y todo el mapa es pared.
func NuevoMapa(ancho, alto int) *Mapa {
	ancho = max(ancho, 0)
	alto = max(alto, 0)
	m := &Mapa{Ancho: ancho, Alto: alto, casillas: make([]Casilla, ancho*alto)}
	for y := 0; y < alto; y++ {
		for x := 0; x < ancho; x++ {
			if x == 0 || y == 0 || x == ancho-1 || y == alto-1 {
				m.casillas[y*ancho+x] = Pared
			}
		}
	}
	return m
}

// MapaDesdeTexto construye un mapa a partir de texto, una línea por fila:
// "#" es pared y "." es suelo. Pensado para escribir pruebas legibles.
//
// Ignora las líneas vacías del principio y del final, y los espacios y
// tabulaciones alrededor de cada línea, para poder escribir el texto con
// sangría. No agrega pared en el borde: lo que se escribe es lo que se obtiene.
// Devuelve error si el texto está vacío, si las filas tienen anchos distintos
// o si hay un carácter que no sea "#" ni ".".
func MapaDesdeTexto(texto string) (*Mapa, error) {
	lineas := strings.Split(texto, "\n")
	for i := range lineas {
		lineas[i] = strings.TrimSpace(lineas[i])
	}
	for len(lineas) > 0 && lineas[0] == "" {
		lineas = lineas[1:]
	}
	for len(lineas) > 0 && lineas[len(lineas)-1] == "" {
		lineas = lineas[:len(lineas)-1]
	}
	if len(lineas) == 0 {
		return nil, errors.New("el texto del mapa está vacío")
	}

	ancho := len(lineas[0])
	m := &Mapa{Ancho: ancho, Alto: len(lineas), casillas: make([]Casilla, 0, ancho*len(lineas))}
	for y, linea := range lineas {
		if len(linea) != ancho {
			return nil, fmt.Errorf("la fila %d mide %d y la primera mide %d", y, len(linea), ancho)
		}
		for x := 0; x < len(linea); x++ {
			switch linea[x] {
			case '#':
				m.casillas = append(m.casillas, Pared)
			case '.':
				m.casillas = append(m.casillas, Suelo)
			default:
				return nil, fmt.Errorf("carácter %q no válido en fila %d, columna %d", linea[x], y, x)
			}
		}
	}
	return m, nil
}

// CasillaEn devuelve la casilla en p. El segundo valor es false si p cae
// fuera del mapa; en ese caso la casilla es Pared, para que quien ignore el
// segundo valor no trate el exterior como suelo.
func (m *Mapa) CasillaEn(p Posicion) (Casilla, bool) {
	if m == nil || p.X < 0 || p.Y < 0 || p.X >= m.Ancho || p.Y >= m.Alto {
		return Pared, false
	}
	return m.casillas[p.Y*m.Ancho+p.X], true
}

// EsTransitable indica si un jugador puede estar parado en p. Solo el suelo
// dentro del mapa es transitable.
func (m *Mapa) EsTransitable(p Posicion) bool {
	c, dentro := m.CasillaEn(p)
	return dentro && c == Suelo
}
