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
	IPerm  [9]byte
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

	superBLoque.TipoDeSistemaDeArchivos = 3

	superBLoque.NumTotalDeInodos = int64(N)
	superBLoque.NumTotalDeBloques = 3 * int64(N)

	superBLoque.NumDeBloquesLibres = (3 * int64(N)) - 14
	superBLoque.NumDeInodosLibres = (int64(N)) - 2

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

	superBLoque.PrimerInodoLibre = 2
	superBLoque.PrimerBloqueLibre = 14

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

	file.Seek(superBLoque.StartBMdeInodos, 0)
	EscribirByteEnBitmap(file)

	file.Seek(superBLoque.StartBMdeBloques, 0)
	EscribirByteEnBitmap(file)

	//creo el inodo raiz
	inodo := Inodo{}
	inodo.IUid = 1
	inodo.IGid = 1
	inodo.ISize = sizeInodo
	inodo.IAtime = fecha //fecha en que se leyo sin modificarlo
	inodo.ICtime = fecha //fecha creacion
	inodo.IMtime = fecha //fecha modificacion
	for i := 0; i < len(inodo.IBlock); i++ {
		inodo.IBlock[i] = -1
	}
	inodo.IBlock[0] = 0 //el primero a puntada al bloqueCarpeta0
	inodo.IType[0] = 0  //indico que es carpeta
	for i := 0; i < len(inodo.IPerm); i++ {
		inodo.IPerm[i] = 7
	}
	//lo escribo
	file.Seek(superBLoque.StartTablaDeInodos, 0)
	err = binary.Write(file, binary.LittleEndian, inodo)
	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println("El inodoRaiz se ha creado y escrito")
	}

	LeerInodo(ruta, superBLoque.StartTablaDeInodos)

	//creo el bloque carpeta
	bloqueCarpeta := BloqueDeCarpeta{}

	for i := 0; i < len(bloqueCarpeta.BContent); i++ {
		bloqueCarpeta.BContent[i].Apuntador = -1
	}

	var nombre [12]byte
	nombreString := "."
	for i := 0; i < len(nombreString); i++ {
		nombre[i] = nombreString[i]
	}
	bloqueCarpeta.BContent[0].Name = nombre
	bloqueCarpeta.BContent[0].Apuntador = 0

	var nombre2 [12]byte
	nombreString2 := ".."
	for i := 0; i < len(nombreString2); i++ {
		nombre2[i] = nombreString2[i]
	}
	bloqueCarpeta.BContent[1].Name = nombre2
	bloqueCarpeta.BContent[1].Apuntador = 0

	//esta
	var nombre23 [12]byte
	nombreString23 := "users.txt"
	for i := 0; i < len(nombreString23); i++ {
		nombre23[i] = nombreString23[i]
	}
	bloqueCarpeta.BContent[2].Name = nombre23
	bloqueCarpeta.BContent[2].Apuntador = 1
	//sssssssssssssssssss
	//lo escribo
	file.Seek(superBLoque.StartTablaDeBloques, 0)
	err = binary.Write(file, binary.LittleEndian, bloqueCarpeta)
	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println("El bloqueCarpeta se ha creado y escrito")
	}
	leerBloqueDeCarpetas(ruta, superBLoque.StartTablaDeBloques)

	EscribirByteBM(ruta, superBLoque.StartBMdeInodos+1, 1)
	inodoUsersTxt(superBLoque, ruta)
	//Buscar(ruta, superBLoque.StartTablaDeInodos, "usersZ.txt", superBLoque)
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

//EscribirByteEnBitmap escribe un byte en el bitmap
func EscribirByteEnBitmap(file *os.File) { // recibe el tamanio del archivo
	tamanio := 1
	var a []byte
	a = make([]byte, tamanio, tamanio)
	for i := 0; i < tamanio; i++ {
		a[i] = 1
	}

	_, err := file.Write(a)

	if err != nil {
		log.Fatal(err)
	}

}

