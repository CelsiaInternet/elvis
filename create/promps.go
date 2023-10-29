package create

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/manifoldco/promptui"
)

func prompCreate() {
	prompt := promptui.Select{
		Label: "What do you want created?",
		Items: []string{"Project", "Microservice", "Modelo", "Rpc"},
	}

	opt, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch opt {
	case 0:
		err := CmdProject.Execute()
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
	case 3:
		err := CmdRpc.Execute()
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
