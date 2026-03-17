package main

// Importación de paquetes
import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// 1.Declaro la constante para el nombre del archivo de salida
const SalidaAuditoria = "auditoria.txt"

// 2. Defino el tipo de dato que representa cada acceso (Un usuario y el nivel de riesgo)
type Acesso struct {
	USUARIO string
	RIESGO  int
	ID      int
}

func main() {

	//Se crea 5 posiciones con datos de prueba (Dos de ellos tienen un riesgo mayor a 8)
	array := [5]Acesso{
		{ID: 001, USUARIO: "Victor", RIESGO: E},  //Seguro
		{ID: 030, USUARIO: "Miguel", RIESGO: 5},  //Advertencia
		{ID: 100, USUARIO: "Seyi", RIESGO: 7},    //Advertencia
		{ID: 007, USUARIO: "Ruben", RIESGO: 9},   //CRÍTICO
		{ID: 064, USUARIO: "Michel", RIESGO: 10}, //CRÍTICO
	}

	//Se usa la constante del archivo de salida
	fmt.Println("Archivo de auditoría:", SalidaAuditoria)
	fmt.Println()

	// 3. Se recorre el array con un bucle for, clasificando riesgos con switch y calculando la media
	var sumaRiesgos int = 0 // Variable para acumular los riesgos

	// Se recorre el array creado con un bucle for, con los riesgos clasificados usando switch
	for i, acceso := range array {
		var categoria string

		switch {
		case acceso.RIESGO >= 1 && acceso.RIESGO <= 3:
			categoria = "Seguro"
		case acceso.RIESGO >= 4 && acceso.RIESGO <= 7:
			categoria = "Advertencia"
		case acceso.RIESGO >= 8 && acceso.RIESGO <= 10:
			categoria = "ATAQUE CRÍTICO"
		default:
			categoria = "Desconocido"
		}

		// Se muestra el resultado del acceso
		fmt.Printf("ID:%d ->[%d] Usuario: %s | Riesgo: %d | Categoría: %s\n", i+1, acceso.ID, acceso.USUARIO, acceso.RIESGO, categoria)

		//Alerta en tiempo real si es un ATAQUE CRÍTICO
		if categoria == "ATAQUE CRÍTICO" {
			fmt.Printf(" [ALERTA], usuario %s ha detectado un ataque crítico con nivel de riesgo:%d\n", acceso.USUARIO, acceso.RIESGO)
		}
		// Acumular los riesgos para calcular la media
		sumaRiesgos += acceso.RIESGO

	}

	fmt.Println()

	// Calcular la media arimética con conversión en decimales
	mediaRiesgos := float64(sumaRiesgos) / float64(len(array))
	fmt.Printf("Media aritmética del nivel de riesgo: %.2f\n", mediaRiesgos)

	// 4. Persistencia y concurrencia
	fmt.Println()
	fmt.Println("Iniciando escritura en archivo de auditoría...")

	// Se utiliza WaitGroup para sincronizar la goroutine
	var wg sync.WaitGroup
	wg.Add(1)

	// Se lanza la función de escritura en una goroutine
	go guardarAuditoria(array, &wg)

	// Se espera a que se complete la goroutine
	wg.Wait()
	fmt.Println("Auditoría completada.")
}

// 4. Persistencia y concuerrencia

func guardarAuditoria(accesos [5]Acesso, wg *sync.WaitGroup) {
	defer wg.Done()

	// Se abre el archivo en modo append (crear si no existe)

	archivo, err := os.OpenFile(SalidaAuditoria, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf(" Error al abrir/crear el archivo de auditoría: %v", err)
		return
	}
	defer archivo.Close()

	// Se configura el log para escribir en el archivo
	logger := log.New(archivo, "", log.LstdFlags)

	// Se registra el inicio de la auditoría para el archivo
	logger.Println("========== INICIO DE AUDITORÍA ==========")

	// Se recorre el array y se escriben los datos en el archivo con formato Sprintf
	for i, acceso := range accesos {
		var categoria string

		switch {
		case acceso.RIESGO >= 1 && acceso.RIESGO <= 3:
			categoria = "Seguro"
		case acceso.RIESGO >= 4 && acceso.RIESGO <= 7:
			categoria = "Advertencia"
		case acceso.RIESGO >= 8 && acceso.RIESGO <= 10:
			categoria = "ATAQUE CRÍTICO"
		default:
			categoria = "Desconocido"
		}

		// Se formatea cada línea con Sprintf y se escribe en el archivo
		registro := fmt.Sprintf("ID:%d ->[%d] Usuario: %-8s | Riesgo: %2d | Categoría: %-15s", i+1, acceso.ID, acceso.USUARIO, acceso.RIESGO, categoria)
		logger.Println(registro)

		// Alerta si es crítico
		if categoria == "ATAQUE CRÍTICO" {
			alertaRegistro := fmt.Sprintf(" [ALERTA], Usuario %s ha detectado un ataque crítico con nivel de riesgo: %d", acceso.USUARIO, acceso.RIESGO)
			logger.Println(alertaRegistro)
		}
	}

	logger.Println("========== FIN DE AUDITORÍA ==========")
	logger.Println("")

	//Tiempo de epsera para la escritura

	time.Sleep(100 * time.Millisecond)

	fmt.Printf("Datos guardados en: %s\n", SalidaAuditoria)
}
