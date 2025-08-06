package operator

import (
	"testing"

	"github.com/denizgursoy/gotouch/internal/compressor"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_cmdExecutor_CompressDirectory(t *testing.T) {
	t.Run("should call compress directory on the executor", func(t *testing.T) {
		instance := GetInstance()

		controller := gomock.NewController(t)
		mockCompressor := compressor.NewMockCompressor(controller)

		sourceDirectory := "sourcedir"

		opts := &CompressDirectoryOptions{
			SourceDirectory: &sourceDirectory,
			TargetDirectory: nil,
			Compressor:      mockCompressor,
			Logger:          logger.NewLogger(),
		}
		mockCompressor.EXPECT().CompressDirectory(gomock.Eq(sourceDirectory), gomock.Eq(""))

		err := instance.CompressDirectory(opts)
		require.Nil(t, err)
	})

	t.Run("should call compress directory on the executor with target", func(t *testing.T) {
		instance := GetInstance()

		controller := gomock.NewController(t)
		mockCompressor := compressor.NewMockCompressor(controller)

		sourceDirectory := "sourcedir"
		targetDirectory := "targetdir"

		opts := &CompressDirectoryOptions{
			SourceDirectory: &sourceDirectory,
			TargetDirectory: &targetDirectory,
			Compressor:      mockCompressor,
			Logger:          logger.NewLogger(),
		}
		mockCompressor.EXPECT().CompressDirectory(gomock.Eq(sourceDirectory), gomock.Eq(targetDirectory))

		err := instance.CompressDirectory(opts)
		require.Nil(t, err)
	})
}
