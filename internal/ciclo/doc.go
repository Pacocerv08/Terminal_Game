// Paquete ciclo ejecuta el ciclo de ticks a intervalo fijo. Es el único que
// modifica el Estado del juego.
//
// Las sesiones nunca tocan el estado: hablan con el ciclo por un único canal
// de Mensaje (unirse, actuar, salir), que conserva el orden de llegada. El
// ciclo les devuelve un Cuadro por tick por el canal de su Suscripcion.
//
// # Una acción por jugador por tick
//
// El ciclo aplica como máximo una acción por jugador en cada tick. Si en un
// mismo tick llegan varias acciones del mismo jugador, conserva solo la
// última y descarta las anteriores. Si el jugador sale en ese tick, su acción
// pendiente se descarta.
//
// # Orden de un tick
//
//  1. Mensajes: se vacía el canal de entrada sin bloquear. Unirse y Salir se
//     atienden en orden de llegada; de cada jugador se recuerda su última
//     acción.
//  2. Acciones de humanos: se aplican las acciones recordadas, solo si
//     partida.Maquina.AccionPermitida las permite en la fase actual.
//  3. Bots: cada bot decide con su Vista y su acción se aplica con su ID
//     como autor, con el mismo filtro de fase. También cuenta como una
//     acción por tick.
//  4. Estado.Avanzar: zona y daño de zona.
//  5. partida.Maquina.Avanzar: transiciones de fase, con el estado ya
//     actualizado.
//  6. Publicación: se envía un Cuadro a cada suscripción sin bloquear. Si el
//     buffer tiene un cuadro viejo sin leer, se reemplaza por el nuevo, así
//     una sesión lenta no frena el juego.
package ciclo
