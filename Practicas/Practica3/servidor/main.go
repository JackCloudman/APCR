package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/dmichael/go-multicast/multicast"
)

var usuarios []string

func handleConnection(conn, connmulti net.Conn) {
	fmt.Printf("Serving %s\n", conn.RemoteAddr().String())
	// Primer mensaje para enviar la lista de usuarios
	d := json.NewDecoder(conn)
	var data Message
	err := d.Decode(&data)
	if err != nil {
		return
	}
	usuarios = append(usuarios, data.User)
	data.Users = usuarios
	sendMessage(conn, data)
	sendMessage(connmulti, data)
	///
	for {
		d := json.NewDecoder(conn)
		var data Message
		err := d.Decode(&data)
		if err != nil {
			break
		}
		fmt.Println(data)
		sendMessage(connmulti, data)
	}
	conn.Close()
}

/*Inicia el socket y se queda a la escucha, lee json y usa la funcion
processCommand para ejecutar el comando que quiere el cliente.*/
func main() {
	my_socket, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer my_socket.Close()
	smulti := start()
	fmt.Println("Servidor iniciado...")
	for {
		conn, err := my_socket.Accept()
		fmt.Println("Nueva conección!")
		if err != nil {
			panic(err)
		}
		go handleConnection(conn, smulti)
	}
}

// Ejemplo paa servidor multicast

const (
	defaultMulticastAddress = "224.1.1.1:5007"
)

func start() net.Conn {
	conn, err := multicast.NewBroadcaster(defaultMulticastAddress)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

func shareMessage(conn net.Conn, m Message) {

}
