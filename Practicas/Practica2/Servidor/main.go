package main

import (
	"encoding/json"
	"fmt"
	"net"
)

// Notense los nombres falsos por si acaso xD
func startLista() {
	crearProducto("Audifonos", "Audifonos bits", "image.jpg", 100, 10, 0)
	crearProducto("Laptop", "Hack book air", "image.jpg", 1000, 2, 10)
	crearProducto("Calculadora", "CASI", "image.jpg", 30, 12, 8)
	crearProducto("WebCam", "WebCam marca logtech", "image.jpg", 30.5, 12, 0)
	crearProducto("Teclado cool master", "Teclado K-202 cool master", "image.jpg", 120, 5, 0)
	crearProducto("Scooter", "Scooter marca Xiaoni", "image.jpg", 200.80, 4, 10)
	crearProducto("Mouse cool master", "Esto es una descripcion generica xD", "image.jpg", 10, 100, 5)
}
func processCommand(data Message, conn net.Conn) {
	if data.Command == "getCatalogo" {
		sendCatalogo(data, conn)
	} else if data.Command == "comprar" {
		status, t := comprarProductos(data.Articulos)
		if status != 200 {
			sendError(conn, "Error al realizar la compra!")
			return
		}
		sendTicket(t, conn)
	} else {
		sendError(conn, "Mensaje erroneo!")
	}
	return
}

/*Inicia el socket y se queda a la escucha, lee json y usa la funcion
processCommand para ejecutar el comando que quiere el cliente.*/
func main() {
	startLista()
	my_socket, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer my_socket.Close()
	fmt.Println("Servidor iniciado...")
	for {
		conn, err := my_socket.Accept()
		fmt.Println("Nueva conección!")
		if err != nil {
			panic(err)
		}
		for {
			d := json.NewDecoder(conn)
			var data Message
			err := d.Decode(&data)
			if err != nil {
				break
			}
			fmt.Println(data)
			processCommand(data, conn)
		}
		conn.Close()
	}
}
