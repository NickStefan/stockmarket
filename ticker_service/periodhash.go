package main

import (
	"github.com/garyburd/redigo/redis"
)

type PeriodHash struct {
	data map[string]*Period
	pool *redis.Pool
}

func (o *PeriodHash) get(key string) *Period {
	return o.data[key]
}

func (o *PeriodHash) set(key string, period *Period) {
	o.data[key] = period
}

func (o *PeriodHash) remove(key string) {
	delete(o.data, key)
}

func NewPeriodHash(pool *redis.Pool, prefix string) *PeriodHash {
	return &PeriodHash{
		pool: pool,
		data: make(map[string]*Period),
	}
}
