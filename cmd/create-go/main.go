package main

import (
	"github.com/celsiainternet/elvis/create"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "create-go"}
	rootCmd.AddCommand(create.Create)
	rootCmd.Execute()
}
