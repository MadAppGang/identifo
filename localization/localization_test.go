package localization_test

import (
	"testing"

	"github.com/madappgang/identifo/v2/localization"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func TestLoadDefaultCatalog(t *testing.T) {
	err := localization.LoadDefaultCatalog()
	require.NoError(t, err)
	peng := message.NewPrinter(language.English)
	assert.Equal(t, "I am test string", peng.Sprintf("test"))

	pukr := message.NewPrinter(language.Ukrainian)
	assert.Equal(t, "Я стрічка для тестів", pukr.Sprintf("test"))
}
