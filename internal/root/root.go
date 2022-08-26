/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package root

import (
	"github.com/denizgursoy/gotouch/internal/compressor"
	"github.com/denizgursoy/gotouch/internal/executor"
	"github.com/denizgursoy/gotouch/internal/lister"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/prompter"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

const (
	FileFlagName = "file"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gotouch",
	Short: "Helps you create new Go Projects",
	Long:  `Tag`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		filePath, inputError := flags.GetString(FileFlagName)
		if inputError != nil {
			log.Fatalln(inputError)
		}

		point := &filePath
		if len(strings.TrimSpace(filePath)) == 0 {
			point = nil
		}

		options := CreateCommandOptions{
			Lister:     lister.GetInstance(),
			Prompter:   prompter.GetInstance(),
			Manager:    manager.GetInstance(),
			Compressor: compressor.GetInstance(),
			Executor:   executor.GetInstance(),
			Path:       point,
		}
		err := CreateNewProject(&options)
		log.Fatalln(err)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP(FileFlagName, "f", "", "input file")
}
