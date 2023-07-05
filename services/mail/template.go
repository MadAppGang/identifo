package mail

import (
	"bufio"
	"bytes"
	"fmt"
)

// extractSubjectAndBody extracts subject and email body from template
// template structure is
// ---
// subject text
// ---
// html body
func extractSubjectAndBody(d []byte) (string, string, error) {
	scanner := bufio.NewScanner(bytes.NewReader(d))
	scanner.Split(split)
	haveSubject := scanner.Scan()
	if !haveSubject {
		return "", "", fmt.Errorf("no subject") // TODO: localized error
	}
	subject := scanner.Text()
	haveBody := scanner.Scan()
	if !haveBody {
		return "", "", fmt.Errorf("no body") // TODO: localized error
	}
	body := scanner.Text()
	return subject, body, nil
}

// TODO: write unit tests
func split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	dataLen := len(data)

	// Return Nothing if at the end of file or no data passed.
	if atEOF && dataLen == 0 {
		return 0, nil, nil
	}

	// Find next separator and return token.
	if i := bytes.Index(data, []byte("---")); i >= 0 {
		return i + 3, data[0:i], nil
	}

	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return dataLen, data, nil
	}

	// Request more data.
	return 0, nil, nil
}
