package funcs

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

//TamanioDelEBRbro asda
func TamanioDelEBRbro() {
	mbr := MBR{}
	fmt.Println(int64(unsafe.Sizeof(mbr)))
	ebr := EBR{}
	fmt.Println(int64(unsafe.Sizeof(ebr)))
}

//EjecutarFDISK ejecuta la instruccion
func EjecutarFDISK(parametros []string, descripciones []string) {
	fmt.Println(" === COMANDO FDISK === ")

	parametrosCorrectos := true
	var tamanio, unidad, ruta, tipo, fit, delete, nombre, add string
	for i := 0; i < len(parametros); i++ {
		switch parametros[i] {
		case "size":
			{
				tamanio = descripciones[i]
				fmt.Println("tamnanio " + tamanio)
			}
		case "unit":
			{
				unidad = descripciones[i]
				unidad = strings.ToLower(unidad)
				fmt.Println("unit " + unidad)
			}
		case "path":
			{
				ruta = descripciones[i]
				fmt.Println("path " + ruta)
			}
		case "type":
			{
				tipo = descripciones[i]
				tipo = strings.ToLower(tipo)
				fmt.Println("type " + tipo)
			}
		case "fit":
			{
				fit = descripciones[i]
				fit = strings.ToLower(fit)
				fmt.Println("fit " + fit)
			}
		case "delete":
			{
				delete = descripciones[i]
				delete = strings.ToLower(delete)
				fmt.Println("delete " + delete)
			}
		case "name":
			{
				nombre = descripciones[i]
				fmt.Println("nombre " + nombre)
			}
		case "add":
			{
				add = descripciones[i]
				fmt.Println("add " + add)
			}
		default:
			{
				fmt.Println("Error en los parametros")
				parametrosCorrectos = false
			}
		}
	}

	//coloco los valores que van por default
	if unidad == "" {
		unidad = "k"
	}
	if tipo == "" {
		tipo = "p"
	}
	if fit == "" {
		fit = "w"
	}

	if add != "" && delete != "" {
		fmt.Println("Los parametros -add y -delete no son compatibles")
	} else if delete != "" && add == "" { // si viene eliminar

		//veo si es una primaria la que se quiere eliminar
		if ExisteParticion(ruta, nombre) == 1 {
			fmt.Println("Confirme la eliminacion de la particion Y/N")
			var confirmacion string
			fmt.Scan(&confirmacion)
			if confirmacion == "Y" || confirmacion == "y" {
				if EliminarPartition(ruta, nombre, delete) == true {
					fmt.Println("La particion se ha eliminado: " + nombre)
				} else {
					fmt.Println("No se elimino la particion")
				}
			}
		} else if ValidarNombreDelEBR(ruta, nombre) == 1 { // si es logica
			fmt.Println("Confirme la eliminacion de la particion Y/N yyyy")
			var confirmacion string
			fmt.Scan(&confirmacion)
			if confirmacion == "Y" || confirmacion == "y" {
				if EliminarPartitionLogica(ruta, nombre, delete) == true {
					fmt.Println("La particion se ha eliminado: " + nombre)
				} else {
					fmt.Println("No se elimino la particion")
				}
			}
		}
	} else if add != "" && delete == "" { // si viene add
		//ADDDDDDDDD

		fmt.Println("\n === COMANDO ADD ===")
		//paso a numero lo que se quiere agregar o quitar
		addNum, _ := strconv.ParseInt(add, 10, 64)
		if unidad == "m" {
			addNum = addNum * 1024 * 1024
		} else if unidad == "k" {
			addNum = addNum * 1024
		}

		//veo que tipo de particion es
		if ExisteParticion(ruta, nombre) == 1 {
			vali := addPartition(ruta, nombre, addNum)
			if vali == 1 {
				fmt.Print("Se agrego el espacio con exito ")
				fmt.Println(addNum)
			} else if vali == 0 {
				fmt.Println("No se puso aumentar o reducir tamanio")
			}
		} else if ValidarNombreDelEBR(ruta, nombre) == 1 {

		} else {
			fmt.Println("No se encontro la particion")
		}

	} else { // CREAR PARTICION
		if parametrosCorrectos == false || nombre == "" || tamanio == "" {
			fmt.Println("Parametros incorrectos")
		} else {
			mbr := MBR{}
			mbr = LeerMBR(ruta) //mbr del disco que voy a modificar
			if tipo == "p" {    // primaria

				nombrecito := ""
				fech := ""
				for i := 0; i < 16; i++ {
					nombrecito += string(mbr.MbrPartition1.PartName[i])
					fech += string(mbr.MbrFechaCreacion[i])
				}
				fmt.Println("El nombre de la particion 1 es: " + nombrecito)
				crearParticion(ruta, nombre, tipo, tamanio, unidad, fit, mbr)
			} else if tipo == "e" { //extendida
				crearParticion(ruta, nombre, tipo, tamanio, unidad, fit, mbr)
			} else if tipo == "l" { //logica
				crearParticionLogica(ruta, mbr, tamanio, unidad, fit, nombre)
			}
		}
	}

}

//LeerMBR va a sacar el mbr del disco que se le mande
func LeerMBR(ruta string) MBR {
	mbr := MBR{}
	file, err := os.Open(ruta)
	defer file.Close()
	if err != nil {
		fmt.Println("Posibles errores en la ruta del disco___")
	} else {
		err := binary.Read(file, binary.LittleEndian, &mbr)
		if err != nil {
			fmt.Println("Error al leer el binario___")
		}
	}
	return mbr
}

