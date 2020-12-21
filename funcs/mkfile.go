package funcs

import (
	"math"
	"strconv"
	"time"
	"unsafe"
)

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

	//ahora enpiezo a llenar el inodo segun el numero de bloquesArchivo necesarios
	total := math.Ceil(float64(size64) / 64)

	for i := 0; i < int(total); i++ {

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
				}
			}
		}
	}
}
