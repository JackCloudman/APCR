package main
import (
  "fmt"
  "encoding/json"
  "net"
  "os"
  "io"
  "github.com/mholt/archiver"
)

const BUFFERSIZE = 1024
const PATH = "temp/"
type Message struct {
  Tipo   string `json:"tipo"`
  Nombres []string `json:"nombres"`
  Sizes []int64 `json:"sizes"`
  Numero int `json:"numero"`
}

func recibir(fileName,ext string,tam int64,conn net.Conn) {
  newFile, err := os.Create(fileName+ext)
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
func uncompress(filename string){
  err := archiver.Unarchive(filename+".zip", PATH+filename)
   if err != nil {
 		panic(err)
 	}
  println("Descompresion completa!")
}
func cleanFiles(filename string){
  _ = os.Remove(filename)
  fmt.Println("Archivos temporales eliminados")
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
      if data.Tipo == "folder"{
        fmt.Println("Recibiendo archivo...")
        recibir(data.Nombres[0],".zip", data.Sizes[0], conn)
        fmt.Println("Archivo recibido! Descomprimiendo...")
        uncompress(data.Nombres[0])
        fmt.Println("Listo!")
        cleanFiles(data.Nombres[0]+".zip")
      }else if data.Tipo == "list"{
        for i := 0; i < data.Numero; i++ {
          recibir(data.Nombres[i], "", data.Sizes[i], conn)
        }
        fmt.Println("Listo!")
      }

		}
		conn.Close()
	}
}
