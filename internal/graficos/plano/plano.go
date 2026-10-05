package plano

import (
	"github.com/Pacocerv08/Terminal_Game/internal/graficos"
	"github.com/Pacocerv08/Terminal_Game/internal/juego"
	"github.com/Pacocerv08/Terminal_Game/internal/partida"
)

var _ graficos.Dibujante = (*Plano)(nil)

// Plano dibuja el mapa visto desde arriba, centrado en el jugador de la vista.
type Plano struct{}

// NuevoPlano crea un dibujante de vista desde arriba.
func NuevoPlano() *Plano {
	return &Plano{}
}

// Dibujar implementa graficos.Dibujante.
func (p *Plano) Dibujar(v juego.Vista, fase partida.Fase, ancho, alto int) string {
	// TODO: calcular la cámara y dibujar casillas, jugadores, botín y zona.
	return ""
}
