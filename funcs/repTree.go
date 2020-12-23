package funcs

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func crearReporteTree(ruta string, id string) {
	estructura := validarElIDRepTree(id)

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

//validarElIdRepInode fsfsd
func validarElIDRepTree(id string) string {
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
			body = generarCuerpoRepTree(mbr.MbrPartition1, rutaDelDisco)

		} else if numP == 2 {
			body = generarCuerpoRepTree(mbr.MbrPartition2, rutaDelDisco)
		} else if numP == 3 {
			body = generarCuerpoRepTree(mbr.MbrPartition3, rutaDelDisco)
		} else if numP == 4 {
			body = generarCuerpoRepTree(mbr.MbrPartition4, rutaDelDisco)
		} else {
			fmt.Println("No se encontro la particion")
		}

	}
	return body
}

func generarCuerpoRepTree(p Partition, ruta string) string {

	Cuerpo := "digraph H {rankdir=\"LR\"; \n"

	sp := LeerSuperBloque(ruta, p.PartStart)

	bytesDelBMInodos := leerBytes(ruta, int(sp.NumTotalDeInodos), sp.StartBMdeInodos)

	fmt.Println("YA SE VA A METER AL FOR")
	for i := 0; i < len(bytesDelBMInodos); i++ {
		if bytesDelBMInodos[i] == 0 {
			//fmt.Print(0)
			//no hago tabla, quiere decir que se borro uno en esa posicion
		} else if bytesDelBMInodos[i] == 1 {
			if int64(i) < sp.PrimerInodoLibre {
				posSeek := calcularPosicionDelInodoEnElArchivo(int64(i), sp)
				inodo := leerElInodo(ruta, posSeek)
				iden := "inode" + strconv.Itoa(i)
				Cuerpo += iden + "[ \n shape=plaintext \n label=< \n <table border='1' \n cellborder='1'> \n <tr><td colspan=\"2\" BGCOLOR=\"green\" >" + "Inodo " + strconv.Itoa(i) + "</td></tr>\n"
				Cuerpo += "<TR><TD>" + "i_uid</TD><TD>" + strconv.FormatInt(inodo.IUid, 10) + "</TD></TR>\n"
				Cuerpo += "<TR><TD>" + "i_gid</TD><TD>" + strconv.FormatInt(inodo.IGid, 10) + "</TD></TR>\n"
				Cuerpo += "<TR><TD>" + "i_size</TD><TD>" + strconv.FormatInt(inodo.ISize, 10) + "</TD></TR>\n"
				fecha := ""
				for i := 0; i < len(inodo.IAtime); i++ {
					if inodo.IAtime[0] != 0 {
						fecha += string(inodo.IAtime[i])
					}
				}
				Cuerpo += "<TR><TD>" + "i_aTime</TD><TD>" + fecha + "</TD></TR>\n"
				fecha = ""
				for i := 0; i < len(inodo.ICtime); i++ {
					if inodo.ICtime[0] != 0 {
						fecha += string(inodo.ICtime[i])
					}
				}
				Cuerpo += "<TR><TD>" + "i_CTime</TD><TD>" + fecha + "</TD></TR>\n"
				fecha = ""
				for i := 0; i < len(inodo.IMtime); i++ {
					if inodo.IMtime[0] != 0 {
						fecha += string(inodo.IMtime[i])
					}
				}
				Cuerpo += "<TR><TD>" + "i_MTime</TD><TD>" + fecha + "</TD></TR>\n"
				for i := 0; i < len(inodo.IBlock); i++ {
					Cuerpo += "<TR><TD>" + "i_block[" + strconv.Itoa(i) + "]</TD><TD " + "port=\"var" + strconv.Itoa(i) + "\">" + strconv.FormatInt(inodo.IBlock[i], 10) + "</TD></TR>\n"
				}
				tipo := ""
				for i := 0; i < len(inodo.IType); i++ {

					g := int(inodo.IType[i])
					tipo += strconv.Itoa(g)

				}
				Cuerpo += "<TR><TD>" + "i_Type</TD><TD>" + tipo + "</TD></TR>\n"
				iPerms := ""
				for i := 0; i < len(inodo.IPerm); i++ {

					g := int(inodo.IPerm[i])
					iPerms += strconv.Itoa(g)

				}
				Cuerpo += "<TR><TD>" + "i_Perms</TD><TD>" + iPerms + "</TD></TR>\n"
				Cuerpo += "</table> >\n ];\n\n"
			}
		}
	}
	//Cuerpo += "\n}"
	Cuerpo += generarCuerpoRepBlockTree(p, ruta)
	Cuerpo += generarEnlacesRepTreeInodos(p, ruta)
	Cuerpo += generarEnlacesRepBlockTree(p, ruta)
	Cuerpo += "\n}"
	//fmt.Println("\n\n\n\n\n========================================================================")
	//fmt.Println(Cuerpo)
	//LeerInodo(ruta, sp.StartTablaDeInodos)
	return Cuerpo
}

