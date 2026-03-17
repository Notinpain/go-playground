// Defino el paquete principal.
package main

// Importo los paquetes necesarios para la entrada/salida, errores, archivos y tiempo.
import (
	"fmt"
	"log"
	"os"
	"time"
)

// CONSTANTE: Mínimo de stock.
const STOCK_MINIMO = 5

// Variable global con el nombre del archivo donde guadaré las alertas.
const ARCHIVO_LOG = "alertas.log"

// ESTRUCTURA para los productos.
type Producto struct {
	ID       int
	Nombre   string
	Cantidad int
	Precio   float64
}

// Función que escribe alertas en un archivo físico.
func guardarAlerta(mensage string) {
	// Abro el archivo, utilizando la flag para crearlo si no existe y la que indica que se escriba al final de este. A parte, indico los permisos 0644 que son el estándar para archivos de lectura/escritura en Windows.
	f, err := os.OpenFile(ARCHIVO_LOG, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	// Manejo de errores.
	if err != nil {
		log.Fatal("[ERROR] Al acceder al archivo de log: ", err)
	}

	// Me aseguro de cerrar el archivo al archivo al terminar la función.
	defer f.Close()

	// Escibo el mensaje con la hora actual para que sea un log "profesional".
	timestamp := time.Now().Format("15:04:05")
	f.WriteString(timestamp + " - " + mensage + "\n")
}

func main() {
	// Array con una lista de 4 productos.
	inventario := [4]Producto{
		{ID: 1, Nombre: "Cajas de Cartón", Cantidad: 33, Precio: 1.50},
		{ID: 2, Nombre: "Cinta de embalar", Cantidad: 1, Precio: 0.80},
		{ID: 3, Nombre: "Palets de Madera", Cantidad: 0, Precio: 15.00},
		{ID: 4, Nombre: "Bolsas de Plástico", Cantidad: 0, Precio: 0.20},
	}

	// Variable para acumular el valor total de todo el stock.
	var valorTotalGeneral float64

	fmt.Println("----- PROCESANDO INVENTARIO LOGÍSTICO -----")
	// Bucles: Recorro el array de productos usando [range]
	for _, p := range inventario {
		// Operadores aritméticos: Calculo el valor del producto actual (Cantidad * Precio).
		// NOTA: Convierto [Cantidad] (int) a (float64) para poder multiplicar correctamente.
		valorProducto := float64(p.Cantidad) * p.Precio
		valorTotalGeneral += valorProducto

		// Variable para almacenar el estado en el que se encuentra el producto (agotado/en strock)
		var estado string

		// SWITCH: Clasificamos el producto según el stock.
		switch {
		case p.Cantidad == 0:
			estado = "AGOTADO"
		case p.Cantidad <= STOCK_MINIMO:
			estado = "PEDIDO URGENTE"
		default:
			estado = "STOCK SUFICIENTE"
		}

		// CONDICIONALES Y CONCURRENCI: Si hay alerta, disparamos la Goroutine.
		if estado == "AGOTADO" || estado == "PEDIDO URGENTE" {
			mensajeAlerta := fmt.Sprintf("ALERTA: %s (ID: %d) tiene estado %s", p.Nombre, p.ID, estado)

			// Lanzo la concurrencia: La escritura no detiene el bucle principal.
			go guardarAlerta(mensajeAlerta)

			fmt.Printf("¡ %s\n", mensajeAlerta)
		} else {
			fmt.Printf(". %s: OK (%d unidades)\n", p.Nombre, p.Cantidad)
		}
	}

	// PAQUETE FMT: Muestro el resumen final con formato de dos decimales.
	fmt.Printf("\nVALOR TOTAL DEL INVENTARIO: %.2f €\n", valorTotalGeneral)

	// ESPERA: Pauso 1 segundo para que las Goroutines terminen de escribir en el archivo.
	fmt.Println("Finalizando procesos de guardado...")
	time.Sleep(1 * time.Second)
	fmt.Println("Sistema cerrado.")
}
