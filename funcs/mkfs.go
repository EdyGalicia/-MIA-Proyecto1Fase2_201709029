package funcs

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"
	"time"
	"unsafe"
)

//SuperBloque struct SP
type SuperBloque struct {
	TipoDeSistemaDeArchivos int64

	//numero total de inodos/bloques
	NumTotalDeInodos  int64
	NumTotalDeBloques int64

	//numero de bloques/inodos libres
	NumDeBloquesLibres int64
	NumDeInodosLibres  int64

	FechaMontada       [16]byte
	FechaDesmontada    [16]byte
	ContadorDeMontadas int64
	SMagic             int64

	//tamanio del inodo/bloque
	TamanioDelInodo  int64
	TamanioDelBloque int64

	//primer inodo/bloque libre
	PrimerInodoLibre  int64
	PrimerBloqueLibre int64

	//inicio Bitmap inodos/bloques
	StartBMdeInodos  int64
	StartBMdeBloques int64

	//inicio de inodos/bloques
	StartTablaDeInodos  int64
	StartTablaDeBloques int64
}

//Inodo info
type Inodo struct {
	IUid   int64
	IGid   int64
	ISize  int64
	IAtime [16]byte
	ICtime [16]byte
	IMtime [16]byte
	IBlock [15]int64
	IType  [1]byte
	IPerm  int64
}

//BloqueDeCarpeta info
type BloqueDeCarpeta struct {
	BContent [4]Contenido
}

//Contenido info
type Contenido struct {
	Name      [12]byte
	Apuntador int32
}

//BloqueDeArchivos info
type BloqueDeArchivos struct {
	Contenido [64]byte
}

//BloqueDeApuntadores info
type BloqueDeApuntadores struct {
	Apuntadores [16]int32
}

//Journaling info
type Journaling struct {
	TipoOperacion [2]byte
	Tipo          [1]byte
	Nombre        [12]byte
	Contenido     [64]byte
	Fecha         [16]byte
	Propietario   [16]byte
	Permisos      int64
}

//TamanioDelSB asda
func TamanioDelSB() {
	mbr := BloqueDeCarpeta{}
	fmt.Println(int64(unsafe.Sizeof(mbr)))
	mbr2 := BloqueDeArchivos{}
	fmt.Println(int64(unsafe.Sizeof(mbr2)))
	mbr3 := BloqueDeApuntadores{}
	fmt.Println(int64(unsafe.Sizeof(mbr3)))

	fmt.Println("\n superbloque")
	mbr4 := SuperBloque{}
	fmt.Println(int64(unsafe.Sizeof(mbr4)))
	fmt.Println("\n inodo")
	mbr5 := Inodo{}
	fmt.Println(int64(unsafe.Sizeof(mbr5)))
	fmt.Println("\n Journaling")
	mbr6 := Journaling{}
	fmt.Println(int64(unsafe.Sizeof(mbr6)))
}

//EjecutarMKFS info
func EjecutarMKFS(parametros []string, descripciones []string) {
	fmt.Println(" === COMANDO MKFS === ")
	var id, tipo string
	for i := 0; i < len(parametros); i++ {
		switch parametros[i] {
		case "id":
			{
				id = descripciones[i]
			}
		case "type":
			{
				tipo = descripciones[i]
			}
		default:
			{
				fmt.Println("Parametros no validos")
			}
		}
	}
	if tipo == "" {
		tipo = "full"
	}
	if id != "" {
		formatearParticion(id)
	} else {
		fmt.Println("Falta el parametro obligatorio -id")
	}
}

func formatearParticion(id string) {
	rutaDelDisco := ""
	nombreDeLaParticion := ""
	for i := 0; i < len(ParticionesMontadas); i++ {
		if id == ParticionesMontadas[i].id {
			rutaDelDisco = ParticionesMontadas[i].ruta
			nombreDeLaParticion = ParticionesMontadas[i].nombre
			break
		}
	}
	if rutaDelDisco != "" && nombreDeLaParticion != "" {
		//aqui va a suceder la magia
		mbr := LeerMBR(rutaDelDisco)

		//me retorna el numero de particion que se va a utilizar
		numP := encontrarPartition(mbr, nombreDeLaParticion)

		if numP == 1 {
			darFormatoInicial(mbr.MbrPartition1, rutaDelDisco)
		} else if numP == 2 {
			darFormatoInicial(mbr.MbrPartition2, rutaDelDisco)
		} else if numP == 3 {
			darFormatoInicial(mbr.MbrPartition3, rutaDelDisco)
		} else if numP == 4 {
			darFormatoInicial(mbr.MbrPartition4, rutaDelDisco)
		} else {
			fmt.Println("No se encontro la particion")
		}

	}
}

