package plano

import (
	"fmt"
	"strings"

	"github.com/Pacocerv08/Terminal_Game/internal/graficos"
	"github.com/Pacocerv08/Terminal_Game/internal/juego"
	"github.com/Pacocerv08/Terminal_Game/internal/partida"
)

var _ graficos.Dibujante = (*Plano)(nil)

// Medidas máximas de pantalla. El cliente informa el tamaño de su terminal, no
// es confiable: lo que pase de aquí se recorta.
const (
	anchoMaximo = 200
	altoMaximo  = 100
)

// Símbolos del dibujo. Cada casilla ocupa dos caracteres ASCII imprimibles,
// así se ve cuadrada en la terminal.
const (
	simboloPared = "##"
	simboloSuelo = ". "
	simboloAgua  = "~~"
	simboloYo    = "@@"
	simboloOtro  = "&&"
	simboloBotin = "**"
	simboloVacio = "  " // fuera del mapa, cuando el mapa es más chico que la pantalla
	anchoCasilla = 2
	lineasFijas  = 1 // la última línea es el indicador
)

// Plano dibuja el mapa visto desde arriba, centrado en el jugador de la vista.
type Plano struct{}

// NuevoPlano crea un dibujante de vista desde arriba.
func NuevoPlano() *Plano {
	return &Plano{}
}

// Dibujar implementa graficos.Dibujante.
//
// La salida tiene exactamente alto líneas, unidas con "\n" y sin salto final,
// y ninguna supera ancho caracteres. Antes se limitan las medidas: las
// negativas valen 0 y las que pasan del máximo se recortan, así que con alto 0
// devuelve "". La última línea es el indicador y el resto es el mapa.
//
// La cámara va centrada en el jugador de la vista. Si el mapa es más grande que
// la ventana, la cámara se detiene en el borde y nunca muestra el exterior; si
// es más chico, el mapa se centra. Con medidas pares el centro cae una casilla
// hacia arriba o a la izquierda.
func (p *Plano) Dibujar(v juego.Vista, fase partida.Fase, ancho, alto int) string {
	ancho = min(max(ancho, 0), anchoMaximo)
	alto = min(max(alto, 0), altoMaximo)
	if alto == 0 {
		return ""
	}

	yo, hayYo := buscarYo(v)
	otros := indexarOtros(v)
	botin := indexarBotin(v)

	ventanaX := ancho / anchoCasilla
	ventanaY := alto - lineasFijas
	anchoMapa, altoMapa := 0, 0
	if v.Mapa != nil {
		anchoMapa, altoMapa = v.Mapa.Ancho, v.Mapa.Alto
	}
	centro := juego.Posicion{X: anchoMapa / 2, Y: altoMapa / 2}
	if hayYo {
		centro = yo.Pos
	}
	inicioX, margenX := ejeCamara(centro.X, ventanaX, anchoMapa)
	inicioY, margenY := ejeCamara(centro.Y, ventanaY, altoMapa)
	columnas := min(anchoMapa, ventanaX)

	lineas := make([]string, 0, alto)
	for fila := 0; fila < ventanaY; fila++ {
		y := inicioY + fila - margenY
		if y < 0 || y >= altoMapa {
			lineas = append(lineas, "")
			continue
		}
		var b strings.Builder
		b.WriteString(strings.Repeat(simboloVacio, margenX))
		for x := inicioX; x < inicioX+columnas; x++ {
			pos := juego.Posicion{X: x, Y: y}
			b.WriteString(simboloEn(v, pos, yo, hayYo, otros, botin))
		}
		lineas = append(lineas, b.String())
	}
	lineas = append(lineas, truncar(indicador(yo, hayYo, fase), ancho))
	return strings.Join(lineas, "\n")
}

// ejeCamara calcula, para un eje, la primera casilla del mapa que se ve y el
// margen en casillas vacías que se deja antes del mapa. Si el mapa cabe en la
// ventana se centra (margen); si no, la ventana se centra en centro sin salirse
// del mapa (inicio).
func ejeCamara(centro, ventana, mapa int) (inicio, margen int) {
	if mapa <= ventana {
		return 0, (ventana - mapa) / 2
	}
	return min(max(centro-ventana/2, 0), mapa-ventana), 0
}

// buscarYo devuelve al jugador de la vista, esté vivo o no.
func buscarYo(v juego.Vista) (juego.JugadorVisible, bool) {
	for _, j := range v.Jugadores {
		if j.ID == v.Yo {
			return j, true
		}
	}
	return juego.JugadorVisible{}, false
}

// indexarOtros indexa por posición a los jugadores vivos que no son el de la
// vista. Si dos comparten casilla gana el de menor ID, para que el resultado
// no dependa del orden del slice.
func indexarOtros(v juego.Vista) map[juego.Posicion]juego.IDJugador {
	otros := make(map[juego.Posicion]juego.IDJugador, len(v.Jugadores))
	for _, j := range v.Jugadores {
		if j.ID == v.Yo || j.Vida <= 0 {
			continue
		}
		if previo, ok := otros[j.Pos]; !ok || j.ID < previo {
			otros[j.Pos] = j.ID
		}
	}
	return otros
}

// indexarBotin indexa por posición el botín que hay en el suelo.
func indexarBotin(v juego.Vista) map[juego.Posicion]struct{} {
	botin := make(map[juego.Posicion]struct{}, len(v.Botin))
	for _, o := range v.Botin {
		botin[o.Pos] = struct{}{}
	}
	return botin
}

// simboloEn devuelve el símbolo de la casilla pos. Prioridad: yo, otro
// jugador, botín y terreno.
func simboloEn(v juego.Vista, pos juego.Posicion, yo juego.JugadorVisible, hayYo bool,
	otros map[juego.Posicion]juego.IDJugador, botin map[juego.Posicion]struct{}) string {
	if hayYo && yo.Vida > 0 && yo.Pos == pos {
		return simboloYo
	}
	if _, ok := otros[pos]; ok {
		return simboloOtro
	}
	if _, ok := botin[pos]; ok {
		return simboloBotin
	}
	c, _ := v.Mapa.CasillaEn(pos)
	switch c {
	case juego.Suelo:
		return simboloSuelo
	case juego.Agua:
		return simboloAgua
	default:
		return simboloPared
	}
}

// indicador es la línea de abajo: vida, posición y fase.
func indicador(yo juego.JugadorVisible, hayYo bool, fase partida.Fase) string {
	vida, pos := 0, "-,-"
	if hayYo {
		vida = max(yo.Vida, 0)
		pos = fmt.Sprintf("%d,%d", yo.Pos.X, yo.Pos.Y)
	}
	return fmt.Sprintf("Vida: %d | Pos: %s | %s", vida, pos, nombreFase(fase))
}

// nombreFase es el texto de la fase en el indicador.
func nombreFase(f partida.Fase) string {
	switch f {
	case partida.SalaDeEspera:
		return "Sala de espera"
	case partida.EnCurso:
		return "En curso"
	case partida.Final:
		return "Final"
	default:
		return "?"
	}
}

// truncar recorta s a como máximo n caracteres, contados por runas y no por
// bytes, para no partir un carácter a la mitad.
func truncar(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}
