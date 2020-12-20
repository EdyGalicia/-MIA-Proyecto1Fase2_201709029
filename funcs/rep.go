package funcs

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// EjecutarREP asd
func EjecutarREP(parametros []string, descripciones []string) {
	fmt.Println(" === COMANDO REP EN EJECUCION === ")
	var nombre string
	var ruta string
	var id string
	hayError := false

	for i := 0; i < len(parametros); i++ {
		switch parametros[i] {
		case "name":
			{
				nombre = descripciones[i]
			}
		case "path":
			{
				ruta = descripciones[i]
			}
		case "id":
			{
				id = descripciones[i]
			}
		default:
			{
				fmt.Println("Hay errores en los parametros")
				hayError = true
			}

		}
	}

	if hayError == true {
		fmt.Println("...")
	} else {
		switch nombre {
		case "mbr":
			{
				fmt.Println("Entro al reporte mbr")
				crearReporteMBR(ruta, id)
			}
		case "disk":
			{
				crearReporteDISK(ruta, id)
			}
		case "inode":
			{
				crearReporteInode(ruta, id)
			}
		case "block":
			{
				crearReporteBlock(ruta, id)
			}
		case "bm_inode":
			{
				crearReporteBmInode(ruta, id)
			}
		case "bm_block":
			{
				crearReporteBmBlock(ruta, id)
			}
		case "sb":
			{
				crearReporteSb(ruta, id)
			}
		default:
			{
				fmt.Println("Nombre de reporte no valido")
			}
		}
	}
}

func crearReporteMBR(ruta string, id string) {

	estructura := genearCuerpoReporteMBR(id)

	aux := strings.Split(ruta, "/")
	nombreV := strings.Split(aux[len(aux)-1], ".") //nombreV [nombre | extension]
	nombre := strings.ReplaceAll(nombreV[0], " ", "")
	directorio := ""
	for i := 0; i < len(aux)-1; i++ {
		directorio += aux[i] + "/"
	}

	err := os.MkdirAll(directorio, 0777)
	if err != nil {
		fmt.Println("Error en la ruta")
	} else {
		file, err1 := os.Create(directorio + "/" + nombre + ".dot")
		defer file.Close()
		if err1 != nil {
			fmt.Println("Error al generar el archivo dot")
		} else {
			_, errr := file.WriteString(estructura)
			if errr != nil {
				fmt.Println("Error al querer escribir")
			} else {
				fmt.Println("el dot fue generado")

				err2 := exec.Command("dot", directorio+nombre+".dot", "-o", directorio+nombre+".jpg", "-Tjpg").Run()
				if err2 != nil {
					fmt.Println("Error al generar el comando en consola")
				} else {
					err3 := exec.Command("xdg-open", directorio+nombre+".jpg").Run()
					if err3 != nil {
						fmt.Println("Error al abrir el reporte")
					} else {
						fmt.Println("Generado correctamente")
					}
				}

			}
		}
	}
}

func genearCuerpoReporteMBR(id string) string {

	rutaDelDisco := ""
	body := "digraph test {\ngraph [ratio=fill];node [label=\"N\", fontsize=15, shape=plaintext]; graph [bb=\"0,0,352,154\"];arset [label=<\n<TABLE ALIGN=\"LEFT\">\n<TR ><TD  BGCOLOR=\"orange\"   colspan=\"2\">MBR</TD></TR>\n"

	for i := 0; i < len(ParticionesMontadas); i++ {
		if id == ParticionesMontadas[i].id {
			rutaDelDisco = ParticionesMontadas[i].ruta
			i = len(ParticionesMontadas)
		}
	}
	if rutaDelDisco != "" {
		mbr := LeerMBR(rutaDelDisco)
		fecha := ""
		for i := 0; i < 16; i++ {
			fecha += string(mbr.MbrFechaCreacion[i])
		}

		body += "<TR><TD>" + "mbr_tamano</TD><TD>" + strconv.FormatInt(mbr.MbrTamanio, 10) + "</TD></TR>\n"
		body += "<TR><TD>" + "mbr_FechaCreacion</TD><TD>" + fecha + "</TD></TR>\n"
		body += "<TR><TD>" + "mbr_DiskSignature</TD><TD>" + strconv.FormatInt(mbr.MbrDiskSignature, 10) + "</TD></TR>\n"

		body += getDatosParticion(mbr.MbrPartition1, 1)
		body += getDatosParticion(mbr.MbrPartition2, 2)
		body += getDatosParticion(mbr.MbrPartition3, 3)
		body += getDatosParticion(mbr.MbrPartition4, 4)

		body += "</TABLE>>, ];\n"

		//aca va las figuras de las logicas, cuando las tenga :(
		body += validarLogica(rutaDelDisco)
	} else {
		fmt.Println("La montura con id: " + id + " no esta montada.")
	}

	body += "}"

	return body
}

