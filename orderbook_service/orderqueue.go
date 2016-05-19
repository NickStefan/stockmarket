package main

import (
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

	_, err := conn.Do("ZADD", queueName, node.Value, node.Lookup)

	if err != nil {
		fmt.Println("orderbook_service: orderqueue enqueue ", err)
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
		rankStart = -1
		rankEnd = -1
	} else {
		rankStart = 0
		rankEnd = 0
	}

	conn.Send("MULTI")
	conn.Send("ZRANGE", queueName, rankStart, rankEnd, "WITHSCORES")
	conn.Send("ZREMRANGEBYRANK", queueName, rankStart, rankEnd)
	values, err := redis.Values(conn.Do("EXEC"))

	var lookup string
	var score float64
	var zremRes int
	var zrangeRes []interface{}
	_, err = redis.Scan(values, &zrangeRes, &zremRes)
	_, err = redis.Scan(zrangeRes, &lookup, &score)
	if err != nil {
		fmt.Println("orderbook_service: orderqueue dequeue ", err)
		return nil
	}

	return &heap.Node{
		Lookup: lookup,
		Value:  score,
	}
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
		rankStart = -1
		rankEnd = -1
	} else {
		rankStart = 0
		rankEnd = 0
	}

	var lookup string
	var score float64
	values, err := redis.Values(conn.Do("ZRANGE", queueName, rankStart, rankEnd, "WITHSCORES"))
	if len(values) < 2 {
		return nil
	}
	_, err = redis.Scan(values, &lookup, &score)

	if err != nil {
		fmt.Println("orderbook_service: orderqueue peek", err)
		return nil
	}

	return &heap.Node{
		Lookup: lookup,
		Value:  score,
	}
}

func (o *OrderQueue) remove(key string) {

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
