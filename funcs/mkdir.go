package funcs

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"unsafe"
)

//EjecutarMKDIR ejecuta el comando mkdir (crear carpetas)
func EjecutarMKDIR(parametros []string, descripciones []string) {
	fmt.Println("\n\n\n\n === EJECUTANDO COMANDO MKDIR === ")
	var id, ruta string
	var p, hayError bool
	for i := 0; i < len(parametros); i++ {
		switch parametros[i] {
		case "path":
			{
				ruta = descripciones[i]
			}
		case "id":
			{
				id = descripciones[i]
			}
		case "p":
			{
				p = true
			}
		default:
			{
				fmt.Println("Error en los parametros")
				hayError = true
			}
		}
	}
	if hayError == false {

		carpetas := strings.Split(ruta, "/")
		//voy a revisar el ID
		validarID(carpetas, id, p)

	}
}

func validarID(carpetas []string, id string, p bool) {

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
			checarDireccion(mbr.MbrPartition1, rutaDelDisco, carpetas, id, p)
		} else if numP == 2 {
			checarDireccion(mbr.MbrPartition2, rutaDelDisco, carpetas, id, p)
		} else if numP == 3 {
			checarDireccion(mbr.MbrPartition3, rutaDelDisco, carpetas, id, p)
		} else if numP == 4 {
			checarDireccion(mbr.MbrPartition4, rutaDelDisco, carpetas, id, p)
		} else {
			fmt.Println("No se encontro la particion")
		}

	}

}

func checarDireccion(partition Partition, ruta string, carpetas []string, id string, p bool) {
	fmt.Println("Ya estoy en checarDireccion")
	//aqui tengo que hacer lo de id de las mierdas que estan montadas
	sp := LeerSuperBloque(ruta, partition.PartStart)

	var pos int64 = 0
	for i := 0; i < len(carpetas); i++ {
		aux := pos
		pos = Buscar(ruta, pos, carpetas[i], sp)
		if pos == -1 {
			//llamo al metodo crear
			fmt.Println("No se encontro la carpeta, se creara")
			buscarEspacioLibreEnBloqueCarpeta(ruta, aux, sp, carpetas[i], partition)
		}
	}
}

//Buscar me va a buscar lo que le mande en origen y me retorna la posicion en el bitmap de inodos
func Buscar(ruta string, posicionEstructura int64, origen string, sp SuperBloque) int64 {

	//tomamos el inodo
	posicionEstructuraArch := calcularPosicionDelInodoEnElArchivo(posicionEstructura, sp)
	inodo := leerElInodo(ruta, posicionEstructuraArch)

	if inodo.IType[0] == 0 { //si es un inodoCarpeta
		var posInodo int64 = -1
		//voy a recorrer los directos
		for i := 0; i < len(inodo.IBlock) && posInodo == -1; i++ {
			if i >= 0 && i <= 12 { //apuntadores directos
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
								//posInodo = calcularPosicionDelInodoEnElArchivo(int64(blCarpeta.BContent[j].Apuntador), sp)
								posInodo = int64(blCarpeta.BContent[j].Apuntador)
								fmt.Println(posInodo)
								return posInodo
							}
						}
					}

					//
				}
			} else if i == 13 {
				if inodo.IBlock[i] != -1 {
					posSeek := calcularPosicionDeBloqueEnElArchivo(inodo.IBlock[i], sp)
					blApuntadores := leerElBloqueDeApuntadores(ruta, posSeek)

					//voy a recorrer los apuntadores
					for i := 0; i < len(blApuntadores.Apuntadores); i++ {
						if blApuntadores.Apuntadores[i] != -1 {

							//tengo que leer bloques de carpetas
							posSeek := calcularPosicionDeBloqueEnElArchivo(int64(blApuntadores.Apuntadores[i]), sp)
							blCarpeta := leerElBloqueCarpeta(ruta, posSeek)

							//recorro los apuntadores del bloque carpeta
							for j := 0; j < len(blCarpeta.BContent); j++ {

								if blCarpeta.BContent[j].Apuntador != -1 {

									temporal := string(blCarpeta.BContent[j].Name[:len(origen)])
									fmt.Println(origen + "-" + temporal + "-")

									if origen == temporal {

										fmt.Println("LO ENCONTRE en un indirecto")
										//posInodo = calcularPosicionDelInodoEnElArchivo(int64(blCarpeta.BContent[j].Apuntador), sp)
										posInodo = int64(blCarpeta.BContent[j].Apuntador)
										fmt.Println(posInodo)
										return posInodo
									}
								}
							}
						}
					}
				}
			} else { //sino, es 14, doble indirecto
				if inodo.IBlock[i] != -1 {
					posSeek := calcularPosicionDeBloqueEnElArchivo(inodo.IBlock[i], sp)
					blApuntadores := leerElBloqueDeApuntadores(ruta, posSeek)

					//voy a recorrer los apuntadores
					for i := 0; i < len(blApuntadores.Apuntadores); i++ {
						if blApuntadores.Apuntadores[i] != -1 {

							//creo el otro bloque de apuntdores
							posSeek2 := calcularPosicionDeBloqueEnElArchivo(int64(blApuntadores.Apuntadores[i]), sp)
							blApuntadores2 := leerElBloqueDeApuntadores(ruta, posSeek2)

							//recorro el bloque de apuntadores mas profundo
							for i := 0; i < len(blApuntadores2.Apuntadores); i++ {
								if blApuntadores2.Apuntadores[i] != -1 {

									//leo el bloque carpeta
									seekCarpeta := calcularPosicionDeBloqueEnElArchivo(int64(blApuntadores2.Apuntadores[i]), sp)
									blCarpeta := leerElBloqueCarpeta(ruta, seekCarpeta)
									//recorro el bloqueCaroeta
									for j := 0; j < len(blCarpeta.BContent); j++ {
										if blCarpeta.BContent[j].Apuntador != -1 {
											temporal := string(blCarpeta.BContent[j].Name[:len(origen)])
											fmt.Println(origen + "-" + temporal + "-")

											if origen == temporal {
												fmt.Println("LO ENCONTRE en un doble indirecto")
												//posInodo = calcularPosicionDelInodoEnElArchivo(int64(blCarpeta.BContent[j].Apuntador), sp)
												posInodo = int64(blCarpeta.BContent[j].Apuntador)
												fmt.Println(posInodo)
												return posInodo
											}

										}
									}

								}
							}
						}
					}
				}
			}
		}
		fmt.Println("No se encontro la carpeta")
		fmt.Println(posInodo)
		return posInodo
	} else if inodo.IType[0] == 1 { //si es un inodoArchivo

	}

	return 0
}

