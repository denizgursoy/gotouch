package root

import (
	"github.com/denizgursoy/gotouch/internal/compressor"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

const (
	DirectoryFlagName = "directory"
)

var (
	packageCommand = &cobra.Command{
		Use:   "package",
		Short: "createYourZip",
		Long:  `Tag`,
		Run: func(cmd *cobra.Command, args []string) {
			flags := cmd.Flags()
			filePath, inputError := flags.GetString(DirectoryFlagName)
			lgr := logger.NewLogger()
			lgr.LogErrorIfExists(inputError)

			if inputError != nil {
				os.Exit(1)
			}

			point := &filePath
			if len(strings.TrimSpace(filePath)) == 0 {
				point = nil
			}

			options := PackageCommandOptions{
				Compressor: compressor.GetInstance(),
				Logger:     lgr,
				Path:       point,
			}
			err := CompressDirectory(&options)
			lgr.LogErrorIfExists(err)
		},
	}
)

type (
	PackageCommandOptions struct {
		Path       *string               `validate:"required,endswith=.yaml,url|file"`
		Compressor compressor.Compressor `validate:"required"`
		Logger     logger.Logger
	}
)

func init() {
	packageCommand.Flags().StringP(DirectoryFlagName, "d", ".", "directory path")
}

func CompressDirectory(opts *PackageCommandOptions) error {
	return opts.Compressor.CompressDirectory(*opts.Path)
}
