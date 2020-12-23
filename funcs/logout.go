package funcs

import "fmt"

//EjecutarLogout cierra sesion
func EjecutarLogout() {
	fmt.Println(" === EJECUTANDO COMANDO LOGOUT ===")
	if UsuarioActual != "" {
		UsuarioActual = ""
		LaRutaDelDisco = ""
		NombreDeLaPartition = ""
		fmt.Println("Se ha cerrado sesion correctamente")
	} else {
		fmt.Println("No hay usuario logueado")
	}
}
