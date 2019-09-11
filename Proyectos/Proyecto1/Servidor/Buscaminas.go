package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Buscaminas struct {
	Matriz [][]int `json:"matriz"`
	F      int     `json:"filas"`
	C      int     `json:"columnas"`
	Minas  int     `json:"minas"`
}

func createBuscaminas(f, c, minas int) Buscaminas {
	b := Buscaminas{}
	b.Matriz = make([][]int, f)
	b.C = c
	b.F = f
	b.Minas = minas
	// Creamos la matriz
	for i := 0; i < f; i++ {
		b.Matriz[i] = make([]int, c)
	}
	// insertamos las bombas
	for i := 0; i < minas; i++ {
		insertarBomba(&b)
	}
	for i := 0; i < f; i++ {
		for j := 0; j < c; j++ {
			if b.Matriz[i][j] == -1 {
				actualizarValores(i, j, &b)
			}
		}
	}
	return b
}
func actualizarValores(rn, c int, b *Buscaminas) {
	// #Row above.
	if rn-1 > -1 {
		r := b.Matriz[rn-1]

		if c-1 > -1 {
			if r[c-1] != -1 {
				r[c-1]++
			}
		}

		if r[c] != -1 {
			r[c]++
		}

		if b.C > c+1 {
			if r[c+1] != -1 {
				r[c+1]++
			}

		}
	}
	//#Same row.
	r := b.Matriz[rn]

	if c-1 > -1 {
		if r[c-1] != -1 {
			r[c-1]++
		}
	}

	if b.C > c+1 {
		if r[c+1] != -1 {
			r[c+1]++
		}
	}
	//#Row below.
	if b.F > rn+1 {
		r := b.Matriz[rn+1]

		if c-1 > -1 {
			if r[c-1] != -1 {
				r[c-1]++
			}
		}

		if r[c] != -1 {
			r[c]++
		}

		if b.C > c+1 {
			if r[c+1] != -1 {
				r[c+1]++
			}

		}
	}
}
func insertarBomba(b *Buscaminas) {
	rand.Seed(time.Now().UTC().UnixNano())
	f := rand.Intn(b.F)
	c := rand.Intn(b.C)
	if b.Matriz[f][c] == 0 {
		b.Matriz[f][c] = -1
	} else {
		insertarBomba(b)
	}
}
func crearFacil() Buscaminas {
	return createBuscaminas(9, 9, 10)
}
func crearMedio() Buscaminas {
	return createBuscaminas(16, 16, 40)
}
func crearDificil() Buscaminas {
	return createBuscaminas(16, 30, 99)
}

func printBuscaminas(b Buscaminas) {
	for i := 0; i < b.F; i++ {
		fmt.Println(b.Matriz[i])
	}
}
