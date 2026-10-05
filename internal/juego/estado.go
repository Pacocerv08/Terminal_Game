package juego

// Accion es algo que un jugador pide hacer. Aplicar recibe al autor porque
// la identidad la asigna la sesión en el servidor, no el contenido de la
// acción.
//
// Aplicar ignora la acción si el autor no existe o no está vivo, y valida los
// campos de la propia acción.
type Accion interface {
	Aplicar(e *Estado, autor IDJugador)
}

// Estado es todo lo que ocurre en la partida. Solo el ciclo de juego lo
// modifica.
type Estado struct {
	Mapa      *Mapa
	Jugadores map[IDJugador]*Jugador
	Botin     []ObjetoEnSuelo // botín en el suelo
	Zona      Zona
	Tick      uint64

	siguienteID IDJugador
}

// NuevoEstado crea un estado vacío sobre el mapa m.
func NuevoEstado(m *Mapa) *Estado {
	return &Estado{
		Mapa:      m,
		Jugadores: make(map[IDJugador]*Jugador),
	}
}

// AgregarJugador crea un jugador con el siguiente ID libre y lo devuelve.
func (e *Estado) AgregarJugador(nombre string) *Jugador {
	// TODO: asignar ID, posición inicial y vida.
	return nil
}

// QuitarJugador saca al jugador id de la partida.
func (e *Estado) QuitarJugador(id IDJugador) {
	// TODO: implementar.
}

// Avanzar mueve el estado un tick: cierre de zona, daño de zona, etc.
func (e *Estado) Avanzar() {
	// TODO: implementar.
}

// Vista es una copia de solo lectura del estado, desde el punto de vista de
// un jugador. No comparte memoria con Estado, salvo el mapa.
//
// El mapa se comparte por puntero porque es inmutable. Si el mapa llega a
// cambiar durante la partida, VistaPara tendrá que copiarlo.
type Vista struct {
	Yo        IDJugador
	Tick      uint64
	Mapa      *Mapa
	Jugadores []JugadorVisible
	Botin     []ObjetoEnSuelo
	Mi        Inventario // copia del inventario propio
	Zona      Zona
}

// VistaPara construye la vista del jugador id. Los slices y el inventario
// son copias completas, sin memoria compartida con Estado.
func (e *Estado) VistaPara(id IDJugador) Vista {
	v := Vista{
		Yo:    id,
		Tick:  e.Tick,
		Mapa:  e.Mapa,
		Botin: append([]ObjetoEnSuelo(nil), e.Botin...),
		Zona:  e.Zona,
	}
	if propio, ok := e.Jugadores[id]; ok {
		v.Mi = propio.Inventario.Copiar()
	}
	// TODO: mostrar solo los jugadores dentro del campo de visión.
	for _, j := range e.Jugadores {
		v.Jugadores = append(v.Jugadores, JugadorVisible{
			ID: j.ID, Nombre: j.Nombre, Pos: j.Pos, Vida: j.Vida,
		})
	}
	return v
}
