package main

import (
	"MIA_Fase2_201709029/funcs"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	//funcs.TamanioDelEBRbro()
	funcs.TipoDeInstruccion("exec -path->\"/home/edygalicia/Documentos/pr.arch\"")
	////funcs.TipoDeInstruccion("mkdisk -size->50 -fit->BF -unit->M -path->\"/home/e discos/Disco1p.dsk\"")
	//funcs.TipoDeInstruccion("montar -size->5 -unit->M -path->\"/home/mis discos/Disco3.dsk\"")
	//funcs.TipoDeInstruccion("rmdisk -path->\"/home/los discos/Disco3.dsk\"")
	//funcs.TipoDeInstruccion("fdisk -path->\"/home/e discos/Disco1p.dsk\" -Size->100 -unit->K -name->part2soy")

	//
	//funcs.TipoDeInstruccion("mount -path->\"/home/Pack1 de discos/Disco2.dsk\" -name->LaPart1D2")
	//funcs.TipoDeInstruccion("rep -path->\"/home/Algun reporte/nuevoRep1.png\" -name->mbr -id->vda1")

	//asd
	//funcs.TipoDeInstruccion("fdisk -path->\"/home/Pack1 de discos/Disco2.dsk\" -delete->fast -name->LaPart2D2")

	//exec -path->/home/edygalicia/Escritorio/cam.arch

	terminal := exec.Command("Clear")
	terminal.Stdout = os.Stdout
	terminal.Run()
	fmt.Println(" =============== FASE 1 EDY GALICIA =================== ")
	scanner := bufio.NewReader(os.Stdin)
	for true {
		fmt.Print("comando: ")
		entrada, err := scanner.ReadString('\n')
		if err == nil {
			if entrada != "" {
				entrada = strings.Replace(entrada, "\n", "", 1)
				funcs.TipoDeInstruccion(entrada)
			} //exec -path->"/home/edygalicia/Escritorio/a.arch"
		}
	}
}
