package funcs

import (
	"fmt"
)

//EjecutarUNMOUNT quita del arreglo una particion montada
func EjecutarUNMOUNT(parametros []string, descripciones []string) {
	fmt.Println("\n === COMANDO UNMOUNT ===")
	montada := descripciones[0]
	encontrada := false
	for i := 0; i < len(ParticionesMontadas); i++ {
		if montada == ParticionesMontadas[i].id {
			montVacia := ParticionMontada{nombre: "", id: "", letra: "", numero: 0, ruta: ""}
			ParticionesMontadas[i] = montVacia

			fmt.Println("La particion: " + montada + " se ha desmontado")
			encontrada = true
			break
		}
	}
	if encontrada == false {
		fmt.Println("La particion: " + montada + " no se encontro")
	}
	listarParticiones()
}
