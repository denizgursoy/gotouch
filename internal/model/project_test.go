package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	projectName     = "test-Project"
	validProjectURL = "https://github.com/denizgursoy/gotouch/graphs/traffic"

	validProjectWithNoCustom = ProjectStructureData{
		Name:      projectName,
		Reference: "",
		URL:       validProjectURL,
		Language:  "go",
	}
)

func TestProjectStructureData_IsValid(t *testing.T) {

	choice1 := "1"
	choice2 := "2"

	fileValid := &File{
		Url:          "https://raw.githubusercontent.com/denizgursoy/go-touch-projects/main/k8s-deployment.yaml",
		Content:      "",
		PathFromRoot: "cmd",
	}

	fileInvalidUrl := &File{
		Url:          "invalidUrl",
		Content:      "",
		PathFromRoot: "cmd",
	}

	fileInvalidHaveBothUrlAndContent := &File{
		Url:          "https://raw.githubusercontent.com/denizgursoy/go-touch-projects/main/k8s-deployment.yaml",
		Content:      "content",
		PathFromRoot: "cmd",
	}

	fileInvalidEmptyPathFromRoot := &File{
		Url:          "",
		Content:      "blablabla",
		PathFromRoot: "",
	}

	choiceValidWithDependencies := &Choice{
		Choice: "choice",
		Resources: Resources{
			Dependencies: []any{choice1, choice2},
			Files:        nil,
		},
	}

	choiceValidOnlyValue := &Choice{
		Choice: "choice",
		Resources: Resources{
			Dependencies: nil,
			Files:        nil,
			Values: map[string]any{
				"key": "value",
			},
		},
	}

	choiceValidWithFile := &Choice{
		Choice: "choice2",
		Resources: Resources{
			Dependencies: nil,
			Files: []*File{
				fileValid,
			},
		},
	}

	choiceInvalidEmptyChoice := &Choice{
		Choice: "",
		Resources: Resources{
			Dependencies: nil,
			Files: []*File{
				fileValid,
			},
		},
	}

	choiceInvalidURL := &Choice{
		Choice: "choice",
		Resources: Resources{
			Dependencies: nil,
			Files: []*File{
				fileInvalidUrl,
			},
		},
	}

	choiceEmptyPathFromRoot := &Choice{
		Choice: "choice",
		Resources: Resources{
			Dependencies: nil,
			Files: []*File{
				fileInvalidEmptyPathFromRoot,
			},
		},
	}

	choiceInvalid1 := &Choice{
		Choice: "choice",
		Resources: Resources{
			Dependencies: nil,
			Files: []*File{
				fileValid, fileInvalidHaveBothUrlAndContent,
			},
		},
	}

	questionValid := &Question{
		Direction:         "Example Question",
		CanSkip:           true,
		CanSelectMultiple: false,
		Choices:           []*Choice{choiceValidWithDependencies, choiceValidWithFile},
	}

	questionInvalidEmptyChoice := &Question{
		Direction:         "Example Question",
		CanSkip:           false,
		CanSelectMultiple: false,
		Choices:           []*Choice{choiceValidWithDependencies, choiceInvalidEmptyChoice},
	}

	questionInvalidEmptyDirection := &Question{
		Direction:         "",
		CanSkip:           true,
		CanSelectMultiple: false,
		Choices:           []*Choice{choiceValidWithDependencies, choiceValidWithFile},
	}

	questionInvalidEmptyChoiceStruct := &Question{
		Direction:         "direction",
		CanSkip:           false,
		CanSelectMultiple: false,
		Choices:           nil,
	}

	questionInvalidWhenCanSelectMultipleTrue := &Question{
		Direction:         "Example Question",
		CanSkip:           true,
		CanSelectMultiple: true,
		Choices:           []*Choice{choiceValidWithFile},
	}

	questionWithChoiceHavingOnlyValueField := &Question{
		Direction: "Example Question",
		CanSkip:   true,
		Choices:   []*Choice{choiceValidOnlyValue},
	}

	questionInvalidWhenCanSkipFalse := &Question{
		Direction:         "Example Question",
		CanSkip:           false,
		CanSelectMultiple: false,
		Choices:           []*Choice{choiceValidWithFile},
	}

	questionInvalidURL := &Question{
		Direction:         "Direction",
		CanSkip:           false,
		CanSelectMultiple: false,
		Choices:           []*Choice{choiceInvalidURL, choiceValidWithDependencies},
	}

	questionEmptyPathFromRoot := &Question{
		Direction:         "direction",
		CanSkip:           false,
		CanSelectMultiple: false,
		Choices:           []*Choice{choiceValidWithFile, choiceValidWithDependencies, choiceEmptyPathFromRoot},
	}

	questionInvalid1 := &Question{
		Direction:         "direction",
		CanSkip:           true,
		CanSelectMultiple: false,
		Choices:           []*Choice{choiceValidWithFile, choiceValidWithDependencies, choiceInvalid1},
	}

	invalidProjectWithEmptyChoice := ProjectStructureData{
		Name:      projectName,
		Reference: "",
		URL:       validProjectURL,
		Questions: []*Question{questionInvalidEmptyChoice},
	}

	invalidProjectWithEmptyDirection := ProjectStructureData{
		Name:      "Example Project",
		Reference: "",
		URL:       validProjectURL,
		Questions: []*Question{questionInvalidEmptyDirection},
	}

	invalidProjectWithEmptyChoiceStruct := ProjectStructureData{
		Name:      "Example Project",
		Reference: "",
		URL:       validProjectURL,
		Questions: []*Question{questionValid, questionInvalidEmptyChoiceStruct},
	}

	invalidProject1 := ProjectStructureData{
		Name:      "invalid project",
		Reference: "",
		URL:       validProjectURL,
		Questions: []*Question{questionValid, questionInvalidWhenCanSelectMultipleTrue},
		Resources: Resources{
			Values: nil,
		},
	}

	validProjectStructureWithChoiceHavingOnlyAValue := ProjectStructureData{
		Name:      "project",
		Reference: "",
		URL:       validProjectURL,
		Questions: []*Question{questionWithChoiceHavingOnlyValueField},
		Resources: Resources{
			Values: nil,
		},
	}

	invalidProject2 := ProjectStructureData{
		Name:      "invalid project",
		Reference: "",
		URL:       validProjectURL,
		Questions: []*Question{questionValid, questionInvalidWhenCanSkipFalse},
		Resources: Resources{
			Values: nil,
		},
	}

	invalidProject3 := ProjectStructureData{
		Name:      "invalid project 3",
		Reference: "",
		URL:       validProjectURL,
		Questions: []*Question{questionValid, questionInvalidURL},
	}

	invalidProject4 := ProjectStructureData{
		Name:      "invalid project 4",
		Reference: "",
		URL:       validProjectURL,
		Questions: []*Question{questionValid, questionEmptyPathFromRoot},
	}

	invalidProject5 := ProjectStructureData{
		Name:      "invalid project 4",
		Reference: "",
		URL:       validProjectURL,
		Questions: []*Question{questionValid, questionInvalid1},
	}

	t.Run("should return ErrProjectNameIsEmpty if url is empty", func(t *testing.T) {
		project := &ProjectStructureData{
			Name:      "",
			Reference: "",
			URL:       "",
		}
		err := project.IsValid()

		require.NotNil(t, err)
		require.ErrorIs(t, err, ErrProjectNameIsEmpty)
	})

	t.Run("should not return ErrProjectURLIsEmpty if url is empty", func(t *testing.T) {
		project := &ProjectStructureData{
			Name:      projectName,
			Reference: "",
			URL:       "",
		}
		err := project.IsValid()

		require.Nil(t, err)
	})

	t.Run("should return ErrProjectURLIsNotValid if url is not valid", func(t *testing.T) {
		project := &ProjectStructureData{
			Name:      projectName,
			Reference: "",
			URL:       "test-url",
		}
		err := project.IsValid()

		require.NotNil(t, err)
		require.ErrorIs(t, err, ErrProjectURLIsNotValid)
	})

	t.Run("should provide Delimiter in correct format", func(t *testing.T) {
		delimiters := []struct {
			Name string
			Case string
			Err  bool
		}{
			{
				Name: "empty test",
				Case: " ",
				Err:  false,
			},
			{
				Name: "correct delimiter",
				Case: "[[   ]]",
				Err:  false,
			},
			{
				Name: "wrong concat value",
				Case: "[[]]",
				Err:  true,
			},
		}

		for _, d := range delimiters {
			project := &ProjectStructureData{
				Name:       projectName,
				URL:        validProjectURL,
				Delimiters: d.Case,
			}
			err := project.IsValid()

			t.Run(d.Name, func(t *testing.T) {
				require.True(t, (err == nil) != d.Err)
				if err != nil {
					require.ErrorAs(t, err, &ErrWrongDelimiterFormat{})
				}
			})
		}
	})

	t.Run("should return no error", func(t *testing.T) {
		err := validProjectWithNoCustom.IsValid()
		require.Nil(t, err)
	})

	t.Run("should return no error if choice has only values,no files and dependencies", func(t *testing.T) {
		err := validProjectStructureWithChoiceHavingOnlyAValue.IsValid()
		require.Nil(t, err)
	})

	t.Run("should return error if choice is empty", func(t *testing.T) {
		err := invalidProjectWithEmptyChoice.IsValid()

		expectedError := &ErrEmptyChoice{}

		require.NotNil(t, err)
		require.ErrorAs(t, err, expectedError)
		require.Equal(t, 0, expectedError.questionIndex)
		require.Equal(t, 1, expectedError.choiceIndex)
		require.Equal(t, "Choice", expectedError.field)
	})

	t.Run("should return error if direction is empty", func(t *testing.T) {
		err := invalidProjectWithEmptyDirection.IsValid()

		expectedError := &ErrEmptyQuestionField{}

		require.NotNil(t, err)
		require.ErrorAs(t, err, expectedError)
		require.Equal(t, 0, expectedError.index)
		require.Equal(t, "Direction", expectedError.field)
	})

	t.Run("should return error if choice is empty", func(t *testing.T) {
		err := invalidProjectWithEmptyChoiceStruct.IsValid()

		expectedError := &ErrEmptyQuestionField{}

		require.NotNil(t, err)
		require.ErrorAs(t, err, expectedError)
		require.Equal(t, 1, expectedError.index)
		require.Equal(t, "Choices", expectedError.field)
	})

	t.Run("should return error if CanSelectMultiple is true when choice length equal 1", func(t *testing.T) {
		err := invalidProject1.IsValid()

		expectedError := &ErrCanSelectMultiple{}

		require.NotNil(t, err)
		require.ErrorAs(t, err, expectedError)
		require.Equal(t, 1, expectedError.index)
	})

	t.Run("should return error if CanSkip is false when choice length equal 1", func(t *testing.T) {
		err := invalidProject2.IsValid()

		expectedError := &ErrCanSkip{}

		require.NotNil(t, err)
		require.ErrorAs(t, err, expectedError)
		require.Equal(t, 1, expectedError.index)
	})

	t.Run("should return error if file URL is invalid", func(t *testing.T) {
		err := invalidProject3.IsValid()

		expectedError := &ErrInvalidURLFile{}

		require.NotNil(t, err)
		require.ErrorAs(t, err, expectedError)
		require.Equal(t, 1, expectedError.questionIndex)
		require.Equal(t, 0, expectedError.choiceIndex)
		require.Equal(t, 0, expectedError.fileIndex)
	})

	t.Run("should return error if file PathFromRoot is empty", func(t *testing.T) {
		err := invalidProject4.IsValid()

		expectedError := &ErrEmptyFileField{}

		require.NotNil(t, err)
		require.ErrorAs(t, err, expectedError)
		require.Equal(t, 1, expectedError.questionIndex)
		require.Equal(t, 2, expectedError.choiceIndex)
		require.Equal(t, 0, expectedError.fileIndex)
		require.Equal(t, "PathFromRoot", expectedError.field)
	})

	t.Run("should return error if file content and url len > 0", func(t *testing.T) {
		err := invalidProject5.IsValid()

		expectedError := &ErrMultipleFieldUrlAndContent{}

		require.NotNil(t, err)
		require.ErrorAs(t, err, expectedError)
		require.Equal(t, 1, expectedError.questionIndex)
		require.Equal(t, 2, expectedError.choiceIndex)
		require.Equal(t, 1, expectedError.fileIndex)
	})
}