func crearParticion(ruta string, nombre string, tipo string, tamanio string, unidad string, fit string, mbr MBR) {
	//veo si existe una particion con el mismo nombre &&
	if ExisteParticion(ruta, nombre) == 0 && ValidarNombreDelEBR(ruta, nombre) == 0 {
		fmt.Println("Paso el if y esta en el metodo de crear particion" + tipo)

		file, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Error al abrir el disco")
		} else {
			//convierto a int64 el tamanio
			tamanio64, _ := strconv.ParseInt(tamanio, 10, 64)

			//tomamos la unidad, si no viene, el default es -> k
			if unidad == "m" {
				tamanio64 = tamanio64 * 1024 * 1024

			} else if unidad == "k" {
				tamanio64 = tamanio64 * 1024
			}

			//capturo la fecha
			today := time.Now()
			var fecha [16]byte
			for i := 0; i < 16; i++ {
				fecha[i] = today.String()[i]
			}

			//si viene una extendida, voy a revisar si ya existe una extendida en el disco
			//retornara 1 si ya existe una extendida
			yaHayExtendida := 0
			if tipo == "e" {
				yaHayExtendida = existeExtendida(ruta)
			}

			if yaHayExtendida == 0 {
				fmt.Println("Aun no existe la particion extendida bro")

				espacioLibre := espacioLibre(ruta)

				if tamanio64 <= espacioLibre {
					fmt.Println("Si hay espacio libre para la particion")

					//procedo a crear la particion
					if mbr.MbrPartition1.PartStatus[0] != 1 {
						mbr.MbrPartition1.PartStatus[0] = 1       //activo la particion
						mbr.MbrPartition1.PartType[0] = tipo[0]   //
						mbr.MbrPartition1.PartFit[0] = fit[0]     //
						elInicio := validarStart(ruta, tamanio64) ////
						mbr.MbrPartition1.PartStart = elInicio    ////
						mbr.MbrPartition1.PartSize = tamanio64    //
						for i := 0; i < len(nombre); i++ {        //
							mbr.MbrPartition1.PartName[i] = nombre[i]
						}

						//veo si pongo el EBR si es una extendida
						if tipo == "e" {

							//hacer el ebr inicial
							mbr.MbrPartition1.PartCont++

							eb := EBR{}
							file.Seek(elInicio, 0)
							var a [1]byte
							a[0] = 2
							eb = EBR{PartStart: elInicio, PartNext: 1 + elInicio + +int64(unsafe.Sizeof(eb)), PartStatus: a}
							err = binary.Write(file, binary.LittleEndian, eb)
							if err != nil {
								log.Fatalln(err)
							}
						}

						file.Seek(0, 0)
						err = binary.Write(file, binary.LittleEndian, mbr)
						if err != nil {
							log.Fatalln(err)
						} else {
							fmt.Println("Particion creada con exito: " + nombre)
						}

					} else if mbr.MbrPartition2.PartStatus[0] != 1 {
						mbr.MbrPartition2.PartStatus[0] = 1       //activo la particion
						mbr.MbrPartition2.PartType[0] = tipo[0]   //
						mbr.MbrPartition2.PartFit[0] = fit[0]     //
						elInicio := validarStart(ruta, tamanio64) ////
						mbr.MbrPartition2.PartStart = elInicio    ////
						mbr.MbrPartition2.PartSize = tamanio64    //
						for i := 0; i < len(nombre); i++ {        //
							mbr.MbrPartition2.PartName[i] = nombre[i]
						}

						//veo si pongo el EBR si es una extendida
						if tipo == "e" {

							//hacer el ebr inicial
							mbr.MbrPartition2.PartCont++

							eb := EBR{}
							file.Seek(elInicio, 0)
							var a [1]byte
							a[0] = 2
							eb = EBR{PartStart: elInicio, PartNext: 1 + elInicio + +int64(unsafe.Sizeof(eb)), PartStatus: a}
							err = binary.Write(file, binary.LittleEndian, eb)
							if err != nil {
								log.Fatalln(err)
							}
						}

						file.Seek(0, 0)
						err = binary.Write(file, binary.LittleEndian, mbr)
						if err != nil {
							log.Fatalln(err)
						} else {
							fmt.Println("Particion creada con exito" + nombre)

						}

					} else if mbr.MbrPartition3.PartStatus[0] != 1 {
						mbr.MbrPartition3.PartStatus[0] = 1       //activo la particion
						mbr.MbrPartition3.PartType[0] = tipo[0]   //
						mbr.MbrPartition3.PartFit[0] = fit[0]     //
						elInicio := validarStart(ruta, tamanio64) ////
						mbr.MbrPartition3.PartStart = elInicio    ////
						mbr.MbrPartition3.PartSize = tamanio64    //
						for i := 0; i < len(nombre); i++ {        //
							mbr.MbrPartition3.PartName[i] = nombre[i]
						}

						//veo si pongo el EBR si es una extendida
						if tipo == "e" {

							//hacer el ebr inicial
							mbr.MbrPartition3.PartCont++

							eb := EBR{}
							file.Seek(elInicio, 0)
							var a [1]byte
							a[0] = 2
							eb = EBR{PartStart: elInicio, PartNext: 1 + elInicio + +int64(unsafe.Sizeof(eb)), PartStatus: a}
							err = binary.Write(file, binary.LittleEndian, eb)
							if err != nil {
								log.Fatalln(err)
							}
						}

						file.Seek(0, 0)
						err = binary.Write(file, binary.LittleEndian, mbr)
						if err != nil {
							log.Fatalln(err)
						} else {
							fmt.Println("Particion creada con exito" + nombre)

						}

					} else if mbr.MbrPartition4.PartStatus[0] != 1 {
						mbr.MbrPartition4.PartStatus[0] = 1       //activo la particion
						mbr.MbrPartition4.PartType[0] = tipo[0]   //
						mbr.MbrPartition4.PartFit[0] = fit[0]     //
						elInicio := validarStart(ruta, tamanio64) ////
						mbr.MbrPartition4.PartStart = elInicio    ////
						mbr.MbrPartition4.PartSize = tamanio64    //
						for i := 0; i < len(nombre); i++ {        //
							mbr.MbrPartition4.PartName[i] = nombre[i]
						}

						//veo si pongo el EBR si es una extendida
						if tipo == "e" {

							//hacer el ebr inicial
							mbr.MbrPartition4.PartCont++

							eb := EBR{}
							file.Seek(elInicio, 0)
							var a [1]byte
							a[0] = 2
							eb = EBR{PartStart: elInicio, PartNext: 1 + elInicio + +int64(unsafe.Sizeof(eb)), PartStatus: a}
							err = binary.Write(file, binary.LittleEndian, eb)
							if err != nil {
								log.Fatalln(err)
							}
						}

						file.Seek(0, 0)
						err = binary.Write(file, binary.LittleEndian, mbr)
						if err != nil {
							log.Fatalln(err)
						} else {
							fmt.Println("Particion creada con exito" + nombre)

						}

					} else {
						fmt.Println("Ya hay 4 particiones")
					}

				} else {
					fmt.Println("No hay espacio suficiente para su particion")
				}

			} else {
				fmt.Println("En el disco ya existe una particion EXTENDIDA")
			}
		}
	} else {
		fmt.Println("Ya existe la particion con ese nombre: " + nombre)
	}
}

//ExisteParticion me retorna 0 si no encuentra la particion (de las 4) con el parametro nombre
func ExisteParticion(ruta string, nombre string) int {

	mbr := MBR{}
	mbr = LeerMBR(ruta)
	//si existe la particion 1
	if mbr.MbrPartition1.PartStatus[0] == 1 { // si esta activa la particion 1
		nombreParticion := ""

		//saco el nombre de la particion
		for i := 0; i < 16; i++ {
			if mbr.MbrPartition1.PartName[i] == 0 {

			} else {
				nombreParticion += string(mbr.MbrPartition1.PartName[i])
			}
		}
		if nombreParticion == nombre { //retorno 1 si encuentro la particion
			return 1
		}
	}

	if mbr.MbrPartition2.PartStatus[0] == 1 {
		nombreParticion := ""

		for i := 0; i < 16; i++ {
			if mbr.MbrPartition2.PartName[i] == 0 {

			} else {
				nombreParticion += string(mbr.MbrPartition2.PartName[i])
			}
		}
		if nombreParticion == nombre {
			return 1
		}
	}

	if mbr.MbrPartition3.PartStatus[0] == 1 {
		nombreParticion := ""

		for i := 0; i < 16; i++ {
			if mbr.MbrPartition3.PartName[i] == 0 {

			} else {
				nombreParticion += string(mbr.MbrPartition3.PartName[i])
			}
		}
		if nombreParticion == nombre {
			return 1
		}
	}

	if mbr.MbrPartition4.PartStatus[0] == 1 {
		nombreParticion := ""

		for i := 0; i < 16; i++ {
			if mbr.MbrPartition4.PartName[i] == 0 {

			} else {
				nombreParticion += string(mbr.MbrPartition4.PartName[i])
			}
		}
		if nombreParticion == nombre {
			return 1
		}
	}

	//si no encuentro el nombre de la particion
	return 0
}

