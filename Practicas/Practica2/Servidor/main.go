package main

import "fmt"

func startLista() {
	crearProducto("Audifonos", "Audifonos beats", 100, 10, 0)
	crearProducto("Laptop Apple", "Mack book air", 1000, 2, 10)
	crearProducto("Calculadora", "CASIO", 30, 12, 8)
}
func main() {
	startLista()
	for _, p := range ListaProductos {
		fmt.Printf("ID:%d\nNombre: %s\nDescripcion:%s\nPrecio:%f\nCantidad:%d\nPromocion:%d\n", p.ID, p.Nombre, p.Descripcion, p.Precio, p.Existencias, p.Promocion)
		fmt.Println("========")
	}
}
