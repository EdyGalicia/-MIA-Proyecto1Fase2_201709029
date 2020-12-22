package funcs

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

//EjecutarMKFILE ejecuta el comando mkdir (crear carpetas)
func EjecutarMKFILE(parametros []string, descripciones []string) {
	fmt.Println("\n\n\n\n === EJECUTANDO COMANDO MK FILE === ")
	var id, ruta, size string
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
		case "size":
			{
				size = descripciones[i]
			}
		default:
			{
				fmt.Println("Error en los parametros")
				hayError = true
			}
		}
	}
	if hayError == false {

		if size == "" {
			size = "0"
		}

		carpetas := strings.Split(ruta, "/")
		//voy a revisar el ID
		validarIDArch(carpetas, id, p, size)

	}
}

func validarIDArch(carpetas []string, id string, p bool, size string) {

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
			checar(mbr.MbrPartition1, rutaDelDisco, carpetas, id, p, size)
		} else if numP == 2 {
			checar(mbr.MbrPartition2, rutaDelDisco, carpetas, id, p, size)
		} else if numP == 3 {
			checar(mbr.MbrPartition3, rutaDelDisco, carpetas, id, p, size)
		} else if numP == 4 {
			checar(mbr.MbrPartition4, rutaDelDisco, carpetas, id, p, size)
		} else {
			fmt.Println("No se encontro la particion")
		}

	}

}

func checar(partition Partition, ruta string, carpetas []string, id string, p bool, size string) {
	fmt.Println("Ya estoy en checar")

	//sp := LeerSuperBloque(ruta, partition.PartStart)
	//var pos int64 = 0
	//buscarEspacioLibreEnBloqueCarpetaArch(ruta, pos, sp, carpetas[1], partition, size)

	//saco el directo de carpetas
	dir := make([]string, (len(carpetas) - 2), (len(carpetas) - 2))
	c := 0
	for i := 1; i < len(carpetas)-1; i++ {
		dir[c] = carpetas[i]
		c++
	}
	archivo := carpetas[len(carpetas)-1]
	fmt.Println(dir)
	letsgo := false
	var aux int64 = 0
	var pos int64 = 0
	for i := 0; i < len(dir); i++ {

		sp := LeerSuperBloque(ruta, partition.PartStart)
		aux = pos
		pos = Buscar(ruta, pos, dir[i], sp)

		if pos == -1 { //si no encuentra lo que buscamos
			if p == true {
				fmt.Println("No se encontro la carpeta, se creara")
				buscarEspacioLibreEnBloqueCarpeta(ruta, aux, sp, dir[i], partition)
				sp = LeerSuperBloque(ruta, partition.PartStart)
				pos = sp.PrimerInodoLibre - int64(1)
				letsgo = true
				//fmt.Println("))))))))))))))))))))))))))))))))))))))))))))))))))))")
				//fmt.Println(pos)
			} else {
				fmt.Println("No se encontro el padre " + dir[i] + " y -p viene falso")
				letsgo = false
				break
			}
		} else {
			letsgo = true
		}
		aux = pos
	}

	if letsgo == true || len(dir) == 0 {
		sp := LeerSuperBloque(ruta, partition.PartStart)
		buscarEspacioLibreEnBloqueCarpetaArch(ruta, aux, sp, archivo, partition, size)
	}
}

