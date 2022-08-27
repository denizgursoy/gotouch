package manager

import "fmt"

type (
	Dependency struct {
		Url     *string
		Version *string
	}
)

func (d *Dependency) String() string {
	return fmt.Sprintf("%s@%s", *d.Url, *d.Version)
}
