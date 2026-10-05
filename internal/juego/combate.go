package juego

// DispararAccion pide que el autor dispare en la dirección Dir.
type DispararAccion struct{ Dir Direccion }

// Aplicar dispara si el autor existe y está vivo.
func (a DispararAccion) Aplicar(e *Estado, autor IDJugador) {
	// TODO: validar Dir y llamar a Estado.ResolverDisparo.
}

// ResolverDisparo calcula el resultado de un disparo del autor en la
// dirección dir: a quién golpea, cuánta vida quita y si lo elimina.
func (e *Estado) ResolverDisparo(autor IDJugador, dir Direccion) {
	// TODO: implementar disparo, daño y eliminación.
}
