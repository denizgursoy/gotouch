// +build unit

package req

import (
	"bytes"
	"errors"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

type RoundTripFunction func(req *http.Request) (*http.Response, error)

// RoundTrip .
func (f RoundTripFunction) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func NewTestClient(fn RoundTripFunction) *http.Client {
	return &http.Client{
		Transport: fn,
	}

}

var (
	testUrl  = "https://github.com/denizgursoy/go-touch-projects/blob/main/properties-2.yaml"
	testPath = "/cmd.txt"
)

func Test_fileTask_Complete(t *testing.T) {

	t.Run("should call create file with the content bytes if url is empty ", func(t *testing.T) {
		type args struct {
			file          model.File
			ExpectedBytes []byte
			client        *http.Client
		}

		onlyContentFile := model.File{
			Url:     "",
			Content: "adsasdsas",
			Path:    testPath,
		}

		urlFileContent := "asds"

		onlyUrlFile := model.File{
			Url:     testUrl,
			Content: "",
			Path:    testPath,
		}

		urlBytes := []byte(urlFileContent)
		testCases := []args{
			{
				file:          onlyContentFile,
				ExpectedBytes: []byte(onlyContentFile.Content),
				client:        &http.Client{},
			},
			{
				file:          onlyUrlFile,
				ExpectedBytes: urlBytes,
				client: NewTestClient(func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						Body: io.NopCloser(bytes.NewReader(urlBytes)),
					}, nil
				}),
			},
		}
		controller := gomock.NewController(t)
		defer controller.Finish()

		for _, arg := range testCases {

			mockManager := manager.NewMockManager(controller)

			mockManager.EXPECT().CreateFile(gomock.Any(), gomock.Eq(arg.file.Path)).DoAndReturn(func(arg1, arg2 interface{}) error {
				closer := arg1.(io.ReadCloser)
				all, err := ioutil.ReadAll(closer)
				if err != nil {
					return err
				}
				if !reflect.DeepEqual(all, arg.ExpectedBytes) {
					return errors.New("does not match")
				}
				return nil
			})

			task := &fileTask{
				File:    arg.file,
				Logger:  logger.NewLogger(),
				Manager: mockManager,
				Client:  arg.client,
			}

			err := task.Complete()
			require.NoError(t, err)
		}

	})

	unexpectedError := errors.New("unexpected error")
	t.Run("should return error if any error is taken while getting file from urls", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockManager := manager.NewMockManager(controller)

		task := &fileTask{
			File: model.File{
				Url:     testUrl,
				Content: "",
				Path:    testPath,
			},
			Logger:  logger.NewLogger(),
			Manager: mockManager,
			Client: NewTestClient(func(req *http.Request) (*http.Response, error) {
				return &http.Response{Body: ioutil.NopCloser(bytes.NewReader([]byte("")))}, unexpectedError
			}),
		}

		err := task.Complete()
		require.Error(t, err)
	})

	t.Run("should return error is file is not created", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockManager := manager.NewMockManager(controller)

		mockManager.EXPECT().CreateFile(gomock.Any(), gomock.Any()).Return(unexpectedError)

		task := &fileTask{
			File: model.File{
				Url:     "",
				Content: "sdsd",
				Path:    testPath,
			},
			Logger:  logger.NewLogger(),
			Manager: mockManager,
			Client:  &http.Client{},
		}

		err := task.Complete()
		require.Error(t, err)
		require.ErrorIs(t, unexpectedError, err)
	})

}
