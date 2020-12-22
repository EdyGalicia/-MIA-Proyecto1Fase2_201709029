package funcs

import (
	"fmt"
	"os"
	"strings"
)

func crearReporteBmBlock(ruta string, id string) {
	estructura := validarElIDRepBmBlock(id)

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
func validarElIDRepBmBlock(id string) string {
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
			body = generarCuerpoRepBmBlock(mbr.MbrPartition1, rutaDelDisco)
		} else if numP == 2 {
			body = generarCuerpoRepBmBlock(mbr.MbrPartition2, rutaDelDisco)
		} else if numP == 3 {
			body = generarCuerpoRepBmBlock(mbr.MbrPartition3, rutaDelDisco)
		} else if numP == 4 {
			body = generarCuerpoRepBmBlock(mbr.MbrPartition4, rutaDelDisco)
		} else {
			fmt.Println("No se encontro la particion")
		}

	}
	return body
}

func generarCuerpoRepBmBlock(p Partition, ruta string) string {

	Cuerpo := ""

	sp := LeerSuperBloque(ruta, p.PartStart)

	bytesDelBMInodos := leerBytes(ruta, int(sp.NumTotalDeBloques), sp.StartBMdeBloques)

	fmt.Println("YA SE VA A METER AL FOR&&&&&")
	count := 0
	for i := 0; i < len(bytesDelBMInodos); i++ {

		if bytesDelBMInodos[i] == 0 {
			if count == 19 {
				Cuerpo += "0\n"
				count = -1
			} else {
				Cuerpo += "0 "
			}
		} else if bytesDelBMInodos[i] == 1 || bytesDelBMInodos[i] == 2 || bytesDelBMInodos[i] == 3 {
			if count == 19 {
				Cuerpo += "1\n"
				count = -1
			} else {
				//if bytesDelBMInodos[i] == 3 {
				//	Cuerpo += "3 "
				//} else {
				Cuerpo += "1 "
				//}

			}
		}
		count++
	}
	fmt.Println(sp.NumTotalDeBloques)
	LeerSuperBloque(ruta, p.PartStart)
	return Cuerpo
}