func leerBytes(ruta string, tam int, seek int64) []byte {
	var a []byte
	aa := int(tam)
	a = make([]byte, tam, tam)
	for i := 0; i < aa; i++ {
		a[i] = 0
	}

	f, err := os.Open(ruta)
	defer f.Close()
	if err != nil {
		log.Fatalln(err)
	} else {

		f.Seek(seek, 0)

		err = binary.Read(f, binary.LittleEndian, a)
		if err != nil {
			log.Fatalln(err)
		}
	}
	return a
}

//LeerInodo lee el inodo de la posicion que le mandemos
func LeerInodo(ruta string, seek int64) {
	fmt.Println(" LEYENDO INODO")
	file, err := os.Open(ruta)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	file.Seek(seek, 0)

	inodo := Inodo{}
	err = binary.Read(file, binary.LittleEndian, &inodo)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(inodo.IUid)
	fmt.Println(inodo.IGid)
	fmt.Println(inodo.ISize)
	fecha := ""
	for i := 0; i < len(inodo.IAtime); i++ {
		fecha += string(inodo.IAtime[i])
	}
	fmt.Println(fecha)
	fecha = ""
	for i := 0; i < len(inodo.ICtime); i++ {
		fecha += string(inodo.ICtime[i])
	}
	fmt.Println(fecha)
	fecha = ""
	for i := 0; i < len(inodo.IMtime); i++ {
		fecha += string(inodo.IMtime[i])
	}
	fmt.Println(fecha)

	for i := 0; i < len(inodo.IBlock); i++ {
		fmt.Println(inodo.IBlock[i])
	}
	fmt.Println(inodo.IType[0])
	perms := ""
	for i := 0; i < len(inodo.IPerm); i++ {
		fmt.Println(inodo.IPerm[i])
	}
	fmt.Println(perms)

	//return inodo
}

func leerBloqueDeArchivos(ruta string, seek int64) BloqueDeArchivos {
	fmt.Println(" LEYENDO BLOQUE DE ARCHIVOS ")

	file, err := os.Open(ruta)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	file.Seek(seek, 0)
	blArchivo := BloqueDeArchivos{}
	err = binary.Read(file, binary.LittleEndian, &blArchivo)
	if err != nil {
		log.Fatalln(err)
	}

	/*cadena := ""
	for i := 0; i < len(blArchivo.Contenido); i++ {
		if blArchivo.Contenido[i] != 0 {
			cadena += string(blArchivo.Contenido[i])
		}
	}
	fmt.Println(cadena)*/

	return blArchivo
}

func leerBloqueDeCarpetas(ruta string, seek int64) BloqueDeCarpeta {
	fmt.Println(" LEYENDO BLOQUE DE CARPETAS")

	file, err := os.Open(ruta)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	file.Seek(seek, 0)
	blCarpeta := BloqueDeCarpeta{}
	err = binary.Read(file, binary.LittleEndian, &blCarpeta)
	if err != nil {
		log.Fatalln(err)
	} else {
		for i := 0; i < len(blCarpeta.BContent); i++ {
			nomb := ""
			for j := 0; j < len(blCarpeta.BContent[i].Name); j++ {
				if blCarpeta.BContent[i].Name[j] != 0 {
					nomb += string(blCarpeta.BContent[i].Name[j])
				}
			}
			fmt.Print("nombre: " + nomb + " apuntador: ")
			fmt.Println(blCarpeta.BContent[i].Apuntador)
		}
	}
	return blCarpeta
}

//EscribirSuperBloque escribe el superbloque
func EscribirSuperBloque(ruta string, p Partition, superBLoque SuperBloque) {
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
}

//EscribirBloqueCarpeta escribe un bloque carpeta en el posicion que le mande
func EscribirBloqueCarpeta(ruta string, seek int64, blCarpeta BloqueDeCarpeta) {
	file, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		fmt.Println("Error al abrir el disco")
	} else {
		file.Seek(seek, 0)
		err = binary.Write(file, binary.LittleEndian, blCarpeta)
		if err != nil {
			log.Fatalln(err)
		} else {
			fmt.Println("El bloqueCarpeta se ha creado y escrito")
		}
	}

}

