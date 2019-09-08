import socket
import threading
import struct
import json
MCAST_GRP = '224.1.1.1'
MCAST_PORT = 5007
IS_ALL_GROUPS = True
MULTICAST_TTL = 2
name = input("Ingresa tu nombre: ")
def sendMessage(sock,message,action):
    message = json.dumps({"user":name,"text":message,"action":action}).encode()
    sock.sendto(message, (MCAST_GRP, MCAST_PORT))
def sender():
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM, socket.IPPROTO_UDP)
    sock.setsockopt(socket.IPPROTO_IP, socket.IP_MULTICAST_TTL, MULTICAST_TTL)
    sendMessage(sock,"","join")
    while True:
        texto = input(">")
        sendMessage(sock,texto,"message")

def processAction(message):
    if message["action"] == "message":
        print("%s:%s\n>"%(message["user"],message["text"]),end="")
    elif message["action"] == "join":
        print("El usuario %s se ha unido!\n>"%(message["user"]),end="")
def listener():
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM, socket.IPPROTO_UDP)
    sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    if IS_ALL_GROUPS:
        # on this port, receives ALL multicast groups
        sock.bind(('', MCAST_PORT))
    else:
        # on this port, listen ONLY to MCAST_GRP
        sock.bind((MCAST_GRP, MCAST_PORT))
    mreq = struct.pack("4sl", socket.inet_aton(MCAST_GRP), socket.INADDR_ANY)

    sock.setsockopt(socket.IPPROTO_IP, socket.IP_ADD_MEMBERSHIP, mreq)

    while True:
        data = sock.recv(10240)
        data = json.loads(data.decode())
        if data["user"] == name:
            pass
        else:
            processAction(data)

threadlist = []
threadlist.append(threading.Thread(target=listener))
threadlist.append(threading.Thread(target=sender))
for t in threadlist:
    t.start()
