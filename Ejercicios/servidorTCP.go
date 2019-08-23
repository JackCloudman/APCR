package main

import (
	"fmt"
	"net"
)

func main() {
	my_socket, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer my_socket.Close()
	for {
		conn, err := my_socket.Accept()
		if err != nil {
			panic(err)
		}
		for {
			bs := make([]byte, 1024)
			n, err := conn.Read(bs)
			fmt.Println(string(bs))
			if err != nil {
				break
			}
			_, err = conn.Write(bs[:n])
			if err != nil {
				break
			}
		}
		conn.Close()
	}
}
