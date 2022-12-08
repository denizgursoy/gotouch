//go:generate mockgen -source=./runner.go -destination=mockCommandrunner.go -package=commandrunner

package commandrunner

import (
	"fmt"
	"github.com/denizgursoy/gotouch/internal/store"
	"os"
	"os/exec"
	"sync"
)

var (
	once      sync.Once
	cmdRunner Runner
)

type (
	Runner interface {
		Run(data *CommandData) error
	}

	runner struct {
		Store store.Store
	}

	CommandData struct {
		WorkingDir *string
		Command    string
		Args       []string
	}
)

func GetInstance(str store.Store) Runner {
	once.Do(func() {
		cmdRunner = &runner{
			Store: str,
		}
	})
	return cmdRunner
}

func (r *runner) Run(data *CommandData) error {
	if data.WorkingDir == nil {
		projectFullPath := r.Store.GetValue(store.ProjectFullPath)
		err := os.Chdir(projectFullPath)
		if err != nil {
			return err
		}
	}
	cmd := exec.Command(data.Command, data.Args...)

	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))
	if err != nil {
		//	log.Printf("Command finished with error: %v", err)
		return err
	}
	return nil
}
