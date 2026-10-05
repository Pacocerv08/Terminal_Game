package juego

import "testing"

// TODO: prueba de que VistaPara no comparte memoria con Estado.
// Debe crear un estado con jugadores, botín e inventario, pedir una vista,
// modificar la vista (Jugadores, Botin y Mi) y comprobar que el estado no
// cambió. También debe hacerlo al revés: cambiar el estado después de pedir
// la vista y comprobar que la vista no cambió. El mapa queda fuera porque se
// comparte a propósito (es inmutable).
func TestVistaParaNoComparteMemoria(t *testing.T) {
	t.Skip("TODO: pendiente de escribir")
}
