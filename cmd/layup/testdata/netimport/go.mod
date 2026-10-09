// The fixture of the test of the package rules (rules_integration_test.go):
// the module path of LAYUP and two breaches of rule 5: an import of net/http,
// and an import of net/smtp in a file behind a build constraint, which also
// breaks the rule of the engine checks in internal/psb.
module github.com/pharzam/layup

go 1.26
