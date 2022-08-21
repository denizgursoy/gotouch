/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package root

import (
	"github.com/denizgursoy/gotouch/internal/lister"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/operation"
	"github.com/denizgursoy/gotouch/internal/prompts"
	"github.com/denizgursoy/gotouch/internal/uncompressor"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gotouch",
	Short: "Helps you create new Go Projects",
	Long:  `Tag`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		options := CreateNewProjectOptions{
			lister:       lister.GetInstance(),
			prompter:     prompts.GetInstance(),
			manager:      manager.GetInstance(),
			uncompressor: uncompressor.GetInstance(),
			executor:     operation.GetInstance(),
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gotouch.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
