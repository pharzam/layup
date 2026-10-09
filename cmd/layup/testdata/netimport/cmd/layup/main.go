// Command layup of the fixture module imports net/http: a breach of rule 5 of
// docs/spec/packages.md, and of no other rule in this package.
package main

import "net/http"

func main() { _ = http.StatusOK }
