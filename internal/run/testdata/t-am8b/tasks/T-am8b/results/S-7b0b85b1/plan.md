size: large
open questions: —
new interfaces: seven ADRs 0010-0016 (indexed in docs/adr/README.md) fix the go-build-only build, a cmd/tally main package plus a count subpackage in one module (cmd-tally-count-subpackage chosen over root-main because go build ./cmd/tally names the binary tally whatever the module path), three exported count functions Lines, Words and Chars over an io.Reader fed once through a teed reader, a hand-written -l/-w/-c parser with combined flags and a literal - name, streaming through a fixed buffer, UTF-8 validated in the stream with the row printed only after EOF, and White_Space tested through unicode.IsSpace (chosen over the unicode.White_Space table: the same code point set with no table lookup to maintain); no Go code, MoSCoW cell or facts record is written
new dependencies: —
tests:
