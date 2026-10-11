package run

import (
	"context"
	"errors"
)

type PhaseSpec struct {
	Step, Role    string
	Harness       []string
	Model         []string
	Models        [][]string
	Base, Records string
	Input         map[string][]byte
	Output        string
	Prompt        func(input map[string]string, output string) []byte
	Keep          func(id string, data []byte) (map[string][]byte, error)
}

type PhaseError struct{ Session, Reason string }

func (e *PhaseError) Error() string { return "the session " + e.Session + ": " + e.Reason }

func (s *Sessions) PhasePair(ctx context.Context, role string, harnesses, models [][]string) ([]string, []string, error) {
	return nil, nil, errors.New("not built")
}

func (s *Sessions) PhaseSession(ctx context.Context, p PhaseSpec) (string, error) {
	return "", errors.New("not built")
}