//ValidarNombreDelEBR retorna 1 si encuentra una particion logica con el mismo nombre
func ValidarNombreDelEBR(ruta string, nombre string) int {

	mbr := MBR{}
	mbr = LeerMBR(ruta)

	if mbr.MbrPartition1.PartType[0] == 'e' && mbr.MbrPartition1.PartStatus[0] == 1 {
		//recorro las particiones logicas
		var i int64
		for i = 1; i < mbr.MbrPartition1.PartCont; i++ {
			ebr := EBR{}
			ebr = LeerEBR(ruta, mbr.MbrPartition1.PartStart+48*i+i)

			//convierto a []bytes el nombre
			var nomB [16]byte
			for i := 0; i < len(nombre); i++ {
				nomB[i] = nombre[i]
			}

			if nomB == ebr.PartName && ebr.PartStatus[0] == 1 {
				return 1
			}
		}
	}

	if mbr.MbrPartition2.PartType[0] == 'e' && mbr.MbrPartition2.PartStatus[0] == 1 {
		var i int64
		for i = 1; i < mbr.MbrPartition2.PartCont; i++ {
			ebr := EBR{}
			ebr = LeerEBR(ruta, mbr.MbrPartition2.PartStart+48*i+i)

			//convierto a []bytes el nombre
			var nomB [16]byte
			for i := 0; i < len(nombre); i++ {
				nomB[i] = nombre[i]
			}

			if nomB == ebr.PartName && ebr.PartStatus[0] == 1 {
				return 1
			}
		}
	}

	if mbr.MbrPartition3.PartType[0] == 'e' && mbr.MbrPartition3.PartStatus[0] == 1 {
		var i int64
		for i = 1; i < mbr.MbrPartition3.PartCont; i++ {
			ebr := EBR{}
			ebr = LeerEBR(ruta, mbr.MbrPartition3.PartStart+48*i+i)

			//convierto a []bytes el nombre
			var nomB [16]byte
			for i := 0; i < len(nombre); i++ {
				nomB[i] = nombre[i]
			}

			if nomB == ebr.PartName && ebr.PartStatus[0] == 1 {
				return 1
			}
		}
	}

	if mbr.MbrPartition4.PartType[0] == 'e' && mbr.MbrPartition4.PartStatus[0] == 1 {
		var i int64
		for i = 1; i < mbr.MbrPartition4.PartCont; i++ {
			ebr := EBR{}
			ebr = LeerEBR(ruta, mbr.MbrPartition4.PartStart+48*i+i)

			//convierto a []bytes el nombre
			var nomB [16]byte
			for i := 0; i < len(nombre); i++ {
				nomB[i] = nombre[i]
			}

			if nomB == ebr.PartName && ebr.PartStatus[0] == 1 {
				return 1
			}
		}
	}

	return 0
}

//LeerEBR lee el ebr del disco. Lo uso cuando voy a buscar una particion por su nombre
func LeerEBR(ruta string, seek int64) EBR {
	file, err := os.Open(ruta)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	file.Seek(seek, 0)

	ebr := EBR{}
	err = binary.Read(file, binary.LittleEndian, &ebr)
	if err != nil {
		log.Fatalln(err)
	}

	return ebr
}

//retorna 1 si encuentra una particion extendida
func existeExtendida(ruta string) int {
	mbr := MBR{}
	mbr = LeerMBR(ruta)

	if mbr.MbrPartition1.PartType[0] == 'e' && mbr.MbrPartition1.PartStatus[0] == 1 {
		return 1
	} else if mbr.MbrPartition2.PartType[0] == 'e' && mbr.MbrPartition2.PartStatus[0] == 1 {
		return 1
	} else if mbr.MbrPartition3.PartType[0] == 'e' && mbr.MbrPartition3.PartStatus[0] == 1 {
		return 1
	} else if mbr.MbrPartition4.PartType[0] == 'e' && mbr.MbrPartition4.PartStatus[0] == 1 {
		return 1
	} else {
		return 0
	}
}

//espacio libre entre para particiones primarias o extendidas
func espacioLibre(ruta string) int64 {
	mbr := MBR{}
	mbr = LeerMBR(ruta)

	var espacioOcupado int64

	//sumo las particiones que esten activas
	if mbr.MbrPartition1.PartStatus[0] == 1 {
		espacioOcupado = espacioOcupado + mbr.MbrPartition1.PartSize
	}
	if mbr.MbrPartition2.PartStatus[0] == 1 {
		espacioOcupado = espacioOcupado + mbr.MbrPartition2.PartSize
	}
	if mbr.MbrPartition3.PartStatus[0] == 1 {
		espacioOcupado = espacioOcupado + mbr.MbrPartition3.PartSize
	}
	if mbr.MbrPartition4.PartStatus[0] == 1 {
		espacioOcupado = espacioOcupado + mbr.MbrPartition4.PartSize
	}

	return mbr.MbrTamanio - espacioOcupado
}

func validarStart(ruta string, tamanio int64) int64 {
	mbr := MBR{}
	mbr = LeerMBR(ruta)
	var inicio int64
	inicio = 250
	var i int64
	for i = 250; i < mbr.MbrTamanio; i++ {
		//todas 0 para que pase. Ponen 1 cuando hay clavo
		if validarRango2(mbr.MbrPartition1, mbr.MbrTamanio, i) == 1 || validarRango2(mbr.MbrPartition2, mbr.MbrTamanio, i) == 1 || validarRango2(mbr.MbrPartition3, mbr.MbrTamanio, i) == 1 || validarRango2(mbr.MbrPartition4, mbr.MbrTamanio, i) == 1 {

		} else {
			if validarRango2(mbr.MbrPartition1, mbr.MbrTamanio, i+tamanio) == 1 || validarRango2(mbr.MbrPartition2, mbr.MbrTamanio, i+tamanio) == 1 || validarRango2(mbr.MbrPartition3, mbr.MbrTamanio, i+tamanio) == 1 || validarRango2(mbr.MbrPartition4, mbr.MbrTamanio, i+tamanio) == 1 {
				//cicla hasta que alcance el espacio a lo largo
			} else {
				inicio = i
				i = mbr.MbrTamanio
				fmt.Print("======= E: INICIO SERA en: ")
				fmt.Println(inicio)
			}
		}
	}

	return inicio
}

func validarRango2(particion Partition, tamanio int64, valor int64) int64 {
	if particion.PartStatus[0] == 1 {
		if validarRango(particion.PartStart, particion.PartStart+particion.PartSize, valor) == 1 {
			return 1
		}
	}
	return 0
}

func validarRango(inicio int64, fin int64, valor int64) int64 {
	if valor >= inicio && valor <= fin {
		return 1
	}
	return 0
}

//==================================================================================

