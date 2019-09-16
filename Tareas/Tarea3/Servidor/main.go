package main

import (
	"encoding/json"
	"fmt"
	"net"
)

var currentGame Ahorcado

func processCommand(data Message, conn net.Conn) {
	if data.Command == "generarJuego" {
		var b Ahorcado
		fmt.Println(data.Dificultad)
		switch data.Dificultad {
		case 0:
			b = crearFacil()
		case 1:
			b = crearMedio()
		case 2:
			b = crearDificil()
		default:
			b = crearFacil()
		}
		fmt.Println("Se ha generado la palabra:", b.Palabra)
		currentGame = b
		data.Tablero = b
		sendTablero(data, conn)
	} else if data.Command == "desbloquear" {
		b := intentar(&currentGame, data.Letra)
		data.Tablero = *b
		sendTablero(data, conn)
	} else {
		sendError(conn, "comando no encontrado")
	}
}
func main() {
	iniciarAhorcado()
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
