package funcs

import (
	"fmt"
	"strconv"
	"strings"
)

//UsuarioActual usario actualmente logueado
var UsuarioActual string = ""

//LaRutaDelDisco guardala la ruta del disco del usuario logueado
var LaRutaDelDisco string = ""

//NombreDeLaPartition guarda el nombre de la particion del usuario logueado
var NombreDeLaPartition string = ""

//NumUsuario UId del usuario (su numero)
var NumUsuario int64

//vamos a creaer un metodo que me saque a todos los usuarios
//quitamos los ?
//hacemos split por /n
//hacemos split por ,
//hacer global un USER y un ID(vda1)

//EjecutarLogin ejecuta el comando mkdir (crear carpetas)
func EjecutarLogin(parametros []string, descripciones []string) {
	fmt.Println("\n\n\n\n === EJECUTANDO COMANDO LOGIN === ")
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
		case "id":
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
		validarIDLogin(id, usr, password)
	} else {
		fmt.Println("Error al intentar hacer el login")
	}
}

func validarIDLogin(id string, usuario string, password string) {

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
			hacerLogin(mbr.MbrPartition1, usuario, password, rutaDelDisco, nombreDeLaParticion)
		} else if numP == 2 {
			hacerLogin(mbr.MbrPartition2, usuario, password, rutaDelDisco, nombreDeLaParticion)
		} else if numP == 3 {
			hacerLogin(mbr.MbrPartition3, usuario, password, rutaDelDisco, nombreDeLaParticion)
		} else if numP == 4 {
			hacerLogin(mbr.MbrPartition4, usuario, password, rutaDelDisco, nombreDeLaParticion)
		} else {
			fmt.Println("No se encontro la particion")
		}

	} else {
		fmt.Println("No hay un usuario logueado")
	}

}

func hacerLogin(p Partition, usuario string, password string, ruta string, nombrePart string) {

	if UsuarioActual == "" { // si no hay usuario logueado
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
			if len(pedazo) == 5 { // si es un pedazo de U	-> [1 U root root 123]
				if pedazo[1] == "U" {
					if pedazo[3] == usuario {
						if pedazo[4] == password {
							fmt.Println("Usuario reconocido " + pedazo[3] + " " + pedazo[4] + "-")
							NumUsuario, _ = strconv.ParseInt(pedazo[0], 10, 64)
							UsuarioActual = pedazo[3]
							LaRutaDelDisco = ruta
							NombreDeLaPartition = nombrePart
							encontrado = true
							break
						} else {
							fmt.Println("PASWORD INCORRECTO")
						}
					}
				}
			}
		}
		if encontrado == false {
			fmt.Println("El usuario no invalido en el sistema de archivos")
		}
	} else {
		fmt.Println("YA EXISTE UN USUARIO LOGUEADO")
	}
}

//Me retorna el numero de grupos y el numero de usuarios
func contarUsGr(elementos []string) (int, int) {
	//def := "1,G,root\n1,U,root,root,123\n"
	contadorGr := 0
	contadorUs := 0
	for i := 0; i < len(elementos[i]); i++ {
		log := elementos[i]
		corte := strings.Split(log, ",")

		if len(corte) == 3 {
			fmt.Println(len(corte))
			contadorGr++
		} else if len(corte) == 5 {
			fmt.Println(len(corte))
			contadorUs++
		}
	}
	return contadorGr, contadorUs
}