func crearParticionLogica(ruta string, mbr MBR, tamanio string, unidad string, fit string, nombre string) {

	//reviso si no existe una particion P, E o L con ese nombre
	if ExisteParticion(ruta, nombre) == 0 && ValidarNombreDelEBR(ruta, nombre) == 0 {

		//checo lo del tamanio de la particion
		tamanio64, _ := strconv.ParseInt(tamanio, 10, 64)
		if unidad == "k" {
			tamanio64 = tamanio64 * 1024
		} else if unidad == "m" {
			tamanio64 = tamanio64 * 1024 * 1024
		}

		if validarEspacioEBR(ruta) >= tamanio64 {

			numDePart, partCont, existeExten := ValidarExtendidaEBR(ruta)
			if existeExten == 1 {
				var seek int64
				file, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
				if err != nil {
					fmt.Println("Error al abrir el disco")
				} else {

					//segun la particion, muevo el seek e incremento el contador de las Logicas
					if numDePart == 1 {
						seek = mbr.MbrPartition1.PartStart + 48*partCont + partCont
						mbr.MbrPartition1.PartCont++
					} else if numDePart == 2 {
						seek = mbr.MbrPartition2.PartStart + 48*partCont + partCont
						mbr.MbrPartition2.PartCont++
					} else if numDePart == 3 {
						seek = mbr.MbrPartition3.PartStart + 48*partCont + partCont
						mbr.MbrPartition3.PartCont++
					} else if numDePart == 4 {
						seek = mbr.MbrPartition4.PartStart + 48*partCont + partCont
						mbr.MbrPartition4.PartCont++
					}

					//escribo el nuevo EBR
					file.Seek(seek, 0)

					//Status
					var a [1]byte
					a[0] = 1

					var fitt [1]byte
					fitt[0] = fit[0]

					var nomb [16]byte
					for i := 0; i < len(nombre); i++ {
						nomb[i] = nombre[i]
					}
					//next seek + 48 + 1
					ebr := EBR{PartStatus: a, PartStart: seek, PartNext: seek + 49, PartName: nomb, PartSize: tamanio64, PartFit: fitt}
					err = binary.Write(file, binary.LittleEndian, ebr)
					if err != nil {
						log.Fatalln(err)
					}

					//escribo el mbr con los nuevos datos de partCont
					file.Seek(0, 0)
					err = binary.Write(file, binary.LittleEndian, mbr)
					if err != nil {
						log.Fatalln(err)
					} else {
						fmt.Println("La particion logica fue creada con exito: " + nombre)
					}
				}
			} else {
				fmt.Println("No se encontro la particion extendida, no se creara la logica")
			}

		} else {
			fmt.Println("La particion logica no se puede crear, espacio insuficiente.")
		}

	} else {
		fmt.Println("Ya existe una particion con ese nombre")
	}
}

//validarEspacioEBR retorna el espacio libre dentro de la extendida
func validarEspacioEBR(ruta string) int64 {
	mbr := MBR{}
	mbr = LeerMBR(ruta)
	var espacioOcupado int64
	//busco la particion extendida, recorro los EBR y sumo sus tamanios, retorno el sobrante
	if mbr.MbrPartition1.PartType[0] == 'e' && mbr.MbrPartition1.PartStatus[0] == 1 {
		var i int64
		for i = 0; i < mbr.MbrPartition1.PartCont; i++ {
			ebr := EBR{}
			ebr = LeerEBR(ruta, mbr.MbrPartition1.PartStart+48*i+i)
			espacioOcupado += ebr.PartSize
		}
		return mbr.MbrPartition1.PartSize - espacioOcupado

	} else if mbr.MbrPartition2.PartType[0] == 'e' && mbr.MbrPartition2.PartStatus[0] == 1 {
		var i int64
		for i = 0; i < mbr.MbrPartition2.PartCont; i++ {
			ebr := EBR{}
			ebr = LeerEBR(ruta, mbr.MbrPartition2.PartStart+48*i+i)
			espacioOcupado += ebr.PartSize
		}
		return mbr.MbrPartition2.PartSize - espacioOcupado

	} else if mbr.MbrPartition3.PartType[0] == 'e' && mbr.MbrPartition3.PartStatus[0] == 1 {
		var i int64
		for i = 0; i < mbr.MbrPartition3.PartCont; i++ {
			ebr := EBR{}
			ebr = LeerEBR(ruta, mbr.MbrPartition3.PartStart+48*i+i)
			espacioOcupado += ebr.PartSize
		}
		return mbr.MbrPartition3.PartSize - espacioOcupado

	} else if mbr.MbrPartition4.PartType[0] == 'e' && mbr.MbrPartition4.PartStatus[0] == 1 {
		var i int64
		for i = 0; i < mbr.MbrPartition4.PartCont; i++ {
			ebr := EBR{}
			ebr = LeerEBR(ruta, mbr.MbrPartition4.PartStart+48*i+i)
			espacioOcupado += ebr.PartSize
		}
		return mbr.MbrPartition4.PartSize - espacioOcupado
	}
	return 0
}

//ValidarExtendidaEBR retorna: #deParticion, cantidadDeLogicas, 1 si encuentra la etendida
//revisa que particion es la extendida
func ValidarExtendidaEBR(ruta string) (int, int64, int) {

	mbr := MBR{}
	mbr = LeerMBR(ruta)

	if mbr.MbrPartition1.PartType[0] == 'e' && mbr.MbrPartition1.PartStatus[0] == 1 {
		//retorno el #1 de particion, la contidad de logicas y un 1 si se encontro la extendida
		return 1, mbr.MbrPartition1.PartCont, 1
	} else if mbr.MbrPartition2.PartType[0] == 'e' && mbr.MbrPartition2.PartStatus[0] == 1 {
		//retorno el #1 de particion, la contidad de logicas y un 1 si se encontro la extendida
		return 2, mbr.MbrPartition2.PartCont, 1
	} else if mbr.MbrPartition3.PartType[0] == 'e' && mbr.MbrPartition3.PartStatus[0] == 1 {
		//retorno el #1 de particion, la contidad de logicas y un 1 si se encontro la extendida
		return 3, mbr.MbrPartition3.PartCont, 1
	} else if mbr.MbrPartition4.PartType[0] == 'e' && mbr.MbrPartition4.PartStatus[0] == 1 {
		//retorno el #1 de particion, la contidad de logicas y un 1 si se encontro la extendida
		return 4, mbr.MbrPartition4.PartCont, 1
	}
	return 0, 0, 0
}

//==================================================================================

