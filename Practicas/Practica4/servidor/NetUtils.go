package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
)

const BUFFERSIZE = 9000
const HOMEFOLDER = "."

var splitregex = regexp.MustCompile(HOMEFOLDER)
var splitargs = regexp.MustCompile(`\?`)
var paramregex = regexp.MustCompile(`.+$`)

func GET(conn net.Conn, request string) {
	f := HOMEFOLDER + rpath.FindString(request)
	data := splitargs.Split(f, 2)
	if len(data) > 1 {
		f = data[0]
		args := ParseParams(data[1])
		mensaje, _ := json.Marshal(args)
		fmt.Println("ARGS", string(mensaje))
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
	var header string
	aux := paramregex.FindAllString(request, -1)
	if len(aux) == 0 {
		aux = append(aux, "")
	}
	args := ParseParams(aux[len(aux)-1])
	mensaje, err := json.Marshal(args)
	if err != nil {
		mensaje = []byte("{}")
		header = makeheader(2, "200", "application/json")
	} else {
		header = makeheader(len(mensaje), "200", "application/json")
	}
	Write(conn, []byte(header))
	Write(conn, mensaje)
}
func PUT(conn net.Conn, request string) {
	aux := paramregex.FindAllString(request, -1)
	if len(aux) == 0 {
		Response400(conn)
		return
	}
	args := ParseParams(aux[len(aux)-1])
	// Campos necesarios para hacer el put
	if _, ok := args["type"]; !ok {
		Response400(conn)
		return
	}
	if _, ok := args["name"]; !ok {
		Response400(conn)
		return
	}
	if _, ok := args["value"]; !ok {
		Response400(conn)
		return
	}
	f := HOMEFOLDER + rpath.FindString(request)
	f, _ = url.QueryUnescape(f)
	_, err := os.Stat(f) // Vemos si existe el archivo/directorio

	if err != nil {
		fmt.Println(err)
		Response404(conn)
		return
	}
	switch args["type"] {
	case "folder":
		fmt.Println("Creando folder!")
		name, _ := url.QueryUnescape(fmt.Sprintf("%s%v", f, args["name"]))
		os.MkdirAll(name, os.ModePerm)
	case "file":
		fmt.Println("Creando archivo")
		name, _ := url.QueryUnescape(fmt.Sprintf("%s%v", f, args["name"]))
		data, _ := url.QueryUnescape(fmt.Sprintf("%v", args["value"]))
		err := ioutil.WriteFile(name, []byte(data), 0755)
		if err != nil {
			fmt.Printf("Unable to write file: %v", err)
		}
	}
	Response200([]byte("ok"), conn, "text/plain")
}
func DELETE(conn net.Conn, request string) {
	aux := paramregex.FindAllString(request, -1)
	if len(aux) == 0 {
		Response400(conn)
		return
	}
	args := ParseParams(aux[len(aux)-1])
	// Campos necesarios para hacer el put
	if _, ok := args["name"]; !ok {
		Response400(conn)
		return
	}
	f := HOMEFOLDER + rpath.FindString(request)
	f, _ = url.QueryUnescape(f)
	_, err := os.Stat(f) // Vemos si existe el archivo/directorio

	if err != nil {
		fmt.Println(err)
		Response200([]byte("ok"), conn, "text/plain")
		return
	}
	name, _ := url.QueryUnescape(fmt.Sprintf("%s%v", f, args["name"]))
	fmt.Println("Borrando: ", name)
	os.RemoveAll(name)
	Response200([]byte("ok"), conn, "text/plain")
}
func makeheader(strlen int, code, ctype string) string {
	dt := time.Now()
	l := strconv.Itoa(strlen)
	header := "HTTP/1.0 " + code + "\nServer: GOLANG Server/1.0 \n" + "Date: " + dt.String() + " \n" + "Content-Type: " + ctype + " \n"
	header += "Content-Length: " + l + "\n\n"
	fmt.Println()
	return header
}
func Response400(conn net.Conn) {
	header := makeheader(15, "400 Bad Request", "text/plain")
	Write(conn, []byte(header))
	Write(conn, []byte("400 Bad Request"))
}
func Response404(conn net.Conn) {
	data, _ := ioutil.ReadFile("404.html")
	header := makeheader(len(data), "404 Not Found", "text/html; charset=utf-8")
	Write(conn, []byte(header))
	Write(conn, data)
}
func Response200(data []byte, conn net.Conn, ctype string) {
	header := makeheader(len(data), "200 ok", ctype)
	Write(conn, []byte(header))
	Write(conn, data)
}
func Response405(conn net.Conn) {
	header := makeheader(18, "405", "text/plain")
	Write(conn, []byte(header))
	Write(conn, []byte("Method Not Allowed"))
}
func SendFile(filepath string, conn net.Conn) {
	data, _ := ioutil.ReadFile(filepath)
	ctype, _, _ := mimetype.DetectFile(filepath)
	fmt.Println("TIPO: ", ctype)
	Response200(data, conn, ctype)

}
func SendDirectory(path string, conn net.Conn) {
	data, err := ioutil.ReadFile(path + "/index.html")
	ctype := "text/html; charset=utf-8"
	if err != nil || len(data) == 0 {
		//data, _ = ioutil.ReadFile("index.html")
		data = getPathInfo(path)
	}
	Response200(data, conn, ctype)
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
func ParseParams(data string) map[string]interface{} {
	dict := map[string]interface{}{}
	if data == "" {
		return dict
	}
	for _, value := range strings.Split(data, "&") {
		mymap := strings.Split(value, "=")
		dict[mymap[0]] = mymap[1]
	}
	return dict
}
