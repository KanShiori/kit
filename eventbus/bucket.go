package eventbus

import "fmt"

type bucket map[string]EventHandler

func newBucket() bucket {
	return make(map[string]EventHandler)
}

func (b bucket) add(h EventHandler) error {
	if _, ok := b[h.Name()]; ok {
		return fmt.Errorf("conflict name")
	}

	b[h.Name()] = h

	return nil
}

func (b bucket) remove(h EventHandler) {
	delete(b, h.Name())
}

func (b bucket) deepcopy() bucket {
	newbucket := make(map[string]EventHandler, len(b))
	for name, handler := range b {
		newbucket[name] = handler
	}
	return newbucket
}
