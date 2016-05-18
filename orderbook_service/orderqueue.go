package main

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/nickstefan/market/orderbook_service/heap"
)

type OrderQueue struct {
	data heap.Heap
	pool *redis.Pool
	env  string
}

func (o *OrderQueue) setEnv(env string) {
	o.env = env
}

func (o *OrderQueue) Enqueue(node *heap.Node) {
	if o.env == "TESTING" {
		o.data.Enqueue(node)
	}

	conn := o.pool.Get()
	defer conn.Close()

	serialized, err := json.Marshal(node)
	// read the ticker off the node, then do correct redis stuff
	_, err = conn.Do("SET", "", serialized)

	if err != nil {
		fmt.Println("TODO: orderbook_service fault tolerance needed; ", err)
	}

}

func (o *OrderQueue) Dequeue() *heap.Node {
	if o.env == "TESTING" {
		return o.data.Dequeue()
	}

	conn := o.pool.Get()
	defer conn.Close()

	serialized, err := redis.Bytes(conn.Do("GET", ""))

	var node *heap.Node
	err = json.Unmarshal(serialized, &node)
	if err != nil {
		fmt.Println("TODO: ticker_service fault tolerance needed; ", err)
	}
	return node
}

func (o *OrderQueue) Peek() *heap.Node {
	if o.env == "TESTING" {
		return o.data.Peek()
	}

	conn := o.pool.Get()
	defer conn.Close()

	serialized, err := redis.Bytes(conn.Do("GET", ""))

	var node *heap.Node
	err = json.Unmarshal(serialized, &node)
	if err != nil {
		fmt.Println("TODO: ticker_service fault tolerance needed; ", err)
	}
	return node

}

func (o *OrderQueue) remove(key string) {
	//if o.env == "TESTING" {
	//delete(o.data, key)
	//return
	//}

	//conn := o.pool.Get()
	//defer conn.Close()

	//conn.Do("DEL", o.prefix+key)
}

func NewOrderQueue(pool *redis.Pool, priority string) *OrderQueue {
	return &OrderQueue{
		pool: pool,
		data: heap.Heap{Priority: priority},
	}
}
