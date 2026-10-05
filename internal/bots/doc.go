// Paquete bots contiene jugadores controlados por el servidor.
//
// Un bot solo recibe una Vista (copia de solo lectura) y devuelve una Accion.
// El ciclo la aplica con el ID del bot como autor, igual que a un humano.
// Los bots corren dentro del ciclo, de forma síncrona, por eso Decidir debe
// ser barato.
package bots