func generarCuerpoRepBlockTree(p Partition, ruta string) string {

	Cuerpo := "\n"

	sp := LeerSuperBloque(ruta, p.PartStart)

	bytesDelBMBloques := leerBytes(ruta, int(sp.NumTotalDeBloques), sp.StartBMdeBloques)

	for i := 0; i < len(bytesDelBMBloques); i++ {
		if bytesDelBMBloques[i] == 0 {
			//fmt.Print(0)
			//no hago tabla, quiere decir que se borro uno en esa posicion
		} else if bytesDelBMBloques[i] == 1 { // si es un bloqueCarpeta
			if int64(i) < sp.PrimerBloqueLibre {
				posSeek := calcularPosicionDeBloqueEnElArchivo(int64(i), sp)
				bloqueCarpeta := leerBloqueDeCarpetas(ruta, posSeek)

				iden := "block" + strconv.Itoa(i)

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

					Cuerpo += "<TR><TD>" + nombre + "</TD><TD port=\"var" + strconv.Itoa(i) + "\">" + ap + "</TD></TR>\n"
				}
				Cuerpo += "</table> >\n ];\n"
			}
		} else if bytesDelBMBloques[i] == 2 { // si es un bloque de apuntadores
			if int64(i) < sp.PrimerBloqueLibre {
				posSeek := calcularPosicionDeBloqueEnElArchivo(int64(i), sp)
				blAp := leerElBloqueDeApuntadores(ruta, posSeek)

				iden := "block" + strconv.Itoa(i)

				Cuerpo += iden + "[ \n shape=plaintext \n label=< \n <table border='1' \n cellborder='1'> \n <tr><td colspan=\"1\" BGCOLOR=\"green\" >" + "Bloque Apuntadores " + strconv.Itoa(i) + "</td></tr>\n"

				for j := 0; j < len(blAp.Apuntadores); j++ {
					ap := strconv.FormatInt(int64(blAp.Apuntadores[j]), 10)
					Cuerpo += "<TR><TD port=\"var" + strconv.Itoa(j) + "\">" + ap + "</TD></TR>\n"
				}
				Cuerpo += "</table> >\n ];\n"
			}
		} else if bytesDelBMBloques[i] == 3 { // si es un bloque de archivos
			if int64(i) < sp.PrimerBloqueLibre {
				posSeek := calcularPosicionDeBloqueEnElArchivo(int64(i), sp)
				bloqueArchivo := leerBloqueDeArchivos(ruta, posSeek)

				iden := "block" + strconv.Itoa(i)

				Cuerpo += iden + "[ \n shape=plaintext \n label=< \n <table border='1' \n cellborder='1'> \n <tr><td colspan=\"1\" BGCOLOR=\"green\" >" + "Bloque Archivo " + strconv.Itoa(i) + "</td></tr>\n"

				cad := ""
				for i := 0; i < len(bloqueArchivo.Contenido); i++ {
					//g := int(bloqueArchivo.Contenido[i])
					if string(bloqueArchivo.Contenido[i]) == "\n" {
						cad += "-"
					} else {
						cad += string(bloqueArchivo.Contenido[i])
					}
				}

				Cuerpo += "<TR><TD>" + cad + "</TD></TR>\n"
				Cuerpo += "</table> >\n ];\n"
			}
		}
	}
	//fmt.Println("\n\n\n\n\n========================================================================")
	//fmt.Println(Cuerpo)
	//LeerInodo(ruta, sp.StartTablaDeInodos)
	return Cuerpo
}

