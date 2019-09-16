package main

import (
	"math/rand"
	"strings"
)

type Ahorcado struct {
	Palabra        string `json:"-"`
	Status         string `json:"status"`
	Longitud       int    `json:"longitud"`
	Dificultad     int    `json:"dificultad"`
	Vidas          int    `json:"vidas"`
	PalabraTablero string `json:"palabra"`
}

var PalabraFacil = []string{}
var PalabraMedio = []string{}
var PalabraDificil = []string{}

func iniciarAhorcado() {
	dificil := []string{"VENTRICULO", "ESCLEROSIS", "MISANTROPIA", "ESCALOFRIANTE"}
	medio := []string{"POKEMON", "GOLANG", "PERRACOS", "HIGO"}
	facil := []string{"CASA", "PERRO", "PALABRA", "PYTHON"}
	PalabraFacil = facil
	PalabraMedio = medio
	PalabraDificil = dificil
}
func crearFacil() Ahorcado {
	a := Ahorcado{}
	a.Dificultad = 0
	a.Vidas = 7
	a.Palabra = PalabraFacil[rand.Intn(len(PalabraFacil))]
	a.Longitud = len(a.Palabra)
	a.PalabraTablero = strings.Repeat("_", a.Longitud)
	a.Status = "JUGANDO"
	return a
}
func crearMedio() Ahorcado {
	a := Ahorcado{}
	a.Dificultad = 1
	a.Vidas = 7
	a.Palabra = PalabraMedio[rand.Intn(len(PalabraMedio))]
	a.Longitud = len(a.Palabra)
	a.PalabraTablero = strings.Repeat("_", a.Longitud)
	a.Status = "JUGANDO"
	return a
}
func crearDificil() Ahorcado {
	a := Ahorcado{}
	a.Dificultad = 1
	a.Vidas = 7
	a.Palabra = PalabraDificil[rand.Intn(len(PalabraMedio))]
	a.Longitud = len(a.Palabra)
	a.PalabraTablero = strings.Repeat("_", a.Longitud)
	a.Status = "JUGANDO"
	return a
}
func intentar(a *Ahorcado, l string) *Ahorcado {
	if strings.ContainsAny(a.Palabra, l) {
		for i, letra := range a.Palabra {
			if string(letra) == l {
				a.PalabraTablero = a.PalabraTablero[:i] + l + a.PalabraTablero[i+1:]
			}
		}
		if !strings.ContainsAny(a.PalabraTablero, "_") {
			a.Status = "GANADO"
		}
	} else {
		a.Vidas--
		if a.Vidas == 0 {
			a.Status = "PERDIDO"
			a.PalabraTablero = a.Palabra
		}
	}
	return a
}
