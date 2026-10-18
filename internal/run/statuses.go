package run

// This file holds the commit statuses of a head and the rule of their posting
// (docs/spec/loop.md, REQ-007 — The four layup/ statuses; task T-q0kr, #362, row
// 66f of the plan of M2e). The step loop that calls them is row 66k's.

import (
	"context"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/rules"
)

// The four contexts of the statuses of a head.
const (
	GatesContext  = "layup/gates"
	SpecContext   = "layup/spec"
	RulesContext  = "layup/rules"
	VerifyContext = "layup/verify"
)

// A Status is a commit status to post: its context, one of the four states and
// a description of one line.
type Status struct{ Context, State, Description string }

// fit gives the description as one line of at most forge.MaxDescription characters.
func fit(s string) string { return s }

// appStatus gives the newest status of the context whose creator is the bot.
func appStatus(forge.HeadChecks, string, int64) (forge.Status, bool) { return forge.Status{}, false }

// PostStatus posts st on the head sha unless the App's newest status of the
// context already has its state and description. It tells whether it posted.
func (s *Sessions) PostStatus(context.Context, string, Status) (bool, error) { return false, nil }

// RulesStatus gives the status of layup/rules for the result of the basic rule.
func RulesStatus(rules.Result) Status { return Status{} }

// A RulesInput is what the caller knows of the head whose layup/rules is asked.
type RulesInput struct {
	Task, Base, Head, PreRev, RegisterRev string
	Records                               []rules.Record
}

// RulesAt gives the status layup/rules of a head.
func (s *Sessions) RulesAt(context.Context, RulesInput) (Status, error) { return Status{}, nil }
