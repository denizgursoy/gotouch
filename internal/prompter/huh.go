package prompter

import (
	"errors"
	"fmt"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// promptRequest represents a prompt to be displayed to the user.
type promptRequest struct {
	form *huh.Form
	done chan error
}

// Messages used by the persistent tea.Program.
type (
	newPromptMsg struct{ req promptRequest }
	formDoneMsg  struct{}
	formAbortMsg struct{}
	shutdownMsg  struct{}
)

// persistentModel manages a single alt-screen session across all prompts.
type persistentModel struct {
	reqCh   chan promptRequest
	current *promptRequest
	form    *huh.Form
}

func newTheme() *huh.Theme {
	theme := huh.ThemeDracula()
	theme.Focused.SelectedPrefix = lipgloss.NewStyle().SetString("[✓] ")
	theme.Focused.UnselectedPrefix = lipgloss.NewStyle().SetString("[ ] ")
	theme.Blurred.SelectedPrefix = lipgloss.NewStyle().SetString("[✓] ")
	theme.Blurred.UnselectedPrefix = lipgloss.NewStyle().SetString("[ ] ")
	return theme
}

func (m persistentModel) Init() tea.Cmd {
	return waitForRequest(m.reqCh)
}

func waitForRequest(ch chan promptRequest) tea.Cmd {
	return func() tea.Msg {
		req, ok := <-ch
		if !ok {
			return shutdownMsg{}
		}
		return newPromptMsg{req: req}
	}
}

func (m persistentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case newPromptMsg:
		m.current = &msg.req
		m.form = msg.req.form
		return m, m.form.Init()

	case formDoneMsg:
		if m.current != nil {
			m.current.done <- nil
			m.current = nil
			m.form = nil
		}
		return m, waitForRequest(m.reqCh)

	case formAbortMsg:
		if m.current != nil {
			m.current.done <- errors.New("prompt cancelled")
			m.current = nil
			m.form = nil
		}
		return m, waitForRequest(m.reqCh)

	case shutdownMsg:
		// Clean up any pending request
		if m.current != nil {
			m.current.done <- errors.New("prompter closed")
			m.current = nil
			m.form = nil
		}
		return m, tea.Quit

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			if m.current != nil {
				m.current.done <- errors.New("interrupted")
				m.current = nil
				m.form = nil
			}
			return m, tea.Quit
		}
		if m.form != nil {
			return m.updateForm(msg)
		}
		return m, nil

	default:
		if m.form != nil {
			return m.updateForm(msg)
		}
		return m, nil
	}
}

func (m persistentModel) updateForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	updatedForm, cmd := m.form.Update(msg)
	m.form = updatedForm.(*huh.Form)

	switch m.form.State {
	case huh.StateCompleted:
		if m.current != nil {
			m.current.done <- nil
			m.current = nil
		}
		m.form = nil
		return m, waitForRequest(m.reqCh)
	case huh.StateAborted:
		if m.current != nil {
			m.current.done <- errors.New("prompt cancelled")
			m.current = nil
		}
		m.form = nil
		return m, waitForRequest(m.reqCh)
	default:
		return m, cmd
	}
}

func (m persistentModel) View() string {
	if m.form != nil {
		return m.form.View()
	}
	return ""
}

// programManager manages the lifecycle of the persistent tea.Program.
var programMgr = &programManager{}

type programManager struct {
	mu      sync.Mutex
	reqCh   chan promptRequest
	program *tea.Program
	started bool
	done    chan struct{}
}

func (pm *programManager) ensureRunning() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if pm.started {
		return
	}

	pm.reqCh = make(chan promptRequest)
	pm.done = make(chan struct{})
	model := persistentModel{reqCh: pm.reqCh}

	pm.program = tea.NewProgram(
		model,
		tea.WithAltScreen(),
	)

	pm.started = true

	go func() {
		_, _ = pm.program.Run()
		close(pm.done)
	}()
}

func (pm *programManager) close() {
	pm.mu.Lock()

	if !pm.started {
		pm.mu.Unlock()
		return
	}

	close(pm.reqCh)
	pm.started = false
	done := pm.done
	pm.program = nil
	pm.mu.Unlock()

	// Wait for bubbletea to restore terminal state (cursor, alt screen, etc.)
	select {
	case <-done:
	case <-time.After(2 * time.Second):
	}
}

func (pm *programManager) sendPrompt(form *huh.Form) error {
	pm.ensureRunning()

	done := make(chan error, 1)
	pm.reqCh <- promptRequest{form: form, done: done}
	return <-done
}

// runPrompt creates a huh.Form with custom submit/cancel commands and sends it
// to the persistent program.
func runPrompt(fields ...huh.Field) error {
	form := huh.NewForm(huh.NewGroup(fields...)).
		WithTheme(newTheme())

	// Override submit/cancel commands to signal our persistent model
	// instead of quitting the program.
	form.SubmitCmd = func() tea.Msg { return formDoneMsg{} }
	form.CancelCmd = func() tea.Msg { return formAbortMsg{} }

	return programMgr.sendPrompt(form)
}

func (s srv) AskForString(direction, initialValue string, validator Validator) (string, error) {
	result := initialValue

	input := huh.NewInput().
		Title(direction).
		Value(&result)

	if validator != nil {
		input.Validate(func(val string) error {
			return validator(val)
		})
	}

	err := runPrompt(input)
	return result, err
}

func (s srv) AskForSelectionFromList(direction string, list []fmt.Stringer) (any, error) {
	count := len(list)

	if count == 0 {
		return nil, EmptyList
	}

	options := make([]huh.Option[string], 0, count)
	lookup := make(map[string]fmt.Stringer)
	for _, item := range list {
		choice := item.String()
		options = append(options, huh.NewOption(choice, choice))
		lookup[choice] = item
	}

	var selected string
	err := runPrompt(
		huh.NewSelect[string]().
			Title(direction + " (select one)").
			Options(options...).
			Value(&selected),
	)

	return lookup[selected], err
}

func (s srv) AskForMultipleSelectionFromList(direction string, list []fmt.Stringer) ([]any, error) {
	count := len(list)

	if count == 0 {
		return nil, EmptyList
	}

	options := make([]huh.Option[string], 0, count)
	lookup := make(map[string]fmt.Stringer)
	for _, item := range list {
		choice := item.String()
		options = append(options, huh.NewOption(choice, choice))
		lookup[choice] = item
	}

	var selected []string
	err := runPrompt(
		huh.NewMultiSelect[string]().
			Title(direction + " (select multiple)").
			Options(options...).
			Value(&selected),
	)

	results := make([]any, 0, len(selected))
	for _, s := range selected {
		results = append(results, lookup[s])
	}

	return results, err
}

func (s srv) AskForYesOrNo(direction string) (bool, error) {
	var result bool
	err := runPrompt(
		huh.NewConfirm().
			Title(direction).
			Inline(false).
			WithButtonAlignment(lipgloss.Left).
			Value(&result),
	)
	return result, err
}

func (s srv) AskForMultilineString(direction, defaultValue, pattern string) (string, error) {
	result := defaultValue
	err := runPrompt(
		huh.NewText().
			Title(direction).
			Value(&result),
	)
	return result, err
}

func (s srv) Close() {
	programMgr.close()
}
