import socket
import json
class Connection():
    def __init__(self,host,port):
        self.HOST = host
        self.PORT = port
        self.s = socket.socket(socket.AF_INET,socket.SOCK_STREAM)
    def recvall(self,sock):
        BUFF_SIZE = 1024 # 4 KiB
        data = b''
        while True:
            part = sock.recv(BUFF_SIZE)
            data += part
            if len(part) < BUFF_SIZE:
                # either 0 or end of data
                break
        return json.loads(data.decode())
    def start(self):
        try:
            self.s.connect((self.HOST,self.PORT))
            return True
        except Exception as e:
            print(e)
            return False
    def sendMessage(self,m):
        self.s.sendall(m.encode())
        return self.recvall(self.s)