func generarEnlacesRepTreeInodos(p Partition, ruta string) string {

	Cuerpo := "\n"

	sp := LeerSuperBloque(ruta, p.PartStart)

	bytesDelBMInodos := leerBytes(ruta, int(sp.NumTotalDeInodos), sp.StartBMdeInodos)

	for i := 0; i < len(bytesDelBMInodos); i++ {
		if bytesDelBMInodos[i] == 0 {
			//fmt.Print(0)
			//no hago tabla, quiere decir que se borro uno en esa posicion
		} else if bytesDelBMInodos[i] == 1 {
			if int64(i) < sp.PrimerInodoLibre {
				posSeek := calcularPosicionDelInodoEnElArchivo(int64(i), sp)
				inodo := leerElInodo(ruta, posSeek)

				iden := "inode" + strconv.Itoa(i)
				for j := 0; j < len(inodo.IBlock); j++ {
					if inodo.IBlock[j] != -1 {
						Cuerpo += iden + ": var" + strconv.Itoa(j) + " -> block" + strconv.FormatInt(inodo.IBlock[j], 10) + ";\n"
					}

				}

			}
		}
	}
	return Cuerpo
}

func generarEnlacesRepBlockTree(p Partition, ruta string) string {

	Cuerpo := "\n"

	sp := LeerSuperBloque(ruta, p.PartStart)

	bytesDelBMBloques := leerBytes(ruta, int(sp.NumTotalDeBloques), sp.StartBMdeBloques)

	for i := 0; i < len(bytesDelBMBloques); i++ {
		if bytesDelBMBloques[i] == 0 {
			//fmt.Print(0)
			//no hago tabla, quiere decir que se borro uno en esa posicion
		} else if bytesDelBMBloques[i] == 1 { // si es un bloqueCarpeta
			if int64(i) < sp.PrimerBloqueLibre {
				posSeek := calcularPosicionDeBloqueEnElArchivo(int64(i), sp)
				bloqueCarpeta := leerBloqueDeCarpetas(ruta, posSeek)

				iden := "block" + strconv.Itoa(i)

				for j := 0; j < len(bloqueCarpeta.BContent); j++ {

					nombre := ""
					for k := 0; k < len(bloqueCarpeta.BContent[j].Name); k++ {
						if bloqueCarpeta.BContent[j].Name[k] != 0 {
							nombre += string(bloqueCarpeta.BContent[j].Name[k])
						}
					}
					if nombre != "." && nombre != ".." {
						if bloqueCarpeta.BContent[j].Apuntador != -1 {
							ap := strconv.FormatInt(int64(bloqueCarpeta.BContent[j].Apuntador), 10)
							Cuerpo += iden + ": var" + strconv.Itoa(j) + " -> inode" + ap + ";\n"
						}
					}

				}
			}
		} else if bytesDelBMBloques[i] == 2 { // si es un bloque de apuntadores
			if int64(i) < sp.PrimerBloqueLibre {
				posSeek := calcularPosicionDeBloqueEnElArchivo(int64(i), sp)
				blAp := leerElBloqueDeApuntadores(ruta, posSeek)

				iden := "block" + strconv.Itoa(i)

				for j := 0; j < len(blAp.Apuntadores); j++ {

					if blAp.Apuntadores[j] != -1 {
						ap := strconv.FormatInt(int64(blAp.Apuntadores[j]), 10)
						Cuerpo += iden + ": var" + strconv.Itoa(j) + " -> block" + ap + ";\n"
					}
				}
			}
		} else if bytesDelBMBloques[i] == 3 { // si es un bloque de archivos

		}
	}
	//fmt.Println("\n\n\n\n\n========================================================================")
	//fmt.Println(Cuerpo)
	//LeerInodo(ruta, sp.StartTablaDeInodos)
	return Cuerpo
}
