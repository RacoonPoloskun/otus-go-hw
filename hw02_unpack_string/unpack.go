package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var i int
	var builder strings.Builder

	explodedString := []rune(str)
	explodedStrLength := len(explodedString)

	for i < explodedStrLength {
		char := explodedString[i]

		if unicode.IsDigit(char) {
			return "", ErrInvalidString
		}

		if (i+1 < explodedStrLength) && (unicode.IsDigit(explodedString[i+1])) {
			i++
			stringNum, err := strconv.Atoi(string(explodedString[i]))
			if err != nil {
				return "", ErrInvalidString
			}

			builder.WriteString(strings.Repeat(string(char), stringNum))
		} else {
			builder.WriteRune(char)
		}

		i++
	}

	return builder.String(), nil
}
