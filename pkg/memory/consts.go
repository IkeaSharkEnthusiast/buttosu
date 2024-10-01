package memory

import "strings"

type ReadError []error

func (r ReadError) Error() string {
	var strs []string

	for _, err := range r {
		strs = append(strs, err.Error())
	}

	return strings.Join(strs, ", ")
}
