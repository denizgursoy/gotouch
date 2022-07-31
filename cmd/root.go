/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/denizgursoy/gotouch/lister"
	"github.com/manifoldco/promptui"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

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
		//prompForProjectAddress()
		//promptForDependencies()
		//promptForProjectStructure()
		// select structure
		projectName := prompForProjectAddress()
		selectedProject := promptForProjectStructure()

		extractFiles(selectedProject, projectName)
	},
}

var projectName = "go-test"

func prompForProjectAddress() string {
	p := promptui.Prompt{
		Label: "Enter Project address",
	}

	projectName, _ = p.Run()
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

func extractFiles(selectedProject *lister.ProjectData, name string) {
	client := http.Client{}
	response, err := client.Get(selectedProject.URL)
	if err != nil {
		println(err)
		return
	}
	filePath2 := filepath.Join(os.TempDir(), filepath.Base(selectedProject.URL))
	println(filePath2)
	create, err := os.Create(filePath2)
	_, err = io.Copy(create, response.Body)
	err = UnGzip(filePath2, projectName+string(filepath.Separator))
	if err != nil {
		log.Fatal(err)
	}

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

func UnGzip(source, target string) error {
	reader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	archive, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer archive.Close()

	target = filepath.Join(target, archive.Name)
	writer, err := os.Create(target)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = io.Copy(writer, archive)
	return err
}

func Untar(tarball, target string) error {
	reader, err := os.Open(tarball)
	if err != nil {
		return err
	}
	defer reader.Close()
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path := filepath.Join(target, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}
	}
	return nil
}
