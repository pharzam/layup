# tally

## A small Go command that counts the lines, words and characters of text files

**Technology stack:** Go, with the standard library only, built with the Go toolchain.

## 1. Purpose

A person who works in a terminal wants to know how long a text file is. The product `tally` reads the files that the person names and prints the number of lines, words and characters of each, and a total.

## 2. Scope

### 2.1 In scope

- The command `tally [-l] [-w] [-c] [file...]`, with zero or more file names.
- The flags `-l` (lines only), `-w` (words only) and `-c` (characters only); with no flag, the three counts.
- Reading standard input when no file name is given.

### 2.2 Out of scope

- Other counts, such as bytes, paragraphs or sentences.
- Reading archives, directories or files that are not text.
- A graphical interface, a network service or a plug-in system.

## 3. Behaviour

1. A line is a run of characters that ends with a line feed. The last line of a file counts when the file does not end with a line feed.
2. A word is a run of characters that are not white space, where white space is a character of the Unicode property `White_Space`.
3. A character is one Unicode code point of a text in UTF-8.
4. The output has one row for each file: the counts that the flags select, then the file name, each separated from the next by one space. With two files or more, the last row is `total`.
5. A file that cannot be read gives a message on standard error, with the file name and the reason, and the exit code 1. The other files are still counted.
6. A file that is not valid UTF-8 gives a message on standard error, with the file name, and the exit code 1.
7. When no error occurs, the exit code is 0.

## 4. Quality targets

| Metric | Target | Measurement |
| ------ | ------ | ----------- |
| Time for a file of 10 MiB | At most 2 seconds on a computer with 4 cores | A test that counts a generated file of 10 MiB and measures the elapsed time |
| Memory for a file of 1 GiB | At most 64 MiB of heap | A test that counts a generated file of 1 GiB and reads the heap size at the end |
| Test coverage of the counting package | At least 90 percent of the statements | The output of `go test -cover` |

## 5. Terms

| Term | Meaning |
| ---- | ------- |
| **UTF-8** | The encoding of Unicode text in one to four bytes for each code point. |
| **MiB** | A mebibyte, 1,048,576 bytes. |
| **GiB** | A gibibyte, 1,073,741,824 bytes. |