//EscribirBloqueApuntadores escribe un bloque de apuntadores en el posicion que le mande
func EscribirBloqueApuntadores(ruta string, seek int64, blAp BloqueDeApuntadores) {
	file, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		fmt.Println("Error al abrir el disco")
	} else {
		file.Seek(seek, 0)
		err = binary.Write(file, binary.LittleEndian, blAp)
		if err != nil {
			log.Fatalln(err)
		} else {
			fmt.Println("El BloqueApuntadores se ha creado y escrito")
		}
	}

}

//EscribirBloqueArchivo escribe un bloque de archivos en el posicion que le mande
func EscribirBloqueArchivo(ruta string, seek int64, blAr BloqueDeArchivos) {
	file, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		fmt.Println("Error al abrir el disco")
	} else {
		file.Seek(seek, 0)
		err = binary.Write(file, binary.LittleEndian, blAr)
		if err != nil {
			log.Fatalln(err)
		} else {
			fmt.Println("El BloqueApuntadores se ha creado y escrito")
		}
	}

}

//EscribirInodo escribe un inodo en la posicion que se le mande (posicion en el archivo)
func EscribirInodo(ruta string, seek int64, inodo Inodo) {

	file, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		fmt.Println("Error al abrir el disco")
	} else {

		file.Seek(seek, 0)
		err = binary.Write(file, binary.LittleEndian, inodo)
		if err != nil {
			log.Fatalln(err)
		} else {
			fmt.Println("El inodo se ha creado y escrito")
		}
	}

}

func inodoUsersTxt(sp SuperBloque, ruta string) {
	inodo := Inodo{}
	inodo.IUid = 1
	inodo.IGid = 1
	inodo.ISize = int64(unsafe.Sizeof(inodo))

	today := time.Now()
	var fecha [16]byte
	for i := 0; i < 16; i++ {
		fecha[i] = today.String()[i]
	}
	inodo.IAtime = fecha //fecha en que se leyo sin modificarlo
	inodo.ICtime = fecha //fecha creacion
	inodo.IMtime = fecha //fecha modificacion

	for i := 0; i < len(inodo.IBlock); i++ {
		inodo.IBlock[i] = -1
	}

	inodo.IType[0] = 1 //indico el tipo. 1 archivo 0 carpeta

	for i := 0; i < len(inodo.IPerm); i++ {
		inodo.IPerm[i] = 7
	}
	seek := calcularPosicionDelInodoEnElArchivo(1, sp)
	EscribirInodo(ruta, seek, inodo)

	def := "1,G,root\n1,U,root,root,123\n"
	preg := "?"
	ct := 0
	for i := 1; i <= 13; i++ {
		inodo.IBlock[ct] = int64(i)
		blArch := BloqueDeArchivos{}

		if i == 1 {
			for j := 0; j < len(blArch.Contenido); j++ {
				if j < len(def) {
					blArch.Contenido[j] = def[j]
				} else {
					blArch.Contenido[j] = preg[0]
				}
			}
			seekP := calcularPosicionDeBloqueEnElArchivo(int64(i), sp)
			EscribirBloqueArchivo(ruta, seekP, blArch)
			EscribirByteBM(ruta, sp.StartBMdeBloques+int64(i), 3)
		} else {
			for j := 0; j < len(blArch.Contenido); j++ {
				blArch.Contenido[j] = preg[0]
			}
			seekP := calcularPosicionDeBloqueEnElArchivo(int64(i), sp)
			EscribirBloqueArchivo(ruta, seekP, blArch)
			EscribirByteBM(ruta, sp.StartBMdeBloques+int64(i), 3)
		}
		ct++
	}
	EscribirInodo(ruta, seek, inodo)
}
