package juego

// Zona es el área segura de la partida. Se va cerrando con el tiempo y quien
// queda fuera recibe daño.
type Zona struct {
	Centro Posicion
	Radio  int
}

// Contiene indica si p está dentro de la zona segura.
func (z *Zona) Contiene(p Posicion) bool {
	// TODO: implementar.
	return false
}

// Reducir cierra la zona un paso.
func (z *Zona) Reducir() {
	// TODO: implementar.
}
