package config

import (
	"fmt"
	"reflect"
)

type Normalizable interface {
	Normalize() error
}

func NormalizeOrDefault[T Normalizable](c T) (T, error) {
	// If input is not nil, normalize it.
	if !reflect.ValueOf(c).IsNil() {
		return c, c.Normalize()
	}

	// If input is nil, new a object and normalize it.

	t := reflect.TypeOf(c)
	if t.Kind() != reflect.Ptr {
		return c, fmt.Errorf("type %s is not a pointer", t)
	}

	n, ok := reflect.New(t.Elem()).Interface().(T)
	if !ok {
		return n, fmt.Errorf("type %s does not implement Normalizable", t)
	}
	return n, n.Normalize()
}
