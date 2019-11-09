package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sort"
	"strconv"
	"time"
)

const BUFFERSIZE = 1024

// Estructura básica para el canal de comunicación entre el cliente (Python) y el servidor (GO)
type Message struct {
	ID      int   `json:"id"`
	Numeros []int `json:"numeros"`
}

func Write(conn net.Conn, message []byte) {
	l := len(message) - 1
	for i := 0; ; {
		if i+BUFFERSIZE >= l {
			conn.Write(message[i : l+1])
			break
		} else {
			conn.Write(message[i : i+BUFFERSIZE])
			i += BUFFERSIZE
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func sendMessage(conn net.Conn, m Message) {
	res, _ := json.Marshal(m)
	Write(conn, res)
}
func startServer(port int, c chan Message) {
	server, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		panic(err)
	}
	defer server.Close()
	fmt.Println("Servidor iniciado en el puerto..", port)
	conn, err := server.Accept()
	if err != nil {
		panic(err)
	}
	var data Message
	d := json.NewDecoder(conn)
	d.Decode(&data)
	c <- data
	conn.Close()
	fmt.Println("Servidor terminado")
}
func startClient(port, ncubeta int, cubeta []int) {
	fmt.Println("Cliente iniciado")
	m := Message{}
	m.ID = ncubeta
	sort.Ints(cubeta)
	m.Numeros = cubeta
	conn, _ := net.Dial("tcp", "localhost:"+strconv.Itoa(port))
	time.Sleep(time.Second)
	fmt.Println("Enviar respuesta a servidor")
	sendMessage(conn, m)
}
