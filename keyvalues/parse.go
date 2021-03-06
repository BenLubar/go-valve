package keyvalues

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"unicode"
)

func consumeSpaces(in *bufio.Reader) (n int64, err error) {
	r := ' '
	var c int
	for unicode.IsSpace(r) {
		n += int64(c)
		r, c, err = in.ReadRune()
		if err != nil {
			return
		}
	}
	err = in.UnreadRune()
	return
}

func readString(in *bufio.Reader) (s string, n int64, err error, special bool) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
				return
			}
			panic(r)
		}
	}()
	var (
		buf    []rune
		quoted bool

		rr = func() rune {
			r, c, err := in.ReadRune()
			n += int64(c)
			if err != nil {
				panic(err)
			}
			return r
		}
	)
	defer func() {
		if !quoted && strings.HasPrefix(s, "[$") && strings.HasSuffix(s, "]") {
			special = true
		}
	}()

	r := rr()

	if r == '"' {
		quoted = true
	} else if r == '/' {
		r = rr()
		var (
			slice []byte
			c     int64
		)
		switch r {
		case '/':
			slice, err = in.ReadSlice('\n')
			n += int64(len(slice))
			if err != nil {
				return
			}
			c, err = consumeSpaces(in)
			n += c
			if err != nil {
				return
			}
			s, c, err, special = readString(in)
			n += c
			return
		case '*':
			for r != '/' && err == nil {
				slice, err = in.ReadSlice('*')
				n += int64(len(slice))
				r = rr()
			}
			if err != nil {
				return
			}
			c, err = consumeSpaces(in)
			n += c
			if err != nil {
				return
			}
			s, c, err, special = readString(in)
			n += c
			return
		default:
			buf = append(buf, '/', r)
		}
	} else if r == '{' || r == '}' {
		s = string([]rune{r})
		special = true
		return
	} else {
		buf = append(buf, r)
	}

	for {
		r = rr()
		if r == '"' && quoted {
			s = string(buf)
			return
		}
		if unicode.IsSpace(r) && !quoted {
			s = string(buf)
			return
		}
		if (r == '"' || r == '{' || r == '}') && !quoted {
			err = in.UnreadRune()
			s = string(buf)
			return
		}
		if r == '\\' && quoted {
			next := rr()
			switch next {
			case '\\':
				// double backslash -> single backslash
			case 'n':
				r = '\n'
			case 'r':
				r = '\r'
			case 't':
				r = '\t'
			case '"':
				r = '"'
			default:
				buf = append(buf, r)
				r = next
			}
		}
		buf = append(buf, r)
	}
}

type stack []*KeyValues

func (s *stack) push(kv *KeyValues) {
	*s = append(*s, kv)
}

func (s stack) peek() *KeyValues {
	if len(s) == 0 {
		panic("stack underflow")
	}
	return s[len(s)-1]
}

func (s *stack) pop() *KeyValues {
	if len(*s) == 0 {
		panic("stack underflow")
	}
	kv := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return kv
}

func (kv *KeyValues) ReadFrom(r io.Reader) (n int64, err error) {
	in := bufio.NewReader(r)

	var last *KeyValues
	var s = stack{kv}

	for err == nil {
		var c int64

		c, err = consumeSpaces(in)
		n += c
		if err != nil {
			if err == io.EOF && s.pop() == kv {
				err = nil
			}
			return
		}

		var key string
		var special bool
		key, c, err, special = readString(in)
		n += c
		if err != nil {
			if err == io.EOF && s.pop() == kv {
				err = nil
			}
			return
		}
		if special {
			switch key {
			case "{":
				err = errors.New("Unexpected '{': expecting '}' or a key")
				return
			case "}":
				last = nil
				s.pop()
				if len(s) == 0 {
					err = errors.New("Unexpected '}': expecting a key")
					return
				}
				continue
			default:
				if last == nil {
					err = errors.New("Unexpected conditional: expecting a key")
					return
				}

				// TODO: better conditionals
				if key == "[$WIN32]" {
					continue
				}

				last.Remove()
				last = nil
				continue
			}
		}

		s.push(s.peek().NewSubKey(key))

		c, err = consumeSpaces(in)
		n += c
		if err != nil {
			return
		}

		var value string
		value, c, err, special = readString(in)
		n += c
		if err != nil {
			return
		}
		if special {
			switch value {
			case "{":
				continue
			case "}":
				err = errors.New("Unexpected '}': expecting '{' or a value")
				return
			default:
				err = errors.New("Unexpected conditional: expecting '{' or a value")
				return
			}
		}

		last = s.pop()
		last.SetValueString(value)
	}

	return
}

var _ io.ReaderFrom = new(KeyValues)
