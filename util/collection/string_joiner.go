package collection

import (
	"fmt"
	"strings"
)

type Joiner struct {
	delimiter string
	prefix    string
	suffix    string
	data      []string
}

func NewJoiner(delimiter, prefix, suffix string) *Joiner {
	return &Joiner{
		delimiter: delimiter,
		prefix:    prefix,
		suffix:    suffix,
		data:      []string{},
	}
}

func (joiner *Joiner) Append(data string) *Joiner {
	if data != "" {
		joiner.data = append(joiner.data, data)
	}
	return joiner
}

func (joiner *Joiner) Size() int {
	return len(joiner.data)
}

func (joiner *Joiner) String() string {
	return fmt.Sprintf("%s%s%s", joiner.prefix, strings.Join(joiner.data, joiner.delimiter), joiner.suffix)
}
