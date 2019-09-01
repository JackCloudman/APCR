from PyQt5 import QtWidgets, uic
import sys
from PyQt5.QtCore import QRegExp
from PyQt5.QtGui import QRegExpValidator, QPixmap
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

        ################################################################

        datas =self.getCatalogo(self.s)
        self.myQListWidget = QListWidget(self)
        #self.myQListWidget.itemDoubleClicked.connect(self.show_message('compraste algo', 'exito!'))
        for articulos in datas['articulos']:
            # Create QCustomQWidget
            myQCustomQWidget = QCustomQWidget()
            myQCustomQWidget.setTextUp(articulos['nombre'])
            myQCustomQWidget.setTextDown(str(articulos['precio']))
            myQCustomQWidget.setItemNum(str(articulos['existencias']))
            myQCustomQWidget.setIcon('icon.png')
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

    def show_message(self,title='jeje',m='exito'):
        QMessageBox.about(self,title,m)


    def recvall(self, sock):
        BUFF_SIZE = 4096 # 4 KiB
        data = b''
        while True:
            part = sock.recv(BUFF_SIZE)
            data += part
            if len(part) < BUFF_SIZE:
                # either 0 or end of data
                break
        return json.loads(data.decode())

    def getCatalogo(self, socket):
        command = '{"command":"getCatalogo"}'
        socket.sendall(command.encode())
        return self.recvall(socket)



app = QtWidgets.QApplication(sys.argv)
window = Ui()
app.exec_()