func buscarEspacioLibreEnBloqueCarpetaArch(ruta string, posicionEstructura int64, sp SuperBloque, nombre string, partition Partition, size string) {
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
							crearInodoBloqueArchivo(ruta, partition, sp, size)

							fmt.Println("Se ha creado el inodoArchivo")
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
					EscribirByteBM(ruta, sp.StartBMdeBloques+sp.PrimerBloqueLibre, 1)
					//ahora actualizo el superBloque
					sp.NumDeBloquesLibres--
					sp.PrimerBloqueLibre++
					EscribirSuperBloque(ruta, partition, sp) //-------------------------------

					//ahora si va todo lo demas
					crearInodoBloqueArchivo(ruta, partition, sp, size)
					fmt.Println("Se ha creado el inodoArchivo")
					return
				}
			} else if i == 13 {
				if inodo.IBlock[i] != -1 {
					posSeek := calcularPosicionDeBloqueEnElArchivo(inodo.IBlock[i], sp)
					bloqueDeApun := leerElBloqueDeApuntadores(ruta, posSeek)

					for i := 0; i < len(bloqueDeApun.Apuntadores); i++ {

						if bloqueDeApun.Apuntadores[i] != -1 {

							posSeek := calcularPosicionDeBloqueEnElArchivo(int64(bloqueDeApun.Apuntadores[i]), sp)
							blCarpeta := leerElBloqueCarpeta(ruta, posSeek)

							//recorro el bloqueCarpeta
							for j := 0; j < len(blCarpeta.BContent); j++ {

								if blCarpeta.BContent[j].Apuntador == -1 {
									//lleno celdas y sobrescribo
									copy(blCarpeta.BContent[j].Name[:], nombre)
									blCarpeta.BContent[j].Apuntador = int32(sp.PrimerInodoLibre) //lo tengo que mandar el crearPack
									EscribirBloqueCarpeta(ruta, posSeek, blCarpeta)

									//creo el nuevo inodo
									crearInodoBloqueArchivo(ruta, partition, sp, size)
									fmt.Println("Se ha creado el inodoArchivo")
									return
								}
							}
						} else if bloqueDeApun.Apuntadores[i] == -1 {
							//tengo que actualizar el bloque de apuntadores
							bloqueDeApun.Apuntadores[i] = int32(sp.PrimerBloqueLibre)
							EscribirBloqueApuntadores(ruta, posSeek, bloqueDeApun)
							//TOCA CREAR EL BLOQUE CARPETA
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
							EscribirByteBM(ruta, sp.StartBMdeBloques+sp.PrimerBloqueLibre, 1)
							//ahora actualizo el superBloque
							sp.NumDeBloquesLibres--
							sp.PrimerBloqueLibre++
							EscribirSuperBloque(ruta, partition, sp) //-------------------------------

							//
							//ahora si va todo lo demas
							//creo el nuevo inodo
							crearInodoBloqueArchivo(ruta, partition, sp, size)

							fmt.Println("Se ha creado el inodoArchivo")
							return
						}
					}

				} else if inodo.IBlock[i] == -1 {
					fmt.Println("\n\n YA ENTRO AL -1 en el elseif ==-1 && == 13")
					//actualizo el inodo
					inodo.IBlock[i] = sp.PrimerBloqueLibre
					EscribirInodo(ruta, posicionEstructuraArch, inodo)

					//=======================================creo el primer bloque de apuntadores
					bloqueApuntadores := BloqueDeApuntadores{}
					for i := 0; i < len(bloqueApuntadores.Apuntadores); i++ {
						bloqueApuntadores.Apuntadores[i] = -1
					}
					bloqueApuntadores.Apuntadores[0] = int32(sp.PrimerBloqueLibre) + 1
					seekBC := calcularPosicionDeBloqueEnElArchivo(sp.PrimerBloqueLibre, sp)
					EscribirBloqueApuntadores(ruta, seekBC, bloqueApuntadores)
					//actualizar los bitmap--------------------------------------------------
					EscribirByteBM(ruta, sp.StartBMdeBloques+sp.PrimerBloqueLibre, 2)
					//ahora actualizo el superBloque
					sp.NumDeBloquesLibres--
					sp.PrimerBloqueLibre++
					EscribirSuperBloque(ruta, partition, sp) //-------------------------------

					//TOCA CREAR EL BLOQUE CARPETA
					//======================================crear el bloqueCrpeta nuevo
					bloqueCarpeta := BloqueDeCarpeta{}
					for i := 0; i < len(bloqueCarpeta.BContent); i++ {
						bloqueCarpeta.BContent[i].Apuntador = -1
					}
					//lleno celdas
					copy(bloqueCarpeta.BContent[0].Name[:], nombre)
					bloqueCarpeta.BContent[0].Apuntador = int32(sp.PrimerInodoLibre)
					seekBC = calcularPosicionDeBloqueEnElArchivo(sp.PrimerBloqueLibre, sp)
					EscribirBloqueCarpeta(ruta, seekBC, bloqueCarpeta)
					//actualizar los bitmap--------------------------------------------------
					EscribirByteBM(ruta, sp.StartBMdeBloques+sp.PrimerBloqueLibre, 1)
					//ahora actualizo el superBloque
					sp.NumDeBloquesLibres--
					sp.PrimerBloqueLibre++
					EscribirSuperBloque(ruta, partition, sp) //-------------------------------

					//
					//ahora si va todo lo demas
					//creo el nuevo inodo
					crearInodoBloqueArchivo(ruta, partition, sp, size)

					fmt.Println("Se ha creado el inodoArchivo")
					return
				}
			} else if i == 14 {
				if inodo.IBlock[i] != -1 {
					posSeek := calcularPosicionDeBloqueEnElArchivo(inodo.IBlock[i], sp)
					bloqueDeApun := leerElBloqueDeApuntadores(ruta, posSeek)

					for i := 0; i < len(bloqueDeApun.Apuntadores); i++ {

						if bloqueDeApun.Apuntadores[i] != -1 {

							posSeek2 := calcularPosicionDeBloqueEnElArchivo(int64(bloqueDeApun.Apuntadores[i]), sp)
							bloqueDeApun2 := leerElBloqueDeApuntadores(ruta, posSeek2)

							for i := 0; i < len(bloqueDeApun2.Apuntadores); i++ {
								if bloqueDeApun2.Apuntadores[i] != -1 {

									posSeek := calcularPosicionDeBloqueEnElArchivo(int64(bloqueDeApun2.Apuntadores[i]), sp)
									blCarpeta := leerElBloqueCarpeta(ruta, posSeek)

									//recorro el bloqueCarpeta
									for j := 0; j < len(blCarpeta.BContent); j++ {

										if blCarpeta.BContent[j].Apuntador == -1 {
											//lleno celdas y sobrescribo
											copy(blCarpeta.BContent[j].Name[:], nombre)
											blCarpeta.BContent[j].Apuntador = int32(sp.PrimerInodoLibre) //lo tengo que mandar el crearPack
											EscribirBloqueCarpeta(ruta, posSeek, blCarpeta)

											//creo el nuevo inodo
											crearInodoBloqueArchivo(ruta, partition, sp, size)
											fmt.Println("Se ha creado el inodoArchivo")
											return
										}
									}
								} else if bloqueDeApun2.Apuntadores[i] == -1 {
									//tengo que actualizar el bloque de apuntadores
									bloqueDeApun2.Apuntadores[i] = int32(sp.PrimerBloqueLibre)
									EscribirBloqueApuntadores(ruta, posSeek2, bloqueDeApun2)
									//TOCA CREAR EL BLOQUE CARPETA
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
									EscribirByteBM(ruta, sp.StartBMdeBloques+sp.PrimerBloqueLibre, 1)
									//ahora actualizo el superBloque
									sp.NumDeBloquesLibres--
									sp.PrimerBloqueLibre++
									EscribirSuperBloque(ruta, partition, sp) //-------------------------------

									//
									//ahora si va todo lo demas
									//creo el nuevo inodo
									crearInodoBloqueArchivo(ruta, partition, sp, size)

									fmt.Println("Se ha creado el inodoArchivo")
									return
								}
							}
						}
					}
				} else if inodo.IBlock[i] == -1 {
					//actualizo el inodo
					inodo.IBlock[i] = sp.PrimerBloqueLibre
					EscribirInodo(ruta, posicionEstructuraArch, inodo)

					//=======================================creo el primer bloque de apuntadores
					bloqueApuntadores := BloqueDeApuntadores{}
					for i := 0; i < len(bloqueApuntadores.Apuntadores); i++ {
						bloqueApuntadores.Apuntadores[i] = -1
					}
					seekBC := calcularPosicionDeBloqueEnElArchivo(sp.PrimerBloqueLibre, sp)
					seekApuntPadre := seekBC //me sirve para reescribir en el for
					EscribirBloqueApuntadores(ruta, seekBC, bloqueApuntadores)

					//actualizar los bitmap--------------------------------------------------
					EscribirByteBM(ruta, sp.StartBMdeBloques+sp.PrimerBloqueLibre, 2)
					//ahora actualizo el superBloque
					sp.NumDeBloquesLibres--
					sp.PrimerBloqueLibre++
					EscribirSuperBloque(ruta, partition, sp) //-------------------------------

					//empiezo a llenarlo de otros bloques de apuntadores
					var aux int64
					for i := 0; i < len(bloqueApuntadores.Apuntadores); i++ {
						bloqueApuntadores.Apuntadores[i] = int32(sp.PrimerBloqueLibre)
						EscribirBloqueApuntadores(ruta, seekApuntPadre, bloqueApuntadores) //lo reescribo

						//creo el bl interno----------------------------------------------------
						blApInterno := BloqueDeApuntadores{}
						for i := 0; i < len(blApInterno.Apuntadores); i++ {
							blApInterno.Apuntadores[i] = -1
						}
						seekEs := calcularPosicionDeBloqueEnElArchivo(sp.PrimerBloqueLibre, sp)
						if i == 0 {
							aux = seekEs // me va a servir para ir a crear el inodo carpeta
						}
						EscribirBloqueApuntadores(ruta, seekEs, blApInterno) //-------------------
						//actualizar los bitmap--------------------------------------------------
						EscribirByteBM(ruta, sp.StartBMdeBloques+sp.PrimerBloqueLibre, 2)
						//ahora actualizo el superBloque
						sp.NumDeBloquesLibres--
						sp.PrimerBloqueLibre++
						EscribirSuperBloque(ruta, partition, sp) //-------------------------------

					}
					//aqui tengo que editar el blApuntador, blCarpeta, y pack inodoCarpeta
					blAp := leerElBloqueDeApuntadores(ruta, aux)
					blAp.Apuntadores[0] = int32(sp.PrimerBloqueLibre)
					EscribirBloqueApuntadores(ruta, aux, blAp)

					//======================================crear el bloqueCrpeta nuevo
					bloqueCarpeta := BloqueDeCarpeta{}
					for i := 0; i < len(bloqueCarpeta.BContent); i++ {
						bloqueCarpeta.BContent[i].Apuntador = -1
					}
					//lleno celdas
					copy(bloqueCarpeta.BContent[0].Name[:], nombre)
					bloqueCarpeta.BContent[0].Apuntador = int32(sp.PrimerInodoLibre)
					seekBC = calcularPosicionDeBloqueEnElArchivo(sp.PrimerBloqueLibre, sp)
					EscribirBloqueCarpeta(ruta, seekBC, bloqueCarpeta)

					//actualizar los bitmap--------------------------------------------------
					EscribirByteBM(ruta, sp.StartBMdeBloques+sp.PrimerBloqueLibre, 1)
					//ahora actualizo el superBloque
					sp.NumDeBloquesLibres--
					sp.PrimerBloqueLibre++
					EscribirSuperBloque(ruta, partition, sp) //-------------------------------

					//va lo del inodo nuevo
					crearInodoBloqueArchivo(ruta, partition, sp, size)

					fmt.Println("Se ha creado el inodoArchivo")
					return
				}
			}
		}
	}
}

