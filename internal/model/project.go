package model

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/denizgursoy/gotouch/internal/langs"
)

type (
	ProjectStructureData struct {
		Name       string      `yaml:"name"`
		Reference  string      `yaml:"reference"`
		URL        string      `yaml:"url"`
		Questions  []*Question `yaml:"questions"`
		Values     interface{} `yaml:"values"`
		Language   string      `yaml:"language"`
		Delimiters string      `yaml:"delimiters"`
	}

	Question struct {
		Direction         string    `yaml:"direction"`
		CanSkip           bool      `yaml:"canSkip"`
		CanSelectMultiple bool      `yaml:"canSelectMultiple"`
		Choices           []*Choice `yaml:"choices"`
	}

	Choice struct {
		Choice       string                 `yaml:"choice"`
		Dependencies []interface{}          `yaml:"dependencies"`
		Files        []*File                `yaml:"files"`
		Values       map[string]interface{} `yaml:"values"`
	}

	File struct {
		Url          string `yaml:"url"`
		Content      string `yaml:"content"`
		PathFromRoot string `yaml:"pathFromRoot"`
	}
)

func (p *ProjectStructureData) String() string {
	projectName := p.Name
	if len(strings.TrimSpace(p.Language)) != 0 {
		projectName += fmt.Sprintf(" ( %s )", strings.Title(p.Language))
	}
	if len(strings.TrimSpace(p.Reference)) != 0 {
		projectName += fmt.Sprintf(" ( %s )", p.Reference)
	}
	return projectName
}

var (
	ErrProjectURLIsEmpty    = errors.New("project url can not be empty")
	ErrProjectNameIsEmpty   = errors.New("project name can not be empty")
	ErrProjectURLIsNotValid = errors.New("project url invalid")
)

func (o *Choice) String() string {
	return o.Choice
}

