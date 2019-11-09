package main

import (
	"fmt"
)

func main() {
	var ncubetas int
	fmt.Println("Numero de cubetas: ")
	_, err := fmt.Scan(&ncubetas)
	if err != nil {
		return
	}
	cubeta := genRandom(4000, 1000)
	cubetas := splitCubetas(cubeta, ncubetas, 1000)
	ch := make(chan Message)
	port := 8000
	// Iniciar servidor y cliente
	for i := 0; i < ncubetas; i++ {
		go startServer(port, ch)
		go startClient(port, i, cubetas[i])
		port++
	}
	// Esperar la respuesta de los servidores
	for i := 0; i < ncubetas; i++ {
		m := <-ch
		fmt.Println("Cubera recibida: ", m.ID)
		cubetas[m.ID] = m.Numeros
	}
	//Combinar las respuestas
	var resultado []int
	for i := 0; i < ncubetas; i++ {
		resultado = append(resultado, cubetas[i]...)
	}
	fmt.Println(resultado)

}
