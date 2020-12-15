package funcs

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

//Partition struct para una particion
type Partition struct {
	PartStatus [1]byte
	PartType   [1]byte
	PartFit    [1]byte
	PartStart  int64
	PartSize   int64
	PartName   [16]byte
	PartCont   int64
}

//MBR struct para el MBR
type MBR struct {
	MbrTamanio       int64
	MbrFechaCreacion [16]byte
	MbrDiskSignature int64
	DiskFit          [1]byte
	MbrPartition1    Partition
	MbrPartition2    Partition
	MbrPartition3    Partition
	MbrPartition4    Partition
}

//EBR struct para el Extended boot record
type EBR struct {
	PartStatus [1]byte
	PartFit    [1]byte
	PartStart  int64
	PartSize   int64
	PartName   [16]byte
	PartNext   int64
}

//EjecutarMKDISK crea los discos
func EjecutarMKDISK(parametros []string, descripciones []string) {
	fmt.Println(" === COMANDO MKDISK ===")
	algoMalo := false
	var ruta, size, unit, fit string
	var nombre, directorio string
	//fmt.Println("Caen los params: ")
	//fmt.Println(parametros)
	for i := 0; i < len(parametros); i++ {
		switch parametros[i] {
		case "path":
			{
				ruta = descripciones[i]
			}
		case "size":
			{
				size = descripciones[i]
			}
		case "unit":
			{
				unit = descripciones[i]
				//fmt.Println(unit)
			}
		case "fit":
			{
				fit = descripciones[i]
				fmt.Println(fit)
			}
		default:
			{
				fmt.Println("Vienen parametros incorrectos")
				algoMalo = true
			}
		}
	}

	if algoMalo == true {
		fmt.Println("Revise los parametos, hay errores.")
	} else {
		if ruta == "" || size == "" { //si no vienen los parametros obligatorios
			fmt.Println("Faltan parametros obligatorios(path,size)")
		} else {

			nombre, directorio = separarNombreRuta(ruta)
			//fmt.Println("El nombre que se saco es: " + nombre)
			//fmt.Println("el directorio que se saco fue: " + directorio)

			//armar el directorio
			auxBool := false
			err := os.MkdirAll(directorio, 0777) //retorna nil si ya existe el directorio
			if err == nil {
				//fmt.Println("Directorio sin novedad")
			} else { //si ya existe la ruta o hay un error
				fmt.Println("no fue posible crear la ruta")
				auxBool = true
			}

			existeDisco := false
			//toca el archivo
			if auxBool == false { //si no hubo problemas con crear la ruta

				var file *os.File
				//vemos si existe el archivo, sino, lo creamos
				if _, err := os.Stat(directorio + nombre); os.IsNotExist(err) {
					//fmt.Println("El archivo no existe, lo vamos a crear")

					file, _ = os.Create(directorio + nombre)
					defer func() {
						if err := file.Close(); err != nil {
							log.Fatalln(err)
						}
					}()
				} else {
					fmt.Println("El disco que desea crear ya existe")
					existeDisco = true
				}

				if existeDisco != true { // si no existe el disco, procedemos a lllenar

					tamanio64, _ := strconv.ParseInt(size, 10, 64)

					if unit == "" || unit == "m" || unit == "M" {
						tamanio64 = tamanio64 * 1024 * 1024
					} else if unit == "k" || unit == "K" {
						tamanio64 = tamanio64 * 1024
					}
					escribirBytes(file, tamanio64)

					today := time.Now()
					var fecha [16]byte
					for i := 0; i < 16; i++ {
						fecha[i] = today.String()[i]
					}

					var b [1]byte
					b[0] = 0
					file.Seek(0, 0)
					p := Partition{PartStatus: b, PartType: b, PartFit: b}
					elMBR := MBR{MbrTamanio: tamanio64, MbrFechaCreacion: fecha, MbrDiskSignature: 10, MbrPartition1: p, MbrPartition2: p, MbrPartition3: p, MbrPartition4: p}
					errR := binary.Write(file, binary.LittleEndian, elMBR)
					if errR != nil {
						log.Fatalln(errR)
					}

					leerArchivo(directorio + nombre)
				}
			}
		}
	}
}

func separarNombreRuta(ruta string) (string, string) {
	aux := strings.Split(ruta, "/")
	nombre := aux[len(aux)-1]
	var directorio string
	directorio = strings.Replace(ruta, nombre, "", 1)
	return nombre, directorio
}

func escribirBytes(file *os.File, tam int64) { // recibe el tamanio del archivo
	if tam != 0 {
		tamanio := int(tam)
		var a []byte
		a = make([]byte, tamanio, tamanio)
		for i := 0; i < tamanio; i++ {
			a[i] = 0
		}

		_, err := file.Write(a)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func leerArchivo(ruta string) {
	file, err := os.Open(ruta)
	if err != nil {
		log.Fatalln(err)
	} else {
		defer file.Close()

		master := MBR{}
		err = binary.Read(file, binary.LittleEndian, &master)
		if err != nil {
			log.Fatalln(err)
		}

		fecha := ""
		for i := 0; i < 16; i++ {
			fecha += string(master.MbrFechaCreacion[i])
		}

		fmt.Println("El disco se ha creado " + fecha)
	}
}
