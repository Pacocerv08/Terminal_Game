package juego

// TipoObjeto es la clase de un objeto del botín.
type TipoObjeto int

const (
	ObjetoPistola TipoObjeto = iota
	ObjetoMunicion
	ObjetoBotiquin
)

// Objeto es una cantidad de objetos de un mismo tipo.
type Objeto struct {
	Tipo     TipoObjeto
	Cantidad int
}

// ObjetoEnSuelo es un objeto del botín que está en una posición del mapa.
type ObjetoEnSuelo struct {
	Objeto
	Pos Posicion
}

// Inventario son los objetos que lleva un jugador, hasta Capacidad.
type Inventario struct {
	Capacidad int
	objetos   []Objeto
}

// Agregar guarda o. Devuelve false si no cabe.
func (i *Inventario) Agregar(o Objeto) bool {
	// TODO: implementar.
	return false
}

// Quitar descarta cantidad unidades del tipo t. Devuelve false si no hay
// suficientes.
func (i *Inventario) Quitar(t TipoObjeto, cantidad int) bool {
	// TODO: implementar.
	return false
}

// Contiene indica si el inventario tiene al menos una unidad del tipo t.
func (i *Inventario) Contiene(t TipoObjeto) bool {
	// TODO: implementar.
	return false
}

// Copiar devuelve un inventario que no comparte memoria con el original.
func (i Inventario) Copiar() Inventario {
	copia := Inventario{Capacidad: i.Capacidad}
	copia.objetos = append([]Objeto(nil), i.objetos...)
	return copia
}

// RecogerAccion pide que el autor recoja el objeto que hay en su casilla.
type RecogerAccion struct{}

// Aplicar recoge el objeto si el autor existe y está vivo.
func (a RecogerAccion) Aplicar(e *Estado, autor IDJugador) {
	// TODO: buscar el objeto en Estado.Botin y pasarlo al inventario.
}