func encontrarPartition(mbr MBR, nombre string) int {
	if mbr.MbrPartition1.PartStatus[0] == 1 { // si esta activa la particion 1
		nombreParticion := ""

		//saco el nombre de la particion
		for i := 0; i < 16; i++ {
			if mbr.MbrPartition1.PartName[i] == 0 {

			} else {
				nombreParticion += string(mbr.MbrPartition1.PartName[i])
			}
		}
		if nombreParticion == nombre { //retorno 1 si encuentro la particion
			return 1
		}
	}

	if mbr.MbrPartition2.PartStatus[0] == 1 {
		nombreParticion := ""

		for i := 0; i < 16; i++ {
			if mbr.MbrPartition2.PartName[i] == 0 {

			} else {
				nombreParticion += string(mbr.MbrPartition2.PartName[i])
			}
		}
		if nombreParticion == nombre {
			return 2
		}
	}

	if mbr.MbrPartition3.PartStatus[0] == 1 {
		nombreParticion := ""

		for i := 0; i < 16; i++ {
			if mbr.MbrPartition3.PartName[i] == 0 {

			} else {
				nombreParticion += string(mbr.MbrPartition3.PartName[i])
			}
		}
		if nombreParticion == nombre {
			return 3
		}
	}

	if mbr.MbrPartition4.PartStatus[0] == 1 {
		nombreParticion := ""

		for i := 0; i < 16; i++ {
			if mbr.MbrPartition4.PartName[i] == 0 {

			} else {
				nombreParticion += string(mbr.MbrPartition4.PartName[i])
			}
		}
		if nombreParticion == nombre {
			return 4
		}
	}
	return 0
}

func darFormatoInicial(p Partition, ruta string) {
	//lo primero que tengo que hacer es encontrar N

	superBLoque := SuperBloque{}
	var sizeSP int64 = int64(unsafe.Sizeof(superBLoque))

	journaling := Journaling{}
	var sizeJour int64 = int64(unsafe.Sizeof(journaling))

	bloque := BloqueDeCarpeta{}
	var sizeBloque int64 = int64(unsafe.Sizeof(bloque))

	inodoCarpeta := Inodo{}
	var sizeInodo int64 = int64(unsafe.Sizeof(inodoCarpeta))

	n := float64(p.PartSize-sizeSP) / float64(sizeJour+1+3+sizeInodo+3*sizeBloque)
	N := math.Floor(n)
	fmt.Print("La N seria bro: ")
	fmt.Println(N)

	//f := SuperBloque{}

	superBLoque.TipoDeSistemaDeArchivos = 0

	superBLoque.NumTotalDeInodos = int64(N)
	superBLoque.NumTotalDeBloques = 3 * int64(N)

	superBLoque.NumDeBloquesLibres = (3 * int64(N)) - 1
	superBLoque.NumDeInodosLibres = (int64(N)) - 1

	today := time.Now()
	var fecha [16]byte
	for i := 0; i < 16; i++ {
		fecha[i] = today.String()[i]
	}
	superBLoque.FechaMontada = fecha
	superBLoque.FechaDesmontada = fecha

	superBLoque.ContadorDeMontadas++
	superBLoque.SMagic = 0xEF53

	superBLoque.TamanioDelInodo = sizeInodo
	superBLoque.TamanioDelBloque = sizeBloque

	superBLoque.PrimerInodoLibre = 1
	superBLoque.PrimerBloqueLibre = 1

	superBLoque.StartBMdeInodos = p.PartStart + sizeSP + int64(N)*sizeJour
	superBLoque.StartBMdeBloques = p.PartStart + sizeSP + int64(N)*sizeJour + int64(N)
	superBLoque.StartTablaDeInodos = p.PartStart + sizeSP + int64(N)*sizeJour + int64(N) + 3*int64(N)
	superBLoque.StartTablaDeBloques = p.PartStart + sizeSP + int64(N)*sizeJour + int64(N) + 3*int64(N) + int64(N)*sizeInodo

	file, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error al abrir el disco")
	} else {
		file.Seek(p.PartStart, 0)
		err = binary.Write(file, binary.LittleEndian, superBLoque)
		if err != nil {
			log.Fatalln(err)
		} else {
			fmt.Println("El superbloque se ha creado y escrito")
		}
	}

	spf := LeerSuperBloque(ruta, p.PartStart)

	fmt.Println(spf.StartBMdeInodos)
	fmt.Println(spf.StartBMdeBloques)
	fmt.Println(spf.StartTablaDeInodos)
	fmt.Println(spf.StartTablaDeBloques)
	fmt.Println()
	TamanioDelSB()
}

//LeerSuperBloque lee
func LeerSuperBloque(ruta string, seek int64) SuperBloque {
	file, err := os.Open(ruta)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	file.Seek(seek, 0)

	SP := SuperBloque{}
	err = binary.Read(file, binary.LittleEndian, &SP)
	if err != nil {
		log.Fatalln(err)
	}

	return SP
}
