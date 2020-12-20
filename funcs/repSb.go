package funcs

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func crearReporteSb(ruta string, id string) {
	estructura := validarElIDRepSb(id)

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
func validarElIDRepSb(id string) string {
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
			body = generarCuerpoRepSb(mbr.MbrPartition1, rutaDelDisco)

		} else if numP == 2 {
			body = generarCuerpoRepSb(mbr.MbrPartition2, rutaDelDisco)
		} else if numP == 3 {
			body = generarCuerpoRepSb(mbr.MbrPartition3, rutaDelDisco)
		} else if numP == 4 {
			body = generarCuerpoRepSb(mbr.MbrPartition4, rutaDelDisco)
		} else {
			fmt.Println("No se encontro la particion")
		}

	}
	return body
}

func generarCuerpoRepSb(p Partition, ruta string) string {

	Cuerpo := "digraph H {rankdir=\"TB\"; \n"

	sp := LeerSuperBloque(ruta, p.PartStart)

	Cuerpo += "Sb" + "[ \n shape=plaintext \n label=< \n <table border='1' \n cellborder='1'> \n <tr><td colspan=\"2\" BGCOLOR=\"green\" >" + "SUPER BLOQUE " + "</td></tr>\n"
	Cuerpo += "<TR><TD>NOMBRE</TD><TD>" + "VALOR" + "</TD></TR>\n"
	Cuerpo += "<TR><TD>" + "Fyle System Type</TD><TD>" + strconv.FormatInt(sp.TipoDeSistemaDeArchivos, 10) + "</TD></TR>\n"
	Cuerpo += "<TR><TD>" + "Inodes count</TD><TD>" + strconv.FormatInt(sp.NumTotalDeInodos, 10) + "</TD></TR>\n"
	Cuerpo += "<TR><TD>" + "Blocks count</TD><TD>" + strconv.FormatInt(sp.NumTotalDeBloques, 10) + "</TD></TR>\n"
	Cuerpo += "<TR><TD>" + "Free blocks count</TD><TD>" + strconv.FormatInt(sp.NumDeBloquesLibres, 10) + "</TD></TR>\n"
	Cuerpo += "<TR><TD>" + "Free inodes count</TD><TD>" + strconv.FormatInt(sp.NumDeInodosLibres, 10) + "</TD></TR>\n"

	//fechas
	fecha := ""
	for i := 0; i < len(sp.FechaMontada); i++ {
		if sp.FechaMontada[i] != 0 {
			fecha += string(sp.FechaMontada[i])
		}
	}
	Cuerpo += "<TR><TD>" + "S Mtime</TD><TD>" + fecha + "</TD></TR>\n"
	fecha = ""
	for i := 0; i < len(sp.FechaDesmontada); i++ {
		if sp.FechaDesmontada[i] != 0 {
			fecha += string(sp.FechaDesmontada[i])
		}
	}
	Cuerpo += "<TR><TD>" + "S Utime</TD><TD>" + fecha + "</TD></TR>\n"

	//contador de cuantas veces ha sido montada
	Cuerpo += "<TR><TD>" + "Mount count</TD><TD>" + strconv.FormatInt(sp.ContadorDeMontadas, 10) + "</TD></TR>\n"
	Cuerpo += "<TR><TD>" + "MAGIC</TD><TD>" + strconv.FormatInt(sp.SMagic, 10) + "</TD></TR>\n"
	Cuerpo += "<TR><TD>" + "Inode Size</TD><TD>" + strconv.FormatInt(sp.TamanioDelInodo, 10) + "</TD></TR>\n"
	Cuerpo += "<TR><TD>" + "Block Size</TD><TD>" + strconv.FormatInt(sp.TamanioDelBloque, 10) + "</TD></TR>\n"

	//
	Cuerpo += "<TR><TD>" + "First Inode free</TD><TD>" + strconv.FormatInt(sp.PrimerInodoLibre, 10) + "</TD></TR>\n"
	Cuerpo += "<TR><TD>" + "First BLock free</TD><TD>" + strconv.FormatInt(sp.PrimerBloqueLibre, 10) + "</TD></TR>\n"

	Cuerpo += "<TR><TD>" + "BM inode start</TD><TD>" + strconv.FormatInt(sp.StartBMdeInodos, 10) + "</TD></TR>\n"
	Cuerpo += "<TR><TD>" + "BM block start</TD><TD>" + strconv.FormatInt(sp.StartBMdeBloques, 10) + "</TD></TR>\n"
	Cuerpo += "<TR><TD>" + "Inode start</TD><TD>" + strconv.FormatInt(sp.StartTablaDeInodos, 10) + "</TD></TR>\n"
	Cuerpo += "<TR><TD>" + "Block start</TD><TD>" + strconv.FormatInt(sp.StartTablaDeBloques, 10) + "</TD></TR>\n"

	Cuerpo += "</table> >\n ];\n"

	Cuerpo += "\n}"
	//fmt.Println("\n\n\n\n\n========================================================================")
	//fmt.Println(Cuerpo)
	//LeerInodo(ruta, sp.StartTablaDeInodos)
	return Cuerpo
}
