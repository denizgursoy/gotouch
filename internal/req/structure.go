package req

import (
	"fmt"
	"github.com/denizgursoy/gotouch/internal/lister"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompts"
	"github.com/denizgursoy/gotouch/internal/queue"
	"github.com/denizgursoy/gotouch/internal/util"
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
	path, err := util.GetBaseName(projectName)
	if err != nil {
		log.Printf("%v", err)
	}

	queue.Extractor.UncompressFromUrl(p.ProjectStructure.URL, path)
	editGoModule(projectName, path)
	return nil
}

func hasGoModule(path string) bool {
	path = fmt.Sprintf("./%s/", path)
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

func editGoModule(projectName, path string) {
	modAvailable := hasGoModule(path)

	dir := fmt.Sprintf("./%s", path)
	err := os.Chdir(dir)
	if err != nil {
		log.Printf("Changed Directory error: %v", err)
	}

	command := "init"

	if modAvailable {
		command = "edit -module"
	}
	command = fmt.Sprintf("go mod %s %s", command, projectName)
	runCommand(command)
}

func runCommand(command string) error {
	//TODO: windows testini yap
	cmd := exec.Command("bash", "-c", command)
	err := cmd.Run()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
		return err
	}
	return nil
}