//EliminarPartition elimina particion de las 4
func EliminarPartition(ruta string, nombre string, delete string) bool {

	mbr := MBR{}
	mbr = LeerMBR(ruta)
	if mbr.MbrPartition1.PartStatus[0] == 1 {

		//le saco el nombre
		nombreAct := ""
		for i := 0; i < 16; i++ {
			if mbr.MbrPartition1.PartName[i] != 0 {
				nombreAct += string(mbr.MbrPartition1.PartName[i])
			}
		}

		//si si es la particion que busco, elimino
		if nombre == nombreAct {
			file, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println("Error al leer el disco")
			} else {
				file.Seek(0, 0)
				if delete == "fast" {
					mbr.MbrPartition1.PartStatus[0] = 3
				} else if delete == "full" {
					mbr.MbrPartition1.PartStatus[0] = 0
					mbr.MbrPartition1.PartType[0] = 0
					mbr.MbrPartition1.PartFit[0] = 0
					mbr.MbrPartition1.PartStart = 0
					mbr.MbrPartition1.PartSize = 0
					var vacio [16]byte
					mbr.MbrPartition1.PartName = vacio
					mbr.MbrPartition1.PartCont = 0
				} else {
					fmt.Println("Parametro delete incorrecto (fast, full)")
					return false
				}
				//escribo el mbr modificado
				err = binary.Write(file, binary.LittleEndian, mbr)

				if err != nil {
					log.Fatalln(err)
				} else {
					return true
				}
			}
			return false
		}
	}
	if mbr.MbrPartition2.PartStatus[0] == 1 {

		//le saco el nombre
		nombreAct := ""
		for i := 0; i < 16; i++ {
			if mbr.MbrPartition2.PartName[i] != 0 {
				nombreAct += string(mbr.MbrPartition2.PartName[i])
			}
		}

		//si si es la particion que busco, elimino
		if nombre == nombreAct {
			file, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println("Error al leer el disco")
			} else {
				file.Seek(0, 0)
				if delete == "fast" {
					mbr.MbrPartition2.PartStatus[0] = 3
				} else if delete == "full" {
					mbr.MbrPartition2.PartStatus[0] = 0
					mbr.MbrPartition2.PartType[0] = 0
					mbr.MbrPartition2.PartFit[0] = 0
					mbr.MbrPartition2.PartStart = 0
					mbr.MbrPartition2.PartSize = 0
					var vacio [16]byte
					mbr.MbrPartition2.PartName = vacio
					mbr.MbrPartition2.PartCont = 0
				} else {
					fmt.Println("Parametro delete incorrecto (fast, full)")
					return false
				}
				//escribo el mbr modificado
				err = binary.Write(file, binary.LittleEndian, mbr)

				if err != nil {
					log.Fatalln(err)
				} else {
					return true
				}
			}
			return false
		}
	}
	if mbr.MbrPartition3.PartStatus[0] == 1 {

		//le saco el nombre
		nombreAct := ""
		for i := 0; i < 16; i++ {
			if mbr.MbrPartition3.PartName[i] != 0 {
				nombreAct += string(mbr.MbrPartition3.PartName[i])
			}
		}

		//si si es la particion que busco, elimino
		if nombre == nombreAct {
			file, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println("Error al leer el disco")
			} else {
				file.Seek(0, 0)
				if delete == "fast" {
					mbr.MbrPartition3.PartStatus[0] = 3
				} else if delete == "full" {
					mbr.MbrPartition3.PartStatus[0] = 0
					mbr.MbrPartition3.PartType[0] = 0
					mbr.MbrPartition3.PartFit[0] = 0
					mbr.MbrPartition3.PartStart = 0
					mbr.MbrPartition3.PartSize = 0
					var vacio [16]byte
					mbr.MbrPartition3.PartName = vacio
					mbr.MbrPartition3.PartCont = 0
				} else {
					fmt.Println("Parametro delete incorrecto (fast, full)")
					return false
				}
				//escribo el mbr modificado
				err = binary.Write(file, binary.LittleEndian, mbr)

				if err != nil {
					log.Fatalln(err)
				} else {
					return true
				}
			}
			return false
		}
	}
	if mbr.MbrPartition4.PartStatus[0] == 1 {

		//le saco el nombre
		nombreAct := ""
		for i := 0; i < 16; i++ {
			if mbr.MbrPartition4.PartName[i] != 0 {
				nombreAct += string(mbr.MbrPartition4.PartName[i])
			}
		}

		//si si es la particion que busco, elimino
		if nombre == nombreAct {
			file, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println("Error al leer el disco")
			} else {
				file.Seek(0, 0)
				if delete == "fast" {
					mbr.MbrPartition4.PartStatus[0] = 3
				} else if delete == "full" {
					mbr.MbrPartition4.PartStatus[0] = 0
					mbr.MbrPartition4.PartType[0] = 0
					mbr.MbrPartition4.PartFit[0] = 0
					mbr.MbrPartition4.PartStart = 0
					mbr.MbrPartition4.PartSize = 0
					var vacio [16]byte
					mbr.MbrPartition4.PartName = vacio
					mbr.MbrPartition4.PartCont = 0
				} else {
					fmt.Println("Parametro delete incorrecto (fast, full)")
					return false
				}
				//escribo el mbr modificado
				err = binary.Write(file, binary.LittleEndian, mbr)

				if err != nil {
					log.Fatalln(err)
				} else {
					return true
				}
			}
			return false
		}
	}

	return false
}

//EliminarPartitionLogica elimina una logica
func EliminarPartitionLogica(ruta string, nombre string, delete string) bool {

	mbr := MBR{}
	mbr = LeerMBR(ruta)
	if mbr.MbrPartition1.PartType[0] == 'e' && mbr.MbrPartition1.PartStatus[0] == 1 {
		var i int64
		for i = 0; i < mbr.MbrPartition1.PartCont; i++ {

			//agarro el ebr
			ebr := EBR{}
			ebr = LeerEBR(ruta, mbr.MbrPartition1.PartStart+48*i+i)

			//obtengo su nombre
			var nombreEBR [16]byte
			for i := 0; i < len(nombre); i++ {
				nombreEBR[i] = nombre[i]
			}

			if nombreEBR == ebr.PartName {
				if delete == "fast" {
					ebr.PartStatus[0] = 3
				} else if delete == "full" {
					ebr.PartStatus[0] = 0
					ebr.PartFit[0] = 0
					ebr.PartStart = 0
					ebr.PartSize = 0
					var vac [16]byte
					ebr.PartName = vac
					ebr.PartNext = 0
				} else {
					fmt.Println("Parametro delete incorrencto (fast/full)")
					return false
				}

				file, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
				//muevo seek al inicio del ebr actual que vamos a modificar
				file.Seek(mbr.MbrPartition1.PartStart+48*i+i, 0)
				if err != nil {
					fmt.Println("Problemas al leer el disco")
				} else {
					err = binary.Write(file, binary.LittleEndian, ebr)
					if err != nil {
						log.Fatalln(err)
						return false
					}
					return true
				}
				return false
			} //if nombre
		} //for
	}
	if mbr.MbrPartition2.PartType[0] == 'e' && mbr.MbrPartition2.PartStatus[0] == 1 {
		var i int64
		for i = 0; i < mbr.MbrPartition2.PartCont; i++ {

			//agarro el ebr
			ebr := EBR{}
			ebr = LeerEBR(ruta, mbr.MbrPartition2.PartStart+48*i+i)

			//obtengo su nombre
			var nombreEBR [16]byte
			for i := 0; i < len(nombre); i++ {
				nombreEBR[i] = nombre[i]
			}

			if nombreEBR == ebr.PartName {
				if delete == "fast" {
					ebr.PartStatus[0] = 3
				} else if delete == "full" {
					ebr.PartStatus[0] = 0
					ebr.PartFit[0] = 0
					ebr.PartStart = 0
					ebr.PartSize = 0
					var vac [16]byte
					ebr.PartName = vac
					ebr.PartNext = 0
				} else {
					fmt.Println("Parametro delete incorrencto (fast/full)")
					return false
				}

				file, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
				//muevo seek al inicio del ebr actual que vamos a modificar
				file.Seek(mbr.MbrPartition2.PartStart+48*i+i, 0)
				if err != nil {
					fmt.Println("Problemas al leer el disco")
				} else {
					err = binary.Write(file, binary.LittleEndian, ebr)
					if err != nil {
						log.Fatalln(err)
						return false
					}
					return true
				}
				return false
			} //if nombre
		} //for
	}
	if mbr.MbrPartition3.PartType[0] == 'e' && mbr.MbrPartition3.PartStatus[0] == 1 {
		var i int64
		for i = 0; i < mbr.MbrPartition3.PartCont; i++ {

			//agarro el ebr
			ebr := EBR{}
			ebr = LeerEBR(ruta, mbr.MbrPartition3.PartStart+48*i+i)

			//obtengo su nombre
			var nombreEBR [16]byte
			for i := 0; i < len(nombre); i++ {
				nombreEBR[i] = nombre[i]
			}

			if nombreEBR == ebr.PartName {
				if delete == "fast" {
					ebr.PartStatus[0] = 3
				} else if delete == "full" {
					ebr.PartStatus[0] = 0
					ebr.PartFit[0] = 0
					ebr.PartStart = 0
					ebr.PartSize = 0
					var vac [16]byte
					ebr.PartName = vac
					ebr.PartNext = 0
				} else {
					fmt.Println("Parametro delete incorrencto (fast/full)")
					return false
				}

				file, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
				//muevo seek al inicio del ebr actual que vamos a modificar
				file.Seek(mbr.MbrPartition3.PartStart+48*i+i, 0)
				if err != nil {
					fmt.Println("Problemas al leer el disco")
				} else {
					err = binary.Write(file, binary.LittleEndian, ebr)
					if err != nil {
						log.Fatalln(err)
						return false
					}
					return true
				}
				return false
			} //if nombre
		} //for
	}

	if mbr.MbrPartition4.PartType[0] == 'e' && mbr.MbrPartition4.PartStatus[0] == 1 {
		var i int64
		for i = 0; i < mbr.MbrPartition4.PartCont; i++ {

			//agarro el ebr
			ebr := EBR{}
			ebr = LeerEBR(ruta, mbr.MbrPartition4.PartStart+48*i+i)

			//obtengo su nombre
			var nombreEBR [16]byte
			for i := 0; i < len(nombre); i++ {
				nombreEBR[i] = nombre[i]
			}

			if nombreEBR == ebr.PartName {
				if delete == "fast" {
					ebr.PartStatus[0] = 3
				} else if delete == "full" {
					ebr.PartStatus[0] = 0
					ebr.PartFit[0] = 0
					ebr.PartStart = 0
					ebr.PartSize = 0
					var vac [16]byte
					ebr.PartName = vac
					ebr.PartNext = 0
				} else {
					fmt.Println("Parametro delete incorrencto (fast/full)")
					return false
				}

				file, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
				//muevo seek al inicio del ebr actual que vamos a modificar
				file.Seek(mbr.MbrPartition4.PartStart+48*i+i, 0)
				if err != nil {
					fmt.Println("Problemas al leer el disco")
				} else {
					err = binary.Write(file, binary.LittleEndian, ebr)
					if err != nil {
						log.Fatalln(err)
						return false
					}
					return true
				}
				return false
			} //if nombre
		} //for
	}
	return false
}

