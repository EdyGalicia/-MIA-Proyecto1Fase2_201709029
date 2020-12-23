package funcs

import (
	"fmt"
	"strconv"
	"strings"
)

//EjecutarMKGRP crear usuarios
func EjecutarMKGRP(parametros []string, descripciones []string) {
	fmt.Println("\n\n\n\n === EJECUTANDO COMANDO MKgrp === ")
	var name string
	var hayError bool

	for i := 0; i < len(parametros); i++ {
		switch parametros[i] {
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
	if hayError == false && name != "" {
		validarPartMKGRP(name)
	} else {
		fmt.Println("Error en mkgrp")
	}
}

func validarPartMKGRP(name string) {

	if UsuarioActual == "root" {
		//aca tenemos que hacer unos cambios por las variables globales
		if LaRutaDelDisco != "" && NombreDeLaPartition != "" {
			//aqui va a suceder la magia
			mbr := LeerMBR(LaRutaDelDisco)

			//me retorna el numero de particion que se va a utilizar
			numP := encontrarPartition(mbr, NombreDeLaPartition)

			if numP == 1 {
				crearGrupo(mbr.MbrPartition1, name, LaRutaDelDisco)
			} else if numP == 2 {
				crearGrupo(mbr.MbrPartition2, name, LaRutaDelDisco)
			} else if numP == 3 {
				crearGrupo(mbr.MbrPartition3, name, LaRutaDelDisco)
			} else if numP == 4 {
				crearGrupo(mbr.MbrPartition4, name, LaRutaDelDisco)
			} else {
				fmt.Println("No se encontro la particion")
			}

		} else {
			fmt.Println("No hay un usuario logueado")
		}
	} else {
		fmt.Println("Necesita ser usuario root para crear grupos")
	}

}

func crearGrupo(p Partition, name string, ruta string) {
	sp := LeerSuperBloque(ruta, p.PartStart)

	//leo el inodo de users.txt
	seekI := calcularPosicionDelInodoEnElArchivo(1, sp)
	inodo := leerElInodo(ruta, seekI)

	cadenaUsers := ""

	for i := 0; i < len(inodo.IBlock)-2; i++ {
		seekB := calcularPosicionDeBloqueEnElArchivo(inodo.IBlock[i], sp)
		blA := leerBloqueDeArchivos(ruta, seekB)
		for j := 0; j < len(blA.Contenido); j++ {
			cadenaUsers += string(blA.Contenido[j])
		}
		//cadenaUsers += " "
	}
	//fmt.Println("La cadena es:\n" + cadenaUsers)
	//fmt.Println(len(cadenaUsers))

	cadenaCortada := ""
	for i := 0; i < len(cadenaUsers); i++ {
		if cadenaUsers[i] != '?' {
			cadenaCortada += string(cadenaUsers[i])
		} else {
			break
		}
	}
	//fmt.Println("La cadena cortada es: \n" + cadenaCortada + "-")
	//fmt.Println(len(cadenaCortada))

	trozos := strings.Split(cadenaCortada, "\n")

	cG, cU := contarUsGr(trozos)
	fmt.Println(cG)
	fmt.Println(cU)
	encontrado := false

	//def := "1,G,root\n1,U,root,root,123\n"
	for i := 0; i < len(trozos); i++ {
		GU := trozos[i]
		pedazo := strings.Split(GU, ",")
		if len(pedazo) == 3 { // si es un pedazo de G	-> [1 U root]
			if pedazo[1] == "G" {
				if pedazo[2] == name {
					fmt.Println("El grupo " + pedazo[2] + " ya existe.")
					encontrado = true
					break
				}
			}
		}
	}
	if encontrado == false {
		fmt.Println("El grupo no existe en el sistema de archivos, se creara")

		//busco donde emepezar a escribir
		nuevoG := strconv.Itoa(cG+1) + ",G," + name + "\n"
		fmt.Println("Nuevo grupo: " + nuevoG + "-")
		fmt.Println(len(nuevoG))

		nuevaC := cadenaCortada
		nuevaC += nuevoG
		q := "?"
		for i := 0; i < 832; i++ {
			if len(nuevaC) < 832 {
				nuevaC += string(q[0])
			}
		}
		//fmt.Println("la mera mera seria: " + nuevaC)
		//fmt.Println(len(nuevaC))

		count := 0
		for j := 0; j < len(inodo.IBlock)-2; j++ {
			seekP := calcularPosicionDeBloqueEnElArchivo(inodo.IBlock[j], sp)

			blArch := BloqueDeArchivos{}
			for i := 0; i < len(blArch.Contenido); i++ {
				blArch.Contenido[i] = nuevaC[count]
				count++
			}
			EscribirBloqueArchivo(ruta, seekP, blArch)
		}
	}
}
