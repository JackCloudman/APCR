package main

import (
	"math/rand"
	"time"
)

func genRandom(tam, nrange int) []int {
	numeros := make([]int, tam)
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < tam; i++ {
		numeros[i] = rand.Intn(nrange)
	}
	return numeros
}
func splitCubetas(cubeta []int, ncubetas, nrange int) [][]int {
	var cubetas [][]int
	cubetas = make([][]int, ncubetas)
	r := nrange / ncubetas
	for _, numero := range cubeta {
		n := numero / r
		if n > ncubetas-1 {
			cubetas[ncubetas-1] = append(cubetas[ncubetas-1], numero)
		} else {
			cubetas[n] = append(cubetas[n], numero)
		}

	}
	return cubetas
}
