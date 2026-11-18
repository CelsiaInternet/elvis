package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/celsiainternet/elvis/agentsguide"
)

/**
* main: Prepara un proyecto que consume elvis: go mod tidy (si hay código) e instala la guía de agentes.
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
* tidy: Ejecuta go mod tidy solo si el módulo ya tiene paquetes Go.
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
