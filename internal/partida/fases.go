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

// PermiteUnirse indica si en la fase actual pueden entrar jugadores nuevos.
//
// TODO: provisional. Hoy solo se puede entrar en la sala de espera; las reglas
// reales de cada fase son el paso 6 del orden de trabajo.
func (m *Maquina) PermiteUnirse() bool {
	return m.fase == SalaDeEspera
}

// AccionPermitida indica si la acción a se puede aplicar en la fase actual.
// Por ejemplo, no se dispara en la sala de espera.
//
// TODO: provisional. Hoy en la sala de espera solo se permite moverse y en
// las demás fases no se permite nada; las reglas reales son el paso 6.
func (m *Maquina) AccionPermitida(a juego.Accion) bool {
	if m.fase == SalaDeEspera {
		_, esMover := a.(juego.MoverAccion)
		return esMover
	}
	return false
}

// Avanzar revisa el estado ya actualizado del tick y hace las transiciones
// que correspondan: iniciar la partida, detectar al último jugador vivo y
// pasar a Final.
func (m *Maquina) Avanzar(e *juego.Estado) {
	// TODO: implementar transiciones.
}
