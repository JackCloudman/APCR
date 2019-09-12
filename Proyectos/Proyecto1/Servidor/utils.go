package main

import (
	"encoding/json"
	"net"
	"time"
)

const BUFFERSIZE = 1024

// Estructura básica para el canal de comunicación entre el cliente (Python) y el servidor (GO)
type Message struct {
	Command    string     `json:"command"`
	Status     int        `json:"status"`
	Dificultad int        `json:"dificultad"`
	Tablero    Buscaminas `json:"buscaminas"`
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
func sendError(conn net.Conn, message string) {
	m := Message{}
	m.Status = 400
	m.Command = message
	res, _ := json.Marshal(m)
	Write(conn, res)
}
func sendTablero(m Message, conn net.Conn) {
	m.Status = 200
	m.Command = "ok"
	res, _ := json.Marshal(m)
	Write(conn, res)
}
