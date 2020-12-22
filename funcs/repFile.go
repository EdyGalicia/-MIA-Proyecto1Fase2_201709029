package funcs

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func crearReporteFile(ruta string, id string, directorioo string) {
	estructura := validarElIDRepFile(id, directorioo)

	aux := strings.Split(ruta, "/")
	nombreV := strings.Split(aux[len(aux)-1], ".") //nombreV [nombre | extension]
	nombre := strings.ReplaceAll(nombreV[0], " ", "")
	directorio := ""
	for i := 0; i < len(aux)-1; i++ {
		directorio += aux[i] + "/"
	}

	err := os.MkdirAll(directorio, 0777)
	if err != nil {
		fmt.Println("Error en la ruta")
	} else {
		file, err1 := os.Create(directorio + "/" + nombre + ".txt")
		defer file.Close()
		if err1 != nil {
			fmt.Println("Error al generar el archivo dot")
		} else {
			_, errr := file.WriteString(string(estructura))
			if errr != nil {
				fmt.Println("Error al querer escribir")
			} else {
				fmt.Println("el dot fue generado")

			}
		}
	}
}

//validarElIDRepBlock me dice en que particion voy a tranajar
func validarElIDRepFile(id string, directorio string) string {
	body := ""
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
			body = generarCuerpoRepFile(mbr.MbrPartition1, rutaDelDisco, directorio)
		} else if numP == 2 {
			body = generarCuerpoRepFile(mbr.MbrPartition2, rutaDelDisco, directorio)
		} else if numP == 3 {
			body = generarCuerpoRepFile(mbr.MbrPartition3, rutaDelDisco, directorio)
		} else if numP == 4 {
			body = generarCuerpoRepFile(mbr.MbrPartition4, rutaDelDisco, directorio)
		} else {
			fmt.Println("No se encontro la particion")
		}

	}
	return body
}

func generarCuerpoRepFile(p Partition, ruta string, directorio string) string {

	dir := strings.Split(directorio, "/")

	var pos int64 = 0
	fmt.Println(dir)
	for i := 1; i < len(dir); i++ {
		sp := LeerSuperBloque(ruta, p.PartStart)
		pos = Buscar(ruta, pos, dir[i], sp)
		if pos == -1 {
			fmt.Println("No se encontro la carpeta/archivo")
			return ""
		}
	}
	return getContenidoDeArchivo(ruta, p, pos)
}

func getContenidoDeArchivo(ruta string, p Partition, numInodo int64) string {

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
					g := int(bloqueArchivo.Contenido[i])
					cad += strconv.Itoa(g)
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
							g := int(bloqueArchivo.Contenido[i])
							cad += strconv.Itoa(g)
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
									g := int(bloqueArchivo.Contenido[i])
									cad += strconv.Itoa(g)
								}
								Cuerpo += cad + " "
							}
						}
					}
				}
			}
		}
	}

	return Cuerpo
}
