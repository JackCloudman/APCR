#Con esta funcion se obtiene todo el json
def recvall(sock):
    BUFF_SIZE = 4096 # 4 KiB
    data = b''
    while True:
        part = sock.recv(BUFF_SIZE)
        data += part
        if len(part) < BUFF_SIZE:
            # either 0 or end of data
            break
    return json.loads(data.decode())
def getCatalogo(socket):
    command = '{"command":"getCatalogo"}'
    s.sendall(command.encode())
    return recvall(s)

def ComprarEjemplo(socket):
    # el comando para comprar es "comprar", debe incluir los articulos a comprar
    #poniendo sus id, existencias representa la cantidad de articulos que quieres
    command = '{"command":"comprar","articulos":[{"id":1,"existencias":2},{"id":3,"existencias":4}]}'
    s.sendall(command.encode())
    ticket = recvall(s)
