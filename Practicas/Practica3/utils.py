import socket
import json
import threading
import struct
import eel
class Connection():
    def __init__(self,MCAST_GRP,MCAST_PORT,NAME):
        self.MCAST_GRP = MCAST_GRP
        self.MCAST_PORT = MCAST_PORT
        self.IS_ALL_GROUPS = True
        self.MULTICAST_TTL = 2
        self.NAME = NAME

    def listener(self):
        sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM, socket.IPPROTO_UDP)
        sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        if self.IS_ALL_GROUPS:
            # on this port, receives ALL multicast groups
            sock.bind(('', self.MCAST_PORT))
        else:
            # on this port, listen ONLY to MCAST_GRP
            sock.bind((self.MCAST_GRP, self.MCAST_PORT))
        mreq = struct.pack("4sl", socket.inet_aton(self.MCAST_GRP), socket.INADDR_ANY)

        sock.setsockopt(socket.IPPROTO_IP, socket.IP_ADD_MEMBERSHIP, mreq)
        while True:
            data = self.recvall(sock)
            print(data)
            if data["user"] == self.NAME:
                pass
            else:
                self.processAction(data)
    def processAction(self,message):
        if message["action"] == "message":
            eel.recvMessage("<b>%s</b>: %s"%(message["user"],message["text"]))
        elif message["action"] == "join":
            eel.recvMessage("<b>Servidor</b>: El usuario %s se ha unido!"%(message["user"]))

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
            self.sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM, socket.IPPROTO_UDP)
            self.sock.setsockopt(socket.IPPROTO_IP, socket.IP_MULTICAST_TTL, self.MULTICAST_TTL)
            self.sendMessage("","join")
            print("iniciando listener...")
            t = threading.Thread(target=self.listener)
            t.start()
            print("iniciado!")
            return True
        except Exception as e:
            print(e)
            return False
    def sendMessage(self,message,action="message"):
        message = json.dumps({"user":self.NAME,"text":message,"action":action}).encode()
        self.sock.sendto(message, (self.MCAST_GRP, self.MCAST_PORT))
