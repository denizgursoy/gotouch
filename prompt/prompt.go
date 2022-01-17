package prompt

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
)

type Definition struct {
	ErrorText string
	Direction string
}

func AskForSelection(definition Definition, options []string) (result string) {
	index := -1
	var err error

	for index < 0 {
		prompt := promptui.Select{
			Label: definition.Direction,
			Items: options,
		}

		index, result, err = prompt.Run()
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	//fmt.Printf("Input: %s\n", result)
	return
}
