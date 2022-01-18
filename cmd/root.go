/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/denizgursoy/gotouch/common"

	"github.com/denizgursoy/gotouch/prompt"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gotouch",
	Short: "Helps you create new Go Projects",
	Long:  `Tag`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("executing the command")
		for _, dependency := range common.AppConfig.Dependencies {
			availableValues := make([]string, 0)

			for _, option := range dependency.Options {
				availableValues = append(availableValues, option.Name)
			}

			_ = prompt.AskForSelection(prompt.Definition{
				//ErrorText: "Please select a HTTP Framework",
				Direction: dependency.Prompt,
			}, availableValues)

			//fmt.Println("User selected ->" + selectedValue)
		}
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