func getDatosParticion(p Partition, numDePart int) string {

	body := ""
	body = "<TR ><TD BGCOLOR=\"cyan\"  colspan=\"2\">Particion" + strconv.Itoa(numDePart) + "</TD></TR>\n"

	if p.PartStatus[0] == 0 || p.PartStatus[0] == 3 {
		body += "<TR><TD>" + "part_status</TD><TD>0</TD></TR>\n"
	} else if p.PartStatus[0] == 1 {
		body += "<TR><TD>" + "part_status</TD><TD>1</TD></TR>\n"
	}

	if p.PartType[0] == 0 {
		body += "<TR><TD>" + "part_type</TD><TD>0</TD></TR>\n"
	} else {
		body += "<TR><TD>" + "part_type</TD><TD>" + string(p.PartType[0]) + "</TD></TR>\n"
	}

	if p.PartFit[0] == 0 {
		body += "<TR><TD>" + "part_fit</TD><TD>0</TD></TR>\n"
	} else {
		body += "<TR><TD>" + "part_fit</TD><TD>" + string(p.PartFit[0]) + "</TD></TR>\n"
	}

	body += "<TR><TD>" + "part_start</TD><TD>" + strconv.FormatInt(p.PartStart, 10) + "</TD></TR>\n"
	body += "<TR><TD>" + "part_size</TD><TD>" + strconv.FormatInt(p.PartSize, 10) + "</TD></TR>\n"
	nom := ""
	for j := 0; j < 16; j++ {
		if p.PartName[j] != 0 {
			nom += string(p.PartName[j])
		}
	}
	body += "<TR><TD>" + "part_name</TD><TD>" + nom + "</TD></TR>\n"
	return body
}

func validarLogica(ruta string) string {
	mbr := MBR{}
	mbr = LeerMBR(ruta)

	if mbr.MbrPartition1.PartType[0] == 'e' && mbr.MbrPartition1.PartStatus[0] == 1 {
		return getPartLogica(mbr.MbrPartition1, ruta, mbr.MbrPartition1.PartStart)
	} else if mbr.MbrPartition2.PartType[0] == 'e' && mbr.MbrPartition2.PartStatus[0] == 1 {
		return getPartLogica(mbr.MbrPartition2, ruta, mbr.MbrPartition2.PartStart)
	} else if mbr.MbrPartition3.PartType[0] == 'e' && mbr.MbrPartition3.PartStatus[0] == 1 {
		return getPartLogica(mbr.MbrPartition3, ruta, mbr.MbrPartition3.PartStart)
	} else if mbr.MbrPartition4.PartType[0] == 'e' && mbr.MbrPartition4.PartStatus[0] == 1 {
		return getPartLogica(mbr.MbrPartition4, ruta, mbr.MbrPartition4.PartStart)
	}
	return ""
}

