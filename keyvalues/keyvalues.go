package keyvalues

import (
	"fmt"
	"strconv"
	"strings"
)

type KeyValues struct {
	name  string
	value string

	next, prev, parent, child *KeyValues
}

func (kv KeyValues) Name() string {
	return kv.name
}

// Returns the string value of this node. If the node is nonexistent or complex
// (one that has subnodes) the default (def) will be returned.
func (kv *KeyValues) String(def string) string {
	if kv == nil || kv.child != nil {
		return def
	}
	return kv.value
}

// Returns the integer value of this node. If the node is nonexistent, complex,
// or unable to be parsed as an integer, the default (def) will be returned.
func (kv *KeyValues) Int(def int64) int64 {
	if kv == nil || kv.child != nil {
		return def
	}
	if i, err := strconv.ParseInt(kv.value, 0, 64); err == nil {
		return i
	}
	return def
}

// Returns the uint64 value of this node. If the node is nonexistent, complex,
// or unable to be parsed as a uint64, the default (def) will be returned.
func (kv *KeyValues) Uint64(def uint64) uint64 {
	if kv == nil || kv.child != nil {
		return def
	}
	if i, err := strconv.ParseUint(kv.value, 0, 64); err == nil {
		return i
	}
	return def
}

// Returns the floating-point value of this node. If the node is nonexistent,
// complex, or unable to be parsed as a float, the default (def) will be returned.
func (kv *KeyValues) Float(def float64) float64 {
	if kv == nil || kv.child != nil {
		return def
	}
	if f, err := strconv.ParseFloat(kv.value, 64); err == nil {
		return f
	}
	return def
}

// Returns the boolean value of this node. The boolean value is equivelent to
// false iff the node has the integer value 0.
func (kv *KeyValues) Bool(def bool) bool {
	if def {
		return kv.Int(1) != 0
	}
	return kv.Int(0) != 0
}

// Sets the value of this node. If this node is complex or nonexistent, this
// method will panic.
func (kv *KeyValues) SetValueString(v string) {
	if kv == nil {
		panic("SetValueString on a nil *KeyValues")
	}
	if kv.child != nil {
		panic("SetValueString on a complex *KeyValues")
	}

	kv.value = v
}

// Sets the value of this node. If this node is complex or nonexistent, this
// method will panic. The integer will be formatted in base 10.
func (kv *KeyValues) SetValueInt(v int64) {
	if kv == nil {
		panic("SetValueInt on a nil *KeyValues")
	}
	if kv.child != nil {
		panic("SetValueInt on a complex *KeyValues")
	}

	kv.value = fmt.Sprint(v)
}

// Sets the value of this node. If this node is complex or nonexistent, this
// method will panic. The uint64 will be formatted as a hexadecimal number
// prefixed by "0x".
func (kv *KeyValues) SetValueUint64(v uint64) {
	if kv == nil {
		panic("SetValueUint64 on a nil *KeyValues")
	}
	if kv.child != nil {
		panic("SetValueUint64 on a complex *KeyValues")
	}

	kv.value = fmt.Sprintf("0x%x", v)
}

// Sets the value of this node. If this node is complex or nonexistent, this
// method will panic.
func (kv *KeyValues) SetValueFloat(v float64) {
	if kv == nil {
		panic("SetValueFloat on a nil *KeyValues")
	}
	if kv.child != nil {
		panic("SetValueFloat on a complex *KeyValues")
	}

	kv.value = fmt.Sprint(v)
}

// Sets the value of this node. If this node is complex or nonexistent, this
// method will panic. The value of the node will be "1" if v is true and "0"
// otherwise.
func (kv *KeyValues) SetValueBool(v bool) {
	if kv == nil {
		panic("SetValueBool on a nil *KeyValues")
	}
	if kv.child != nil {
		panic("SetValueBool on a complex *KeyValues")
	}

	if v {
		kv.value = "1"
	} else {
		kv.value = "0"
	}
}

// Returns the first subkey (if any) that has a name equal to the argument under
// Unicode case-folding.
func (kv *KeyValues) SubKey(name string) *KeyValues {
	if kv == nil {
		return nil
	}
	for ch := kv.child; ch != nil; ch = ch.next {
		if strings.EqualFold(ch.name, name) {
			return ch
		}
	}
	return nil
}

// Creates, appends, and returns a new subkey. If the current node (the subkey's
// parent) is nil, this method will panic.
func (kv *KeyValues) NewSubKey(name string) *KeyValues {
	if kv == nil {
		panic("Call to NewSubKey on a nil *KeyValues")
	}

	ch := &KeyValues{name: name}
	kv.Append(ch)
	return ch
}

// Appends the given child node to this node. If the current node (the subkey's
// parent) is nil, this method will panic. This method is a no-op on a nil child
// if the parent is valid.
func (kv *KeyValues) Append(child *KeyValues) {
	if kv == nil {
		panic("Call to Append on a nil *KeyValues")
	}

	if child == nil {
		return
	}

	if child.parent != nil || child.next != nil || child.prev != nil {
		panic("Call to Append with a *KeyValues already in a heirarchy")
	}

	var prev *KeyValues
	p := &kv.child
	for *p != nil {
		prev = *p
		p = &(*p).next
	}

	*p = child
	child.prev = prev
	child.parent = kv
}

func (kv *KeyValues) Each(f func(*KeyValues)) {
	if kv == nil {
		panic("Call to Each on a nil *KeyValues")
	}

	ch := kv.child
	for ch != nil {
		next := ch.next
		f(ch)
		ch = next
	}
}

func (kv *KeyValues) Remove() {
	if kv == nil {
		return
	}

	if kv.next != nil {
		kv.next.prev = kv.prev
	}
	if kv.prev != nil {
		kv.prev.next = kv.next
	}
	if kv.parent != nil && kv.parent.child == kv {
		kv.parent.child = kv.next
	}
	kv.parent, kv.next, kv.prev = nil, nil, nil
}
