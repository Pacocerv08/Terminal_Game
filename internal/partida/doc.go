// Paquete partida define las fases de la partida (sala de espera, en curso y
// final) como una máquina de estados.
//
// La máquina decide qué acciones se permiten en cada fase y cuándo se pasa de
// una fase a otra. Solo el ciclo de juego la invoca, después de avanzar el
// estado en cada tick.
package partida
