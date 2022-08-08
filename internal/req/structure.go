package req

import (
	"fmt"
	"github.com/denizgursoy/gotouch/internal/lister"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompts"
	"github.com/denizgursoy/gotouch/internal/uncompressor"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type (
	ProjectStructureRequirement struct {
		ProjectsData []*lister.ProjectStructureData
	}

	projectStructureTask struct {
		ProjectStructure *lister.ProjectStructureData
	}
)

func (p ProjectStructureRequirement) AskForInput() (model.Task, error) {

	instance := prompts.GetInstance()

	projectList := make([]*prompts.ListOption, 0)

	for _, project := range p.ProjectsData {
		displayText := fmt.Sprintf("%s (%s)", project.Name, project.Reference)
		projectList = append(projectList, &prompts.ListOption{
			DisplayText: displayText,
			ReturnVal:   project,
		})
	}

	selected := instance.AskForSelectionFromList("select project type", projectList).(*lister.ProjectStructureData)
	return projectStructureTask{
		ProjectStructure: selected,
	}, nil
}

func (p projectStructureTask) Complete(previousResponse interface{}) interface{} {
	projectName := previousResponse.(string)
	ex := uncompressor.GetInstance()
	ex.UncompressFromUrl(p.ProjectStructure.URL, projectName)
	modAvailable := checkGoModule(projectName)
	editGoModule(modAvailable, projectName)
	return nil
}

func checkGoModule(projectName string) bool {
	path := fmt.Sprintf("./%s/", projectName)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.Name() == "go.mod" {
			return true
		}
	}
	return false
}

func editGoModule(modAvailable bool, projectName string) {
	dir := fmt.Sprintf("./%s", projectName)

	err := os.Chdir(dir)
	if err != nil {
		log.Printf("Changed Directory error: %v", err)
	}

	command := "init"

	if modAvailable {
		command = "edit"
	}
	command = fmt.Sprintf("go mod %s %s", command, projectName)
	runCommand(command)
}

func runCommand(command string) error {
	cmd := exec.Command("bash", "-c", command)
	err := cmd.Run()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
		return err
	}
	return nil
}