//==================================================================================
func addPartition(ruta string, nombre string, valor int64) int {

	mbr := MBR{}
	mbr = LeerMBR(ruta)
	if mbr.MbrPartition1.PartStatus[0] == 1 {
		nom := ""
		for i := 0; i < 16; i++ {
			if mbr.MbrPartition1.PartName[i] != 0 {
				nom += string(mbr.MbrPartition1.PartName[i])
			}
		}

		if nom == nombre {
			nuevoFin := (mbr.MbrPartition1.PartStart + mbr.MbrPartition1.PartSize) + valor
			tamanioDePartition := (mbr.MbrPartition1.PartStart + mbr.MbrPartition1.PartSize)

			if nuevoFin > tamanioDePartition { // hay que agregar

				//con uno que tenga problema, retorno 0
				if verRangos(mbr.MbrPartition1, nuevoFin) == 1 || verRangos(mbr.MbrPartition2, nuevoFin) == 1 || verRangos(mbr.MbrPartition3, nuevoFin) == 1 || verRangos(mbr.MbrPartition4, nuevoFin) == 1 {
					fmt.Println("No cabe, se pasa")
					return 0
				}
				//ahora valido que no se haya quedado una particion adentro
				if quedoDentro(mbr.MbrPartition2, mbr.MbrPartition1.PartStart, nuevoFin) == 1 || quedoDentro(mbr.MbrPartition3, mbr.MbrPartition1.PartStart, nuevoFin) == 1 || quedoDentro(mbr.MbrPartition4, mbr.MbrPartition1.PartStart, nuevoFin) == 1 {
					fmt.Println("Se quedo una adentro, asi que no crecera")
					return 0
				}
				//todo bien, escribimos
				//editando mbr
				file, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
				if err != nil {
					fmt.Println("error al abrir disco")
					return 0
				}
				file.Seek(0, 0)
				mbr.MbrPartition1.PartSize += valor
				err = binary.Write(file, binary.LittleEndian, mbr)

				if err != nil {
					log.Fatalln(err)
					return 0
				}

				//aqui tengo que ver lo nuevo

				//si alguno da 1 q
				validarDeleteFast(mbr.MbrPartition2, 2, nuevoFin, valor, ruta, mbr)
				validarDeleteFast(mbr.MbrPartition3, 3, nuevoFin, valor, ruta, mbr)
				validarDeleteFast(mbr.MbrPartition4, 4, nuevoFin, valor, ruta, mbr)

				//
				return 1
				//
			} else if nuevoFin < tamanioDePartition { // hay que quitar

				if nuevoFin > mbr.MbrPartition1.PartStart { // para que no quede negativo
					//editando mbr
					file, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
					if err != nil {
						fmt.Println("error al abrir disco")
						return 0
					}
					file.Seek(0, 0)
					mbr.MbrPartition1.PartSize += valor
					err = binary.Write(file, binary.LittleEndian, mbr)

					if err != nil {
						log.Fatalln(err)
						return 0
					}
					return 1
				}
				fmt.Sprintln("Quedaria negativo")
				return 0

			}
			return 0
		}
	}

	if mbr.MbrPartition2.PartStatus[0] == 1 {
		nom := ""
		for i := 0; i < 16; i++ {
			if mbr.MbrPartition2.PartName[i] != 0 {
				nom += string(mbr.MbrPartition2.PartName[i])
			}
		}

		if nom == nombre {
			nuevoFin := (mbr.MbrPartition2.PartStart + mbr.MbrPartition2.PartSize) + valor
			tamanioDePartition := (mbr.MbrPartition2.PartStart + mbr.MbrPartition2.PartSize)

			if nuevoFin > tamanioDePartition { // hay que agregar

				//con uno que tenga problema, retorno 0
				if verRangos(mbr.MbrPartition1, nuevoFin) == 1 || verRangos(mbr.MbrPartition2, nuevoFin) == 1 || verRangos(mbr.MbrPartition3, nuevoFin) == 1 || verRangos(mbr.MbrPartition4, nuevoFin) == 1 {
					fmt.Println("No cabe, se pasa")
					return 0
				}
				//ahora valido que no se haya quedado una particion adentro
				if quedoDentro(mbr.MbrPartition1, mbr.MbrPartition2.PartStart, nuevoFin) == 1 || quedoDentro(mbr.MbrPartition3, mbr.MbrPartition2.PartStart, nuevoFin) == 1 || quedoDentro(mbr.MbrPartition4, mbr.MbrPartition2.PartStart, nuevoFin) == 1 {
					fmt.Println("Se quedo una adentro, asi que no crecera")
					return 0
				}
				//todo bien, escribimos
				//editando mbr
				file, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
				if err != nil {
					fmt.Println("error al abrir disco")
					return 0
				}
				file.Seek(0, 0)
				mbr.MbrPartition2.PartSize += valor
				err = binary.Write(file, binary.LittleEndian, mbr)

				if err != nil {
					log.Fatalln(err)
					return 0
				}

				//aqui tengo que ver lo nuevo

				//si alguno da 1 q
				validarDeleteFast(mbr.MbrPartition1, 1, nuevoFin, valor, ruta, mbr)
				validarDeleteFast(mbr.MbrPartition3, 3, nuevoFin, valor, ruta, mbr)
				validarDeleteFast(mbr.MbrPartition4, 4, nuevoFin, valor, ruta, mbr)

				//
				return 1
				//
			} else if nuevoFin < tamanioDePartition { // hay que quitar

				if nuevoFin > mbr.MbrPartition2.PartStart { // para que no quede negativo
					//editando mbr
					file, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
					if err != nil {
						fmt.Println("error al abrir disco")
						return 0
					}
					file.Seek(0, 0)
					mbr.MbrPartition2.PartSize += valor
					err = binary.Write(file, binary.LittleEndian, mbr)

					if err != nil {
						log.Fatalln(err)
						return 0
					}
					return 1
				}
				fmt.Sprintln("Quedaria negativo")
				return 0

			}
			return 0
		}
	}

	if mbr.MbrPartition3.PartStatus[0] == 1 {
		nom := ""
		for i := 0; i < 16; i++ {
			if mbr.MbrPartition3.PartName[i] != 0 {
				nom += string(mbr.MbrPartition3.PartName[i])
			}
		}

		if nom == nombre {
			nuevoFin := (mbr.MbrPartition3.PartStart + mbr.MbrPartition3.PartSize) + valor
			tamanioDePartition := (mbr.MbrPartition3.PartStart + mbr.MbrPartition3.PartSize)

			if nuevoFin > tamanioDePartition { // hay que agregar

				fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
				fmt.Println(valor)
				//con uno que tenga problema, retorno 0
				if verRangos(mbr.MbrPartition1, nuevoFin) == 1 || verRangos(mbr.MbrPartition2, nuevoFin) == 1 || verRangos(mbr.MbrPartition3, nuevoFin) == 1 || verRangos(mbr.MbrPartition4, nuevoFin) == 1 {
					fmt.Println("No cabe, se pasa")
					return 0
				}
				//ahora valido que no se haya quedado una particion adentro
				if quedoDentro(mbr.MbrPartition1, mbr.MbrPartition3.PartStart, nuevoFin) == 1 || quedoDentro(mbr.MbrPartition2, mbr.MbrPartition3.PartStart, nuevoFin) == 1 || quedoDentro(mbr.MbrPartition4, mbr.MbrPartition3.PartStart, nuevoFin) == 1 {
					fmt.Println("Se quedo una adentro, asi que no crecera")
					return 0
				}
				//todo bien, escribimos
				//editando mbr
				file, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
				if err != nil {
					fmt.Println("error al abrir disco")
					return 0
				}
				file.Seek(0, 0)
				mbr.MbrPartition3.PartSize += valor
				err = binary.Write(file, binary.LittleEndian, mbr)

				if err != nil {
					log.Fatalln(err)
					return 0
				}

				//aqui tengo que ver lo nuevo

				//si alguno da 1 q
				validarDeleteFast(mbr.MbrPartition1, 1, nuevoFin, valor, ruta, mbr)
				validarDeleteFast(mbr.MbrPartition2, 2, nuevoFin, valor, ruta, mbr)
				validarDeleteFast(mbr.MbrPartition4, 4, nuevoFin, valor, ruta, mbr)

				//
				return 1
				//
			} else if nuevoFin < tamanioDePartition { // hay que quitar

				if nuevoFin > mbr.MbrPartition3.PartStart { // para que no quede negativo
					//editando mbr
					file, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
					if err != nil {
						fmt.Println("error al abrir disco")
						return 0
					}
					file.Seek(0, 0)
					mbr.MbrPartition3.PartSize += valor
					err = binary.Write(file, binary.LittleEndian, mbr)

					if err != nil {
						log.Fatalln(err)
						return 0
					}
					return 1
				}
				fmt.Sprintln("Quedaria negativo")
				return 0

			}
			fmt.Println("================ no paso ninguno de los ifs =============")
			return 0
		}
	}

	if mbr.MbrPartition4.PartStatus[0] == 1 {
		nom := ""
		for i := 0; i < 16; i++ {
			if mbr.MbrPartition4.PartName[i] != 0 {
				nom += string(mbr.MbrPartition4.PartName[i])
			}
		}

		if nom == nombre {
			nuevoFin := (mbr.MbrPartition4.PartStart + mbr.MbrPartition4.PartSize) + valor
			tamanioDePartition := (mbr.MbrPartition4.PartStart + mbr.MbrPartition4.PartSize)

			if nuevoFin > tamanioDePartition { // hay que agregar

				//con uno que tenga problema, retorno 0
				if verRangos(mbr.MbrPartition1, nuevoFin) == 1 || verRangos(mbr.MbrPartition2, nuevoFin) == 1 || verRangos(mbr.MbrPartition3, nuevoFin) == 1 || verRangos(mbr.MbrPartition4, nuevoFin) == 1 {
					fmt.Println("No cabe, se pasa")
					return 0
				}
				//ahora valido que no se haya quedado una particion adentro
				if quedoDentro(mbr.MbrPartition1, mbr.MbrPartition4.PartStart, nuevoFin) == 1 || quedoDentro(mbr.MbrPartition2, mbr.MbrPartition4.PartStart, nuevoFin) == 1 || quedoDentro(mbr.MbrPartition3, mbr.MbrPartition4.PartStart, nuevoFin) == 1 {
					fmt.Println("Se quedo una adentro, asi que no crecera")
					return 0
				}
				//todo bien, escribimos
				//editando mbr
				file, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
				if err != nil {
					fmt.Println("error al abrir disco")
					return 0
				}
				file.Seek(0, 0)
				mbr.MbrPartition4.PartSize += valor
				err = binary.Write(file, binary.LittleEndian, mbr)

				if err != nil {
					log.Fatalln(err)
					return 0
				}

				//aqui tengo que ver lo nuevo

				//
				validarDeleteFast(mbr.MbrPartition1, 1, nuevoFin, valor, ruta, mbr)
				validarDeleteFast(mbr.MbrPartition2, 2, nuevoFin, valor, ruta, mbr)
				validarDeleteFast(mbr.MbrPartition3, 3, nuevoFin, valor, ruta, mbr)

				//
				return 1
				//
			} else if nuevoFin < tamanioDePartition { // hay que quitar

				if nuevoFin > mbr.MbrPartition4.PartStart { // para que no quede negativo
					//editando mbr
					file, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
					if err != nil {
						fmt.Println("error al abrir disco")
						return 0
					}
					file.Seek(0, 0)
					mbr.MbrPartition4.PartSize += valor
					err = binary.Write(file, binary.LittleEndian, mbr)

					if err != nil {
						log.Fatalln(err)
						return 0
					}
					return 1
				}
				fmt.Sprintln("Quedaria negativo")
				return 0

			}
			return 0
		}
	}

	//
	return 0
}

