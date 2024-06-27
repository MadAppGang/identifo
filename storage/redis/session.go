package redis

import (
	"encoding/json"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/madappgang/identifo/v2/logging"
	"github.com/madappgang/identifo/v2/model"
)

// RedisSessionStorage is a Redis-backed storage for admin sessions.
type RedisSessionStorage struct {
	logger *slog.Logger
	client redis.Cmdable
	prefix string
}

// NewSessionStorage creates new Redis session storage.
func NewSessionStorage(
	logger *slog.Logger,
	settings model.RedisDatabaseSettings,
) (model.SessionStorage, error) {
	var addr, password string
	var db int

	if settings.Address == "" {
		addr = defaultRedisAddress
	} else {
		addr = settings.Address
	}

	if settings.Password == "" {
		password = defaultRedisPassword
	} else {
		password = settings.Password
	}

	var client redis.Cmdable

	if settings.Cluster {
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    strings.Split(settings.Address, ","),
			Password: password,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		})
	}

	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	p := strings.TrimSpace(settings.Prefix)
	if p != "" && !strings.HasSuffix(settings.Prefix, ":") {
		p = p + ":"
	}

	return &RedisSessionStorage{
		logger: logger,
		client: client,
		prefix: p,
	}, nil
}

// GetSession fetches session by ID.
func (r *RedisSessionStorage) GetSession(id string) (model.Session, error) {
	var session model.Session

	key := r.prefix + id
	bs, err := r.client.Get(key).Bytes()
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

	key := r.prefix + session.ID
	err = r.client.SetNX(key, bs, time.Until(time.Unix(session.ExpirationTime, 0))).Err()
	return err
}

// DeleteSession deletes session from the storage.
func (r *RedisSessionStorage) DeleteSession(id string) error {
	key := r.prefix + id
	count, err := r.client.Del(key).Result()
	if count == 0 {
		r.logger.Warn("Tried to delete non existent session",
			"key", key,
			logging.FieldError, err)
	}

	return err
}

// ProlongSession sets new duration for the existing session.
func (r *RedisSessionStorage) ProlongSession(id string, newDuration model.SessionDuration) error {
	session, err := r.GetSession(id)
	if err != nil {
		return err
	}

	session.ExpirationTime = time.Now().Add(newDuration.Duration).Unix()

	bs, err := json.Marshal(session)
	if err != nil {
		return err
	}

	key := r.prefix + session.ID
	err = r.client.SetXX(key, bs, newDuration.Duration).Err()
	return err
}

func (r *RedisSessionStorage) Close() {
	if c, ok := r.client.(io.Closer); ok {
		c.Close()
	}
}
