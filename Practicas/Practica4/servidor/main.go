package main

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
)

var usuarios []string
var rmethod = regexp.MustCompile(`(GET)|(POST)`)
var rpath = regexp.MustCompile(`[/?][\w.,@?^=%&:/~+#-]*`)

/*Revisa si la peticion es metodo POST o metodo GET*/
func handleConnection(conn net.Conn) {
	// Mensaje
	line, err := bufio.NewReader(conn).ReadString('\n')
	fmt.Println(line)

	if err != nil {
		return
	}
	method := rmethod.FindString(string(line)) // Buscamos el metodo
	if method == "POST" {
		POST(conn, string(line))
	} else if method == "GET" {
		GET(conn, string(line))
	}
	if err != nil {
		return
	}
	conn.Close()
}

/*Inicia el socket y lo ejecute en una gorutina*/
func main() {
	server, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}
	defer server.Close()
	fmt.Println("Servidor iniciado...")
	for {
		conn, err := server.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(conn)
	}
}
