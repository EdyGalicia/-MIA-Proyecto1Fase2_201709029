package funcs

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func ejecutarExec(ruta string) {
	fmt.Println(" === EJECUTANDO COMANDO EXEC ===")

	archivo, err := os.Open(ruta)
	if err != nil {
		fmt.Println("Error al abrir el archivo")
	} else {

		defer archivo.Close()
		scanner := bufio.NewScanner(archivo)

		i := 1
		for scanner.Scan() {
			linea := scanner.Text()

			if len(linea) != 0 { //si trae algo la linea
				if linea[0] == '#' {
					fmt.Println("\n\nlineaC" + strconv.Itoa(i) + " " + linea)
				} else {
					fmt.Println("\n\nlinea" + strconv.Itoa(i) + " " + linea)
					TipoDeInstruccion(linea)
				}
			} else {
				//fmt.Println("linea vacia")
			}
			i++
		}
		fmt.Println(" === COMANDO EXEC TERMINADO === ")
	}
}