func getPartLogica(p Partition, ruta string, inicio int64) string {
	var i int64
	body := "tbl[label=<<TABLE ALIGN=\"LEFT\">\n"

	for i = 1; i < p.PartCont; i++ {
		ebr := EBR{}
		//PartStart+48*i+i
		ebr = LeerEBR(ruta, inicio+48*i+i)
		//el primer ebr tiene un status = 2
		if ebr.PartStatus[0] != 0 {
			body += "<TR ><TD BGCOLOR=\"orange\"  colspan=\"2\">EBR" + strconv.FormatInt(i, 10) + "</TD></TR>\n"

			//escribo el status de la particion
			if ebr.PartStatus[0] == 3 {
				body += "<TR><TD>" + "part_status</TD><TD>0</TD></TR>\n"
			} else if ebr.PartStatus[0] == 1 {
				body += "<TR><TD>" + "part_status</TD><TD>1</TD></TR>\n"
			}

			//escribo el fit de la particion
			if ebr.PartFit[0] == 0 {
				body += "<TR><TD>" + "part_fit</TD><TD>0</TD></TR>\n"
			} else {
				body += "<TR><TD>" + "part_fit</TD><TD>" + string(ebr.PartFit[0]) + "</TD></TR>\n"
			}

			//escribo el start
			body += "<TR><TD>" + "part_start</TD><TD>" + strconv.FormatInt(ebr.PartStart, 10) + "</TD></TR>\n"

			//escribo el part size
			body += "<TR><TD>" + "part_size</TD><TD>" + strconv.FormatInt(ebr.PartSize, 10) + "</TD></TR>\n"

			//escribo el nombre-----------------------------------------------
			nom := ""
			for j := 0; j < 16; j++ {
				if ebr.PartName[j] != 0 {
					nom += string(ebr.PartName[j])
				}
			}
			body += "<TR><TD>" + "part_name</TD><TD>" + nom + "</TD></TR>\n"
			//----------------------------------------------------------------

			//escribo el partNext
			if i == p.PartCont-1 { //cuando llegue a la ultima particion logica
				body += "<TR><TD>" + "part_next</TD><TD>-1</TD></TR>\n"
			} else {
				body += "<TR><TD>" + "part_next</TD><TD>" + strconv.FormatInt(ebr.PartNext, 10) + "</TD></TR>\n"
			}
		}
	}
	body += "</TABLE>>, ];\n"
	return body
}

//El otro reporte
func crearReporteDISK(ruta string, id string) {

	estructura := genearCuerpoReporteDISK(id)

	aux := strings.Split(ruta, "/")
	nombreV := strings.Split(aux[len(aux)-1], ".") //nombreV [nombre | extension]
	nombre := strings.ReplaceAll(nombreV[0], " ", "")
	directorio := ""
	for i := 0; i < len(aux)-1; i++ {
		directorio += aux[i] + "/"
	}

	err := os.MkdirAll(directorio, 0777)
	if err != nil {
		fmt.Println("Error en la ruta")
	} else {
		file, err1 := os.Create(directorio + "/" + nombre + ".dot")
		defer file.Close()
		if err1 != nil {
			fmt.Println("Error al generar el archivo dot")
		} else {
			_, errr := file.WriteString(estructura)
			if errr != nil {
				fmt.Println("Error al querer escribir")
			} else {
				fmt.Println("el dot fue generado")

				err2 := exec.Command("dot", directorio+nombre+".dot", "-o", directorio+nombre+".jpg", "-Tjpg").Run()
				if err2 != nil {
					fmt.Println("Error al generar el comando en consola")
				} else {
					err3 := exec.Command("xdg-open", directorio+nombre+".jpg").Run()
					if err3 != nil {
						fmt.Println("Error al abrir el reporte")
					} else {
						fmt.Println("Generado correctamente")
					}
				}

			}
		}
	}
}

