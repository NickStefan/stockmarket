package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type TickAggregator struct {
	db *mgo.Database
}

func (t *TickAggregator) setDB(db *mgo.Database) {
	t.db = db
}

// periodType ex: 1min, 5min, 1day
// periodsRange ex: 5, 10, 20
// ticker ex: 'STOCK'

func NewGroup(period string, num int) bson.M {
	return bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"year":      bson.M{"$year": "$date"},
				"dayOfYear": bson.M{"$dayOfYear": "$date"},
				"interval": bson.M{
					"$subtract": []bson.M{
						{"$minute": "$date"},
						{"$mod": []interface{}{
							map[string]string{"$minute": "$date"},
							num,
						}},
					},
				},
			},
			"high":   bson.M{"$max": "$high"},
			"low":    bson.M{"$min": "$low"},
			"open":   bson.M{"$first": "$open"},
			"close":  bson.M{"$last": "$close"},
			"volume": bson.M{"$sum": "$volume"},
		},
	}
}

func (t *TickAggregator) query() {
	c := t.db.C("ticks")
	pipe := c.Pipe([]bson.M{
		{"$match": bson.M{"ticker": "STOCK"}},
		{"$sort": bson.M{"date": -1}}, // sorting seems off, open of 11 first period is wrong
		NewGroup("minute", 2),
		{"$sort": bson.M{"interval": -1}},
		{"$project": bson.M{
			"_id":      false,
			"interval": "$_id.interval",
			"open":     "$open",
			"close":    "$close",
			"high":     "$high",
			"low":      "$low",
			"ticker":   "$ticker",
			"volume":   "$volume",
		}},
	})
	results := []bson.M{}
	err := pipe.All(&results)
	if err != nil {
		fmt.Println("TODO: fault tolerance needed; ", err)
	}
	fmt.Println(results)
}
