package file

import (
	"errors"
	"fmt"
	"io"
)

type char struct {
	val     rune
	escaped bool
}

func isGlob(str string) (bool, error) {
	pos := 0

	nextChar := func(at int) (char, int, error) {
		if str[at] != '\\' {
			val := str[at]
			return char{val: rune(val), escaped: false}, 1, nil
		} else {
			at++
			if at >= len(str) {
				return char{}, 0, fmt.Errorf("error invalid escape")
			}
			val := str[at]
			return char{val: rune(val), escaped: true}, 2, nil
		}
	}

	peek := func() (char, error) {
		char, _, err := nextChar(pos)
		return char, err
	}

	peekUntil := func(terminator rune) ([]char, error) {
		var chars []char
		for pos := pos; pos < len(str); {
			char, consumed, err := nextChar(pos)
			if err != nil {
				return nil, err
			}
			pos += consumed
			chars = append(chars, char)
			if !char.escaped && char.val == terminator {
				return chars, nil
			}
		}
		return nil, io.EOF
	}

	for pos < len(str) {
		next, consumed, err := nextChar(pos)
		if err != nil {
			return false, err
		}
		pos += consumed
		switch next.val {
		case '*':
			if !next.escaped {
				return true, nil
			}
		case '?':
			if !next.escaped {
				return true, nil
			}
		case '[':
			if !next.escaped {
				after, err := peek()
				if err != nil {
					if !errors.Is(err, io.EOF) {
						return false, err
					}
				} else {
					if after.escaped || after.val != ' ' {
						return true, nil
					}
				}
			}
		case '{':
			if !next.escaped {
				chars, err := peekUntil('}')
				if err != nil {
					if !errors.Is(err, io.EOF) {
						return false, err
					}
				} else {
					for _, c := range chars {
						if !c.escaped && c.val == ',' {
							return true, nil
						}
					}
				}
			}
		}
	}

	return false, nil
}
