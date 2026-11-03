package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/celsiainternet/elvis/agentsguide"
)

/**
* main: Prepara un proyecto que consume elvis: completa el go.mod con
* go mod tidy (solo si ya hay codigo) e instala la guia de agentes.
* No fija versiones: tidy toma las que exige elvis y nunca las baja.
**/
func main() {
	goMod, err := os.ReadFile("go.mod")
	if err != nil {
		fmt.Println("No se encontro go.mod en el directorio actual; ejecuta primero: go mod init <modulo> && go get github.com/celsiainternet/elvis")
		os.Exit(1)
	}

	ok := tidy()

	if err := agentsguide.Install(); err != nil {
		fmt.Printf("Advertencia: no se pudo instalar la guia de agentes - %v\n", err)
		ok = false
	} else {
		fmt.Println("Guia de agentes de Elvis instalada en .claude/guides/elvis/ (referenciada desde CLAUDE.md y AGENTS.md)")
	}

	if strings.Contains(string(goMod), "github.com/celsiainternet/jdb ") {
		fmt.Println("Este proyecto usa jdb: instala tambien su guia con go run github.com/celsiainternet/jdb/cmd/install")
	}

	if !ok {
		os.Exit(1)
	}
}

/**
* tidy: Ejecuta go mod tidy si el modulo ya tiene paquetes Go. En un
* proyecto vacio tidy quitaria elvis del go.mod (nada la importa todavia),
* y despues cmd/create no se podria ejecutar.
* @return bool
**/
func tidy() bool {
	if !hasPackages() {
		fmt.Println("El proyecto aun no tiene codigo: se omite go mod tidy (el generador lo ejecuta al crear el proyecto)")
		return true
	}

	fmt.Println("Ejecutando go mod tidy...")
	out, err := exec.Command("go", "mod", "tidy").CombinedOutput()
	if err != nil {
		fmt.Printf("go mod tidy fallo: %v\n%s\n", err, out)
		return false
	}

	fmt.Println("go.mod actualizado")
	return true
}

/**
* hasPackages: Indica si el modulo actual tiene al menos un paquete Go.
* @return bool
**/
func hasPackages() bool {
	out, err := exec.Command("go", "list", "-e", "-f", "{{if .GoFiles}}{{.ImportPath}}{{end}}", "./...").Output()
	if err != nil {
		return false
	}

	return strings.TrimSpace(string(out)) != ""
}
