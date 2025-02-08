//go:generate mockgen -source=./prompter.go -destination=mockPrompter.go -package=prompter

package prompter

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	once      = sync.Once{}
	prompter  Prompter
	EmptyList = errors.New("options can not be empty")
)

type (
	Prompter interface {
		AskForString(direction string, validator Validator) (string, error)
		AskForSelectionFromList(direction string, list []fmt.Stringer) (any, error)
		AskForMultipleSelectionFromList(direction string, list []fmt.Stringer) ([]any, error)
		AskForYesOrNo(direction string) (bool, error)
		AskForMultilineString(direction, defaultValue, pattern string) (string, error)
	}

	ListOption struct {
		DisplayText string
		ReturnVal   any
	}

	Validator func(any) error
)

func GetInstance() Prompter {
	once.Do(func() {
		prompter = &terminalPrompter{}
	})
	return prompter
}

type terminalPrompter struct {
	prompt        string
	choices       []fmt.Stringer
	selectedIndex int
	inputBuffer   string
	state         interfaceState
	err           error
}

type interfaceState int

const (
	stateIdle interfaceState = iota
	stateAskString
	stateAskSelection
	stateAskMultipleSelection
	stateAskYesNo
	stateAskMultiline
)

func (tp *terminalPrompter) AskForString(direction string, validator Validator) (string, error) {
	tp.state = stateAskString
	tp.prompt = direction
	tp.inputBuffer = ""
	return tp.run()
}

func (tp *terminalPrompter) AskForSelectionFromList(direction string, list []fmt.Stringer) (any, error) {
	tp.state = stateAskSelection
	tp.prompt = direction
	tp.choices = list
	tp.selectedIndex = 0
	return tp.runSelect()
}

func (tp *terminalPrompter) AskForMultipleSelectionFromList(direction string, list []fmt.Stringer) ([]any, error) {
	tp.state = stateAskMultipleSelection
	tp.prompt = direction
	tp.choices = list
	multiSelect, err := tp.runMultiSelect()
	anyArray := make([]any, 0)
	for _, stringer := range multiSelect {
		anyArray = append(anyArray, stringer)
	}
	return anyArray, err
}

func (tp *terminalPrompter) AskForYesOrNo(direction string) (bool, error) {
	tp.state = stateAskYesNo
	tp.prompt = direction
	return tp.runYesNo()
}

func (tp *terminalPrompter) AskForMultilineString(direction, defaultValue, pattern string) (string, error) {
	tp.state = stateAskMultiline
	tp.prompt = direction
	tp.inputBuffer = defaultValue
	return tp.run()
}

func (tp *terminalPrompter) Init() tea.Cmd {
	return nil
}

func (tp *terminalPrompter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			tp.err = fmt.Errorf("operation cancelled")
			return tp, tea.Quit
		case "enter":
			return tp, tea.Quit
		case "up":
			if tp.state == stateAskSelection && tp.selectedIndex > 0 {
				tp.selectedIndex--
			}
		case "down":
			if tp.state == stateAskSelection && tp.selectedIndex < len(tp.choices)-1 {
				tp.selectedIndex++
			}
		default:
			tp.inputBuffer = msg.String()
		}
	}
	return tp, nil
}

func (tp *terminalPrompter) View() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s\n", tp.prompt))
	switch tp.state {
	case stateAskString, stateAskMultiline:
		b.WriteString(fmt.Sprintf("Input: %s", tp.inputBuffer))
	case stateAskSelection:
		for i, choice := range tp.choices {
			cursor := " "
			if i == tp.selectedIndex {
				cursor = ">"
			}
			b.WriteString(fmt.Sprintf("%s %d. %s\n", cursor, i+1, choice))
		}
	case stateAskMultipleSelection:
		// Display logic can be similar to stateAskSelection
	case stateAskYesNo:
		b.WriteString("Yes or No? (y/n)")
	}
	return b.String()
}

func (tp *terminalPrompter) run() (string, error) {
	p := tea.NewProgram(tp)
	if err := p.Start(); err != nil {
		return "", err
	}
	return tp.inputBuffer, tp.err
}

func (tp *terminalPrompter) runSelect() (fmt.Stringer, error) {
	p := tea.NewProgram(tp)
	if err := p.Start(); err != nil {
		return nil, err
	}
	return tp.choices[tp.selectedIndex], tp.err
}

func (tp *terminalPrompter) runMultiSelect() ([]fmt.Stringer, error) {
	// Implement multi-select logic here
	return nil, nil
}

func (tp *terminalPrompter) runYesNo() (bool, error) {
	p := tea.NewProgram(tp)
	if err := p.Start(); err != nil {
		return false, err
	}
	return true, nil // Adjust logic based on captured input
}
