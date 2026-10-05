package partida

import "github.com/Pacocerv08/Terminal_Game/internal/juego"

// Fase es la etapa en la que está la partida.
type Fase int

const (
	SalaDeEspera Fase = iota
	EnCurso
	Final
)

// Maquina lleva la fase actual de la partida y decide las transiciones.
type Maquina struct{ fase Fase }

// NuevaMaquina crea una máquina en la sala de espera.
func NuevaMaquina() *Maquina {
	return &Maquina{fase: SalaDeEspera}
}

// Fase devuelve la fase actual.
func (m *Maquina) Fase() Fase {
	return m.fase
}

// AccionPermitida indica si la acción a se puede aplicar en la fase actual.
// Por ejemplo, no se dispara en la sala de espera.
func (m *Maquina) AccionPermitida(a juego.Accion) bool {
	// TODO: decidir según la fase y el tipo de acción.
	return false
}

// Avanzar revisa el estado ya actualizado del tick y hace las transiciones
// que correspondan: iniciar la partida, detectar al último jugador vivo y
// pasar a Final.
func (m *Maquina) Avanzar(e *juego.Estado) {
	// TODO: implementar transiciones.
}
