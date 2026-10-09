//go:build layupfixture

// This file is behind a build constraint, so go list gives none of its
// imports: net/smtp, which depends on net and crypto/tls, breaks rule 5 and
// the rule of the engine checks.
package psb

import "net/smtp"

var _ = smtp.PlainAuth
