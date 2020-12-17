package funcs

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"unsafe"
)

func prueba() {
	//hola := Buscar()
	fmt.Println("ya")
}

//Buscar me va a buscar lo que le mande en origen
func Buscar(ruta string, posicionEstructura int64, origen string, sp SuperBloque) int64 {

	//tomamos el inodo
	inodo := leerElInodo(ruta, posicionEstructura)

	if inodo.IType[0] == 0 { //si es un inodoCarpeta
		var posicionBitmap int64 = -1
		//voy a recorrer los directos
		for i := 0; i < len(inodo.IBlock)-2 && posicionBitmap == -1; i++ {
			if inodo.IBlock[i] != -1 {

				posSeek := calcularPosicionDeBloqueEnElArchivo(inodo.IBlock[i], sp)
				blCarpeta := leerElBloqueCarpeta(ruta, posSeek)

				//recorro el bloqueCarpeta
				for j := 0; j < len(blCarpeta.BContent); j++ {
					if blCarpeta.BContent[j].Apuntador != -1 {
						temporal := string(blCarpeta.BContent[j].Name[:len(origen)])
						fmt.Println(origen + "-" + temporal + "-")
						if origen == temporal {
							fmt.Println("LO ENCONTRE")
							posInodo := calcularPosicionDelInodoEnElArchivo(int64(blCarpeta.BContent[j].Apuntador), sp)
							fmt.Println(posInodo)
							return posInodo
						}
					}
				}
			}
		}
		fmt.Println("No se encontro en los directos")
		return posicionBitmap
	} else if inodo.IType[0] == 1 { //si es un inodoArchivo

	}

	return 0
}

func leerElInodo(ruta string, seek int64) Inodo {

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
	return inodo
}

func leerElBloqueCarpeta(ruta string, seek int64) BloqueDeCarpeta {
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
	}
	return blCarpeta
}

//le mando la posicion en el bitmap
func calcularPosicionDeBloqueEnElArchivo(posicion int64, superBloque SuperBloque) int64 {

	bloqueC := BloqueDeCarpeta{}
	pos := superBloque.StartTablaDeBloques + posicion*int64(unsafe.Sizeof(bloqueC))

	return pos
}

//le mando la posicion en el bitmap
func calcularPosicionDelInodoEnElArchivo(posicion int64, superBloque SuperBloque) int64 {

	inodo := Inodo{}
	pos := superBloque.StartTablaDeInodos + posicion*int64(unsafe.Sizeof(inodo))

	return pos
}
