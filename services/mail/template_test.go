package mail_test

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/madappgang/identifo/v2/services/mail"
	"github.com/stretchr/testify/assert"
)

func TestSplit(t *testing.T) {
	template := `---
This is a subject.
---
<html>
	<body>
		<h1>I am a template body!</h1>
	</body>
</html>
`
	scanner := bufio.NewScanner(bytes.NewReader([]byte(template)))
	scanner.Split(mail.Split)
	assert.True(t, scanner.Scan())
	assert.Equal(t, "This is a subject.", scanner.Text())
	assert.True(t, scanner.Scan())
	txt := scanner.Text()
	assert.False(t, scanner.Scan())
	assert.Contains(t, txt, "<html>")
	assert.Contains(t, txt, "I am a template body!")
	expected := `<html>
	<body>
		<h1>I am a template body!</h1>
	</body>
</html>
`
	assert.Equal(t, expected, txt)
}

func TestSplitTrailingSeparator(t *testing.T) {
	template := `---
This is a subject.
---l;sdakjfl;asdjfl;skadjfkl;sdj     
<html>
	<body>
		<h1>I am a template body!</h1>
	</body>
</html>
----
`
	scanner := bufio.NewScanner(bytes.NewReader([]byte(template)))
	scanner.Split(mail.Split)
	assert.True(t, scanner.Scan())
	assert.Equal(t, "This is a subject.", scanner.Text())
	assert.True(t, scanner.Scan())
	txt := scanner.Text()
	assert.False(t, scanner.Scan())
	assert.Contains(t, txt, "<html>")
	assert.Contains(t, txt, "I am a template body!")
	expected := `<html>
	<body>
		<h1>I am a template body!</h1>
	</body>
</html>`
	assert.Equal(t, expected, txt)
}
