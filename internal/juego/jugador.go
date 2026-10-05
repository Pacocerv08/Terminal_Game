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

// Mover intenta desplazar al jugador una casilla en la dirección d. Devuelve
// false si la casilla destino no es transitable.
func (j *Jugador) Mover(d Direccion, m *Mapa) bool {
	// TODO: implementar movimiento con colisión.
	return false
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
type MoverAccion struct{ Dir Direccion }

// Aplicar mueve al autor si existe y está vivo.
func (a MoverAccion) Aplicar(e *Estado, autor IDJugador) {
	// TODO: validar Dir y llamar a Jugador.Mover.
}
