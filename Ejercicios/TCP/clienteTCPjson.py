import socket
import json
host = "localhost"
port = 8080

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

def main():
    with socket.socket(socket.AF_INET,socket.SOCK_STREAM) as s:
        s.connect((host,port))
        while True:
            m = input(">>")
            if m == "!close":
                exit(code=0)
            s.sendall(m.encode())
            data = recvall(s)
            print(data)

if __name__ == '__main__':
    main()
