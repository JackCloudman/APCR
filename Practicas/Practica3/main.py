from __future__ import print_function	# For Py2/3 compatibility
import eel
from utils import Connection


if __name__ == '__main__':
    eel.init('gui')
    s = Connection("224.1.1.1",5007,"JACK")
    print("Iniciando conexion!")
    if not s.start():
        print("Error al establecer conexion!")
        exit(code=1)

    @eel.expose                         # Decorador para eel (Exponemos la funcion a javascript)
    def sendCommand(message,action):
        s.sendMessage(message,action) # Enviamos el mensaje
        return "ok" # Regresamos a Javascript el resultado de la consulta!

    eel.start('index.html')    # Start
