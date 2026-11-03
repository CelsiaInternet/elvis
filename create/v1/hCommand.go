package create

import (
	"fmt"

	"github.com/celsiainternet/elvis/strs"
	"github.com/celsiainternet/elvis/utility"
	"github.com/spf13/cobra"
)

var CmdProject = &cobra.Command{
	Use:   "project [name author schema]",
	Short: "Create a new project with a base microservice.",
	Long:  "Template project including README, .env, .gitignore and a base microservice with folders cmd, deployments, pkg, rest and test.",
	Run: func(cmd *cobra.Command, args []string) {
		packageName, err := moduleName()
		if err != nil {
			fmt.Println(err)
			return
		}

		name, err := PrompStr("Name", true)
		if err != nil {
			fmt.Printf("\tPrompt failed: %v\n", err)
			return
		}

		author, err := PrompStr("Author", true)
		if err != nil {
			fmt.Printf("\tPrompt failed: %v\n", err)
			return
		}

		schema, err := PrompStr("Schema", false)
		if err != nil {
			fmt.Printf("\tPrompt failed: %v\n", err)
			return
		}

		err = MkProject(packageName, name, author, schema)
		if err != nil {
			fmt.Printf("\tCommand failed: %v\n", err)
			return
		}
	},
}

var CmdMicro = &cobra.Command{
	Use:   "micro [name schema]",
	Short: "Create a microservice inside the current project.",
	Long:  "Template microservice including folders cmd, deployments, internal, pkg, scripts, test and www, with the .go files required for making a microservice.",
	Run: func(cmd *cobra.Command, args []string) {
		packageName, err := moduleName()
		if err != nil {
			fmt.Println(err)
			return
		}

		name, err := PrompStr("Name", true)
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		schema, err := PrompStr("Schema", false)
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		err = MkMicroservice(packageName, name, schema)
		if err != nil {
			fmt.Printf("Command failed %v\n", err)
			return
		}
	},
}

var CmdModelo = &cobra.Command{
	Use:   "modelo [name modelo, schema]",
	Short: "Create model to microservice.",
	Long:  "Template model to microservice include function handler model.",
	Run: func(cmd *cobra.Command, args []string) {
		packageName, err := PrompStr("Package", true)
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		modelo, err := PrompStr("Model", true)
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		schema, err := PrompStr("Schema", false)
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		err = MkMolue(packageName, modelo, schema)
		if err != nil {
			fmt.Printf("Command failed %v\n", err)
			return
		}

		title := strs.Titlecase(modelo)
		message := strs.Format(`Copia las rutas que estan al final de pkg/%s/h%s.go en Routes() de router.go.`, packageName, title)
		if len(schema) > 0 {
			message = strs.Format(`%s Si el paquete se creo sin schema, agrega la llamada a initModels(c.Db) en Controller.Init.`, message)
		}
		fmt.Println(message)
	},
}

var CmdRpc = &cobra.Command{
	Use:   "rpc [package]",
	Short: "Create rpc model to microservice.",
	Long:  "Template rpc model to microservice include function handler model.",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := PrompStr("Package", true)
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		err = MkRpc(name)
		if err != nil {
			fmt.Printf("Command failed %v\n", err)
			return
		}
	},
}

/**
* moduleName: Lee el modulo del go.mod del directorio actual; el generador
* lo usa para los imports internos del proyecto.
* @return string, error
**/
func moduleName() (string, error) {
	result, err := utility.GoMod("module")
	if err != nil || result == "" {
		return "", fmt.Errorf("no se encontro go.mod con module en el directorio actual; ejecuta primero: go mod init <modulo>")
	}

	return result, nil
}
