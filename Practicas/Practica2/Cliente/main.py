from PyQt5 import QtWidgets, uic
import sys
from PyQt5.QtCore import QRegExp, QByteArray, qUncompress
from PyQt5.QtGui import QRegExpValidator, QPixmap, QImage
from PyQt5.QtWidgets import QMessageBox, QFileDialog, QListWidget, QWidget, QVBoxLayout, QHBoxLayout,QLabel, QListWidgetItem,QPushButton
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
        self.productos = []
        self.tickets = []
        # Validator
        reg_host = QRegExp("^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$")
        input_validator = QRegExpValidator(reg_host, self.HOST)
        self.HOST.setValidator(input_validator)

        reg_port = QRegExp("^([0-9]{1,4}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$")
        input_validator = QRegExpValidator(reg_port, self.PORT)
        self.PORT.setValidator(input_validator)
        self.BTN_CONN.clicked.connect(self.connect)
        self.BTN_COMPRAR.clicked.connect(self.comprar)
        self.BTN_TICKETS.clicked.connect(self.mostrartickets)
        #self.BTN_SFILE.clicked.connect(self.openFile)
        #self.BTN_SFOLDER.clicked.connect(self.openFolder)
        #self.BTN_UPLOAD.clicked.connect(self.upload)
        #self.BTN_LIST.clicked.connect(self.path_list)
        self.send = None
        self.show() # Show the GUI
    def comprar(self):
        bandera = False
        articulos = []
        for p in self.productos:
            if p.ItemNumSpin.value()>0:
                articulos.append({"id":p.getId(),"existencias":p.getNum()})
        if not articulos:
            print("Compra al menos 1")
            return
        ticket = self.reques_buy(articulos)
        self.tickets.append(ticket)
        self.reloadArticulos()
    def reloadArticulos(self):
        articulos = self.getCatalogo(self.s)["articulos"]
        for a in self.productos:
            a.setItemNum("Disponibilidad:%d"%articulos[a.getId()]["existencias"])
            a.resetSpinbox()
    def mostrartickets(self):
        self.list = QListWidget()
        if self.tickets:
            for t in self.tickets:
                texto = "Numero de ticket:%d\nFecha:%s\n"%(t["ticketid"],t["fechacompra"])
                texto+= "Articulo\t\tCANT\tPRECIO\tDESC(%)\tIMPORTE\n"
                for index,a in enumerate(t["articulos"]):
                    texto+="%s\t\t%d\t%.2f\t%d\t%.2f\n"%(a["nombre"][:10],a["existencias"],a["precio"],a["promocion"],t["preciostotal"][index])
                texto+="Total: %.2f\nUsted ahorro:%.2f"%(t["total"],t["descuentototal"])
                self.list.addItem(texto)
        else:
            self.list.addItem("No hay tickets!")
        self.list.show()
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

        datas =self.getCatalogo(self.s)
        self.myQListWidget = self.myList
        #self.myQListWidget.itemDoubleClicked.connect(self.show_message('compraste algo', 'exito!'))
        for articulos in datas['articulos']:
            # Create QCustomQWidget
            str_image =articulos['imagen']
            image = self.decode_thumbnail(str_image)
            myQCustomQWidget = QCustomQWidget()
            myQCustomQWidget.setTextUp('Articulo: '+articulos['nombre'])
            myQCustomQWidget.setTextDown('$ '+str(articulos['precio']))
            myQCustomQWidget.setItemNum('Disponibilidad:'+str(articulos['existencias']))
            myQCustomQWidget.ItemNumSpin.setMaximum(articulos['existencias'])

            myQCustomQWidget.setIcon(QPixmap.fromImage(image))
            myQCustomQWidget.setId(articulos['id'])
            self.productos.append(myQCustomQWidget)
            # Create QListWidgetItem
            myQListWidgetItem = QListWidgetItem(self.myQListWidget)
            # Set size hint
            myQListWidgetItem.setSizeHint(myQCustomQWidget.sizeHint())
            # Add QListWidgetItem into QListWidget
            self.myQListWidget.addItem(myQListWidgetItem)
            self.myQListWidget.setItemWidget(myQListWidgetItem, myQCustomQWidget)
            self.setCentralWidget(self.myQListWidget)
        self.BTN_COMPRAR.setEnabled(True)
        self.BTN_TICKETS.setEnabled(True)
        #self.myQListWidget.show()
        #self.myQListWidget.itemDoubleClicked.connect(self.items_buyied)

    '''def items_buyied(self, item):

        item_list=item.listWidget()
        my_customWidget=item_list.itemWidget(item)
        print(my_customWidget.getTextName())
        print(my_customWidget.getPrice())
        print(my_customWidget.getExistence())
        if self.reques_buy(self.s, my_customWidget.getId()):
            show_message('Compra Exitosa!', 'haz comprado:'+ my_customWidget.getTextName())'''

    def reques_buy(self,articulos):
        # el comando para comprar es "comprar", debe incluir los articulos a comprar
        #poniendo sus id, existencias representa la cantidad de articulos que quieres
        command = {"command":"comprar","articulos":articulos}
        command = json.dumps(command)
        self.s.sendall(command.encode())
        response = self.recvall(self.s)
        if response["status"] != 200:
            print("Error!")
        return response["ticket"]


    def show_message(self,title='jeje',m='exito'):
        QMessageBox.about(self,title,m)


    def recvall(self, sock):
        BUFF_SIZE = 1024 # 4 KiB
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

    def decode_thumbnail(self,str_image):
       bytearray = QByteArray.fromBase64(str_image.encode())
       return QImage.fromData(bytearray, 'JPG')



app = QtWidgets.QApplication(sys.argv)
window = Ui()
app.exec_()
