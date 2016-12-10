package keyvalues

import (
	"fmt"
	"io"
	"strings"
)

var escaper = strings.NewReplacer("\\", "\\\\", "\n", "\\n", "\r", "\\r", "\t", "\\t", "\"", "\\\"")

func (kv *KeyValues) WriteTo(w io.Writer) (n int64, err error) {
	kv.Each(func(ch *KeyValues) {
		if err != nil {
			return
		}
		var c int64
		c, err = ch.writeIndented(w, 0)
		n += c
	})
	return
}

func (kv *KeyValues) writeIndented(w io.Writer, indent int) (n int64, err error) {
	indentString := strings.Repeat("\t", indent)
	c, err := fmt.Fprintf(w, "%s\"%s\" ", indentString, escaper.Replace(kv.name))
	n += int64(c)
	if err != nil {
		return
	}

	if kv.child == nil {
		c, err = fmt.Fprintf(w, "\"%s\"\n", escaper.Replace(kv.value))
		n += int64(c)
		if err != nil {
			return
		}
	} else {
		c, err = fmt.Fprintf(w, "{\n")
		n += int64(c)
		if err != nil {
			return
		}

		kv.Each(func(ch *KeyValues) {
			if err != nil {
				return
			}
			var c int64
			c, err = ch.writeIndented(w, indent+1)
			n += c
		})

		c, err = fmt.Fprintf(w, "%s}\n", indentString)
		n += int64(c)
		if err != nil {
			return
		}
	}
	return
}

var _ io.WriterTo = new(KeyValues)
