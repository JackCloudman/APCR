package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/mholt/archiver"
)

const BUFFERSIZE = 1024
const PATH = "temp/"

type Message struct {
	Tipo    string   `json:"tipo"`
	Nombres []string `json:"nombres"`
	Sizes   []int64  `json:"sizes"`
	Numero  int      `json:"numero"`
	Command string   `json:"command"`
	Path    string   `json:"path"`
}

type PathInfo struct {
	Path    string   `json:"path"`
	Nombres []string `json:"nombres"`
	Tipo    []int    `json:"tipo"`
}

func recibir(fileName, ext string, tam int64, conn net.Conn) {
	newFile, err := os.Create(PATH + fileName + ext)
	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	var receivedBytes int64
	receivedBytes = 0
	for {
		if (tam - receivedBytes) < BUFFERSIZE {
			io.CopyN(newFile, conn, (tam - receivedBytes))
			break
		}
		io.CopyN(newFile, conn, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
	}
	fmt.Println("Archivo guardado!")
}
func uncompress(filename string) {
	err := archiver.Unarchive(PATH+filename+".zip", PATH+filename)
	if err != nil {
		panic(err)
	}
	println("Descompresion completa!")
}
func cleanFiles(filename string) {
	_ = os.Remove(PATH + filename)
	fmt.Println("Archivos temporales eliminados")
}

func upload(data Message, conn net.Conn) {
	if data.Tipo == "folder" {
		fmt.Println("Recibiendo archivo...")
		recibir(data.Nombres[0], ".zip", data.Sizes[0], conn)
		fmt.Println("Archivo recibido! Descomprimiendo...")
		uncompress(data.Nombres[0])
		fmt.Println("Listo!")
		cleanFiles(data.Nombres[0] + ".zip")
	} else if data.Tipo == "list" {
		for i := 0; i < data.Numero; i++ {
			recibir(data.Nombres[i], "", data.Sizes[i], conn)
		}
		fmt.Println("Listo!")
	}
}
func getPathInfo(path string) []byte {
	pi := &PathInfo{}
	files, err := ioutil.ReadDir(PATH + path)
	if err != nil {
		log.Fatal(err)
	}
	pi.Path = path
	var t int
	for _, f := range files {
		pi.Nombres = append(pi.Nombres, f.Name())
		if f.IsDir() {
			t = 1
		} else {
			t = 0
		}
		pi.Tipo = append(pi.Tipo, t)
		fmt.Println(f.Name(), f.IsDir())
	}
	res, _ := json.Marshal(pi)
	println(string(res))
	return res
}
func main() {
	my_socket, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer my_socket.Close()
	fmt.Println("Servidor iniciado...")
	for {
		conn, err := my_socket.Accept()
		fmt.Println("Nueva conección!")
		if err != nil {
			panic(err)
		}
		for {
			d := json.NewDecoder(conn)
			var data Message
			err := d.Decode(&data)
			if err != nil {
				break
			}
			fmt.Println(data)
			if data.Command == "descargar" {
				//download(data, conn)
			} else if data.Command == "subir" {
				fmt.Println("Subir!")
				upload(data, conn)
			} else if data.Command == "getPathInfo" {
				s := getPathInfo(data.Path)
				conn.Write(s)
			}

		}
		conn.Close()
	}
}
