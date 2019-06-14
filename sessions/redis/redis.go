package sessions

import (
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/madappgang/identifo/model"
)

const (
	defaultRedisAddress  = "localhost:6379"
	defaultRedisPassword = ""
	defaultRedisDB       = 0
)

// RedisSessionStorage is a Redis-backed storage for admin sessions.
type RedisSessionStorage struct {
	client *redis.Client
}

// NewSessionStorage creates new Redis session storage.
func NewSessionStorage(settings model.RedisSettings) (model.SessionStorage, error) {
	var addr, password string
	var db int

	if settings.Address == "" {
		addr = defaultRedisAddress
	}
	if settings.Password == "" {
		password = defaultRedisPassword
	}
	if settings.DB == 0 {
		db = defaultRedisDB
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	return &RedisSessionStorage{client: client}, nil
}

// GetSession fetches session by ID.
func (r *RedisSessionStorage) GetSession(id string) (model.Session, error) {
	var session model.Session

	bs, err := r.client.Get(id).Bytes()
	if err != nil {
		return session, err
	}

	err = json.Unmarshal(bs, &session)
	return session, err
}

// InsertSession inserts session to the storage.
func (r *RedisSessionStorage) InsertSession(session model.Session) error {
	bs, err := json.Marshal(session)
	if err != nil {
		return err
	}

	err = r.client.SetNX(session.ID, bs, time.Until(session.ExpirationDate)).Err()
	return err
}

// DeleteSession deletes session from the storage.
func (r *RedisSessionStorage) DeleteSession(id string) error {
	count, err := r.client.Del(id).Result()
	if count == 0 {
		log.Println("Tried to delete nonexistent session:", id)
	}

	return err
}

// ProlongSession sets new duration for the existing session.
func (r *RedisSessionStorage) ProlongSession(id string, newDuration model.SessionDuration) error {
	session, err := r.GetSession(id)
	if err != nil {
		return err
	}

	session.ExpirationDate = time.Now().Add(newDuration.Duration)

	bs, err := json.Marshal(session)
	if err != nil {
		return err
	}

	err = r.client.SetXX(session.ID, bs, newDuration.Duration).Err()
	return err
}
