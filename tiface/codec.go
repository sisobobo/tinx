package tiface

import "github.com/sisobobo/tinx/tpkg/bufio"

type Codec interface {
	Decode(reader *bufio.Reader) (interface{}, error)
	Encode(interface{}) ([]byte, error)
}
