package funcs

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func crearReporteBlock(ruta string, id string) {
	estructura := validarElIDRepBlock(id)

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
		file, err1 := os.Create(directorio + "/" + nombre + ".dot")
		defer file.Close()
		if err1 != nil {
			fmt.Println("Error al generar el archivo dot")
		} else {
			_, errr := file.WriteString(string(estructura))
			if errr != nil {
				fmt.Println("Error al querer escribir")
			} else {
				fmt.Println("el dot fue generado")

				err2 := exec.Command("dot", directorio+nombre+".dot", "-o", directorio+nombre+".svg", "-Tsvg").Run()
				if err2 != nil {
					fmt.Println("Error al generar el comando en consola")
				} else {
					err3 := exec.Command("xdg-open", directorio+nombre+".svg").Run()
					if err3 != nil {
						fmt.Println("Error al abrir el reporte")
					} else {
						fmt.Println("Generado correctamente")
					}
				}

			}
		}
	}
}

//validarElIDRepBlock me dice en que particion voy a tranajar
func validarElIDRepBlock(id string) string {
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
			body = generarCuerpoRepBlock(mbr.MbrPartition1, rutaDelDisco)

		} else if numP == 2 {
			body = generarCuerpoRepBlock(mbr.MbrPartition2, rutaDelDisco)
		} else if numP == 3 {
			body = generarCuerpoRepBlock(mbr.MbrPartition3, rutaDelDisco)
		} else if numP == 4 {
			body = generarCuerpoRepBlock(mbr.MbrPartition4, rutaDelDisco)
		} else {
			fmt.Println("No se encontro la particion")
		}

	}
	return body
}

func generarCuerpoRepBlock(p Partition, ruta string) string {

	Cuerpo := "digraph H {rankdir=\"TB\"; \n"

	sp := LeerSuperBloque(ruta, p.PartStart)

	bytesDelBMBloques := leerBytes(ruta, int(sp.NumTotalDeBloques), sp.StartBMdeBloques)

	fmt.Println("YA SE VA A METER AL FOR")
	for i := 0; i < len(bytesDelBMBloques); i++ {
		if bytesDelBMBloques[i] == 0 {
			//fmt.Print(0)
			//no hago tabla, quiere decir que se borro uno en esa posicion
		} else if bytesDelBMBloques[i] == 1 {
			if int64(i) < sp.PrimerBloqueLibre {
				posSeek := calcularPosicionDeBloqueEnElArchivo(int64(i), sp)
				bloqueCarpeta := leerBloqueDeCarpetas(ruta, posSeek)

				iden := "node" + strconv.Itoa(i)

				Cuerpo += iden + "[ \n shape=plaintext \n label=< \n <table border='1' \n cellborder='1'> \n <tr><td colspan=\"2\" BGCOLOR=\"green\" >" + "Bloque Carpeta " + strconv.Itoa(i) + "</td></tr>\n"
				Cuerpo += "<TR><TD>" + "b_Name</TD><TD>" + "b_inodo" + "</TD></TR>\n"

				for i := 0; i < len(bloqueCarpeta.BContent); i++ {
					//saco el nombre
					nombre := ""
					for j := 0; j < len(bloqueCarpeta.BContent[i].Name); j++ {
						if bloqueCarpeta.BContent[i].Name[j] != 0 {
							nombre += string(bloqueCarpeta.BContent[i].Name[j])
						}
					}
					ap := strconv.FormatInt(int64(bloqueCarpeta.BContent[i].Apuntador), 10)

					Cuerpo += "<TR><TD>" + nombre + "</TD><TD>" + ap + "</TD></TR>\n"
				}
				Cuerpo += "</table> >\n ];\n"
			}
		} else if bytesDelBMBloques[i] == 2 { // si es un bloque de apuntadores
			if int64(i) < sp.PrimerBloqueLibre {
				posSeek := calcularPosicionDeBloqueEnElArchivo(int64(i), sp)
				blAp := leerElBloqueDeApuntadores(ruta, posSeek)
				iden := "node" + strconv.Itoa(i)

				Cuerpo += iden + "[ \n shape=plaintext \n label=< \n <table border='1' \n cellborder='1'> \n <tr><td colspan=\"1\" BGCOLOR=\"green\" >" + "Bloque Apuntadores " + strconv.Itoa(i) + "</td></tr>\n"

				for j := 0; j < len(blAp.Apuntadores); j++ {
					ap := strconv.FormatInt(int64(blAp.Apuntadores[j]), 10)
					Cuerpo += "<TR><TD>" + ap + "</TD></TR>\n"
				}
				Cuerpo += "</table> >\n ];\n"
			}
		} else if bytesDelBMBloques[i] == 3 { // si es un bloque de archivos
			if int64(i) < sp.PrimerBloqueLibre {
				posSeek := calcularPosicionDeBloqueEnElArchivo(int64(i), sp)
				bloqueArchivo := leerBloqueDeArchivos(ruta, posSeek)

				iden := "node" + strconv.Itoa(i)

				Cuerpo += iden + "[ \n shape=plaintext \n label=< \n <table border='1' \n cellborder='1'> \n <tr><td colspan=\"1\" BGCOLOR=\"green\" >" + "Bloque Archivo " + strconv.Itoa(i) + "</td></tr>\n"

				cad := ""
				for i := 0; i < len(bloqueArchivo.Contenido); i++ {
					g := int(bloqueArchivo.Contenido[i])
					cad += strconv.Itoa(g)
				}

				Cuerpo += "<TR><TD>" + cad + "</TD></TR>\n"
				Cuerpo += "</table> >\n ];\n"
			}
		}
	}
	Cuerpo += "\n}"
	//fmt.Println("\n\n\n\n\n========================================================================")
	//fmt.Println(Cuerpo)
	//LeerInodo(ruta, sp.StartTablaDeInodos)
	return Cuerpo
}
