package session

import "errors"

func MakeInput(d Dir, files map[string][]byte) error { return errors.New("not built") }

func ReadOutput(d Dir, name string) ([]byte, error) { return nil, errors.New("not built") }
