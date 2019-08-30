package main

type Message struct {
	Command string `json:"command"`
}
type Articulo struct {
	ID          int     `json:"id"`
	Nombre      string  `json:"nombre"`
	Precio      float64 `json:"precio"`
	Existencias int     `json:"existencias"`
	Descripcion string  `json:"descripcion"`
	Promocion   int     `json:"promocion"`
}

var ListaProductos = []*Articulo{}

func crearProducto(nombre, descripcion string, precio float64, existencias, promocion int) {
	a := new(Articulo)
	a.Nombre = nombre
	a.Descripcion = descripcion
	a.Precio = precio
	a.Existencias = existencias
	a.Promocion = promocion
	a.ID = len(ListaProductos)
	ListaProductos = append(ListaProductos, a)
}
