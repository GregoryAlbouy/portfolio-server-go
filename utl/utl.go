package utl

import (
	"errors"
	"strconv"
	"strings"
)

// SplitInt64 splits a string and converts the resulting values to []int64.
func SplitInt64(str string, sep string) ([]int64, error) {
	strs := strings.Split(str, sep)
	nbs := []int64{}
	for _, str := range strs {
		nb, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, errors.New("Cannot convert to []int64")
		}
		nbs = append(nbs, nb)
	}

	return nbs, nil
}
