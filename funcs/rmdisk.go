package funcs

import (
	"fmt"
	"os"
)

//EjecutarRMDISK : elimina el disco
func EjecutarRMDISK(ruta string) {
	fmt.Println(" === COMANDO RMDISK === ")

	var opc string
	fmt.Println("Desea eliminar el disco? Y/N")
	fmt.Scanf("%s\n", &opc)

	if opc == "Y" || opc == "y" {
		err := os.Remove(ruta)
		if err != nil {
			fmt.Println("Error al eliminar el disco")
		}

		fmt.Println(" === COMANDO RMDISK TERMINADO ===")
	}
}
