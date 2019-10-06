package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"regexp"
	"time"
)

var usuarios []string
var rmethod = regexp.MustCompile(`(GET)|(POST)`)
var rpath = regexp.MustCompile(`[/?][\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-]`)

func GET(conn net.Conn, request string) {
	data, _ := ioutil.ReadFile("index.html")
	header := makeheader(len(data))
	conn.Write([]byte(header))
	conn.Write(data)
}
func POST(conn net.Conn, request string) {
	mensaje := "METODO POST!"
	header := makeheader(len(mensaje))
	conn.Write([]byte(header))
	conn.Write([]byte(mensaje))
}
func makeheader(strlen int) string {
	dt := time.Now()
	header := "HTTP/1.0 200 ok\n" + "Server: Axel Server/1.0 \n" + "Date: " + dt.String() + " \n" + "Content-Type: text/html \n"
	header += "Content-Length: " + string(strlen) + "\n\n"
	return header
}
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

/*Inicia el socket y se queda a la escucha, lee json y usa la funcion
processCommand para ejecutar el comando que quiere el cliente.*/
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
