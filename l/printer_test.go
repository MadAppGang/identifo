package l_test

import (
	"testing"

	"github.com/madappgang/identifo/v2/l"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestPrinter(t *testing.T) {
	printer, err := l.NewPrinter("en")
	require.NoError(t, err)
	require.NotNil(t, printer)

	assert.Equal(t, "I am test string", printer.SD(l.Test))
	assert.Equal(t, "I am test string", printer.S(language.English, l.Test))
	assert.Equal(t, "I am test string", printer.SL("en", l.Test))
	assert.Equal(t, "I am test string", printer.SL("uu", l.Test))
	assert.Equal(t, "Я стрічка для тестів", printer.SL("ukr", l.Test))
	assert.Equal(t, "Я стрічка для тестів", printer.SL("uk;q=0.9, en;q=0.8, de;q=0.7, *;q=0.5", l.Test))
}
