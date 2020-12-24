package funcs

import (
	"fmt"
	"strings"
)

//Herramienta separa las cosas
func Herramienta(entrada string) (string, []string, []string) {
	var elementos []string
	vecEntrada := strings.Split(entrada, "#") //separo comentarios
	if len(vecEntrada) == 1 {
		elementos = strings.Split(entrada, " -")
	} else {
		elementos = strings.Split(vecEntrada[0], " -")
	}

	//aqui voy a guardar parametro y descripcion
	parametros := make([]string, (len(elementos) - 1), (len(elementos) - 1))
	descripciones := make([]string, (len(elementos) - 1), (len(elementos) - 1))
	instruccion := strings.ToLower(elementos[0]) //guardo la instruccion
	aux := ""
	indice := 0
	for i := 1; i < len(elementos); i++ {
		aux = elementos[i]
		pareja := strings.Split(aux, "->") //pareja[0]=parametro, pareja[1]=desc

		if len(parametros) != 0 { //si vinieron parametros
			parametros[indice] = strings.ToLower(pareja[0])

			if strings.ToLower(pareja[0]) == "p" {
				descripciones[indice] = ""
			} else {
				auxE := strings.Replace(pareja[1], "”", "", 2)
				descripciones[indice] = strings.Replace(auxE, "\"", "", 2)
			}

			indice++
		}
	}
	//fmt.Println(instruccion)
	//fmt.Println(parametro)
	//fmt.Println(descripcion)

	return instruccion, parametros, descripciones
}

//TipoDeInstruccion ve que onda con eso
func TipoDeInstruccion(entrada string) {
	instruccion, param, descrip := Herramienta(entrada)
	//fmt.Println(param)
	//fmt.Println(descrip)
	switch instruccion {
	case "exec":
		{
			if len(param) == 1 && len(descrip) == 1 {
				if param[0] == "path" {
					ejecutarExec(descrip[0])
				} else {
					fmt.Println("Es necesario el path unicamente")
				}
			}
		}
	case "pause":
		{
			EjecutarPAUSE()
		}
	case "mkdisk":
		{
			if param != nil && descrip != nil {
				EjecutarMKDISK(param, descrip)
			}
		}
	case "rmdisk":
		{
			if param[0] == "path" {
				EjecutarRMDISK(descrip[0])
			} else {
				fmt.Println("Necesita el parametro -path")
			}
		}
	case "fdisk":
		{
			EjecutarFDISK(param, descrip)
		}
	case "mount":
		{
			EjecutarMOUNT(param, descrip)
		}
	case "unmount":
		{
			EjecutarUNMOUNT(param, descrip)
		}
	case "rep":
		{
			fmt.Println("el REP")
			EjecutarREP(param, descrip)
		}
	case "mkfs":
		{
			EjecutarMKFS(param, descrip)
		}
	case "mkdir":
		{
			fmt.Println("capto el mkdir")
			fmt.Println(param)
			fmt.Println(descrip)
			EjecutarMKDIR(param, descrip)
		}
	case "mkfile":
		{
			EjecutarMKFILE(param, descrip)
		}
	case "login":
		{
			EjecutarLogin(param, descrip)
		}
	case "logout":
		{
			EjecutarLogout()
		}
	case "mkgrp":
		{
			EjecutarMKGRP(param, descrip)
		}
	case "mkusr":
		{
			EjecutarMKUSR(param, descrip)
		}
	case "cat":
		{
			EjecutarCat(param, descrip)
		}
	case "ren":
		{
			EjecutarRen(param, descrip)
		}
	case "rem":
		{
			fmt.Println("No se termino comando rem")
		}
	default:
		{
			fmt.Println("El comando es incorrecto.")
		}
	}
}

//ParticionMontada obj que se agrega a al arreglo cuando monto
type ParticionMontada struct {
	nombre string
	id     string
	letra  string
	numero int
	ruta   string
}

//ParticionesMontadas arreglo de particiones montadas
var ParticionesMontadas = [100]ParticionMontada{}