//retorna 1 si el valor esta dentro del rango. [valor] ret 1
func verRangos(p Partition, valor int64) int {

	if p.PartStatus[0] == 1 {
		inicio := p.PartStart
		final := p.PartStart + p.PartSize
		if valor >= inicio && valor <= final {
			n := ""
			for i := 0; i < 16; i++ {
				n += string(p.PartName[i])
			}
			fmt.Println("-------------------------------------------@@@@ nombre " + n)
			return 1
		}
	}

	return 0
}

func quedoDentro(p Partition, Vi int64, Vf int64) int {
	if p.PartStatus[0] == 1 {
		inicio := p.PartStart
		final := p.PartStart + p.PartSize
		if inicio >= Vi && inicio <= Vf && final >= Vi && final <= Vf {
			//quiere decir que si esta dentro
			return 1
		}
	}
	return 0
}

//retorno 1 si el nuevoFinal esta dentro de una particion borrada con FAST
func validarDeleteFast(p Partition, numPart int, nuevoF int64, add int64, ruta string, mbr MBR) {

	if p.PartStatus[0] == 3 {
		inicio := p.PartStart
		final := p.PartStart + p.PartSize
		if nuevoF >= inicio && nuevoF <= final {
			//quiere decir que si esta dentro de una particion con delete fast
			//entonces hay que quitarle a la fast size, lo que se le agrego

			if numPart == 1 {
				//editando mbr
				file, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
				if err != nil {
					fmt.Println("error al abrir disco")
					return
				}
				file.Seek(0, 0)
				mbr.MbrPartition1.PartSize = mbr.MbrPartition1.PartSize - add
				err = binary.Write(file, binary.LittleEndian, mbr)

				if err != nil {
					log.Fatalln(err)
					return
				}
			} else if numPart == 2 {
				//editando mbr
				file, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
				if err != nil {
					fmt.Println("error al abrir disco")
					return
				}
				file.Seek(0, 0)
				mbr.MbrPartition2.PartSize = mbr.MbrPartition2.PartSize - add
				err = binary.Write(file, binary.LittleEndian, mbr)

				if err != nil {
					log.Fatalln(err)
					return
				}
			} else if numPart == 3 {
				//editando mbr
				file, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
				if err != nil {
					fmt.Println("error al abrir disco")
					return
				}
				file.Seek(0, 0)
				mbr.MbrPartition3.PartSize = mbr.MbrPartition3.PartSize - add
				err = binary.Write(file, binary.LittleEndian, mbr)

				if err != nil {
					log.Fatalln(err)
					return
				}
			} else if numPart == 4 {
				//editando mbr
				file, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
				if err != nil {
					fmt.Println("error al abrir disco")
					return
				}
				file.Seek(0, 0)
				mbr.MbrPartition4.PartSize = mbr.MbrPartition4.PartSize - add
				err = binary.Write(file, binary.LittleEndian, mbr)

				if err != nil {
					log.Fatalln(err)
					return
				}
			}
		}
	}

}

