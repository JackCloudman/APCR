import socket,json
host = "localhost"
port = 8080
buffersize = 1024
def download(s,name,filesize):
    f = open(name, 'wb')
    myfile = s.recv(1024) #Descargamos el archivo
    totalRecv = len(myfile)
    while True:
        if totalRecv>=filesize:
            f.write(myfile[:-(totalRecv-filesize)])
            f.close()
            break
        f.write(myfile)
        myfile = s.recv(1024)
        totalRecv += len(myfile)
    print ('Descarga completa')
    f.close()
def main():
    with socket.socket(socket.AF_INET,socket.SOCK_STREAM) as s:
        s.connect((host,port))
        while True:
            m = input("Archivo a descargar: ")
            s.sendall(json.dumps({"command":"descargar","path":".","nombres":[m]}).encode()) # Peticion de descarga
            data = s.recv(1024)
            data = json.loads(data.decode(errors='ignore')) # Leemos la respuesta
            name = data["nombres"][0]
            filesize = data["sizes"][0]
            download(s,name,filesize)

if __name__ == "__main__":
    main()
