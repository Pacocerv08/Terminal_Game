package red

import (
	"context"
	"errors"
	"os"

	"github.com/Pacocerv08/Terminal_Game/internal/ciclo"
	"github.com/Pacocerv08/Terminal_Game/internal/entrada"
	"github.com/Pacocerv08/Terminal_Game/internal/graficos"
)

// Nombres de las variables de entorno de la configuración.
const (
	VariableDireccion = "TERMINAL_GAME_DIRECCION"
	VariableClaveHost = "TERMINAL_GAME_CLAVE_HOST"
)

// Config son los datos de red del servidor. Nada de esto va fijo en el código.
type Config struct {
	Direccion     string // dirección y puerto donde escucha, por ejemplo ":2222"
	RutaClaveHost string // archivo con la clave del host SSH
}

// ConfigDesdeEntorno lee la configuración de las variables de entorno. Los
// valores que falten quedan vacíos; Validar detecta cuáles.
func ConfigDesdeEntorno() Config {
	return Config{
		Direccion:     os.Getenv(VariableDireccion),
		RutaClaveHost: os.Getenv(VariableClaveHost),
	}
}

// Validar devuelve un error si falta algún dato obligatorio.
func (c Config) Validar() error {
	if c.Direccion == "" {
		return errors.New("falta la dirección: use el argumento -direccion o la variable " + VariableDireccion)
	}
	if c.RutaClaveHost == "" {
		return errors.New("falta la clave del host: use el argumento -clave o la variable " + VariableClaveHost)
	}
	return nil
}

// Servidor atiende las conexiones SSH y crea una Sesion por jugador.
type Servidor struct {
	cfg       Config
	mensajes  chan<- ciclo.Mensaje
	dibujante graficos.Dibujante
	traductor entrada.Traductor
}

// NuevoServidor crea el servidor. Las sesiones envían sus mensajes al ciclo
// por mensajes.
func NuevoServidor(cfg Config, mensajes chan<- ciclo.Mensaje, d graficos.Dibujante, t entrada.Traductor) *Servidor {
	return &Servidor{cfg: cfg, mensajes: mensajes, dibujante: d, traductor: t}
}

// Iniciar escucha conexiones hasta que ctx se cancele.
func (s *Servidor) Iniciar(ctx context.Context) error {
	// TODO: levantar el servidor SSH con Wish y crear una Sesion por conexión.
	return nil
}
