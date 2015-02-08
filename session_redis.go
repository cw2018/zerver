package zerver

import (
	"github.com/cosiner/golib/encoding"
	"github.com/cosiner/gomodule/redis"
)

// redisStore is a session store use redis
type redisStore struct {
	store *redis.RedisStore
}

// NewRedisStore create a session store background on redis
func NewRedisStore() SessionStore {
	return new(redisStore)
}

// Init init redis store, config like maxidle=*&idletimeout=*&addr=*
func (rstore *redisStore) Init(conf string) (err error) {
	if rstore.store == nil {
		rstore.store, err = redis.NewRedisStore(conf)
	}
	return
}

// Destroy destroy redis store, release resources
func (rstore *redisStore) Destroy() {
	rstore.store.Destroy()
}

// IsExist check whether given id of node is exist
func (rstore *redisStore) IsExist(id string) bool {
	exist, _ := rstore.store.IsExist(id)
	return exist
}

// Save save values with given id and lifetime
func (rstore *redisStore) Save(id string, values Values, lifetime int64) {
	if lifetime != 0 {
		if bs, err := encoding.GOBEncode(values); err == nil {
			go rstore.store.SetWithExpire(id, bs, lifetime)
		}
	}
}

// Get return values of given id
func (rstore *redisStore) Get(id string) (vals Values) {
	if bs, err := redis.ToBytes(rstore.store.Get(id)); err == nil && len(bs) != 0 {
		vals = make(Values)
		encoding.GOBDecode(bs, &vals)
	}
	return
}

// Rename move values exist in old id to new id
func (rstore *redisStore) Rename(oldId string, newId string) {
	rstore.store.Update("RENAME", oldId, newId)
}
