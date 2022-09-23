// +build unit

package model

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProjectStructureData_IsValid(t *testing.T) {
	projectName := "test-Project"
	validProjectURL := "https://github.com/denizgursoy/gotouch/graphs/traffic"

	validProjectWithNoCustom := ProjectStructureData{
		Name:      projectName,
		Reference: "",
		URL:       validProjectURL,
	}

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
		Choice:       "choice",
		Dependencies: []*string{&choice1, &choice2},
		Files:        nil,
	}

	choiceValidWithFile := &Choice{
		Choice:       "choice2",
		Dependencies: nil,
		Files:        []*File{fileValid},
	}

	choiceInvalidEmptyChoice := &Choice{
		Choice:       "",
		Dependencies: nil,
		Files:        []*File{fileValid},
	}

	choiceInvalidURL := &Choice{
		Choice:       "choice",
		Dependencies: nil,
		Files:        []*File{fileInvalidUrl},
	}

	choiceEmptyPathFromRoot := &Choice{
		Choice:       "choice",
		Dependencies: nil,
		Files:        []*File{fileInvalidEmptyPathFromRoot},
	}

	choiceInvalid1 := &Choice{
		Choice:       "choice",
		Dependencies: nil,
		Files:        []*File{fileValid, fileInvalidHaveBothUrlAndContent},
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
		Values:    nil,
	}

	invalidProject2 := ProjectStructureData{
		Name:      "invalid project",
		Reference: "",
		URL:       validProjectURL,
		Questions: []*Question{questionValid, questionInvalidWhenCanSkipFalse},
		Values:    nil,
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

	t.Run("should return ErrProjectURLIsEmpty if url is empty", func(t *testing.T) {
		project := &ProjectStructureData{
			Name:      projectName,
			Reference: "",
			URL:       "",
		}
		err := project.IsValid()

		require.NotNil(t, err)
		require.ErrorIs(t, err, ErrProjectURLIsEmpty)
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

	t.Run("should return no error", func(t *testing.T) {
		err := validProjectWithNoCustom.IsValid()
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