func crearInodoBloqueArchivo(ruta string, partition Partition, sp SuperBloque, size string) {

	//tamanio del archivo
	size64, _ := strconv.ParseInt(size, 10, 64)

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

	inodo.IType[0] = 1 //indico el tipo. 1 archivo 0 carpeta

	for i := 0; i < len(inodo.IPerm); i++ {
		inodo.IPerm[i] = 7
	}

	//lo escribo
	seekInodo := calcularPosicionDelInodoEnElArchivo(sp.PrimerInodoLibre, sp)
	EscribirInodo(ruta, seekInodo, inodo)
	//actualizo el bitmap de inodos
	EscribirByteBM(ruta, sp.StartBMdeInodos+sp.PrimerInodoLibre, 1)
	sp.NumDeInodosLibres--
	sp.PrimerInodoLibre++
	EscribirSuperBloque(ruta, partition, sp)

	fmt.Println("PROCESO A LLENARLO")
	//ahora enpiezo a llenar el inodo segun el numero de bloquesArchivo necesarios
	total := math.Ceil(float64(size64) / 64)

	for i := 0; i < int(total); i++ {
		buscarMenos1(ruta, seekInodo, partition)
	}
}

//busca un -1 para insertar bloques de archivos
func buscarMenos1(ruta string, seekInodo int64, partition Partition) {

	inodo := leerElInodo(ruta, seekInodo)
	sp := LeerSuperBloque(ruta, partition.PartStart)

	//voy a recorrer el inodo para ver donde lo pongo
	for i := 0; i < len(inodo.IBlock); i++ {

		if i >= 0 && i <= 12 {

			if inodo.IBlock[i] == -1 {

				//actualizo el inodo
				inodo.IBlock[i] = sp.PrimerBloqueLibre
				EscribirInodo(ruta, seekInodo, inodo)

				//escribo el bloque archivo
				bloqueArchivo := BloqueDeArchivos{}
				var ct byte = 0
				for i := 0; i < len(bloqueArchivo.Contenido); i++ {
					if ct <= 9 {
						bloqueArchivo.Contenido[i] = ct
					} else {
						ct = 0
						bloqueArchivo.Contenido[i] = ct
					}
					ct++
				}

				posSeek := calcularPosicionDeBloqueEnElArchivo(sp.PrimerBloqueLibre, sp)
				EscribirBloqueArchivo(ruta, posSeek, bloqueArchivo)
				//actualizar los bitmap--------------------------------------------------
				EscribirByteBM(ruta, sp.StartBMdeBloques+sp.PrimerBloqueLibre, 3)
				//ahora actualizo el superBloque
				sp.NumDeBloquesLibres--
				sp.PrimerBloqueLibre++
				EscribirSuperBloque(ruta, partition, sp) //-------------------------------

				cad := ""
				for i := 0; i < len(bloqueArchivo.Contenido); i++ {
					g := int(bloqueArchivo.Contenido[i])
					cad += strconv.Itoa(g)
				}
				fmt.Println(cad)
				fmt.Println("Se ha colocado el bloqueArchivo en un directooooooooooooooooooo")
				return
			}
		} else if i == 13 {
			if inodo.IBlock[i] != -1 {
				//leo bloque de apuntadores
				posSeek := calcularPosicionDeBloqueEnElArchivo(inodo.IBlock[i], sp)
				bloqueDeApun := leerElBloqueDeApuntadores(ruta, posSeek)

				for i := 0; i < len(bloqueDeApun.Apuntadores); i++ {
					if bloqueDeApun.Apuntadores[i] == -1 {
						//tengo que actualizar el bloque de apuntadores
						bloqueDeApun.Apuntadores[i] = int32(sp.PrimerBloqueLibre)
						EscribirBloqueApuntadores(ruta, posSeek, bloqueDeApun)

						//toca crear el bloqueArchivo
						//escribo el bloque archivo
						bloqueArchivo := BloqueDeArchivos{}
						var ct byte = 0
						for i := 0; i < len(bloqueArchivo.Contenido); i++ {
							if ct <= 9 {
								bloqueArchivo.Contenido[i] = ct
							} else {
								ct = 0
								bloqueArchivo.Contenido[i] = ct
							}
							ct++
						}

						posSeek := calcularPosicionDeBloqueEnElArchivo(sp.PrimerBloqueLibre, sp)
						EscribirBloqueArchivo(ruta, posSeek, bloqueArchivo)
						//actualizar los bitmap--------------------------------------------------
						EscribirByteBM(ruta, sp.StartBMdeBloques+sp.PrimerBloqueLibre, 3)
						//ahora actualizo el superBloque
						sp.NumDeBloquesLibres--
						sp.PrimerBloqueLibre++
						EscribirSuperBloque(ruta, partition, sp) //-------------------------------

						fmt.Println("Se ha colocado el bloqueArchivo en un indirecto !=-1")
						return
					}
				}

			} else if inodo.IBlock[i] == -1 {
				//actualizo el inodo
				inodo.IBlock[i] = sp.PrimerBloqueLibre
				EscribirInodo(ruta, seekInodo, inodo)
				//crear bloque de apuntadores
				bloqueApuntadores := BloqueDeApuntadores{}
				for i := 0; i < len(bloqueApuntadores.Apuntadores); i++ {
					bloqueApuntadores.Apuntadores[i] = -1
				}
				bloqueApuntadores.Apuntadores[0] = int32(sp.PrimerBloqueLibre) + 1
				seekBC := calcularPosicionDeBloqueEnElArchivo(sp.PrimerBloqueLibre, sp)
				EscribirBloqueApuntadores(ruta, seekBC, bloqueApuntadores)
				//actualizar los bitmap--------------------------------------------------
				EscribirByteBM(ruta, sp.StartBMdeBloques+sp.PrimerBloqueLibre, 2)
				//ahora actualizo el superBloque
				sp.NumDeBloquesLibres--
				sp.PrimerBloqueLibre++
				EscribirSuperBloque(ruta, partition, sp) //-------------------------------

				//creo el bloqueArchivo
				//escribo el bloque archivo
				bloqueArchivo := BloqueDeArchivos{}
				var ct byte = 0
				for i := 0; i < len(bloqueArchivo.Contenido); i++ {
					if ct <= 9 {
						bloqueArchivo.Contenido[i] = ct
					} else {
						ct = 0
						bloqueArchivo.Contenido[i] = ct
					}
					ct++
				}

				posSeek := calcularPosicionDeBloqueEnElArchivo(sp.PrimerBloqueLibre, sp)
				EscribirBloqueArchivo(ruta, posSeek, bloqueArchivo)
				//actualizar los bitmap--------------------------------------------------
				EscribirByteBM(ruta, sp.StartBMdeBloques+sp.PrimerBloqueLibre, 3)
				//ahora actualizo el superBloque
				sp.NumDeBloquesLibres--
				sp.PrimerBloqueLibre++
				EscribirSuperBloque(ruta, partition, sp) //-------------------------------

				fmt.Println("Se ha colocado el bloqueArchivo en un indirecto == -1")
				return
			}
		} else if i == 14 {
			if inodo.IBlock[i] != -1 {
				posSeek := calcularPosicionDeBloqueEnElArchivo(inodo.IBlock[i], sp)
				bloqueDeApun := leerElBloqueDeApuntadores(ruta, posSeek)

				for i := 0; i < len(bloqueDeApun.Apuntadores); i++ {
					if bloqueDeApun.Apuntadores[i] != -1 {
						posSeek2 := calcularPosicionDeBloqueEnElArchivo(int64(bloqueDeApun.Apuntadores[i]), sp)
						bloqueDeApun2 := leerElBloqueDeApuntadores(ruta, posSeek2)

						for i := 0; i < len(bloqueDeApun2.Apuntadores); i++ {
							if bloqueDeApun2.Apuntadores[i] == -1 {
								//tengo que actualizar el bloque de apuntadores
								bloqueDeApun2.Apuntadores[i] = int32(sp.PrimerBloqueLibre)
								EscribirBloqueApuntadores(ruta, posSeek2, bloqueDeApun2)

								//TOCA CREAR EL BLOQUE archivo
								//escribo el bloque archivo
								bloqueArchivo := BloqueDeArchivos{}
								var ct byte = 0
								for i := 0; i < len(bloqueArchivo.Contenido); i++ {
									if ct <= 9 {
										bloqueArchivo.Contenido[i] = ct
									} else {
										ct = 0
										bloqueArchivo.Contenido[i] = ct
									}
									ct++
								}

								posSeek := calcularPosicionDeBloqueEnElArchivo(sp.PrimerBloqueLibre, sp)
								EscribirBloqueArchivo(ruta, posSeek, bloqueArchivo)
								//actualizar los bitmap--------------------------------------------------
								EscribirByteBM(ruta, sp.StartBMdeBloques+sp.PrimerBloqueLibre, 3)
								//ahora actualizo el superBloque
								sp.NumDeBloquesLibres--
								sp.PrimerBloqueLibre++
								EscribirSuperBloque(ruta, partition, sp) //-------------------------------

								fmt.Println("Se ha colocado el bloqueArchivo en un indirecto dobble != -1")
								return
							}
						}
					}
				}
			} else if inodo.IBlock[i] == -1 {
				//actualizo el inodo
				inodo.IBlock[i] = sp.PrimerBloqueLibre
				EscribirInodo(ruta, seekInodo, inodo)

				//=======================================creo el primer bloque de apuntadores
				bloqueApuntadores := BloqueDeApuntadores{}
				for i := 0; i < len(bloqueApuntadores.Apuntadores); i++ {
					bloqueApuntadores.Apuntadores[i] = -1
				}
				seekBC := calcularPosicionDeBloqueEnElArchivo(sp.PrimerBloqueLibre, sp)
				seekApuntPadre := seekBC //me sirve para reescribir en el for
				EscribirBloqueApuntadores(ruta, seekBC, bloqueApuntadores)
				//actualizar los bitmap--------------------------------------------------
				EscribirByteBM(ruta, sp.StartBMdeBloques+sp.PrimerBloqueLibre, 2)
				//ahora actualizo el superBloque
				sp.NumDeBloquesLibres--
				sp.PrimerBloqueLibre++
				EscribirSuperBloque(ruta, partition, sp) //-------------------------------

				//empiezo a llenarlo de otros bloques de apuntadores
				var aux int64
				for i := 0; i < len(bloqueApuntadores.Apuntadores); i++ {
					bloqueApuntadores.Apuntadores[i] = int32(sp.PrimerBloqueLibre)
					EscribirBloqueApuntadores(ruta, seekApuntPadre, bloqueApuntadores) //lo reescribo

					//creo el bl interno----------------------------------------------------
					blApInterno := BloqueDeApuntadores{}
					for i := 0; i < len(blApInterno.Apuntadores); i++ {
						blApInterno.Apuntadores[i] = -1
					}
					seekEs := calcularPosicionDeBloqueEnElArchivo(sp.PrimerBloqueLibre, sp)
					if i == 0 {
						aux = seekEs // me va a servir para ir a crear el inodo carpeta
					}
					EscribirBloqueApuntadores(ruta, seekEs, blApInterno) //-------------------
					//actualizar los bitmap--------------------------------------------------
					EscribirByteBM(ruta, sp.StartBMdeBloques+sp.PrimerBloqueLibre, 2)
					//ahora actualizo el superBloque
					sp.NumDeBloquesLibres--
					sp.PrimerBloqueLibre++
					EscribirSuperBloque(ruta, partition, sp) //-------------------------------

				}
				//aqui tengo que editar el blApuntador, blCarpeta, y pack inodoCarpeta
				blAp := leerElBloqueDeApuntadores(ruta, aux)
				blAp.Apuntadores[0] = int32(sp.PrimerBloqueLibre)
				EscribirBloqueApuntadores(ruta, aux, blAp)

				//creo el bloquearchivos
				//escribo el bloque archivo
				bloqueArchivo := BloqueDeArchivos{}
				var ct byte = 0
				for i := 0; i < len(bloqueArchivo.Contenido); i++ {
					if ct <= 9 {
						bloqueArchivo.Contenido[i] = ct
					} else {
						ct = 0
						bloqueArchivo.Contenido[i] = ct
					}
					ct++
				}

				posSeek := calcularPosicionDeBloqueEnElArchivo(sp.PrimerBloqueLibre, sp)
				EscribirBloqueArchivo(ruta, posSeek, bloqueArchivo)
				//actualizar los bitmap--------------------------------------------------
				EscribirByteBM(ruta, sp.StartBMdeBloques+sp.PrimerBloqueLibre, 3)
				//ahora actualizo el superBloque
				sp.NumDeBloquesLibres--
				sp.PrimerBloqueLibre++
				EscribirSuperBloque(ruta, partition, sp) //-------------------------------

				fmt.Println("Se ha colocado el bloqueArchivo en un indirecto dobble == -1")
				return
			}
		}
	}
}
