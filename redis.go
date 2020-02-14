package main

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/satori/go.uuid"
)

type RedisLock struct {
	lockKey string
	value   string
	rd      redis.Conn
	timeout int
}

//保证原子性（redis是单线程），避免del删除了，其他client获得的lock
var delScript = redis.NewScript(1, `
if redis.call("get", KEYS[1]) == ARGV[1] then
return redis.call("del", KEYS[1])
else
return 0
end`)

func NewGetRedis(key string, rd redis.Conn) RedisLock {
	return RedisLock{
		key,
		"",
		rd,
		1,
	}
}
func (this *RedisLock) Lock() error {
	uid, err := uuid.NewV4()
	if err != nil {
		return errors.New("uuid err:" + err.Error())
	}
	this.value = uid.String()

	lockReply, err := redis.String(this.rd.Do("SET", this.lockKey, this.value, "ex", this.timeout, "nx"))
	if err != nil {
		return errors.New("redis fail")
	}
	if lockReply == "OK" {
		return nil
	} else {
		return errors.New("lock fail")
	}
}

func (this *RedisLock) Unlock() error {
	_, err := delScript.Do(this.rd, this.lockKey, this.value)
	return err
}
func (this *RedisLock) GetValue(key string) (string, error) {
	reply, err := redis.String(this.rd.Do("get", this.lockKey))
	if err != nil {
		return "", errors.New("redis fail")
	}
	return reply, nil
}
