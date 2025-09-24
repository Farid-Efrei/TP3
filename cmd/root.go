package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "loganalyzer",
	Short: "Analyse des logs en parallèle avec Export JSON",
	Long:  "loganalyzer - outil CLI pour analyser plusieurs fichiers de logs en parallèle (goroutines).",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	fmt.Print("Initialisation de l'application loganalyzer\n")
}
