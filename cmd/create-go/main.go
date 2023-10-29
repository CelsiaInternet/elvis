package main

import (
	"fmt"
	"runtime"

	"github.com/cgalvisleon/elvis/create"
	"github.com/spf13/cobra"
)

func main() {
	if runtime.GOOS == "windows" {
		fmt.Println("This program is not aplicable for windows machine.")
	} else {
		var rootCmd = &cobra.Command{Use: "create-go"}
		rootCmd.AddCommand(create.Create)
		rootCmd.Execute()
	}
}
