package funcs

import (
	"fmt"
	"strings"
)

//EjecutarRen renombra un archivo o carpeta
func EjecutarRen(parametros []string, descripciones []string) {
	fmt.Println("\n\n\n\n === EJECUTANDO COMANDO REN === ")
	var path, name string
	var hayError bool

	for i := 0; i < len(parametros); i++ {
		switch parametros[i] {
		case "path":
			{
				path = descripciones[i]
			}
		case "name":
			{
				name = descripciones[i]
			}
		default:
			{
				fmt.Println("Error en los parametros")
				hayError = true
			}
		}
	}
	if hayError == false && path != "" && name != "" {
		validarIDRen(path, name)
	} else {
		fmt.Println("Error al intentar hacer el login")
	}
}

func validarIDRen(path string, renombre string) {

	//aca tenemos que hacer unos cambios por las variables globales
	if LaRutaDelDisco != "" && NombreDeLaPartition != "" {
		//aqui va a suceder la magia
		mbr := LeerMBR(LaRutaDelDisco)

		//me retorna el numero de particion que se va a utilizar
		numP := encontrarPartition(mbr, NombreDeLaPartition)

		if numP == 1 {
			checarRen(mbr.MbrPartition1, LaRutaDelDisco, path, renombre)
		} else if numP == 2 {
			checarRen(mbr.MbrPartition2, LaRutaDelDisco, path, renombre)
		} else if numP == 3 {
			checarRen(mbr.MbrPartition3, LaRutaDelDisco, path, renombre)
		} else if numP == 4 {
			checarRen(mbr.MbrPartition4, LaRutaDelDisco, path, renombre)
		} else {
			fmt.Println("No se encontro la particion")
		}

	} else {
		fmt.Println("No hay un usuario logueado")
	}

}

func checarRen(partition Partition, ruta string, path string, renombre string) {
	var pos int64 = 0
	carpetas := strings.Split(path, "/")

	for i := 1; i < len(carpetas); i++ {
		sp := LeerSuperBloque(ruta, partition.PartStart)
		//fmt.Println("))))))))))))))))))))))))))))))))))))))))))))))))))))voy a ir a buscar" + carpetas[i] + " en")
		//fmt.Println(pos)
		if i == len(carpetas)-1 {
			pos = BuscarRen(ruta, pos, carpetas[i], sp, renombre, 1) //1 si lo debe cambiar
		} else {
			pos = BuscarRen(ruta, pos, carpetas[i], sp, renombre, 0) //0 si no lo cambia
		}

		if pos == -1 {
			fmt.Println("No se ubico la ruta para hacer comando REN")
			break
		}

	}
}

//BuscarRen me va a buscar lo que le mande en origen y me retorna la posicion en el bitmap de inodos
func BuscarRen(ruta string, posicionEstructura int64, origen string, sp SuperBloque, nuevoNombre string, ban int) int64 {

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

								if ban == 1 {
									for i := 0; i < len(blCarpeta.BContent[j].Name); i++ {
										blCarpeta.BContent[j].Name[i] = 0
									}
									copy(blCarpeta.BContent[j].Name[:], nuevoNombre)
									EscribirBloqueCarpeta(ruta, posSeek, blCarpeta)
								}

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

										if ban == 1 {
											//
											for i := 0; i < len(blCarpeta.BContent[j].Name); i++ {
												blCarpeta.BContent[j].Name[i] = 0
											}
											copy(blCarpeta.BContent[j].Name[:], nuevoNombre)
											EscribirBloqueCarpeta(ruta, posSeek, blCarpeta)
											//
										}

										return posInodo
									}
								}
							}
						}
					}
				}
			} else if i == 14 { //sino, es 14, doble indirecto
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

												if ban == 1 {
													//
													for i := 0; i < len(blCarpeta.BContent[j].Name); i++ {
														blCarpeta.BContent[j].Name[i] = 0
													}
													copy(blCarpeta.BContent[j].Name[:], nuevoNombre)
													EscribirBloqueCarpeta(ruta, seekCarpeta, blCarpeta)
													//
												}

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
