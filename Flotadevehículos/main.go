package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

// CONSTANTE: Kilometraje máximo antes de revisión
const KM_MAXIMO = 10000

// Nombre del archivo de log
const ARCHIVO_LOG = "flota.log"

// Estructura para los vehículos
type Vehiculo struct {
	ID          int
	Nombre      string
	Kilometraje int
	CostePorKm  float64
}

// Función que escribe alertas en un archivo físico
func guardarAlerta(mensaje string) {
	f, err := os.OpenFile(ARCHIVO_LOG, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("[ERROR] Al acceder al archivo de log: ", err)
	}
	defer f.Close()

	timestamp := time.Now().Format("15:04:05")
	f.WriteString(timestamp + " - " + mensaje + "\n")
}

func main() {
	// Array con 4 vehículos
	flota := [4]Vehiculo{
		{ID: 1, Nombre: "Furgoneta A", Kilometraje: 500, CostePorKm: 0.15},
		{ID: 2, Nombre: "Camión B", Kilometraje: 9800, CostePorKm: 0.45},
		{ID: 3, Nombre: "Furgoneta C", Kilometraje: 10500, CostePorKm: 0.15},
		{ID: 4, Nombre: "Moto D", Kilometraje: 10000, CostePorKm: 0.08},
	}

	var costeTotalGeneral float64

	fmt.Println("----- PROCESANDO FLOTA DE VEHÍCULOS -----")

	for _, v := range flota {
		// Calcular coste del vehículo
		costeVehiculo := float64(v.Kilometraje) * v.CostePorKm
		costeTotalGeneral += costeVehiculo

		// Clasificar el vehículo según kilometraje
		var estado string
		switch {
		case v.Kilometraje >= KM_MAXIMO:
			estado = "REVISION URGENTE"
		case v.Kilometraje >= KM_MAXIMO-1000:
			estado = "PROXIMAMENTE REVISION"
		default:
			estado = "EN BUEN ESTADO"
		}

		// Si necesita atención, lanzar goroutine
		if estado != "EN BUEN ESTADO" {
			mensaje := fmt.Sprintf("ALERTA: %s (ID: %d) - %s", v.Nombre, v.ID, estado)
			go guardarAlerta(mensaje)
			fmt.Printf("¡ %s\n", mensaje)
		} else {
			fmt.Printf(". %s: OK (%d km)\n", v.Nombre, v.Kilometraje)
		}
	}

	fmt.Printf("\nCOSTE TOTAL DE LA FLOTA: %.2f €\n", costeTotalGeneral)
	fmt.Println("Guardando registros...")
	time.Sleep(1 * time.Second)
	fmt.Println("Sistema cerrado.")
}
