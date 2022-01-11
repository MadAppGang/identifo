package mem

import (
	"github.com/madappgang/identifo/v2/model"
)

type ConnectionTester struct{}

// NewConnectionTester creates a memory connection tester, which never has errors.

func NewConnectionTester() model.ConnectionTester {
	return &ConnectionTester{}
}

func (ct *ConnectionTester) Connect() error {
	return nil
}
