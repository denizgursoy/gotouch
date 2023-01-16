package template

import (
	"bytes"
	"fmt"
	"io"
	textTemplate "text/template"

	"github.com/Masterminds/sprig/v3"
)

type Template struct {
	template *textTemplate.Template
}

func New() *Template {
	return &Template{
		template: textTemplate.New("txt"),
	}
}

func (t *Template) Reset() {
	t.template = textTemplate.New("txt")
}

func (t *Template) SetDelims(left, right string) {
	if left == "" {
		left = "{{"
	}

	if right == "" {
		right = "}}"
	}

	t.template.Delims(left, right)
}

func (t *Template) SetSprigFuncs() {
	t.template.Funcs(sprig.TxtFuncMap())
}

func (t *Template) Execute(v any, content string) (string, error) {
	var b bytes.Buffer
	// Execute the template and write the output to the buffer
	if err := textTemplate.Must(t.template.Parse(content)).Execute(&b, v); err != nil {
		return "", fmt.Errorf("execute error: %w", err)
	}

	return b.String(), nil
}

func (t *Template) ExecuteContent(writer io.Writer, v any, content []byte) error {
	// Execute the template and write the output to the buffer
	txtTemplate, err := t.template.Parse(string(content))
	if err != nil {
		return err
	}
	if err := txtTemplate.Execute(writer, v); err != nil {
		return fmt.Errorf("ExecuteContent error: %w", err)
	}

	return nil
}
