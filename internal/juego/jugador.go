package juego

// IDJugador identifica a un jugador dentro de una partida. Lo asigna el
// ciclo al unirse el jugador; nunca viene de lo que escribe el cliente.
type IDJugador int

// Jugador es un participante de la partida, sea humano o bot.
type Jugador struct {
	ID         IDJugador
	Nombre     string
	Pos        Posicion
	Vida       int
	Inventario Inventario
}

// EstaVivo indica si el jugador sigue en la partida.
func (j *Jugador) EstaVivo() bool {
	return j.Vida > 0
}

// Mover desplaza al jugador una casilla en la dirección d. Devuelve false y
// no lo mueve si d no es válida o si la casilla destino no es transitable
// (pared o fuera del mapa). Solo mira el terreno: que otro jugador ocupe la
// casilla lo comprueba MoverAccion.Aplicar, que conoce el Estado.
func (j *Jugador) Mover(d Direccion, m *Mapa) bool {
	if !d.EsValida() {
		return false
	}
	destino := j.Pos.Desplazar(d)
	if !m.EsTransitable(destino) {
		return false
	}
	j.Pos = destino
	return true
}

// JugadorVisible es lo que otros jugadores pueden saber de un jugador.
// Es una copia, sin inventario.
type JugadorVisible struct {
	ID     IDJugador
	Nombre string
	Pos    Posicion
	Vida   int
}

// MoverAccion pide mover al autor una casilla en la dirección Dir.
//
// Regla de empate: si dos jugadores quieren la misma casilla en un tick, gana
// la acción que se aplica primero; la otra encuentra la casilla ocupada y no
// se mueve. El orden de aplicación lo define quien llama a Aplicar y debe ser
// determinista (por ejemplo, por ID de jugador o por orden de llegada), nunca
// el recorrido de un map.
type MoverAccion struct{ Dir Direccion }

// Aplicar mueve al autor una casilla. No hace nada si Dir no es válida, si el
// autor no existe o no está vivo, si el destino es pared o está fuera del
// mapa, o si otro jugador vivo ocupa el destino. Un jugador muerto no ocupa
// su casilla.
func (a MoverAccion) Aplicar(e *Estado, autor IDJugador) {
	if !a.Dir.EsValida() {
		return
	}
	j, ok := e.Jugadores[autor]
	if !ok || !j.EstaVivo() {
		return
	}
	if e.hayJugadorVivoEn(j.Pos.Desplazar(a.Dir)) {
		return
	}
	j.Mover(a.Dir, e.Mapa)
}