func genearCuerpoReporteDISK(id string) string {

	rutaDelDisco := ""
	body := "digraph test {\ngraph [ratio=fill];node [label=\"N\", fontsize=15, shape=plaintext]; graph [bb=\"0,0,352,154\"];arset [label=<\n<TABLE ALIGN=\"LEFT\">\n"

	for i := 0; i < len(ParticionesMontadas); i++ {
		if id == ParticionesMontadas[i].id {
			rutaDelDisco = ParticionesMontadas[i].ruta
			i = len(ParticionesMontadas)
		}
	}
	if rutaDelDisco != "" {
		mbr := LeerMBR(rutaDelDisco)
		body += "<TR >"

		//bloquecito del MBR
		body += "<td BGCOLOR=\"green\"  rowspan=\"2\">MBR <br/> " + "</td>"

		body += generarParticionesPE(mbr.MbrPartition1, mbr.MbrTamanio, rutaDelDisco)
		body += generarParticionesPE(mbr.MbrPartition2, mbr.MbrTamanio, rutaDelDisco)
		body += generarParticionesPE(mbr.MbrPartition3, mbr.MbrTamanio, rutaDelDisco)
		body += generarParticionesPE(mbr.MbrPartition4, mbr.MbrTamanio, rutaDelDisco)

		//obtengo el porcentaje de espacio disponible en el disco
		espacioDisponible := espacioDisponibleEnElDisco(rutaDelDisco)
		porcentajeEspDisp := (float64(espacioDisponible) / float64(mbr.MbrTamanio)) * 100
		body += "<td BGCOLOR=\"cadetblue1\" rowspan=\"2\">Libre <br/> " + strconv.FormatFloat(porcentajeEspDisp, 'f', 2, 64) + "%</td>"

		body += "</TR>\n"

		//checo si existe Extendida
		_, _, existeExtendida := ValidarExtendidaEBR(rutaDelDisco)
		if existeExtendida == 1 {
			body += "<tr>\n"
			body += validarLogicaDisk(rutaDelDisco)
			tam1 := validarEspacioEBR(rutaDelDisco)
			tam := ValidarTamEBR(rutaDelDisco)
			porcentaje := (float64(tam1) / float64(tam)) * 100
			body += "<td BGCOLOR=\"cadetblue1\" >Libre <br/> " + strconv.FormatFloat(porcentaje, 'f', 2, 64) + "%</td>"
			body += "</tr>\n"
		}

		//
		body += "</TABLE>>, ];\n"
	} else {
		fmt.Println("La montura con id: " + id + " no esta montada.")
	}

	body += "}"

	return body
}

func generarParticionesPE(p Partition, tamMax int64, ruta string) string {
	body := ""

	if p.PartStatus[0] == 1 {
		if p.PartType[0] == 'e' || p.PartType[0] == 'E' {
			porcentaje := (float64(p.PartSize) / float64(tamMax)) * 100
			body += "<td BGCOLOR=\"chartreuse\" colspan=\"" + strconv.FormatInt(p.PartCont*2-1, 10) + "\">Extendida <br/> " + strconv.FormatFloat(porcentaje, 'f', 2, 64) + "%</td>"

		} else if p.PartType[0] == 'p' || p.PartType[0] == 'P' {
			porcentaje := (float64(p.PartSize) / float64(tamMax)) * 100
			body += "<td BGCOLOR=\"green\"  rowspan=\"2\">Primaria <br/> " + strconv.FormatFloat(porcentaje, 'f', 2, 64) + "%</td>"

		}
	} else if p.PartStatus[0] == 3 {
		porcentaje := (float64(p.PartSize) / float64(tamMax)) * 100
		body += "<td BGCOLOR=\"cadetblue1\"  rowspan=\"2\">Libre <br/> " + strconv.FormatFloat(porcentaje, 'f', 2, 64) + "%</td>"

	}

	return body
}

