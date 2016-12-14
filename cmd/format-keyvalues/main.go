package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/BenLubar/go-valve/keyvalues"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
)

const bom = "\ufeff"

func main() {
	var kv keyvalues.KeyValues

	var enc encoding.Encoding

	var initial [len(bom)]byte
	n, err := io.ReadFull(os.Stdin, initial[:])
	if err == io.EOF || err == io.ErrUnexpectedEOF {
		err = nil
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if n >= 2 && initial[0] == 0xFE && initial[1] == 0xFF {
		enc = unicode.UTF16(unicode.BigEndian, unicode.UseBOM)
	} else if n >= 2 && initial[0] == 0xFF && initial[1] == 0xFE {
		enc = unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)
	} else if n >= 3 && initial[0] == bom[0] && initial[1] == bom[1] && initial[2] == bom[2] {
		n = 0
		_, err = os.Stdout.WriteString(bom)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		enc = unicode.UTF8
	} else {
		enc = unicode.UTF8
	}

	_, err = kv.ReadFrom(enc.NewDecoder().Reader(io.MultiReader(bytes.NewReader(initial[:n]), os.Stdin)))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	_, err = kv.WriteTo(enc.NewEncoder().Writer(os.Stdout))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
