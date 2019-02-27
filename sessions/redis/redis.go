package sessions

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/madappgang/identifo/model"
)

const (
	// RedisAddress is a host:port Redis network address.
	RedisAddress = "REDIS_ADDRESS"
	// RedisPassword is a password to connect to Redis.
	RedisPassword = "REDIS_PASSWORD"
	// RedisDB is an enumerator for database to be selected after connecting to the Redis server.
	RedisDB = "REDIS_DB"
)

type redisStorage struct {
	client *redis.Client
}

// NewSessionStorageFromEnv creates new Redis session storage getting all settings from env.
func NewSessionStorageFromEnv() (model.SessionStorage, error) {
	addr := os.Getenv(RedisAddress)
	password := os.Getenv(RedisPassword)
	db, err := strconv.Atoi(os.Getenv(RedisDB))

	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	return &redisStorage{client: client}, err
}

func (r *redisStorage) GetSession(id string) (model.Session, error) {
	var session model.Session

	bs, err := r.client.Get(id).Bytes()
	if err != nil {
		return session, err
	}

	err = json.Unmarshal(bs, &session)
	return session, err
}

func (r *redisStorage) InsertSession(session model.Session) error {
	bs, err := json.Marshal(session)
	if err != nil {
		return err
	}

	err = r.client.SetNX(session.ID, bs, time.Until(session.ExpirationDate)).Err()
	return err
}

func (r *redisStorage) DeleteSession(id string) error {
	count, err := r.client.Del(id).Result()
	if count == 0 {
		log.Println("Tried to delete nonexistent session:", id)
	}

	return err
}

func (r *redisStorage) ProlongSession(id string, newDuration time.Duration) error {
	session, err := r.GetSession(id)
	if err != nil {
		return err
	}

	session.ExpirationDate = time.Now().Add(newDuration)

	bs, err := json.Marshal(session)
	if err != nil {
		return err
	}

	err = r.client.SetXX(session.ID, bs, newDuration).Err()
	return err
}