// ADD logica
func addLogica(ruta string, nombre string, valor int64) int {
	mbr := MBR{}
	mbr = LeerMBR(ruta)

	if mbr.MbrPartition1.PartType[0] == 'e' && mbr.MbrPartition1.PartStatus[0] == 1 {

		//recorro las logicas
		var i int64
		for i = 1; i < mbr.MbrPartition1.PartCont; i++ {
			ebr := EBR{}
			ebr = LeerEBR(ruta, mbr.MbrPartition1.PartStart+48*i+i)

			//convierto el nombre
			var nomb [16]byte
			for i := 0; i < len(nombre); i++ {
				nomb[i] = nombre[i]
			}

			if nomb == ebr.PartName {
				nuevo := ebr.PartSize + valor
				if nuevo > 0 {
					ebr.PartSize = nuevo
					file, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
					file.Seek(mbr.MbrPartition1.PartStart+48*i+i, 0)
					if err != nil {
						fmt.Println("Error al abrir el disco")
					} else {
						err = binary.Write(file, binary.LittleEndian, ebr)
						if err != nil {
							log.Fatalln(err)
							return 0
						}
						return 1
					}
				}
				return 0
			}
		}
	}

	if mbr.MbrPartition2.PartType[0] == 'e' && mbr.MbrPartition2.PartStatus[0] == 1 {

		//recorro las logicas
		var i int64
		for i = 1; i < mbr.MbrPartition2.PartCont; i++ {
			ebr := EBR{}
			ebr = LeerEBR(ruta, mbr.MbrPartition2.PartStart+48*i+i)

			//convierto el nombre
			var nomb [16]byte
			for i := 0; i < len(nombre); i++ {
				nomb[i] = nombre[i]
			}

			if nomb == ebr.PartName {
				nuevo := ebr.PartSize + valor
				if nuevo > 0 {
					ebr.PartSize = nuevo
					file, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
					file.Seek(mbr.MbrPartition2.PartStart+48*i+i, 0)
					if err != nil {
						fmt.Println("Error al abrir el disco")
					} else {
						err = binary.Write(file, binary.LittleEndian, ebr)
						if err != nil {
							log.Fatalln(err)
							return 0
						}
						return 1
					}
				}
				return 0
			}
		}
	}

	if mbr.MbrPartition3.PartType[0] == 'e' && mbr.MbrPartition3.PartStatus[0] == 1 {

		//recorro las logicas
		var i int64
		for i = 1; i < mbr.MbrPartition3.PartCont; i++ {
			ebr := EBR{}
			ebr = LeerEBR(ruta, mbr.MbrPartition3.PartStart+48*i+i)

			//convierto el nombre
			var nomb [16]byte
			for i := 0; i < len(nombre); i++ {
				nomb[i] = nombre[i]
			}

			if nomb == ebr.PartName {
				nuevo := ebr.PartSize + valor
				if nuevo > 0 {
					ebr.PartSize = nuevo
					file, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
					file.Seek(mbr.MbrPartition3.PartStart+48*i+i, 0)
					if err != nil {
						fmt.Println("Error al abrir el disco")
					} else {
						err = binary.Write(file, binary.LittleEndian, ebr)
						if err != nil {
							log.Fatalln(err)
							return 0
						}
						return 1
					}
				}
				return 0
			}
		}
	}

	if mbr.MbrPartition4.PartType[0] == 'e' && mbr.MbrPartition4.PartStatus[0] == 1 {

		//recorro las logicas
		var i int64
		for i = 1; i < mbr.MbrPartition4.PartCont; i++ {
			ebr := EBR{}
			ebr = LeerEBR(ruta, mbr.MbrPartition4.PartStart+48*i+i)

			//convierto el nombre
			var nomb [16]byte
			for i := 0; i < len(nombre); i++ {
				nomb[i] = nombre[i]
			}

			if nomb == ebr.PartName {
				nuevo := ebr.PartSize + valor
				if nuevo > 0 {
					ebr.PartSize = nuevo
					file, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
					file.Seek(mbr.MbrPartition4.PartStart+48*i+i, 0)
					if err != nil {
						fmt.Println("Error al abrir el disco")
					} else {
						err = binary.Write(file, binary.LittleEndian, ebr)
						if err != nil {
							log.Fatalln(err)
							return 0
						}
						return 1
					}
				}
				return 0
			}
		}
	}

	return 0
}
