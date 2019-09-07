import socket
host = "localhost"
port = 8080
with socket.socket(socket.AF_INET,socket.SOCK_STREAM) as s:
    s.connect((host,port))
    while True:
        m = input(">>")
        if m == "!close":
            exit(code=0)
        s.sendall(m.encode())
        data = s.recv(1024)
        print(data.decode())
