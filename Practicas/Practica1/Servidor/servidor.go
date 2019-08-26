package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

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
	Sizes   []int64  `json:"sizes"`
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
func compress(filename string) int64 {
	err := archiver.Archive([]string{PATH + filename}, PATH+filename+".zip")
	if err != nil {
		panic(err)
	}
	println("Compresion completa!")
	fileStat, _ := os.Stat(PATH + filename + ".zip")
	return fileStat.Size()
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
		pi.Sizes = append(pi.Sizes, f.Size())
		fmt.Println(f.Name(), f.IsDir(), f.Size())
	}
	res, _ := json.Marshal(pi)
	println(string(res))
	return res
}
func download(data Message, conn net.Conn) {
	fi := &PathInfo{}
	fileStat, err := os.Stat(PATH + data.Nombres[0])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("File Name:", fileStat.Name())        // Base name of the file
	fmt.Println("Size:", fileStat.Size())             // Length in bytes for regular files
	fmt.Println("Permissions:", fileStat.Mode())      // File mode bits
	fmt.Println("Last Modified:", fileStat.ModTime()) // Last modification time
	fmt.Println("Is Directory: ", fileStat.IsDir())   // Abbreviation for Mode().IsDir()
	if fileStat.IsDir() {
		fi.Tipo = []int{1}
		fi.Sizes = []int64{compress(data.Nombres[0])}
		fi.Nombres = []string{fileStat.Name() + ".zip"}
	} else {
		fi.Tipo = []int{0}
		fi.Sizes = []int64{fileStat.Size()}
		fi.Nombres = []string{fileStat.Name()}
	}
	fi.Path = "."
	res, _ := json.Marshal(fi)
	fmt.Println(string(res))
	conn.Write(res) //Enviamos el .json con la informacion del Archivo
	time.Sleep(100 * time.Millisecond)
	file, err := os.Open(PATH + fi.Nombres[0])
	if err != nil {
		log.Fatal(err)
	}
	sendBuffer := make([]byte, BUFFERSIZE)
	fmt.Println("Start sending file!")
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		conn.Write(sendBuffer)
	}
	fmt.Println("Archivo enviado!")
	if fi.Tipo[0] == 1 {
		cleanFiles(fi.Nombres[0])
	}

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
				download(data, conn)
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
