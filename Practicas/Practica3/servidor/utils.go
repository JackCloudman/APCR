package main

import (
	"encoding/json"
	"net"
	"time"
)

const BUFFERSIZE = 1024

// Estructura básica para el canal de comunicación entre el cliente (Python) y el servidor (GO)
type Message struct {
	User   string   `json:"user"`
	Text   string   `json:"text"`
	Action string   `json:"action"`
	TypeM  string   `json:"type_message"`
	Users  []string `json:"users"`
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
func sendError(conn net.Conn, message string) {
	m := Message{}
	res, _ := json.Marshal(m)
	Write(conn, res)
}
