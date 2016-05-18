package main

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/nickstefan/market/orderbook_service/heap"
	"strings"
)

type OrderQueue struct {
	data map[string]*heap.Heap
	pool *redis.Pool
	env  string
}

func (o *OrderQueue) setEnv(env string) {
	o.env = env
}

func (o *OrderQueue) Enqueue(queueName string, node *heap.Node) {
	if o.env == "TESTING" {
		o.data[queueName].Enqueue(node)
		return
	}

	conn := o.pool.Get()
	defer conn.Close()

	//serialized, err := json.Marshal(node)
	// read the ticker off the node, then do correct redis stuff
	_, err := conn.Do("ZADD", queueName, node.Value, node.Lookup) // serialized)

	if err != nil {
		fmt.Println("TODO: orderbook_service fault tolerance needed; ", err)
	}
}

func (o *OrderQueue) Dequeue(queueName string) *heap.Node {
	if o.env == "TESTING" {
		return o.data[queueName].Dequeue()
	}

	conn := o.pool.Get()
	defer conn.Close()

	var rankStart int
	var rankEnd int
	if true == strings.HasPrefix(queueName, "BUY") {
		rankStart = 0
		rankEnd = 0
	} else {
		rankStart = -1
		rankEnd = -1
	}
	serialized, err := redis.Bytes(conn.Do("ZREMRANGEBYRANK", queueName, rankStart, rankEnd))

	var node *heap.Node
	err = json.Unmarshal(serialized, &node)
	if err != nil {
		fmt.Println("TODO: ticker_service fault tolerance needed; ", err)
	}
	return node
}

func (o *OrderQueue) Peek(queueName string) *heap.Node {
	if o.env == "TESTING" {
		return o.data[queueName].Peek()
	}

	conn := o.pool.Get()
	defer conn.Close()

	var rankStart int
	var rankEnd int
	if true == strings.HasPrefix(queueName, "BUY") {
		rankStart = 0
		rankEnd = 0
	} else {
		rankStart = -1
		rankEnd = -1
	}
	serialized, err := redis.Bytes(conn.Do("ZRANGEBYRANK", queueName, rankStart, rankEnd))

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

func NewOrderQueue(pool *redis.Pool) *OrderQueue {
	// only used for testing OrderBook logic without having to run redis
	var data = make(map[string]*heap.Heap)
	data["BUYSTOCK"] = &heap.Heap{Priority: "max"}
	data["SELLSTOCK"] = &heap.Heap{Priority: "min"}

	return &OrderQueue{
		pool: pool,
		data: data,
	}
}
