package mail

import (
	"bufio"
	"bytes"
	"errors"
)

var (
	ErrorNoSubject = errors.New("no subject")
	ErrorNoBody    = errors.New("no body")
)

// extractSubjectAndBody extracts subject and email body from template
// template structure is
// ---
// subject text
// ---
// html body
func ExtractSubjectAndBody(d []byte) (string, string, error) {
	scanner := bufio.NewScanner(bytes.NewReader(d))
	scanner.Split(Split)
	haveSubject := scanner.Scan()
	if !haveSubject {
		return "", "", ErrorNoSubject
	}
	subject := scanner.Text()
	haveBody := scanner.Scan()
	if !haveBody {
		return "", "", ErrorNoBody
	}
	body := scanner.Text()
	return subject, body, nil
}

// Split split "---" separator
// separator should be in the beginning of the line
// separator in the end of the file is not required
// separator in the beginning of the file is not required
// any symbols in the same line after separator are ignored
func Split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	dataLen := len(data)

	// Return Nothing if at the end of file or  no data passed.
	if atEOF && dataLen == 0 {
		return 0, nil, nil
	}

	i := 0
	sl := len("---")
	// the buffer not starts from separator
	if bytes.Index(data, []byte("---")) != 0 {
		i = bytes.Index(data, []byte("\n---"))
		sl = len("\n---")
	}
	if i >= 0 {
		// let's find the end of the line
		nli := bytes.Index(data[i+sl:], []byte("\n"))
		// there is no newline, the separator is the last line of the buffer
		if nli < 0 {
			// the separator is last in the file, return what is before it
			if atEOF {
				return i, data[:i], nil
			}
			// Request more data
			return 0, nil, nil
		}
		sl += nli + 1

		// there is nothing before, just skipping the separator and seek next
		if i == 0 {
			a, t, e := Split(data[i+sl:], atEOF)
			return a + i + sl, t, e
		}
		return i + sl, data[:i], nil
	}

	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return dataLen, data, nil
	}

	// Request more data.
	return 0, nil, nil
}
