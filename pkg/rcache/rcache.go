package rcache

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"

	"gopkg.in/inconshreveable/log15.v2"

	"sourcegraph.com/sourcegraph/sourcegraph/pkg/env"
)

const (
	// dataVersion is used for releases that change type struture for
	// data that may already be cached. Increasing this number will
	// change the key prefix that is used for all hash keys,
	// effectively resetting the cache at the same time the new code
	// is deployed.
	dataVersion = "v1"
)

// Cache implements httpcache.Cache
type Cache struct {
	keyPrefix  string
	ttlSeconds int
}

// New creates a redis backed Cache
func New(keyPrefix string) *Cache {
	return &Cache{
		keyPrefix: keyPrefix,
	}
}

// NewWithTTL creates a redis backed Cache which expires values after
// ttlSeconds.
func NewWithTTL(keyPrefix string, ttlSeconds int) *Cache {
	return &Cache{
		keyPrefix:  keyPrefix,
		ttlSeconds: ttlSeconds,
	}
}

// Get implements httpcache.Cache.Get
func (r *Cache) Get(key string) ([]byte, bool) {
	c := pool.Get()
	defer c.Close()

	b, err := redis.Bytes(c.Do("GET", r.rkeyPrefix()+key))
	if err != nil && err != redis.ErrNil {
		log15.Warn("failed to execute redis command", "cmd", "GET", "error", err)
	}

	return b, err == nil
}

// Delete implements httpcache.Cache.Set
func (r *Cache) Set(key string, b []byte) {
	c := pool.Get()
	defer c.Close()

	if r.ttlSeconds == 0 {
		_, err := c.Do("SET", r.rkeyPrefix()+key, b)
		if err != nil {
			log15.Warn("failed to execute redis command", "cmd", "SET", "error", err)
		}
	} else {
		_, err := c.Do("SETEX", r.rkeyPrefix()+key, r.ttlSeconds, b)
		if err != nil {
			log15.Warn("failed to execute redis command", "cmd", "SETEX", "error", err)
		}
	}
}

// Delete implements httpcache.Cache.Delete
func (r *Cache) Delete(key string) {
	c := pool.Get()
	defer c.Close()

	_, err := c.Do("DEL", r.rkeyPrefix()+key)
	if err != nil {
		log15.Warn("failed to execute redis command", "cmd", "DEL", "error", err)
	}
}

func (r *Cache) Keys(pattern string) []string {
	c := pool.Get()
	defer c.Close()

	prefix := r.rkeyPrefix()
	keys, err := redis.Strings(c.Do("KEYS", prefix+pattern))
	if err != nil {
		log15.Warn("failed to execute redis command", "cmd", "KEYS", "error", err)
	}
	for i := range keys {
		keys[i] = keys[i][len(prefix):]
	}
	return keys
}

// rkeyPrefix generates the actual key prefix we use on redis.
func (r *Cache) rkeyPrefix() string {
	return fmt.Sprintf("%s:%s:", globalPrefix, r.keyPrefix)
}

// SetupForTest adjusts the globalPrefix and clears it out. You will have
// conflicts if you do `t.Parallel()`
func SetupForTest(name string) {
	globalPrefix = "__test__" + name
	// Make mutex fails faster
	mutexTries = 1
	c := pool.Get()
	defer c.Close()
	_, err := c.Do("EVAL", `local keys = redis.call('keys', ARGV[1])
if #keys > 0 then
	return redis.call('del', unpack(keys))
else
	return ''
end`, 0, globalPrefix+":*")
	if err != nil {
		log15.Error("Could not clear test prefix", "name", name, "globalPrefix", globalPrefix, "error", err)
	}
}

var (
	pool         *redis.Pool
	globalPrefix string
)

var redisMasterEndpoint = env.Get("REDIS_MASTER_ENDPOINT", "redis-cache:6379", "redis used for caches")

func init() {
	globalPrefix = dataVersion

	pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisMasterEndpoint)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
