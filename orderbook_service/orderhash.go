package main

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

type OrderHash struct {
	data   map[string]*Order
	pool   *redis.Pool
	env    string
	prefix string
}

func (o *OrderHash) setEnv(env string) {
	o.env = env
}

func (o *OrderHash) get(key string) *Order {
	if o.env == "TESTING" {
		return o.data[key]
	}

	conn := o.pool.Get()
	defer conn.Close()

	serialized, err := redis.Bytes(conn.Do("GET", o.prefix+key))

	var order *Order
	err = json.Unmarshal(serialized, &order)
	if err != nil {
		fmt.Println("TODO: ticker_service fault tolerance needed; ", err)
	}
	if o.prefix == "minute" {
		//fmt.Println("get", order)
	}
	return order
}

func (o *OrderHash) set(key string, order *Order) {
	if o.env == "TESTING" {
		o.data[key] = order
		return
	}

	conn := o.pool.Get()
	defer conn.Close()

	serialized, err := json.Marshal(order)
	_, err = conn.Do("SET", o.prefix+key, serialized)
	if o.prefix == "minute" {
		//fmt.Println("res", res)
	}
	if err != nil {
		fmt.Println("TODO: ticker_service fault tolerance needed; ", err)
	}
}

func (o *OrderHash) remove(key string) {
	if o.env == "TESTING" {
		delete(o.data, key)
		return
	}

	conn := o.pool.Get()
	defer conn.Close()

	conn.Do("DEL", o.prefix+key)
}

func NewOrderHash(pool *redis.Pool, prefix string) *OrderHash {
	return &OrderHash{
		pool:   pool,
		prefix: prefix,
		data:   make(map[string]*Order),
	}
}
