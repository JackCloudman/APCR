import socket	#for sockets
import sys	#for exit

# create dgram udp socket
try:
	s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
except socket.error:
	print('Failed to create socket')
	sys.exit()

host = 'localhost';
port = 8080;

while(1) :
	msg = input('Enter message to send : ')

	try :
		#Set the whole string
		s.sendto(msg.encode(), (host, port))

		# receive data from client (data, addr)
		d = s.recvfrom(1024)
		print(d)

	except Exception as e:
		print(e)
		sys.exit()
