package funcs

import (
	"fmt"
	"strconv"
	"strings"
)

//EjecutarMKUSR crear usuarios
func EjecutarMKUSR(parametros []string, descripciones []string) {
	fmt.Println("\n\n\n\n === EJECUTANDO COMANDO MkUSR === ")
	var usr, password, id string
	var hayError bool

	for i := 0; i < len(parametros); i++ {
		switch parametros[i] {
		case "usr":
			{
				usr = descripciones[i]
			}
		case "pwd":
			{
				password = descripciones[i]
			}
		case "grp":
			{
				id = descripciones[i]
			}
		default:
			{
				fmt.Println("Error en los parametros")
				hayError = true
			}
		}
	}
	if hayError == false && usr != "" && password != "" && id != "" {
		validarIDmkusr(id, usr, password)
	} else {
		fmt.Println("Error al intentar hacer el mkusr")
	}
}

//el id es del grupo al que va a pertenecer
func validarIDmkusr(id string, usuario string, password string) {

	if UsuarioActual == "root" {
		//aca tenemos que hacer unos cambios por las variables globales
		if LaRutaDelDisco != "" && NombreDeLaPartition != "" {
			//aqui va a suceder la magia
			mbr := LeerMBR(LaRutaDelDisco)

			//me retorna el numero de particion que se va a utilizar
			numP := encontrarPartition(mbr, NombreDeLaPartition)

			if numP == 1 {
				crearUsuario(mbr.MbrPartition1, usuario, LaRutaDelDisco, id, password)
			} else if numP == 2 {
				crearUsuario(mbr.MbrPartition2, usuario, LaRutaDelDisco, id, password)
			} else if numP == 3 {
				crearUsuario(mbr.MbrPartition3, usuario, LaRutaDelDisco, id, password)
			} else if numP == 4 {
				crearUsuario(mbr.MbrPartition4, usuario, LaRutaDelDisco, id, password)
			} else {
				fmt.Println("No se encontro la particion")
			}

		} else {
			fmt.Println("No hay un usuario logueado")
		}
	} else {
		fmt.Println("Debe estar logueado como usaurio root para poder crear usuarios")
	}

}

func crearUsuario(p Partition, user string, ruta string, grupo string, pass string) {
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
	Gencontrado := false
	Uencontrado := false

	//buscamos si existe el grupo
	//def := "1,G,root\n1,U,root,root,123\n"
	for i := 0; i < len(trozos); i++ {
		GU := trozos[i]
		pedazo := strings.Split(GU, ",")
		if len(pedazo) == 3 { // si es un pedazo de G	-> [1 U root]
			if pedazo[1] == "G" {
				if pedazo[2] == grupo {
					fmt.Println("El grupo " + pedazo[2] + " si existe, ahora revisamos usuarios")
					Gencontrado = true
					break
				}
			}
		}
	}

	//
	for i := 0; i < len(trozos) && Gencontrado == true; i++ {
		GU := trozos[i]
		pedazo := strings.Split(GU, ",")
		if len(pedazo) == 5 { // si es un pedazo de U	-> [1 U root root 123]
			if pedazo[1] == "U" {
				if pedazo[3] == user {
					fmt.Println("El usuario " + pedazo[3] + " ya existe")
					Uencontrado = true
				}
			}
		}
	}
	//

	//si se encontro su grupo y es un usuario nuevo...
	if Gencontrado == true && Uencontrado == false {
		fmt.Println("El grupo existe y el usuario no, encontes si se puede crear el USER")

		if len(pass) <= 10 && len(user) <= 10 {
			//busco donde emepezar a escribir
			nuevoU := strconv.Itoa(cU+1) + ",U," + grupo + "," + user + "," + pass + "\n"
			//fmt.Println("Nuevo User: " + nuevoU + "-")
			//fmt.Println(len(nuevoU))

			nuevaC := cadenaCortada
			nuevaC += nuevoU
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
		} else {
			fmt.Println("Se exceden los limetes de caracteres en la password y/o usuario")
		}
	} else {
		fmt.Println("El grupo que desea asignar no existe O el usuario ya existe")
	}
}
