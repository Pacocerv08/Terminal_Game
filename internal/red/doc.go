// Paquete red es el servidor SSH y el manejo de sesiones.
//
// Cada sesión se une al ciclo, traduce las teclas del jugador en acciones, las
// valida y las envía al ciclo como Mensaje; nunca toca el estado. Por cada
// Cuadro que recibe, dibuja con un graficos.Dibujante y escribe en la
// terminal del jugador.
//
// La dirección y la clave del servidor no se escriben fijas: se leen de
// variables de entorno o de argumentos (ver Config). El servidor SSH usará
// Wish y la interfaz de terminal usará Bubble Tea; se agregan al módulo cuando
// se implemente este paquete.
package red
