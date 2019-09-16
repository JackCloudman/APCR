import socket,json,os
host = "localhost"
port = 8080
estados = [
'''
 ___
 |
 |
 |
_|_''',
'''
 ___
 |  o
 |
 |
_|_
''',
'''
 ___
 |  o
 |  |
 |
_|_
''',
'''
 ___
 |  o
 | -|
 |
_|_
''',
'''
 ___
 |  o
 | -|-
 |
_|_
''',
'''
 ___
 |  o
 | -|-
 | (
_|_
''',
'''
 ___
 |  o
 | -|-
 | ( )
_|_
'''
]
def recvall(sock):
    BUFF_SIZE = 1024 # 4 KiB
    data = b''
    while True:
        part = sock.recv(BUFF_SIZE)
        data += part
        if len(part) < BUFF_SIZE:
            # either 0 or end of data
            break
    return json.loads(data.decode())

def main():
    with socket.socket(socket.AF_INET,socket.SOCK_STREAM) as s:
        s.connect((host,port))
        while True:
            m = input("Selecciona una dificultad(0: facil, 1:medio, 2:dificil): ")
            if m == "!close":
                exit(code=0)
            m = '{"command":"generarJuego","dificultad":%d}'%int(m)
            s.sendall(m.encode())
            data = recvall(s)
            while data["ahorcado"]["status"] == "JUGANDO":
                a = data["ahorcado"]
                printTablero(a["palabra"],a["vidas"])
                m = input("Dame una letra: ")
                m = '{"command":"desbloquear","letra":"%s"}'%m.upper()
                s.sendall(m.encode())
                data = recvall(s)
            if data["ahorcado"]["status"] == "PERDIDO":
                os.system("clear")
                print("HAS PERDIDOOO! :(\n La palabra era: "+data["ahorcado"]["palabra"])
            else:
                os.system("clear")
                print("HAS GANADOOO! :D  La palabra era: "+data["ahorcado"]["palabra"])
            play = input("¿Volver a jugar? y/n")
            if play == "n":exit(code=0)

def printTablero(palabra,vidas):
    os.system("clear")
    print(estados[7-vidas])
    print(palabra)

if __name__ == '__main__':
    main()