func TestProjectStructureText(t *testing.T) {
	t.Run("should display reference it exists", func(t *testing.T) {
		data := ProjectStructureData{
			Name:      "x",
			Reference: "y",
		}
		expected := data.Name + " ( " + data.Reference + " )"
		require.Equal(t, expected, data.String())
	})

	t.Run("should not display if reference does not exist", func(t *testing.T) {
		data := ProjectStructureData{
			Name:      "x",
			Reference: "   ",
		}
		expected := data.Name
		require.Equal(t, expected, data.String())
	})

	t.Run("should display if language exists", func(t *testing.T) {
		data := ProjectStructureData{
			Name:      "x",
			Language:  "go",
			Reference: "   ",
		}
		expected := data.Name + " ( Go )"
		require.Equal(t, expected, data.String())
	})

	t.Run("should display all fields", func(t *testing.T) {
		data := ProjectStructureData{
			Name:      "x",
			Language:  "go",
			Reference: "github.com",
		}
		expected := data.Name + " ( Go ) ( github.com )"
		require.Equal(t, expected, data.String())
	})
}

func TestDependencyTest(t *testing.T) {
	t.Run("should validate go dependency", func(t *testing.T) {
		custom := validProjectWithNoCustom
		choices := make([]*Choice, 0)

		goDependencyChoice := &Choice{
			Choice: "asdsd",
			Resources: Resources{
				Dependencies: []any{"asdsad"},
			},
		}

		choices = append(choices, goDependencyChoice)

		custom.Questions = append(custom.Questions, &Question{
			Direction: "asdasdas",
			CanSkip:   true,
			Choices:   choices,
		})

		err := custom.IsValid()
		require.Nil(t, err)
	})

	t.Run("should return false if not go dependency", func(t *testing.T) {
		custom := validProjectWithNoCustom
		choices := make([]*Choice, 0)

		type wrongType struct {
			name    string
			version string
		}

		wrongGoDependencyChoice := &Choice{
			Choice: "asdsd",
			Resources: Resources{
				Dependencies: []any{
					"asdsa",
					wrongType{
						name:    "x",
						version: "y",
					},
				},
			},
		}

		choices = append(choices, wrongGoDependencyChoice)

		custom.Questions = append(custom.Questions, &Question{
			Direction: "asdasdas",
			CanSkip:   true,
			Choices:   choices,
		})

		err := custom.IsValid()
		require.NotNil(t, err)

		expectedError := &ErrWrongDependencyFormat{}

		require.ErrorAs(t, err, expectedError)
		require.Equal(t, 0, expectedError.questionIndex)
		require.Equal(t, 0, expectedError.choiceIndex)
		require.Equal(t, 1, expectedError.dependencyIndex)
	})
}
