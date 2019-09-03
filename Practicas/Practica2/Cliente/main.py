from PyQt5 import QtWidgets, uic
import sys
from PyQt5.QtCore import QRegExp, QByteArray, qUncompress
from PyQt5.QtGui import QRegExpValidator, QPixmap, QImage
from PyQt5.QtWidgets import QMessageBox, QFileDialog, QListWidget, QWidget, QVBoxLayout, QHBoxLayout,QLabel, QListWidgetItem
import socket
import shutil
from os.path import basename
from os import getcwd,stat,remove
import json,time
from utils import QCustomQWidget



class Ui(QtWidgets.QMainWindow):
    def __init__(self):
        super(Ui, self).__init__() # Call the inherited classes __init__ method
        uic.loadUi('gui.ui', self) # Load the .ui file
        # Validator
        reg_host = QRegExp("^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$")
        input_validator = QRegExpValidator(reg_host, self.HOST)
        self.HOST.setValidator(input_validator)

        reg_port = QRegExp("^([0-9]{1,4}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$")
        input_validator = QRegExpValidator(reg_port, self.PORT)
        self.PORT.setValidator(input_validator)

        self.BTN_CONN.clicked.connect(self.connect)
        #self.BTN_SFILE.clicked.connect(self.openFile)
        #self.BTN_SFOLDER.clicked.connect(self.openFolder)
        #self.BTN_UPLOAD.clicked.connect(self.upload)
        #self.BTN_LIST.clicked.connect(self.path_list)
        self.send = None
        self.show() # Show the GUI
    def connect(self):
        h = self.HOST.text()
        p = self.PORT.text()
        if h==p=="":
            return
        try:
            self.s = socket.socket(socket.AF_INET,socket.SOCK_STREAM)
            self.s.connect((h,int(p)))
        except Exception as e:
            self.show_message("Error!","No se ha podido establecer comunicación")
            print(e)
            return
        # Buttons
        self.BTN_CONN.setEnabled(False)
        self.HOST.setEnabled(False)
        self.PORT.setEnabled(False)
        self.STATUS.setChecked(True)
        #self.BTN_SFILE.setEnabled(True)
        #self.BTN_SFOLDER.setEnabled(True)
        self.show_message("Exito!","Conexion establecida")
        #str_image = b'''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'''
        ################################################################

        datas =self.getCatalogo(self.s)
        self.myQListWidget = QListWidget(self)
        #self.myQListWidget.itemDoubleClicked.connect(self.show_message('compraste algo', 'exito!'))
        for articulos in datas['articulos']:
            # Create QCustomQWidget
            str_image =articulos['imagen']
            image = self.decode_thumbnail(str_image)
            myQCustomQWidget = QCustomQWidget()
            myQCustomQWidget.setTextUp('Articulo: '+articulos['nombre'])
            myQCustomQWidget.setTextDown('$ '+str(articulos['precio']))
            myQCustomQWidget.setItemNum('Disponibilidad:'+str(articulos['existencias']))

            myQCustomQWidget.setIcon(QPixmap.fromImage(image))
            myQCustomQWidget.setId=articulos['id']
            # Create QListWidgetItem
            myQListWidgetItem = QListWidgetItem(self.myQListWidget)
            # Set size hint
            myQListWidgetItem.setSizeHint(myQCustomQWidget.sizeHint())
            # Add QListWidgetItem into QListWidget
            self.myQListWidget.addItem(myQListWidgetItem)
            self.myQListWidget.setItemWidget(myQListWidgetItem, myQCustomQWidget)
            self.setCentralWidget(self.myQListWidget)

        self.myQListWidget.show()
        self.myQListWidget.itemDoubleClicked.connect(self.items_buyied)

    def items_buyied(self, item):

        item_list=item.listWidget()
        my_customWidget=item_list.itemWidget(item)
        print(my_customWidget.getTextName())
        print(my_customWidget.getPrice())
        print(my_customWidget.getExistence())
        if self.reques_buy(self.s, my_customWidget.getId()):
            show_message('Compra Exitosa!', 'haz comprado:'+ my_customWidget.getTextName())

    def reques_buy(self, socket, idItem):
        # el comando para comprar es "comprar", debe incluir los articulos a comprar
        #poniendo sus id, existencias representa la cantidad de articulos que quieres
        command = '{"command":"comprar","articulos":[{"id":%s,"existencias":1}]}'%idItem
        socket.sendall(command.encode())
        ticket = self.recvall(socket)


    def show_message(self,title='jeje',m='exito'):
        QMessageBox.about(self,title,m)


    def recvall(self, sock):
        BUFF_SIZE = 1024 # 4 KiB
        data = b''
        while True:
            part = sock.recv(BUFF_SIZE)
            data += part
            print(len(part))
            if len(part) < BUFF_SIZE:
                # either 0 or end of data
                break
        print(data)
        return json.loads(data.decode())

    def getCatalogo(self, socket):
        command = '{"command":"getCatalogo"}'
        socket.sendall(command.encode())
        return self.recvall(socket)

    def decode_thumbnail(self,str_image):
       bytearray = QByteArray.fromBase64(str_image)
       bytearray = qUncompress(bytearray)
       return QImage.fromData(bytearray, 'PNG')



app = QtWidgets.QApplication(sys.argv)
window = Ui()
app.exec_()
