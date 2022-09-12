package model

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type (
	ProjectStructureData struct {
		Name      string      `yaml:"name"`
		Reference string      `yaml:"reference"`
		URL       string      `yaml:"url"`
		Questions []*Question `yaml:"questions"`
		Values    interface{} `yaml:"values"`
	}

	Question struct {
		Direction         string    `yaml:"direction"`
		CanSkip           bool      `yaml:"canSkip"`
		CanSelectMultiple bool      `yaml:"canSelectMultiple"`
		Options           []*Option `yaml:"options"`
	}

	Option struct {
		Choice       string                 `yaml:"choice"`
		Dependencies []*string              `yaml:"dependencies"`
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
	return fmt.Sprintf("%s (%s)", p.Name, p.Reference)
}

var (
	ErrProjectURLIsEmpty    = errors.New("project url can not be empty")
	ErrProjectNameIsEmpty   = errors.New("project name can not be empty")
	ErrProjectURLIsNotValid = errors.New("project url invalid")
)

func (o *Option) String() string {
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

	if len(p.Questions) > 0 {
		for i, q := range p.Questions {
			if len(strings.TrimSpace(q.Direction)) == 0 {
				return ErrEmptyQuestionField{
					index: i,
					field: "Direction",
				}
			}
			if len(q.Options) == 0 {
				return ErrEmptyQuestionField{
					index: i,
					field: "Options",
				}
			} else if len(q.Options) == 1 && q.CanSelectMultiple == true {
				return ErrCanSelectMultiple{index: i}
			} else if len(q.Options) == 1 && q.CanSkip == false {
				return ErrCanSkip{index: i}
			}
			for j, option := range q.Options {
				if len(strings.TrimSpace(option.Choice)) == 0 {
					return ErrEmptyOption{
						questionIndex: i,
						optionIndex:   j,
						field:         "Choice",
					}
				}

				if len(option.Files) == 0 && len(option.Dependencies) == 0 {
					return ErrEmptyFileAndDependency{
						questionIndex: i,
						optionIndex:   j,
					}
				}

				if len(option.Files) > 0 {
					for k, file := range option.Files {

						if len(strings.TrimSpace(file.PathFromRoot)) == 0 {
							return ErrEmptyFileField{
								questionIndex: i,
								optionIndex:   j,
								fileIndex:     k,
								field:         "PathFromRoot",
							}
						}

						if len(strings.TrimSpace(file.Url)) == 0 && len(strings.TrimSpace(file.Content)) == 0 {
							return ErrEmptyUrlAndContent{
								questionIndex: i,
								optionIndex:   j,
								fileIndex:     k,
							}
						}

						if len(strings.TrimSpace(file.Url)) > 0 && len(strings.TrimSpace(file.Content)) > 0 {
							return ErrMultipleFieldUrlAndContent{
								questionIndex: i,
								optionIndex:   j,
								fileIndex:     k,
							}
						}

						if len(file.Content) == 0 {
							fileUrl := strings.TrimSpace(file.Url)

							if len(fileUrl) == 0 {
								return ErrEmptyFileField{
									questionIndex: i,
									optionIndex:   j,
									fileIndex:     k,
									field:         "URL",
								}
							}

							if _, err := url.ParseRequestURI(fileUrl); err != nil {
								return ErrInvalidURLFile{
									questionIndex: i,
									optionIndex:   j,
									fileIndex:     k,
								}
							}
						} else {
							content := strings.TrimSpace(file.Content)

							if len(content) == 0 {
								return ErrEmptyFileField{
									questionIndex: i,
									optionIndex:   j,
									fileIndex:     k,
									field:         "Content",
								}
							}
						}
					}
				}
			}
		}
	}

	return nil
}

type (
	ErrEmptyQuestionField struct {
		index int
		field string
	}
	ErrCanSelectMultiple struct {
		index int
	}
	ErrCanSkip struct {
		index int
	}
	ErrEmptyOption struct {
		questionIndex int
		optionIndex   int
		field         string
	}
	ErrEmptyFileAndDependency struct {
		questionIndex int
		optionIndex   int
	}
	ErrEmptyUrlAndContent struct {
		questionIndex int
		optionIndex   int
		fileIndex     int
	}
	ErrMultipleFieldUrlAndContent struct {
		questionIndex int
		optionIndex   int
		fileIndex     int
	}
	ErrEmptyFileField struct {
		questionIndex int
		optionIndex   int
		fileIndex     int
		field         string
	}
	ErrInvalidURLFile struct {
		questionIndex int
		optionIndex   int
		fileIndex     int
	}
)

func (e ErrEmptyQuestionField) Error() string {
	return fmt.Sprintf("%d. %s is empty. %s can not be empty", e.index+1, e.field, e.field)
}

func (e ErrCanSelectMultiple) Error() string {
	return fmt.Sprintf("%d. question 'CanSelectMultiple' field is 'true'. This field can not be true if there are less than two options", e.index+1)
}

func (e ErrCanSkip) Error() string {
	return fmt.Sprintf("%d. question 'CanSkip' field is 'false'. This field can not be false if there is only one option", e.index+1)
}

func (e ErrEmptyOption) Error() string {
	return fmt.Sprintf("%d. question, %d. option %s is empty. %s can not be empty", e.questionIndex+1, e.optionIndex+1, e.field, e.field)
}

func (e ErrEmptyFileAndDependency) Error() string {
	return fmt.Sprintf("%d. question %d. option do not have both file and dependency. Option must have at least one file or dependency", e.questionIndex+1, e.optionIndex+1)
}

func (e ErrEmptyUrlAndContent) Error() string {
	return fmt.Sprintf("%d. question %d. option %d. file do not have both URL and Content. File must have at least one URL or Content", e.questionIndex+1, e.optionIndex+1, e.fileIndex+1)
}

func (e ErrMultipleFieldUrlAndContent) Error() string {
	return fmt.Sprintf("%d. question %d. option %d. file have both URL and Content. File can not have both URL and Content at the same time", e.questionIndex+1, e.optionIndex+1, e.fileIndex+1)
}

func (e ErrEmptyFileField) Error() string {
	return fmt.Sprintf("%d. question, %d. option, %d. file %s is empty. %s can not be empty", e.questionIndex+1, e.optionIndex+1, e.fileIndex+1, e.field, e.field)
}

func (e ErrInvalidURLFile) Error() string {
	return fmt.Sprintf("%d. question, %d. option, %d. file URL is invalid.", e.questionIndex+1, e.optionIndex+1, e.fileIndex+1)
}
