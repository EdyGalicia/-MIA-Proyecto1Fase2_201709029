package funcs

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func crearReporteInode(ruta string, id string) {
	estructura := validarElIDRepInode(id)

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
func validarElIDRepInode(id string) string {
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
			body = generarCuerpoRepInode(mbr.MbrPartition1, rutaDelDisco)

		} else if numP == 2 {
			body = generarCuerpoRepInode(mbr.MbrPartition2, rutaDelDisco)
		} else if numP == 3 {
			body = generarCuerpoRepInode(mbr.MbrPartition3, rutaDelDisco)
		} else if numP == 4 {
			body = generarCuerpoRepInode(mbr.MbrPartition4, rutaDelDisco)
		} else {
			fmt.Println("No se encontro la particion")
		}

	}
	return body
}

func generarCuerpoRepInode(p Partition, ruta string) string {

	Cuerpo := "digraph H {rankdir=\"TB\"; \n"

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
				iden := "node" + strconv.Itoa(i)
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
					Cuerpo += "<TR><TD>" + "i_block[" + strconv.Itoa(i) + "]</TD><TD>" + strconv.FormatInt(inodo.IBlock[i], 10) + "</TD></TR>\n"
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
