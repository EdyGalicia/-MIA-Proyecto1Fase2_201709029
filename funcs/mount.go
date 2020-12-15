package funcs

import (
	"fmt"
	"strconv"
)

//EjecutarMOUNT monta una particion
func EjecutarMOUNT(parametros []string, descripciones []string) {
	fmt.Println("\n === EJECUTANDO COMANDO MOUNT === ")
	var name, path string
	hayErrores := false
	for i := 0; i < len(parametros); i++ {
		switch parametros[i] {
		case "name":
			{
				name = descripciones[i]
				//fmt.Println("name -------" + name + "-")
			}
		case "path":
			{
				path = descripciones[i]
				//fmt.Println("ruta ------" + path + "-")
			}
		default:
			{
				fmt.Println("Hay error en los parametros")
				hayErrores = true
			}
		}
	}
	if hayErrores == true {
		fmt.Println("...")
	} else {
		//si encuentra una P o E o si encuentra una logica; con ese nombre retornan un 1
		if ExisteParticion(path, name) == 1 || ValidarNombreDelEBR(path, name) == 1 {
			if yaEstaMontada(path, name) == true {
				fmt.Println("La particion que desea montar YA esta montada")
			} else {
				// va a ver si existe la ruta y si si existe, retorna su letra, sino ""
				letra := verificarLetraExistente(path)
				numero := 0

				if letra == "" {
					letra = obtenerLetraDisponile()
				}

				numero = obtenerNumero(letra)

				for i := 0; i < len(ParticionesMontadas); i++ {
					if ParticionesMontadas[i].nombre == "" {
						montada := ParticionMontada{}
						montada.nombre = name
						montada.id = "vd" + letra + strconv.Itoa(numero)
						montada.letra = letra
						montada.numero = numero
						montada.ruta = path
						ParticionesMontadas[i] = montada
						i = len(ParticionesMontadas)
					}
				}
				fmt.Println("LA PARTICION SE HA MONTADO")
				listarParticiones()
			}
		} else {
			fmt.Println("La particion no existe: " + name)
		}
	}
}

func yaEstaMontada(ruta string, nombre string) bool {

	for i := 0; i < len(ParticionesMontadas); i++ {
		if nombre == ParticionesMontadas[i].nombre && ruta == ParticionesMontadas[i].ruta {
			//true si la encontro
			return true
		}
	}

	return false
}

func verificarLetraExistente(ruta string) string {
	for i := 0; i < len(ParticionesMontadas); i++ {
		if ruta == ParticionesMontadas[i].ruta {
			return ParticionesMontadas[i].letra
		}
	}
	return ""
}

func obtenerLetraDisponile() string {
	coleccion := [26]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	letra := ""
	for i := 0; i < len(coleccion); i++ {
		var bandera bool
		bandera = false
		for b := 0; b < len(ParticionesMontadas); b++ {
			if ParticionesMontadas[b].letra == coleccion[i] {
				//paro
				b = len(ParticionesMontadas)
				bandera = true
			}
		}
		if bandera == false {
			letra = coleccion[i]
			i = len(coleccion)
		}
	}
	return letra
}

func obtenerNumero(letra string) int {
	aux := 0
	for i := 1; i < len(ParticionesMontadas); i++ {
		if numeroDIsponible(letra, i) == 0 {

		} else {
			aux = i
			i = len(ParticionesMontadas)
		}
	}
	return aux
}

func numeroDIsponible(letra string, numero int) int {
	for i := 0; i < len(ParticionesMontadas); i++ {
		if letra == ParticionesMontadas[i].letra && numero == ParticionesMontadas[i].numero {
			return 0
		}
	}
	return 1
}

func listarParticiones() {
	fmt.Println("\n === LISTADO DE PARTICIOENS MONTADAS")
	for i := 0; i < len(ParticionesMontadas); i++ {
		if ParticionesMontadas[i].nombre != "" {
			fmt.Println("Nombre: " + ParticionesMontadas[i].nombre + ", Id " + ParticionesMontadas[i].id + ", Ruta: " + ParticionesMontadas[i].ruta)
		}
	}
}
