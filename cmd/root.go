/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/denizgursoy/gotouch/common"
	"github.com/denizgursoy/gotouch/extractor"
	"github.com/denizgursoy/gotouch/lister"
	"github.com/manifoldco/promptui"
	"os"

	"github.com/denizgursoy/gotouch/prompt"

	"github.com/spf13/cobra"
)

var (
	ex = extractor.GetInstance()
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gotouch",
	Short: "Helps you create new Go Projects",
	Long:  `Tag`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		projectName := prompForProjectAddress()
		projectStructure := promptForProjectStructure()
		ex.Extract(projectStructure.URL, projectName)
	},
}

func prompForProjectAddress() string {
	p := promptui.Prompt{
		Label: "Enter Project name",
	}

	projectName, _ := p.Run()
	return projectName
}

func promptForDependencies() {
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
}

func promptForProjectStructure() *lister.ProjectData {

	options := make([]string, 0)
	projectLister := lister.GetInstance()

	projects := projectLister.GetDefaultProjects()
	for _, project := range projects {
		projectString := fmt.Sprintf("%s (%s)", project.Name, project.Reference)
		options = append(options, projectString)
	}

	selectedName := prompt.AskForSelection(prompt.Definition{
		Direction: "Select the project structure",
	}, options)

	var selectedProject *lister.ProjectData

	for _, project := range projects {
		projectString := fmt.Sprintf("%s (%s)", project.Name, project.Reference)
		if projectString == selectedName {
			selectedProject = project
			break
		}
	}
	return selectedProject
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