func buscarEspacioLibreEnBloqueCarpeta(ruta string, posicionEstructura int64, sp SuperBloque, nombre string, partition Partition) {
	//tomamos el inodo
	posicionEstructuraArch := calcularPosicionDelInodoEnElArchivo(posicionEstructura, sp)
	inodo := leerElInodo(ruta, posicionEstructuraArch)
	//iasd := Inodo{}
	if inodo.IType[0] == 0 {
		for i := 0; i < len(inodo.IBlock); i++ {
			if i >= 0 && i <= 12 { //busco en los directos
				if inodo.IBlock[i] != -1 {

					posSeek := calcularPosicionDeBloqueEnElArchivo(inodo.IBlock[i], sp)
					blCarpeta := leerElBloqueCarpeta(ruta, posSeek)

					//recorro el bloqueCarpeta
					for j := 0; j < len(blCarpeta.BContent); j++ {

						if blCarpeta.BContent[j].Apuntador == -1 {
							//lleno celdas
							copy(blCarpeta.BContent[j].Name[:], nombre)
							blCarpeta.BContent[j].Apuntador = int32(sp.PrimerInodoLibre) //lo tengo que mandar el crearPack
							EscribirBloqueCarpeta(ruta, posSeek, blCarpeta)

							//creo el nuevo inodo
							nuevoInodoBloqueCarpeta(ruta, posicionEstructura, sp) // posicion estructura es el padre

							//actualizar los bitmap----------------------------------------------
							EscribirByteBM(ruta, sp.StartBMdeInodos+sp.PrimerInodoLibre)
							EscribirByteBM(ruta, sp.StartBMdeBloques+sp.PrimerBloqueLibre)

							//ahora actualizo el superBloque
							sp.NumDeInodosLibres--
							sp.NumDeBloquesLibres--
							sp.PrimerInodoLibre++
							sp.PrimerBloqueLibre++
							EscribirSuperBloque(ruta, partition, sp) //-----------------------------

							fmt.Println("Se ha creado la carpeta /home")
							return
						}
					}

					//
				} else if inodo.IBlock[i] == -1 {
					//creo el bloqueCarpeta nuevo y despues lleno celda
					//actualizo el padre
					inodo.IBlock[i] = sp.PrimerBloqueLibre
					EscribirInodo(ruta, posicionEstructuraArch, inodo)
					//======================================crear el bloqueCrpeta nuevo
					bloqueCarpeta := BloqueDeCarpeta{}
					for i := 0; i < len(bloqueCarpeta.BContent); i++ {
						bloqueCarpeta.BContent[i].Apuntador = -1
					}
					//lleno celdas
					copy(bloqueCarpeta.BContent[0].Name[:], nombre)
					bloqueCarpeta.BContent[0].Apuntador = int32(sp.PrimerInodoLibre)
					seekBC := calcularPosicionDeBloqueEnElArchivo(sp.PrimerBloqueLibre, sp)
					EscribirBloqueCarpeta(ruta, seekBC, bloqueCarpeta)

					//actualizar los bitmap--------------------------------------------------
					EscribirByteBM(ruta, sp.StartBMdeBloques+sp.PrimerBloqueLibre)
					//ahora actualizo el superBloque
					sp.NumDeBloquesLibres--
					sp.PrimerBloqueLibre++
					EscribirSuperBloque(ruta, partition, sp) //-------------------------------

					//
					//ahora si va todo lo demas
					//creo el nuevo inodo
					nuevoInodoBloqueCarpeta(ruta, posicionEstructura, sp) // posicion estructura es el padre

					//actualizar los bitmap----------------------------------------------
					EscribirByteBM(ruta, sp.StartBMdeInodos+sp.PrimerInodoLibre)
					EscribirByteBM(ruta, sp.StartBMdeBloques+sp.PrimerBloqueLibre)

					//ahora actualizo el superBloque
					sp.NumDeInodosLibres--
					sp.NumDeBloquesLibres--
					sp.PrimerInodoLibre++
					sp.PrimerBloqueLibre++
					EscribirSuperBloque(ruta, partition, sp) //-----------------------------

					fmt.Println("Se ha creado una nuevecita")
					LeerInodo(ruta, sp.StartTablaDeInodos)
					fmt.Println("de bloques")
					f := leerBytes(ruta, 5, sp.StartBMdeBloques)
					for i := 0; i < 5; i++ {
						if f[i] == 0 {
							fmt.Print(0)
						} else if f[i] == 1 {
							fmt.Print(1)
						}
					}
					fmt.Println("de de inodos")
					ff := leerBytes(ruta, 5, sp.StartBMdeInodos)
					for i := 0; i < 5; i++ {
						if ff[i] == 0 {
							fmt.Print(0)
						} else if ff[i] == 1 {
							fmt.Print(1)
						}
					}
					fmt.Println()

					vb := leerBloqueDeCarpetas(ruta, sp.StartTablaDeBloques)
					if vb != bloqueCarpeta {

					}
					return
				}
			}
		}
	}
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

func leerElBloqueDeApuntadores(ruta string, seek int64) BloqueDeApuntadores {
	fmt.Println(" LEYENDO BLOQUE DE APUNTADORES")

	file, err := os.Open(ruta)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	file.Seek(seek, 0)
	blApuntadores := BloqueDeApuntadores{}
	err = binary.Read(file, binary.LittleEndian, &blApuntadores)
	if err != nil {
		log.Fatalln(err)
	}
	for i := 0; i < len(blApuntadores.Apuntadores); i++ {
		fmt.Println(blApuntadores.Apuntadores[i])
	}
	return blApuntadores
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

func nuevoInodoBloqueCarpeta(ruta string, padre int64, sp SuperBloque) {

	//lo creo
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

	inodo.IBlock[0] = sp.PrimerBloqueLibre //esto debe ser el primerBloqueLibre

	inodo.IType[0] = 0 //indico que es carpeta
	for i := 0; i < len(inodo.IPerm); i++ {
		inodo.IPerm[i] = 7
	}
	//lo escribo
	seekInodo := calcularPosicionDelInodoEnElArchivo(sp.PrimerInodoLibre, sp)
	EscribirInodo(ruta, seekInodo, inodo)

	//======================================crear el bloqueCrpeta nuevo
	bloqueCarpeta := BloqueDeCarpeta{}
	for i := 0; i < len(bloqueCarpeta.BContent); i++ {
		bloqueCarpeta.BContent[i].Apuntador = -1
	}
	//el mismo
	var nombre [12]byte
	nombreString := "."
	for i := 0; i < len(nombreString); i++ {
		nombre[i] = nombreString[i]
	}
	bloqueCarpeta.BContent[0].Name = nombre
	bloqueCarpeta.BContent[0].Apuntador = int32(sp.PrimerInodoLibre)
	//el padre
	var nombre2 [12]byte
	nombreString2 := ".."
	for i := 0; i < len(nombreString2); i++ {
		nombre2[i] = nombreString2[i]
	}
	bloqueCarpeta.BContent[1].Name = nombre2
	bloqueCarpeta.BContent[1].Apuntador = int32(padre)

	seekBC := calcularPosicionDeBloqueEnElArchivo(sp.PrimerBloqueLibre, sp)
	EscribirBloqueCarpeta(ruta, seekBC, bloqueCarpeta)
}

//EscribirByteBM escribe un byte en el bitmap
func EscribirByteBM(ruta string, seek int64) { // recibe el tamanio del archivo

	file, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error al abrir el disco en EscribirByte")
	} else {
		file.Seek(seek, 0)
		tamanio := 1
		var a []byte
		a = make([]byte, tamanio, tamanio)
		for i := 0; i < tamanio; i++ {
			a[i] = 1
		}

		_, err = file.Write(a)

		if err != nil {
			log.Fatal(err)
		}
	}

}
