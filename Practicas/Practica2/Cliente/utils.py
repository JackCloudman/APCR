from PyQt5 import QtWidgets, uic
from PyQt5.QtGui import QRegExpValidator, QPixmap
from PyQt5.QtWidgets import QMessageBox, QFileDialog, QListWidget, QWidget, QVBoxLayout, QHBoxLayout,QLabel, QListWidgetItem


class QCustomQWidget(QWidget):

    def __init__ (self, parent = None):
        super(QCustomQWidget, self).__init__(parent)

        self.setWindowTitle("My Awesome App")
        self.textQVBoxLayout = QVBoxLayout()
        self.textUpQLabel    = QLabel()
        self.textDownQLabel  = QLabel()
        self.textItemsnum =QLabel()
        self.textQVBoxLayout.addWidget(self.textUpQLabel)
        self.textQVBoxLayout.addWidget(self.textDownQLabel)
        self.textQVBoxLayout.addWidget(self.textItemsnum)
        self.allQHBoxLayout  = QHBoxLayout()
        self.iconQLabel      = QLabel()
        self.allQHBoxLayout.addWidget(self.iconQLabel, 0)
        self.allQHBoxLayout.addLayout(self.textQVBoxLayout, 1)
        self.setLayout(self.allQHBoxLayout)
        # setStyleSheet
        self.textUpQLabel.setStyleSheet('''
        font: bold 20px;
        color: rgb(0, 0, 255);
        ''')
        self.textDownQLabel.setStyleSheet('''
        font: bold 20px;
        color: rgb(255, 0, 0);
        ''')
        self.textItemsnum.setStyleSheet('''
        font: bold 20px;
        color: rgb(255, 0, 0);
        ''')

    def setTextUp(self, text):
            self.textUpQLabel.setText(text)

    def setTextDown (self, text):
            self.textDownQLabel.setText(text)

    def setIcon (self, imagePath):
            self.iconQLabel.setPixmap(QPixmap(imagePath))

    def setItemNum(self,text):
        self.textItemsnum.setText(text)

    def getTextName(self):
         return self.textUpQLabel.text()


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
