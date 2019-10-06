package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/gabriel-vasile/mimetype"
)

const BUFFERSIZE = 9000
const HOMEFOLDER = "."

var splitregex = regexp.MustCompile(HOMEFOLDER)
var splitargs = regexp.MustCompile(`\?`)

func GET(conn net.Conn, request string) {
	f := HOMEFOLDER + rpath.FindString(request)
	data := splitargs.Split(f, 2)
	if len(data) > 1 {
		f = data[0]
		args := data[1]
		fmt.Println("ARGS", args)
	}
	f, _ = url.QueryUnescape(f)
	fmt.Println(f)
	fi, err := os.Stat(f) // Vemos si existe el archivo/directorio

	if err != nil {
		fmt.Println(err)
		Response404(conn)
		return
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		SendDirectory(f, conn)
	case mode.IsRegular():
		SendFile(f, conn)
	}
}
func POST(conn net.Conn, request string) {
	mensaje := "METODO POST!"
	header := makeheader(len(mensaje), "200", "text/html; charset=utf-8")
	Write(conn, []byte(header))
	Write(conn, []byte(mensaje))
}
func makeheader(strlen int, code, ctype string) string {
	dt := time.Now()
	l := strconv.Itoa(strlen)
	header := "HTTP/1.0 " + code + "\nServer: GOLANG Server/1.0 \n" + "Date: " + dt.String() + " \n" + "Content-Type: " + ctype + " \n"
	header += "Content-Length: " + l + "\n\n"
	fmt.Println()
	return header
}
func Response404(conn net.Conn) {
	data, _ := ioutil.ReadFile("404.html")
	header := makeheader(len(data), "404 Not Found", "text/html; charset=utf-8")
	Write(conn, []byte(header))
	Write(conn, data)
}
func Response202(data []byte, conn net.Conn, ctype string) {
	fmt.Println("LEN DATA:", len(data))
	header := makeheader(len(data), "202 ok", ctype)
	Write(conn, []byte(header))
	Write(conn, data)
}
func SendFile(filepath string, conn net.Conn) {
	data, _ := ioutil.ReadFile(filepath)
	ctype, _, _ := mimetype.DetectFile(filepath)
	fmt.Println("TIPO: ", ctype)
	Response202(data, conn, ctype)

}
func SendDirectory(path string, conn net.Conn) {
	data, err := ioutil.ReadFile(path + "/index.html")
	ctype := "text/html; charset=utf-8"
	if err != nil || len(data) == 0 {
		//data, _ = ioutil.ReadFile("index.html")
		data = getPathInfo(path)
	}
	Response202(data, conn, ctype)
}
func Write(conn net.Conn, message []byte) {
	l := len(message) - 1
	for i := 0; ; {
		if i+BUFFERSIZE >= l {
			conn.Write(message[i : l+1])
			break
		} else {
			conn.Write(message[i : i+BUFFERSIZE])
			i += BUFFERSIZE
		}
	}
}
