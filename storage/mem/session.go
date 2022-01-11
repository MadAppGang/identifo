package mem

import (
	"sync"
	"time"

	"github.com/madappgang/identifo/v2/model"
)

type memoryStorage struct {
	sync.Mutex
	sessions map[string]model.Session
}

// NewSessionStorage creates an in-memory session storage.
func NewSessionStorage() model.SessionStorage {
	return &memoryStorage{
		sessions: make(map[string]model.Session),
	}
}

func (m *memoryStorage) GetSession(id string) (model.Session, error) {
	session, ok := m.sessions[id]
	if !ok {
		return model.Session{}, model.ErrorNotFound
	}

	return session, nil
}

func (m *memoryStorage) InsertSession(session model.Session) error {
	m.Lock()
	defer m.Unlock()

	m.sessions[session.ID] = session
	return nil
}

func (m *memoryStorage) DeleteSession(id string) error {
	m.Lock()
	defer m.Unlock()

	delete(m.sessions, id)
	return nil
}

func (m *memoryStorage) ProlongSession(id string, newDuration model.SessionDuration) error {
	m.Lock()
	defer m.Unlock()

	session, ok := m.sessions[id]
	if !ok {
		return model.ErrorNotFound
	}

	session.ExpirationTime = time.Now().Add(newDuration.Duration).Unix()

	m.sessions[session.ID] = session
	return nil
}
