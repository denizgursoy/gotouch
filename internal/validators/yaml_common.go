package validators

import (
	"gopkg.in/yaml.v3"
	"io"
)

func isYaml(r io.Reader) bool {
	node := new(yaml.Node)
	decoder := yaml.NewDecoder(r)
	return decoder.Decode(node) == nil
}
