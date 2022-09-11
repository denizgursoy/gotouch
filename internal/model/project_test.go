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
		Url:             "https://raw.githubusercontent.com/denizgursoy/go-touch-projects/main/k8s-deployment.yaml",
		Content:         "",
		TargetDirectory: "cmd",
	}

	fileInvalidUrl := &File{
		Url:             "invalidUrl",
		Content:         "",
		TargetDirectory: "cmd",
	}

	fileInvalidHaveBothUrlAndContent := &File{
		Url:             "https://raw.githubusercontent.com/denizgursoy/go-touch-projects/main/k8s-deployment.yaml",
		Content:         "content",
		TargetDirectory: "cmd",
	}

	fileInvalidEmptyTargetDirectory := &File{
		Url:             "",
		Content:         "blablabla",
		TargetDirectory: "",
	}

	optionValidWithDependencies := &Option{
		Choice:       "choice",
		Dependencies: []*string{&choice1, &choice2},
		Files:        nil,
	}

	optionValidWithFile := &Option{
		Choice:       "choice2",
		Dependencies: nil,
		Files:        []*File{fileValid},
	}

	optionInvalidEmptyChoice := &Option{
		Choice:       "",
		Dependencies: nil,
		Files:        []*File{fileValid},
	}

	optionInvalidURL := &Option{
		Choice:       "choice",
		Dependencies: nil,
		Files:        []*File{fileInvalidUrl},
	}

	optionEmptyTargetDirectory := &Option{
		Choice:       "choice",
		Dependencies: nil,
		Files:        []*File{fileInvalidEmptyTargetDirectory},
	}

	optionInvalid1 := &Option{
		Choice:       "choice",
		Dependencies: nil,
		Files:        []*File{fileValid, fileInvalidHaveBothUrlAndContent},
	}

	questionValid := &Question{
		Direction:         "Example Question",
		CanSkip:           true,
		CanSelectMultiple: false,
		Options:           []*Option{optionValidWithDependencies, optionValidWithFile},
	}

	questionInvalidEmptyChoice := &Question{
		Direction:         "Example Question",
		CanSkip:           false,
		CanSelectMultiple: false,
		Options:           []*Option{optionValidWithDependencies, optionInvalidEmptyChoice},
	}

	questionInvalidEmptyDirection := &Question{
		Direction:         "",
		CanSkip:           true,
		CanSelectMultiple: false,
		Options:           []*Option{optionValidWithDependencies, optionValidWithFile},
	}

	questionInvalidEmptyOption := &Question{
		Direction:         "direction",
		CanSkip:           false,
		CanSelectMultiple: false,
		Options:           nil,
	}

	questionInvalidWhenCanSelectMultipleTrue := &Question{
		Direction:         "Example Question",
		CanSkip:           true,
		CanSelectMultiple: true,
		Options:           []*Option{optionValidWithFile},
	}

	questionInvalidWhenCanSkipFalse := &Question{
		Direction:         "Example Question",
		CanSkip:           false,
		CanSelectMultiple: false,
		Options:           []*Option{optionValidWithFile},
	}

	questionInvalidURL := &Question{
		Direction:         "Direction",
		CanSkip:           false,
		CanSelectMultiple: false,
		Options:           []*Option{optionInvalidURL, optionValidWithDependencies},
	}

	questionEmptyTargetDirectory := &Question{
		Direction:         "direction",
		CanSkip:           false,
		CanSelectMultiple: false,
		Options:           []*Option{optionValidWithFile, optionValidWithDependencies, optionEmptyTargetDirectory},
	}

	questionInvalid1 := &Question{
		Direction:         "direction",
		CanSkip:           true,
		CanSelectMultiple: false,
		Options:           []*Option{optionValidWithFile, optionValidWithDependencies, optionInvalid1},
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

	invalidProjectWithEmptyOption := ProjectStructureData{
		Name:      "Example Project",
		Reference: "",
		URL:       validProjectURL,
		Questions: []*Question{questionValid, questionInvalidEmptyOption},
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
		Questions: []*Question{questionValid, questionEmptyTargetDirectory},
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

		expectedError := &ErrEmptyOption{}

		require.NotNil(t, err)
		require.ErrorAs(t, err, expectedError)
		require.Equal(t, 0, expectedError.questionIndex)
		require.Equal(t, 1, expectedError.optionIndex)
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

	t.Run("should return error if option is empty", func(t *testing.T) {
		err := invalidProjectWithEmptyOption.IsValid()

		expectedError := &ErrEmptyQuestionField{}

		require.NotNil(t, err)
		require.ErrorAs(t, err, expectedError)
		require.Equal(t, 1, expectedError.index)
		require.Equal(t, "Options", expectedError.field)
	})

	t.Run("should return error if CanSelectMultiple is true when option length equal 1", func(t *testing.T) {
		err := invalidProject1.IsValid()

		expectedError := &ErrCanSelectMultiple{}

		require.NotNil(t, err)
		require.ErrorAs(t, err, expectedError)
		require.Equal(t, 1, expectedError.index)
	})

	t.Run("should return error if CanSkip is false when option length equal 1", func(t *testing.T) {
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
		require.Equal(t, 0, expectedError.optionIndex)
		require.Equal(t, 0, expectedError.fileIndex)
	})

	t.Run("should return error if file Target Directory is empty", func(t *testing.T) {
		err := invalidProject4.IsValid()

		expectedError := &ErrEmptyFileField{}

		require.NotNil(t, err)
		require.ErrorAs(t, err, expectedError)
		require.Equal(t, 1, expectedError.questionIndex)
		require.Equal(t, 2, expectedError.optionIndex)
		require.Equal(t, 0, expectedError.fileIndex)
		require.Equal(t, "Target Directory", expectedError.field)
	})

	t.Run("should return error if file content and url len > 0", func(t *testing.T) {
		err := invalidProject5.IsValid()

		expectedError := &ErrMultipleFieldUrlAndContent{}

		require.NotNil(t, err)
		require.ErrorAs(t, err, expectedError)
		require.Equal(t, 1, expectedError.questionIndex)
		require.Equal(t, 2, expectedError.optionIndex)
		require.Equal(t, 1, expectedError.fileIndex)
	})
}
