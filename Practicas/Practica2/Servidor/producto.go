package main

// Estructura básica de un producto
type Articulo struct {
	ID          int     `json:"id"`
	Nombre      string  `json:"nombre"`
	Precio      float64 `json:"precio"`
	Existencias int     `json:"existencias"`
	Descripcion string  `json:"descripcion"`
	Promocion   int     `json:"promocion"`
	Imagen      []byte  `json:"imagen"`
}
type Ticket struct {
	ID                    int         `json:"ticketid"`
	ArticulosPreciostotal []float64   `json:"preciostotal"`
	Articulos             []*Articulo `json:"articulos"`
	Fecha                 string      `json:"fechacompra"`
	Total                 float64     `json:"total"`
	DescuentoTotal        float64     `json:"descuentototal"`
}

// Variable global que simula la BD .-.s
var ListaProductos = []*Articulo{}
var ListaTickets = []Ticket{}

// Funcion que añade productos
func crearProducto(nombre, descripcion, image string, precio float64, existencias, promocion int) {
	a := new(Articulo)
	a.Nombre = nombre
	a.Descripcion = descripcion
	a.Precio = precio
	a.Existencias = existencias
	a.Promocion = promocion
	a.ID = len(ListaProductos)
	a.Imagen = ImageTobase64(image)
	ListaProductos = append(ListaProductos, a)
}
func generarTicket(articulos []*Articulo) Ticket {
	t := Ticket{}
	t.ID = len(ListaTickets) + 1
	t.Total = 0
	t.DescuentoTotal = 0
	for _, a := range articulos {
		ListaProductos[a.ID].Existencias -= a.Existencias
		articulo := *ListaProductos[a.ID]
		articulo.Imagen = nil
		articulo.Descripcion = ""
		articulo.Existencias = a.Existencias
		descuento := articulo.Precio * float64(articulo.Promocion) / 100
		t.DescuentoTotal += (descuento * float64(a.Existencias))
		atotal := float64(a.Existencias) * (articulo.Precio - descuento)
		t.Total += atotal
		t.ArticulosPreciostotal = append(t.ArticulosPreciostotal, atotal)
		t.Articulos = append(t.Articulos, &articulo)
	}
	ListaTickets = append(ListaTickets, t)
	return t
}
func comprarProductos(articulos []*Articulo) (int, Ticket) {
	t := Ticket{}
	for _, a := range articulos {
		if ListaProductos[a.ID].Existencias-a.Existencias < 0 {
			return 400, t
		}
	}
	return 200, generarTicket(articulos)
}
