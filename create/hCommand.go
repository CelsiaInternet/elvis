package create

import (
	"fmt"

	utl "github.com/cgalvisleon/elvis/utilities"
	"github.com/spf13/cobra"
)

var CmdProject = &cobra.Command{
	Use:   "micro [name author schema, schema_var]",
	Short: "Create project base type microservice.",
	Long:  "Template project to microservice include folder cmd, deployments, pkg, rest, test and web, with files .go required for making a microservice.",
	Run: func(cmd *cobra.Command, args []string) {
		packageName, err := utl.ModuleName()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		name, err := PrompStr("Name")
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		author, err := PrompStr("Author")
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		schema, err := PrompStr("Schema")
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		err = MkPMicroservice(packageName, name, author, schema)
		if err != nil {
			fmt.Printf("Command failed %v\n", err)
			return
		}
	},
}

var CmdMicro = &cobra.Command{
	Use:   "micro [name schema, schema_var]",
	Short: "Create project base type microservice.",
	Long:  "Template project to microservice include folder cmd, deployments, pkg, rest, test and web, with files .go required for making a microservice.",
	Run: func(cmd *cobra.Command, args []string) {
		packageName, err := utl.ModuleName()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		name, err := PrompStr("Name")
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		schema, err := PrompStr("Schema")
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
		name, err := PrompStr("Package")
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		modelo, err := PrompStr("Model")
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		schema, err := PrompStr("Schema")
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		err = MkMolue(name, modelo, schema)
		if err != nil {
			fmt.Printf("Command failed %v\n", err)
			return
		}

		title := utl.Titlecase(name)
		message := fmt.Sprintf(`Remember, including the router, that it is on the bottom of the h%s.go, in routers section of the router.go file`, title)
		fmt.Println(message)
	},
}

var CmdRpc = &cobra.Command{
	Use:   "rpc [name]",
	Short: "Create rpc model to microservice.",
	Long:  "Template rpc model to microservice include function handler model.",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := PrompStr("Package")
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
