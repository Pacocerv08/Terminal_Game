package juego

import "sort"

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

// VidaInicial es la vida con la que entra un jugador a la partida.
const VidaInicial = 100

// hayJugadorVivoEn indica si algún jugador vivo está parado en p.
func (e *Estado) hayJugadorVivoEn(p Posicion) bool {
	for _, j := range e.Jugadores {
		if j.Pos == p && j.EstaVivo() {
			return true
		}
	}
	return false
}

// AgregarJugador crea un jugador con el siguiente ID libre (desde 1) y lo
// coloca en una casilla transitable y sin jugador vivo. Lo devuelve, o nil si
// no hay ninguna casilla libre; en ese caso no consume ningún ID.
//
// La casilla es la primera que se encuentra recorriendo el mapa por filas, de
// arriba abajo y de izquierda a derecha, así el resultado es determinista.
//
// TODO: tomar la primera casilla libre es provisional; después se repartirá
// a los jugadores con azar controlado por semilla.
func (e *Estado) AgregarJugador(nombre string) *Jugador {
	if e.Mapa == nil {
		return nil
	}
	for y := 0; y < e.Mapa.Alto; y++ {
		for x := 0; x < e.Mapa.Ancho; x++ {
			p := Posicion{X: x, Y: y}
			if !e.Mapa.EsTransitable(p) || e.hayJugadorVivoEn(p) {
				continue
			}
			e.siguienteID++
			j := &Jugador{ID: e.siguienteID, Nombre: nombre, Pos: p, Vida: VidaInicial}
			e.Jugadores[j.ID] = j
			return j
		}
	}
	return nil
}

// QuitarJugador saca al jugador id de la partida. Si no existe, no hace nada.
// Su casilla queda libre y su ID no se vuelve a asignar.
func (e *Estado) QuitarJugador(id IDJugador) {
	delete(e.Jugadores, id)
}

// Avanzar mueve el estado un tick: cuenta el tick y, más adelante, cierra la
// zona y aplica su daño.
func (e *Estado) Avanzar() {
	e.Tick++
	// TODO: cierre de zona y daño de zona.
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
	// El recorrido de un map es aleatorio: se ordena por ID para que la vista
	// sea determinista.
	sort.Slice(v.Jugadores, func(a, b int) bool { return v.Jugadores[a].ID < v.Jugadores[b].ID })
	return v
}
