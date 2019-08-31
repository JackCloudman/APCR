package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

// Estructura básica para el canal de comunicación entre el cliente (Python) y el servidor (GO)
type Message struct {
	Command   string      `json:"command"`
	Articulos []*Articulo `json:"articulos"`
	Status    int         `json:"status"`
	Ticket    Ticket      `json:"ticket"`
}

func sendError(conn net.Conn, message string) {
	m := Message{}
	m.Status = 400
	m.Command = message
	res, _ := json.Marshal(m)
	conn.Write(res)
}

func ImageTobase64(path string) []byte {
	// Open file on disk.
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)
	return content
}
func sendTicket(t Ticket, conn net.Conn) {
	m := Message{}
	m.Status = 200
	m.Command = "ok"
	m.Ticket = t
	res, _ := json.Marshal(m)
	conn.Write(res)
}

//
func sendCatalogo(data Message, conn net.Conn) {
	data.Articulos = ListaProductos
	data.Status = 200
	data.Command = "ok"
	res, _ := json.Marshal(data)
	fmt.Println(string(res))
	conn.Write(res)
}
