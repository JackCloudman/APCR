import socket
import json
from buscaminas import Buscaminas
host = "localhost"
port = 8080
class Juego():
    def __init__(self,HOST,PORT):
        self.HOST = HOST
        self.PORT = PORT
        self.sock = socket.socket(socket.AF_INET,socket.SOCK_STREAM)
        self.sock.connect((HOST,PORT))
        print("Conectado al servidor!")
    def play(self):
        dificultad = int(input("Dame un dificultad (facil:0,medio:1,dificultad:2)"))
        b,f,c,minas = self.getTablero(dificultad)
        for i in range(f):
            for j in range(c):
                if b[i][j] == -1:
                    b[i][j] = "*"
        self.buscaminas = Buscaminas(b,f,c,minas)
        self.buscaminas.start()


    def recvall(self):
        sock = self.sock
        BUFF_SIZE = 1024 # 4 KiB
        data = b''
        while True:
            part = sock.recv(BUFF_SIZE)
            data += part
            if len(part) < BUFF_SIZE:
                # either 0 or end of data
                break
        return json.loads(data.decode())
    def getTablero(self,dificultad):
        s = self.sock
        m = {"command":"generarJuego","dificultad":dificultad}
        m = json.dumps(m)
        s.sendall(m.encode())
        data = self.recvall()["buscaminas"]
        return data["matriz"],data["filas"],data["columnas"],data["minas"]

def main():
    b = Juego(host,port)
    while True:
        b.play()
        op = input("Volver a jugar? y/n")
        if op =="y":
            pass
        else:
            exit(0)
if __name__ == '__main__':
    main()