func validarLogicaDisk(ruta string) string {
	mbr := MBR{}
	mbr = LeerMBR(ruta)

	if mbr.MbrPartition1.PartType[0] == 'e' && mbr.MbrPartition1.PartStatus[0] == 1 {
		return obtenerLogicaDisk(mbr.MbrPartition1, ruta, mbr.MbrPartition1.PartSize, mbr.MbrPartition1.PartStart)
	} else if mbr.MbrPartition2.PartType[0] == 'e' && mbr.MbrPartition2.PartStatus[0] == 1 {
		return obtenerLogicaDisk(mbr.MbrPartition2, ruta, mbr.MbrPartition2.PartSize, mbr.MbrPartition2.PartStart)
	} else if mbr.MbrPartition3.PartType[0] == 'e' && mbr.MbrPartition3.PartStatus[0] == 1 {
		return obtenerLogicaDisk(mbr.MbrPartition3, ruta, mbr.MbrPartition3.PartSize, mbr.MbrPartition3.PartStart)
	} else if mbr.MbrPartition4.PartType[0] == 'e' && mbr.MbrPartition4.PartStatus[0] == 1 {
		return obtenerLogicaDisk(mbr.MbrPartition4, ruta, mbr.MbrPartition4.PartSize, mbr.MbrPartition4.PartStart)
	}

	return ""
}

func obtenerLogicaDisk(Pa Partition, ruta string, tam int64, start int64) string {

	var i int64
	Cuerpo := ""

	for i = 1; i < Pa.PartCont; i++ {
		eb := EBR{}
		eb = LeerEBR(ruta, start+48*i+i)
		if eb.PartStatus[0] == 1 {
			Cuerpo += "<td BGCOLOR=\"chartreuse\">EBR</td>"
			porcentaje := (float64(eb.PartSize) / float64(tam)) * 100

			Cuerpo += "<td BGCOLOR=\"greenyellow\" >Logica <br/> " + strconv.FormatFloat(porcentaje, 'f', 2, 64) + "%</td>"

		} else if eb.PartStatus[0] == 3 {
			Cuerpo += "<td BGCOLOR=\"chartreuse\" ></td>"
			porcentaje := (float64(eb.PartSize) / float64(tam)) * 100
			Cuerpo += "<td BGCOLOR=\"cadetblue1\" >Libre <br/> " + strconv.FormatFloat(porcentaje, 'f', 2, 64) + "%</td>"
		}

	}

	return Cuerpo

}

func espacioDisponibleEnElDisco(ruta string) int64 {

	mbr := MBR{}
	mbr = LeerMBR(ruta)
	var espacioOcupado int64
	if mbr.MbrPartition1.PartStatus[0] != 0 {
		espacioOcupado = espacioOcupado + mbr.MbrPartition1.PartSize
	}
	if mbr.MbrPartition2.PartStatus[0] != 0 {
		espacioOcupado = espacioOcupado + mbr.MbrPartition2.PartSize
	}
	if mbr.MbrPartition3.PartStatus[0] != 0 {
		espacioOcupado = espacioOcupado + mbr.MbrPartition3.PartSize
	}
	if mbr.MbrPartition4.PartStatus[0] != 0 {
		espacioOcupado = espacioOcupado + mbr.MbrPartition4.PartSize
	}
	//retorno el espacio libre
	return (mbr.MbrTamanio - espacioOcupado)
}

//ValidarTamEBR : espacio del ebr
func ValidarTamEBR(ruta string) int64 {
	mbr := MBR{}
	mbr = LeerMBR(ruta)

	if mbr.MbrPartition1.PartType[0] == 'e' && mbr.MbrPartition1.PartStatus[0] == 1 {

		return mbr.MbrPartition1.PartSize
	} else if mbr.MbrPartition2.PartType[0] == 'e' && mbr.MbrPartition2.PartStatus[0] == 1 {

		return mbr.MbrPartition2.PartSize
	} else if mbr.MbrPartition3.PartType[0] == 'e' && mbr.MbrPartition3.PartStatus[0] == 1 {

		return mbr.MbrPartition3.PartSize
	} else if mbr.MbrPartition4.PartType[0] == 'e' && mbr.MbrPartition4.PartStatus[0] == 1 {

		return mbr.MbrPartition4.PartSize
	}

	return 0
}
