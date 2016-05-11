package main

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

type PeriodHash struct {
	data   map[string]*Period
	pool   *redis.Pool
	env    string
	prefix string
}

func (o *PeriodHash) setEnv(env string) {
	o.env = env
}

func (o *PeriodHash) get(key string) *Period {
	if o.env == "TESTING" {
		return o.data[key]
	}

	conn := o.pool.Get()
	defer conn.Close()

	serialized, err := redis.Bytes(conn.Do("GET", o.prefix+key))

	var period *Period
	err = json.Unmarshal(serialized, &period)
	if err != nil {
		fmt.Println("TODO: ticker_service fault tolerance needed; ", err)
	}
	return period
}

func (o *PeriodHash) set(key string, period *Period) {
	if o.env == "TESTING" {
		o.data[key] = period
		return
	}

	conn := o.pool.Get()
	defer conn.Close()

	serialized, err := json.Marshal(period)
	_, err = conn.Do("SET", o.prefix+key, serialized)
	if err != nil {
		fmt.Println("TODO: ticker_service fault tolerance needed; ", err)
	}
}

func (o *PeriodHash) remove(key string) {
	if o.env == "TESTING" {
		delete(o.data, key)
		return
	}

	conn := o.pool.Get()
	defer conn.Close()

	conn.Do("DEL", o.prefix+key)
}

func NewPeriodHash(pool *redis.Pool, prefix string) *PeriodHash {
	return &PeriodHash{
		pool:   pool,
		prefix: prefix,
		data:   make(map[string]*Period),
	}
}
