package create

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/cgalvisleon/elvis/utilities"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func prompCreate() {
	prompt := promptui.Select{
		Label: "What do you want created?",
		Items: []string{"Project", "Microservice", "Modelo"},
	}

	opt, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch opt {
	case 0:
		err := CmdPMicro.Execute()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
	case 1:
		err := CmdMicro.Execute()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
	case 2:
		err := CmdModelo.Execute()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
	}
}

func prompStr(label string) (string, error) {
	validate := func(input string) error {
		if len(input) == 0 {
			return errors.New(fmt.Sprintf("Invalid %s", label))
		}

		return nil
	}

	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		return "", err
	}

	return result, nil
}

func prompInt(label string) (int, error) {
	validate := func(input string) error {
		_, err := strconv.Atoi(input)
		if err != nil {
			return errors.New(fmt.Sprintf("Invalid %s", label))
		}

		return nil
	}

	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}

	value, err := prompt.Run()

	if err != nil {
		return 0, err
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return result, nil
}

var CmdCreate = &cobra.Command{
	Use:   "create",
	Short: "You can created Microservice.",
	Long:  "Template project to create microservice include required folders and basic files.",
	Run: func(cmd *cobra.Command, args []string) {
		prompCreate()
	},
}

var CmdPMicro = &cobra.Command{
	Use:   "micro [name author schema, schema_var]",
	Short: "Create project base type microservice.",
	Long:  "Template project to microservice include folder cmd, deployments, pkg, rest, test and web, with files .go required for making a microservice.",
	Run: func(cmd *cobra.Command, args []string) {
		packageName, err := utilities.ModuleName()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		name, err := prompStr("Name")
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		author, err := prompStr("Author")
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		schema, err := prompStr("Schema")
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
		packageName, err := utilities.ModuleName()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		name, err := prompStr("Name")
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		schema, err := prompStr("Schema")
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
		name, err := prompStr("Package")
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		modelo, err := prompStr("Model")
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		schema, err := prompStr("Schema")
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		err = MkMolue(name, modelo, schema)
		if err != nil {
			fmt.Printf("Command failed %v\n", err)
			return
		}
	},
}
