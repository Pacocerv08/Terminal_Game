# Terminal_Game

## Qué es

Battle royale multijugador en tiempo real que se juega en terminal.

- Los jugadores entran con un solo comando `ssh`, sin instalar nada.
- Servidor autoritativo: toda la lógica y todo el dibujo ocurren en el servidor. No existe programa cliente.
- Vista desde arriba, con una cámara que sigue a cada jugador.
- El juego avanza con un ciclo de ticks a intervalo fijo.

## Tecnología

- Lenguaje: Go. La versión mínima está en `go.mod`.
- Wish: servidor SSH.
- Bubble Tea: interfaz de terminal.
- Módulo: `github.com/Pacocerv08/Terminal_Game`

## Mapa de paquetes

| Ruta | Responsabilidad |
| --- | --- |
| `cmd/servidor` | Punto de entrada. Une las piezas y arranca el servidor. |
| `internal/juego` | Estado y reglas: mapa, jugadores, combate, botín, zona. |
| `internal/partida` | Fases de la partida (sala de espera, en curso, final) como máquina de estados. |
| `internal/bots` | Jugadores controlados por el servidor. |
| `internal/ciclo` | Ciclo de ticks a intervalo fijo. |
| `internal/red` | Servidor SSH y manejo de sesiones. |
| `internal/entrada` | Traduce teclas en acciones del juego. |
| `internal/graficos` | Interfaz de dibujo. |
| `internal/graficos/plano` | Vista desde arriba con cámara. |

## Reglas de arquitectura

1. `internal/juego` solo usa la biblioteca estándar de Go. No importa ningún otro paquete del proyecto ni bibliotecas de red o de terminal.
2. Las dependencias apuntan hacia `juego`, nunca al revés. No hay importaciones circulares.
3. Solo el ciclo de juego modifica el estado. Las sesiones envían acciones por canales; nunca tocan el estado directamente.
4. `graficos` se usa mediante una interfaz, para poder cambiar de vista sin tocar `juego`.
5. El servidor valida toda entrada que llega de un jugador antes de aplicarla.
6. Direcciones y puertos no se escriben fijos en el código. Se leen de argumentos o de variables de entorno.

## Nombres

- Carpetas, paquetes, tipos, funciones y comentarios en español.
- Sin acentos ni ñ en nombres de carpetas, paquetes e identificadores.
- Cada nombre describe exactamente su contenido. Si es la parte de gráficos, se llama `graficos`.

## Pruebas y verificación

- Todo cambio en `internal/juego` incluye pruebas.
- Antes de dar una tarea por terminada, estos tres comandos deben terminar sin errores:

```
go build ./...
go vet ./...
go test ./...
```

- Para ejecutar el servidor en local: `go run ./cmd/servidor`

## Git

- Nunca usar `git add .` ni `git add -A`. Los archivos se agregan por ruta explícita.
- Nunca subir `.DS_Store`, binarios compilados ni archivos de editor.
- No hacer commit ni push sin que se pida.
- Cada tarea nueva va en su propia rama. No se sube directo a `main`, salvo la estructura inicial.

## Forma de trabajo

- Antes de escribir código, mostrar el plan y esperar aprobación.
- Al terminar, explicar los cambios por bloques: qué hace cada parte y por qué se hizo así.
- Si una tarea obliga a romper una regla de este archivo, detenerse y preguntar.
- No modificar este archivo sin que se pida.

## Orden de trabajo

1. Estructura inicial: paquetes, tipos e interfaces, sin lógica.
2. `juego`: mapa grande y movimiento con colisión, con pruebas.
3. `graficos/plano` con cámara, `red` y `ciclo`: varios jugadores moviéndose por ssh.
4. Despliegue en un VPS con Ubuntu.
5. Combate: disparo, vida y eliminación.
6. Fases de la partida y ganador.
7. Zona que se cierra.
8. Botín e inventario.
9. Bots.

Pendiente de decidir: si el juego tendrá construcción.
