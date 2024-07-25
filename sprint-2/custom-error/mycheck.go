//go:build !solution

package mycheck

import (
	"errors"
	"strings"
	"unicode"
)

var (
	errFoundNumbers  error = errors.New("found numbers")
	errLineIsTooLong error = errors.New("line is too long")
	errNoTwoSpaces   error = errors.New("no two spaces")
)

const (
	minLengthLongLine    = 20
	requiredNumberSpaces = 2
)

type stringsErrors []error

func (s stringsErrors) Error() string {
	var mess []string
	for _, err := range s {
		mess = append(mess, err.Error())
	}
	return strings.Join(mess, ";")
}

func MyCheck(input string) error {
	var (
		isFoundNumbers bool
		length         int
		countSpaces    int

		err stringsErrors
	)

	for _, symb := range input {
		length++
		if unicode.IsDigit(symb) {
			isFoundNumbers = true
		}
		if unicode.IsSpace(symb) {
			countSpaces++
		}
	}

	if isFoundNumbers {
		err = append(err, errFoundNumbers)
	}
	if length > minLengthLongLine {
		err = append(err, errLineIsTooLong)
	}
	if countSpaces != requiredNumberSpaces {
		err = append(err, errNoTwoSpaces)
	}

	return err
}
