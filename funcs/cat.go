package funcs

import (
	"fmt"
	"strings"
)

//EjecutarCat lee archivos
func EjecutarCat(parametros []string, descripciones []string) {
	fmt.Println(" === EJECUTANDO COMANDO CAT ===")
	hayError := true
	for i := 0; i < len(parametros); i++ {
		param := parametros[i]
		if param[0] == 'f' && param[1] == 'i' && param[2] == 'l' && param[3] == 'e' {
			hayError = false
		} else {
			hayError = true
			break
		}
	}
	if hayError == true {
		fmt.Println("Error en los parametros del cat")
	} else {
		validarIDCat(parametros, descripciones)
	}
}

func validarIDCat(param []string, desc []string) {

	//aca tenemos que hacer unos cambios por las variables globales
	if LaRutaDelDisco != "" && NombreDeLaPartition != "" {
		//aqui va a suceder la magia
		mbr := LeerMBR(LaRutaDelDisco)

		//me retorna el numero de particion que se va a utilizar
		numP := encontrarPartition(mbr, NombreDeLaPartition)

		if numP == 1 {
			hacerCat(mbr.MbrPartition1, LaRutaDelDisco, param, desc)
		} else if numP == 2 {
			hacerCat(mbr.MbrPartition2, LaRutaDelDisco, param, desc)
		} else if numP == 3 {
			hacerCat(mbr.MbrPartition3, LaRutaDelDisco, param, desc)
		} else if numP == 4 {
			hacerCat(mbr.MbrPartition4, LaRutaDelDisco, param, desc)
		} else {
			fmt.Println("No se encontro la particion")
		}

	} else {
		fmt.Println("No hay un usuario logueado")
	}

}

func hacerCat(partition Partition, ruta string, par []string, desc []string) {

	cuerpo := ""

	for j := 0; j < len(desc); j++ {
		paso := true

		carpetas := strings.Split(desc[j], "/")

		var pos int64 = 0
		for i := 1; i < len(carpetas); i++ {
			sp := LeerSuperBloque(ruta, partition.PartStart)

			pos = Buscar(ruta, pos, carpetas[i], sp)
			if pos == -1 {
				fmt.Println("No se encontro la carpeta, se creara")
				paso = false
				break
			} else {
				fmt.Println("Si se encontro el dir, seguimos buscando")
			}
		}

		if paso == false {
			fmt.Println("Alguna parte no se encontro, no se podra realizar el CAT")
			cuerpo = ""
			break
		} else {
			fmt.Println(" === FIle satisfactorio")
			//llamamos al metodo que lee
			cad := getContenidoDeArchivoLeer(ruta, partition, pos)
			cuerpo += cad
			fmt.Println(cad)
		}
	}

	fmt.Println("Resultado del CAT: ")
	fmt.Println(cuerpo)

}

func getContenidoDeArchivoLeer(ruta string, p Partition, numInodo int64) string {

	Cuerpo := ""
	sp := LeerSuperBloque(ruta, p.PartStart)

	posSeek := calcularPosicionDelInodoEnElArchivo(numInodo, sp)
	inodo := leerElInodo(ruta, posSeek)

	for i := 0; i < len(inodo.IBlock); i++ {
		if i >= 0 && i <= 12 {
			if inodo.IBlock[i] != -1 {
				posSeek := calcularPosicionDeBloqueEnElArchivo(inodo.IBlock[i], sp)
				bloqueArchivo := leerBloqueDeArchivos(ruta, posSeek)
				cad := ""
				for i := 0; i < len(bloqueArchivo.Contenido); i++ {
					//g := int(bloqueArchivo.Contenido[i])
					//cad += strconv.Itoa(g)
					if string(bloqueArchivo.Contenido[i]) == "?" {

					} else {
						cad += string(bloqueArchivo.Contenido[i])
					}

				}
				Cuerpo += cad + " "
			}
		} else if i == 13 {
			if inodo.IBlock[i] != -1 {
				posSeek := calcularPosicionDeBloqueEnElArchivo(inodo.IBlock[i], sp)
				bloqueDeApun := leerElBloqueDeApuntadores(ruta, posSeek)

				for i := 0; i < len(bloqueDeApun.Apuntadores); i++ {
					if bloqueDeApun.Apuntadores[i] != -1 {
						posSeek := calcularPosicionDeBloqueEnElArchivo(int64(bloqueDeApun.Apuntadores[i]), sp)
						bloqueArchivo := leerBloqueDeArchivos(ruta, posSeek)
						cad := ""
						for i := 0; i < len(bloqueArchivo.Contenido); i++ {
							//g := int(bloqueArchivo.Contenido[i])
							//cad += strconv.Itoa(g)
							if string(bloqueArchivo.Contenido[i]) == "?" {

							} else {
								cad += string(bloqueArchivo.Contenido[i])
							}
						}
						Cuerpo += cad + " "
					}
				}
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
								bloqueArchivo := leerBloqueDeArchivos(ruta, posSeek)
								cad := ""
								for i := 0; i < len(bloqueArchivo.Contenido); i++ {
									//g := int(bloqueArchivo.Contenido[i])
									//cad += strconv.Itoa(g)
									if string(bloqueArchivo.Contenido[i]) == "?" {

									} else {
										cad += string(bloqueArchivo.Contenido[i])
									}
								}
								Cuerpo += cad + " "
							}
						}
					}
				}
			}
		}
	}
	fmt.Println("hola bb" + Cuerpo)
	return Cuerpo + "\n"
}
