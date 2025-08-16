package model

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/denizgursoy/gotouch/internal/auth"
	"github.com/denizgursoy/gotouch/internal/langs"
)

const (
	GitUploadPackContentType   = "application/x-git-upload-pack-advertisement"
	GitDiscoveryReferencesPath = "/info/refs?service=git-upload-pack"
)

type (
	ProjectStructureData struct {
		Resources `yaml:",inline"`
		Name      string `yaml:"name"`
		// InitialModuleName is the default value of the module name. It can be populated to suggest path to creator.
		// It can be something like github.com/my-company/project-x and user only changes the project-x part and does not,
		// have to remember full path.
		InitialModuleName string `yaml:"initialModuleName"`
		// Reference is displayed on the screen for reference during the project selection.
		Reference string `yaml:"reference"`
		// URL is either the address of the git repository or the zip file in the HTTP.
		// if the URL starts with @git, then it checkouts with SSH connection.
		URL string `yaml:"url"`
		// LocalPath stores the relative or absolute path of either compressed file or the directory,
		// that stores the template project
		LocalPath string `yaml:"localPath"`
		// Branch is the git branch to be checkout.
		Branch    string      `yaml:"branch"`
		Questions []*Question `yaml:"questions"`
		Language  string      `yaml:"language"`
		// Delimiters are the values that will be used instead of default go template delimiters. It should be written,
		// like [[ ]], left delimiter and right delimiter separated by a space
		Delimiters string `yaml:"delimiters"`
	}

	Question struct {
		Direction         string    `yaml:"direction"`
		CanSkip           bool      `yaml:"canSkip"`
		CanSelectMultiple bool      `yaml:"canSelectMultiple"`
		Choices           []*Choice `yaml:"choices"`
	}

	Choice struct {
		Choice    string `yaml:"choice"`
		Resources `yaml:",inline"`
	}

	Resources struct {
		Dependencies []any          `yaml:"dependencies"`
		Files        []*File        `yaml:"files"`
		Values       map[string]any `yaml:"values"`
		CustomValues map[string]any `yaml:"customValues"`
	}

	File struct {
		Url          string `yaml:"url"`
		Content      string `yaml:"content"`
		PathFromRoot string `yaml:"pathFromRoot"`
	}

	HttpRequester interface {
		Do(req *http.Request) (*http.Response, error)
	}
)

func (p *ProjectStructureData) String() string {
	projectName := p.Name
	if len(strings.TrimSpace(p.Language)) != 0 {
		englishCaser := cases.Title(language.English)
		projectNameTitle := englishCaser.String(p.Language)
		projectName += fmt.Sprintf(" ( %s )", projectNameTitle)
	}
	if len(strings.TrimSpace(p.Reference)) != 0 {
		projectName += fmt.Sprintf(" ( %s )", p.Reference)
	}
	return projectName
}

var (
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
	if !strings.HasPrefix(projectUrl, "git@") {
		if len(projectUrl) != 0 {
			if _, err := url.ParseRequestURI(projectUrl); err != nil {
				return ErrProjectURLIsNotValid
			}
		}
	}
	localFilePath := strings.TrimSpace(p.LocalPath)
	if len(localFilePath) > 0 {
		// TODO
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

func (p *ProjectStructureData) IsGit(ctx context.Context, requester HttpRequester) (bool, error) {
	if requester == nil {
		requester = auth.NewAuthenticatedHTTPClient()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s", p.URL, GitDiscoveryReferencesPath), nil)
	if err != nil {
		return false, err
	}

	resp, err := requester.Do(req)
	if err != nil {
		return false, nil
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return false, nil
	}

	if resp.Header.Get("Content-Type") != GitUploadPackContentType {
		return false, nil
	}

	return true, nil
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
	} else if len(q.Choices) == 1 && q.CanSelectMultiple {
		return ErrCanSelectMultiple{
			projectName: p.Name,
			index:       questionIndex,
		}
	} else if len(q.Choices) == 1 && !q.CanSkip {
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

func (p *ProjectStructureData) validateChoice(choice *Choice, questionIndex, choiceIndex int) error {
	if len(strings.TrimSpace(choice.Choice)) == 0 {
		return ErrEmptyChoice{
			projectName:   p.Name,
			questionIndex: questionIndex,
			choiceIndex:   choiceIndex,
			field:         "Choice",
		}
	}

	if len(choice.Files) == 0 && len(choice.Dependencies) == 0 && len(choice.Values) == 0 && len(choice.CustomValues) == 0 {
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
	return fmt.Sprintf("%s's %d. question %d. choice does not have any files,values,custom values or dependencies. Choice must have at least one file, value,custom value or dependency", e.projectName, e.questionIndex+1, e.choiceIndex+1)
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
