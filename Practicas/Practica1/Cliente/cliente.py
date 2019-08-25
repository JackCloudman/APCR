from PyQt5 import QtWidgets, uic
import sys
from PyQt5.QtCore import QRegExp
from PyQt5.QtGui import QRegExpValidator
from PyQt5.QtWidgets import QMessageBox, QFileDialog, QListWidget
import socket
import shutil
from os.path import basename
from os import getcwd,stat,remove
import json,time

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
        self.BTN_SFILE.clicked.connect(self.openFile)
        self.BTN_SFOLDER.clicked.connect(self.openFolder)
        self.BTN_UPLOAD.clicked.connect(self.upload)
        self.BTN_LIST.clicked.connect(self.path_list)
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
        self.BTN_SFILE.setEnabled(True)
        self.BTN_SFOLDER.setEnabled(True)
        self.show_message("Exito!","Conexion establecida")
    def openFile(self):
        options = QFileDialog.Options()
        options |= QFileDialog.DontUseNativeDialog
        files, _ = QFileDialog.getOpenFileNames(self,"Select files", "","All Files (*)", options=options)
        if files:
            self.BTN_UPLOAD.setEnabled(True)
            self.FILE_PATH.setText(str(files))
            self.send = files
        else:
            self.BTN_UPLOAD.setEnabled(False)
            self.send = None
            self.FILE_PATH.setText("")
    def openFolder(self):
        options = QFileDialog.Options()
        options |= QFileDialog.DontUseNativeDialog
        folder= QFileDialog.getExistingDirectory(self, 'Select directory',options=options)
        if folder:
            self.BTN_UPLOAD.setEnabled(True)
            self.FILE_PATH.setText(folder)
            self.send = folder
        else:
            self.BTN_UPLOAD.setEnabled(False)
            self.FILE_PATH.setText("")
            self.send = None
    def upload(self):

        '''
            Si lo que vamos a enviar es una carpeta
        '''
        if type(self.send) == str:
            name = basename(self.send)
            if not self.compress(self.send,name):
                self.show_message("Error!","Ha ocurrido un error al preparar el archivo para su envio!")
                return
            print("compress success!")
            message = json.dumps({"command":"subir","tipo":"folder","nombres":[name],"numero":1,"sizes":[self.getSize("temp/%s.zip"%name)]})
            self.s.sendall(message.encode())
            time.sleep(0.1)

            path = "temp/%s.zip"%(basename(self.send))
            self.sendFile(path)
            try:
                remove(path)
            except Exception as e:print(e)
            print("Datos enviados!")
        elif type(self.send) == list:
            message = json.dumps({"command":"subir","tipo":"list","nombres":[basename(x) for x in self.send],"numero":len(self.send),"sizes":[self.getSize(x) for x in self.send]})
            self.s.sendall(message.encode())
            time.sleep(0.05)
            for f in self.send:
                self.sendFile(f)
        else:
            return
    def sendFile(self,path):
        f = open(path,"rb")
        l = f.read(1024)
        bytes_enviados = 0
        while (l):
            self.s.send(l)
            l = f.read(1024)
            bytes_enviados+=1024
        f.close()

    def path_list(self):
        self.s.sendall(json.dumps({"command":"getPathInfo","path":"."}).encode())
        data = self.s.recv(1024)
        json_data = json.loads(data.decode())
        self.list = QListWidget()
        for file in json_data['nombres']:
            len_file = len(file)
            self.list.addItem(str(file) + '     '+str(len_file)+' bytes')

        self.list.show()
        self.list.setSelectionMode(
            QtWidgets.QAbstractItemView.ExtendedSelection
        )
        self.list.itemDoubleClicked.connect(self.items_dwnl)

    def downloads(self,s,name,filesize):
        f = open(name, 'wb')
        myfile = s.recv(1024) #Descargamos el archivo
        totalRecv = len(myfile)
        while True:
            if totalRecv>=filesize:
                f.write(myfile[:-(totalRecv-filesize)])
                f.close()
                break
            f.write(myfile)
            myfile = s.recv(1024)
            totalRecv += len(myfile)
        print ('Descarga completa')
        f.close()

    def items_dwnl(self,item):
        item =item.text()
        print(item)
        take_name = QRegExp("^\S*")
        take_name.indexIn(item)
        input_val= take_name.cap();
        self.s.sendall(json.dumps({"command":"descargar","path":".","nombres":[str(input_val)]}).encode())
        print(json.dumps({"command":"descargar","path":".","nombres":[str(input_val)]}).encode())# Peticion de descarga
        data = self.s.recv(1024)
        data = json.loads(data.decode(errors='ignore')) # Leemos la respuesta
        name = data["nombres"][0]
        filesize = data["sizes"][0]
        self.downloads(self.s,name,filesize)

    def compress(self,path,name):
        try:
            shutil.make_archive("temp/"+name, "zip",path)
            return True
        except:
            return False


    def show_message(self,title,m):
        QMessageBox.about(self, title,m)
    def getSize(self,filename):
        st = stat(filename)
        return st.st_size

app = QtWidgets.QApplication(sys.argv)
window = Ui()
app.exec_()
