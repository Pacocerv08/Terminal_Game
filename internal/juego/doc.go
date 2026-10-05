// Paquete juego contiene el estado y las reglas del juego: mapa, jugadores,
// combate, botín y zona.
//
// Solo usa la biblioteca estándar de Go. No importa ningún otro paquete del
// proyecto ni bibliotecas de red o de terminal: todos los demás paquetes
// dependen de este, nunca al revés.
//
// Solo el ciclo de juego modifica un Estado. Quien necesite leerlo (por
// ejemplo para dibujar) recibe una Vista, que es una copia de solo lectura.
package juego
