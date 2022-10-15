package operator

import (
	"github.com/denizgursoy/gotouch/internal/compressor"
	"github.com/denizgursoy/gotouch/internal/logger"
)

type (
	CompressDirectoryOptions struct {
		SourceDirectory *string               `validate:"required,dir"`
		TargetDirectory *string               `validate:"omitempty,dir"`
		Compressor      compressor.Compressor `validate:"required"`
		Logger          logger.Logger
	}
)

func (c *operator) CompressDirectory(opts *CompressDirectoryOptions) error {
	targetDirectory := ""
	if opts.TargetDirectory != nil {
		targetDirectory = *opts.TargetDirectory
	}
	return opts.Compressor.CompressDirectory(*opts.SourceDirectory, targetDirectory)
}
