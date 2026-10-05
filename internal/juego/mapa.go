package juego

// Posicion es una coordenada entera en el mapa.
type Posicion struct{ X, Y int }

// Direccion es uno de los cuatro sentidos en los que se puede mover o disparar.
type Direccion int

const (
	Arriba Direccion = iota
	Abajo
	Izquierda
	Derecha
)

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
	casillas    []Casilla
}

// NuevoMapa crea un mapa del tamaño indicado.
func NuevoMapa(ancho, alto int) *Mapa {
	// TODO: reservar las casillas y generar el terreno.
	return &Mapa{Ancho: ancho, Alto: alto}
}

// CasillaEn devuelve la casilla en p. El segundo valor es false si p cae
// fuera del mapa.
func (m *Mapa) CasillaEn(p Posicion) (Casilla, bool) {
	// TODO: implementar.
	return Suelo, false
}

// EsTransitable indica si un jugador puede estar parado en p.
func (m *Mapa) EsTransitable(p Posicion) bool {
	// TODO: implementar.
	return false
}
