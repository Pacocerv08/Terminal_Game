// Paquete ciclo ejecuta el ciclo de ticks a intervalo fijo. Es el único que
// modifica el Estado del juego.
//
// Las sesiones nunca tocan el estado: hablan con el ciclo por un único canal
// de Mensaje, que conserva el orden de llegada. Los mensajes solo se crean
// con Unirse, Actuar y Salir; el valor cero de Mensaje se ignora. El ciclo
// les devuelve un Cuadro por tick por el canal de su Suscripcion.
//
// # Unirse, actuar y salir
//
// Unirse envía el mensaje y espera la respuesta del ciclo. La respuesta no
// tiene buffer: el ciclo la entrega con un select contra la cancelación del
// contexto de Unirse, y si no se pudo entregar retira al jugador que acababa
// de agregar. Salir intenta el envío antes de mirar su contexto, que debe ser
// el del servidor y no el de la sesión.
//
// # Una acción por jugador por tick
//
// El ciclo aplica como máximo una acción por jugador en cada tick. Si en un
// mismo tick llegan varias acciones del mismo jugador, conserva solo la
// última y descarta las anteriores. Si el jugador sale en ese tick, su acción
// pendiente se descarta. Solo actúan los jugadores con suscripción.
//
// # Orden de un tick
//
// El tick ejecuta un solo paso, sin leer el reloj (ver Ciclo.tick):
//
//  1. Mensajes: se atienden los que había en la cola al empezar el tick, sin
//     bloquear y en orden de llegada. Unirse y Salir se resuelven en el acto;
//     de cada jugador se recuerda su última acción.
//  2. Bots: cada bot decide con la Vista del estado de inicio de tick, y su
//     acción se suma a las recordadas.
//  3. Aplicación: todas las acciones se aplican en el orden de
//     ordenDeAplicacion, y solo las que partida.Maquina.AccionPermitida
//     permite en la fase actual. Humanos y bots van mezclados, así que nadie
//     gana siempre un desempate.
//  4. Estado.Avanzar: cuenta el tick, zona y daño de zona.
//  5. partida.Maquina.Avanzar: transiciones de fase, con el estado ya
//     actualizado.
//  6. Publicación: se envía un Cuadro a cada suscripción sin bloquear. Si el
//     buffer tiene un cuadro viejo sin leer, se reemplaza por el nuevo, así
//     una sesión lenta no frena el juego.
//
// # Orden de aplicación y desempates
//
// Si dos jugadores quieren lo mismo en un tick (la misma casilla, por
// ejemplo), gana quien se aplica primero. El orden sale de ordenar los IDs de
// menor a mayor y barajarlos con un generador sembrado con (semilla, tick):
// es determinista, no depende de ningún map y no favorece siempre a los
// mismos jugadores. La semilla se fija con Ciclo.FijarSemilla.
//
// # Cierre
//
// Ejecutar termina cuando se cancela su contexto y cierra todas las
// suscripciones. Un Unirse que quedó en la cola sin atender solo termina
// cuando se cancela su propio contexto, así que ese contexto debe derivarse
// del contexto del servidor.
package ciclo