func (p *ProjectStructureData) IsValid() error {
	if len(strings.TrimSpace(p.Name)) == 0 {
		return ErrProjectNameIsEmpty
	}

	projectUrl := strings.TrimSpace(p.URL)
	if len(projectUrl) == 0 {
		return ErrProjectURLIsEmpty
	}

	if _, err := url.ParseRequestURI(projectUrl); err != nil {
		return ErrProjectURLIsNotValid
	}

	delimiters := strings.Fields(p.Delimiters)
	if len(delimiters) != 0 && len(delimiters) != 2 {
		return ErrWrongDelimiterFormat{projectName: p.Name}
	}

	if len(p.Questions) > 0 {
		for questionIndex, q := range p.Questions {
			err := p.validateQuestion(q, questionIndex)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *ProjectStructureData) validateQuestion(q *Question, questionIndex int) error {
	if len(strings.TrimSpace(q.Direction)) == 0 {
		return ErrEmptyQuestionField{
			projectName: p.Name,
			index:       questionIndex,
			field:       "Direction",
		}
	}
	if len(q.Choices) == 0 {
		return ErrEmptyQuestionField{
			projectName: p.Name,
			index:       questionIndex,
			field:       "Choices",
		}
	} else if len(q.Choices) == 1 && q.CanSelectMultiple == true {
		return ErrCanSelectMultiple{
			projectName: p.Name,
			index:       questionIndex,
		}
	} else if len(q.Choices) == 1 && q.CanSkip == false {
		return ErrCanSkip{
			projectName: p.Name,
			index:       questionIndex,
		}
	}
	for choiceIndex, choice := range q.Choices {
		err := p.validateChoice(choice, questionIndex, choiceIndex)
		if err != nil {
			return err
		}

	}
	return nil
}

func (p *ProjectStructureData) validateChoice(choice *Choice, questionIndex int, choiceIndex int) error {
	if len(strings.TrimSpace(choice.Choice)) == 0 {
		return ErrEmptyChoice{
			projectName:   p.Name,
			questionIndex: questionIndex,
			choiceIndex:   choiceIndex,
			field:         "Choice",
		}
	}

	if len(choice.Files) == 0 && len(choice.Dependencies) == 0 && len(choice.Values) == 0 {
		return ErrEmptyFileAndDependency{
			projectName:   p.Name,
			questionIndex: questionIndex,
			choiceIndex:   choiceIndex,
		}
	}

	if len(choice.Files) > 0 {
		for k, file := range choice.Files {

			if len(strings.TrimSpace(file.PathFromRoot)) == 0 {
				return ErrEmptyFileField{
					projectName:   p.Name,
					questionIndex: questionIndex,
					choiceIndex:   choiceIndex,
					fileIndex:     k,
					field:         "PathFromRoot",
				}
			}

			if len(strings.TrimSpace(file.Url)) == 0 && len(strings.TrimSpace(file.Content)) == 0 {
				return ErrEmptyUrlAndContent{
					projectName:   p.Name,
					questionIndex: questionIndex,
					choiceIndex:   choiceIndex,
					fileIndex:     k,
				}
			}

			if len(strings.TrimSpace(file.Url)) > 0 && len(strings.TrimSpace(file.Content)) > 0 {
				return ErrMultipleFieldUrlAndContent{
					projectName:   p.Name,
					questionIndex: questionIndex,
					choiceIndex:   choiceIndex,
					fileIndex:     k,
				}
			}

			if len(file.Content) == 0 {
				fileUrl := strings.TrimSpace(file.Url)

				if len(fileUrl) == 0 {
					return ErrEmptyFileField{
						projectName:   p.Name,
						questionIndex: questionIndex,
						choiceIndex:   choiceIndex,
						fileIndex:     k,
						field:         "URL",
					}
				}

				if _, err := url.ParseRequestURI(fileUrl); err != nil {
					return ErrInvalidURLFile{
						projectName:   p.Name,
						questionIndex: questionIndex,
						choiceIndex:   choiceIndex,
						fileIndex:     k,
					}
				}
			} else {
				content := strings.TrimSpace(file.Content)

				if len(content) == 0 {
					return ErrEmptyFileField{
						projectName:   p.Name,
						questionIndex: questionIndex,
						choiceIndex:   choiceIndex,
						fileIndex:     k,
						field:         "Content",
					}
				}
			}
		}
	}

	if len(choice.Dependencies) > 0 {

		for dependencyIndex, dependency := range choice.Dependencies {
			if err := langs.GetChecker(p.Language, nil, nil, nil).CheckDependency(dependency); err != nil {
				return ErrWrongDependencyFormat{
					projectName:     p.Name,
					questionIndex:   questionIndex,
					choiceIndex:     choiceIndex,
					dependencyIndex: dependencyIndex,
				}
			}
		}
	}
	return nil
}

type (
	ErrEmptyQuestionField struct {
		projectName string
		index       int
		field       string
	}
	ErrCanSelectMultiple struct {
		projectName string
		index       int
	}
	ErrCanSkip struct {
		projectName string
		index       int
	}
	ErrEmptyChoice struct {
		projectName   string
		questionIndex int
		choiceIndex   int
		field         string
	}
	ErrEmptyFileAndDependency struct {
		projectName   string
		questionIndex int
		choiceIndex   int
	}
	ErrEmptyUrlAndContent struct {
		projectName   string
		questionIndex int
		choiceIndex   int
		fileIndex     int
	}
	ErrMultipleFieldUrlAndContent struct {
		projectName   string
		questionIndex int
		choiceIndex   int
		fileIndex     int
	}
	ErrWrongDependencyFormat struct {
		projectName     string
		questionIndex   int
		choiceIndex     int
		dependencyIndex int
	}
	ErrEmptyFileField struct {
		projectName   string
		questionIndex int
		choiceIndex   int
		fileIndex     int
		field         string
	}
	ErrInvalidURLFile struct {
		projectName   string
		questionIndex int
		choiceIndex   int
		fileIndex     int
	}
	ErrWrongDelimiterFormat struct {
		projectName string
	}
)

func (e ErrEmptyQuestionField) Error() string {
	return fmt.Sprintf("%s's %d. %s is empty. %s can not be empty", e.projectName, e.index+1, e.field, e.field)
}

func (e ErrCanSelectMultiple) Error() string {
	return fmt.Sprintf("%s's %d. question 'CanSelectMultiple' field is 'true'. This field can not be true if there are less than two choices", e.projectName, e.index+1)
}

func (e ErrCanSkip) Error() string {
	return fmt.Sprintf("%s's %d. question 'CanSkip' field is 'false'. This field can not be false if there is only one choice", e.projectName, e.index+1)
}

func (e ErrEmptyChoice) Error() string {
	return fmt.Sprintf("%s's %d. question, %d. choice %s is empty. %s can not be empty", e.projectName, e.questionIndex+1, e.choiceIndex+1, e.field, e.field)
}

func (e ErrEmptyFileAndDependency) Error() string {
	return fmt.Sprintf("%s's %d. question %d. choice does not have any files,values or dependencies. Choice must have at least one file, value or dependency", e.projectName, e.questionIndex+1, e.choiceIndex+1)
}

func (e ErrEmptyUrlAndContent) Error() string {
	return fmt.Sprintf("%s's %d. question %d. choice %d. file do not have both URL and Content. File must have at least one URL or Content", e.projectName, e.questionIndex+1, e.choiceIndex+1, e.fileIndex+1)
}

func (e ErrWrongDependencyFormat) Error() string {
	return fmt.Sprintf("%s's %d. question %d. choice %d. dependecies' format is not correct", e.projectName, e.questionIndex+1, e.choiceIndex+1, e.dependencyIndex+1)
}

func (e ErrMultipleFieldUrlAndContent) Error() string {
	return fmt.Sprintf("%s's %d. question %d. choice %d. file have both URL and Content. File can not have both URL and Content at the same time", e.projectName, e.questionIndex+1, e.choiceIndex+1, e.fileIndex+1)
}

func (e ErrEmptyFileField) Error() string {
	return fmt.Sprintf("%s's %d. question, %d. choice, %d. file %s is empty. %s can not be empty", e.projectName, e.questionIndex+1, e.choiceIndex+1, e.fileIndex+1, e.field, e.field)
}

func (e ErrInvalidURLFile) Error() string {
	return fmt.Sprintf("%s's %d. question, %d. choice, %d. file URL is invalid.", e.projectName, e.questionIndex+1, e.choiceIndex+1, e.fileIndex+1)
}

func (e ErrWrongDelimiterFormat) Error() string {
	return fmt.Sprintf("%s's delimiter must be seperated by space as '[[ ]]'", e.projectName)
}
