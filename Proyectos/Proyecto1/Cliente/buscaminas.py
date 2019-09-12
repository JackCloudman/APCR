'''
Codigo modificado por @JackCloudman
Buscaminas en socket, codigo basado en el tutorial
de @ThomasS1, con las modificaciones para funcionar con un
buscaminas de MxN con R minas.
'''

import random, time, copy,os
from termcolor import cprint
class Buscaminas():
    def __init__(self,b,f,c,minas):
        self.b = b
        self.f = f
        self.c = c
        self.minas = minas
    def start(self):
        print('''
        MAIN MENU
        =========

        -> Para instrucciones escribe I
        -> Para jugar escribe P
        ''')

        choice = input('Escribe aquí: ').upper()

        if choice == 'I':
            #replit.clear()
            os.system("clear")

            #Prints instructions.
            print(open('instrucciones.txt', 'r').read())

            input('Presiona [Enter] para jugar')

        elif choice != 'P':
            os.system("clear")
            self.start()

        #The solution grid.

        #Iniciamos la matriz de tamaño personalizado
        self.k = [[" " for j in range(self.c)] for i in range(self.f)]

        self.printBoard(self.k)

        #Start timer
        self.startTime = time.time()

        #The game begins!
        self.play()
    def printBoard(self,b):
        os.system("clear")
        print("\t  "+"   ".join([str(x) for x in range(self.c)]))
        print('\t╔═══'+'╦═══'*(self.c-2)+'╦═══╗')
        for fila in range(0,self.f):
            print(fila,'\t║'," ║ ".join([str(x) for x in b[fila]]),"║")
            if fila != self.f-1:
                print('\t╠═══'+'╬═══'*(self.c-2)+'╬═══╣')
        print('\t╚═══'+'╩═══'*(self.c-2)+'╩═══╝')
    def play(self):
        b = self.b
        k = self.k
        #Player chooses square.
        f, c = self.choose()
        #Gets the value at that location.
        v = b[f][c]
        #If you hit a bomb, it ends the game.
        if v == '*':
            self.printBoard(b)
            print('Perdiste!')
            #Print timer result.
            print('Tiempo: ' + str(round(time.time() - self.startTime)) + 's')
            return
            #Offer to play again.
        #Puts that value into the known grid (k).
        k[f][c] = v
        #Runs checkZeros() if that value is a 0.
        if v == 0:
            self.checkZeros(f, c)
        self.printBoard(k)
        #Checks to see if you have won.
        squaresLeft = 0
        for x in range(self.f):
            row = k[x]
            squaresLeft += row.count(' ')
            squaresLeft += row.count('⚐')
        if squaresLeft == self.minas:
            self.printBoard(b)
            print('Ganaste!!')
            #Print timer result.
            print('Time: ' + str(round(time.time() - self.startTime)) + 's')
            return
            #Offer to play again.
        #Repeats!
        self.play()
        return
    def choose(self):
        b = self.b
        k = self.k
        startTime = self.startTime
        #Variables 'n stuff.
        columNumber = [str(x) for x in range(self.c)]
        numbers = [str(x) for x in range(self.f)]
        #Loop in case of invalid entry.
        while True:
            try:
                chosen = input('Escoge un cuadro (fila,columna) ej. 4,3 o pon una marca (ej. m4,3): ').lower()
                #Checks for valid square.
                if chosen[0] == "m":
                    f,c = chosen[1:].split(",")
                    f = int(f)
                    c = int(c)
                    print("fila %d columna: %d"%(f,c))
                    self.marker(f, c)
                else:
                    f,c = chosen.split(",")
                    f = int(f)
                    c = int(c)
                    print("fila %d columna: %d"%(f,c))
                    return f,c
            except:
                pass
    def marker(self,f,c):
        self.k[f][c] = '⚐'
        self.printBoard(self.k)
    def zeroProcedure(self,f, c):
        b = self.b
        k = self.k

        #Row above
        if f-1 > -1:
            row = k[f-1]
            if c-1 > -1: row[c-1] = b[f-1][c-1]
            row[c] = b[f-1][c]
            if self.c > c+1: row[c+1] = b[f-1][c+1]

        #Same row
        row = k[f]
        if c-1 > -1: row[c-1] = b[f][c-1]
        if self.c > c+1: row[c+1] = b[f][c+1]

        #Row below
        if self.f > f+1:
            row = k[f+1]
            if c-1 > -1: row[c-1] = b[f+1][c-1]
            row[c] = b[f+1][c]
            if self.c > c+1: row[c+1] = b[f+1][c+1]
    def checkZeros(self,f, c):
        k = self.k
        b = self.b
        oldGrid = copy.deepcopy(k)
        self.zeroProcedure(f, c)
        if oldGrid == k:
            return
        while True:
            oldGrid = copy.deepcopy(k)
            for x in range (self.f):
                for y in range (self.c):
                    if k[x][y] == 0:
                        self.zeroProcedure(x, y)
            if oldGrid == k:
                return
