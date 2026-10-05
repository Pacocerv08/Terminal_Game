// Comando servidor es el punto de entrada: une las piezas y arranca el
// servidor.
package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/Pacocerv08/Terminal_Game/internal/ciclo"
	"github.com/Pacocerv08/Terminal_Game/internal/entrada"
	"github.com/Pacocerv08/Terminal_Game/internal/graficos/plano"
	"github.com/Pacocerv08/Terminal_Game/internal/juego"
	"github.com/Pacocerv08/Terminal_Game/internal/partida"
	"github.com/Pacocerv08/Terminal_Game/internal/red"
)

// TODO: mover el tamaño del mapa, el intervalo y el tamaño del canal a la
// configuración cuando se implemente el mapa grande.
const (
	anchoMapa         = 200
	altoMapa          = 200
	intervaloTick     = 100 * time.Millisecond
	capacidadMensajes = 256
)

func main() {
	// Los argumentos tienen prioridad sobre las variables de entorno.
	cfg := red.ConfigDesdeEntorno()
	flag.StringVar(&cfg.Direccion, "direccion", cfg.Direccion, "dirección y puerto donde escucha, o "+red.VariableDireccion)
	flag.StringVar(&cfg.RutaClaveHost, "clave", cfg.RutaClaveHost, "archivo con la clave del host SSH, o "+red.VariableClaveHost)
	flag.Parse()
	if err := cfg.Validar(); err != nil {
		log.Fatal(err)
	}

	estado := juego.NuevoEstado(juego.NuevoMapa(anchoMapa, altoMapa))
	maquina := partida.NuevaMaquina()
	mensajes := make(chan ciclo.Mensaje, capacidadMensajes)

	bucle := ciclo.Nuevo(intervaloTick, estado, maquina, mensajes)
	// TODO: la semilla debe ser aleatoria al arrancar el servidor. Mientras no
	// se llame a FijarSemilla vale 0, así que el orden de desempate de cada
	// partida es siempre el mismo y los jugadores podrían predecirlo.
	servidor := red.NuevoServidor(cfg, mensajes, plano.NuevoPlano(), entrada.TraductorTeclado{})

	ctx, parar := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer parar()

	go bucle.Ejecutar(ctx)
	if err := servidor.Iniciar(ctx); err != nil {
		log.Fatal(err)
	}
}
