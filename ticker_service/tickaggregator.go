package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Query struct {
	TickerName   string
	PeriodName   string
	PeriodNumber int
	Periods      int
	StartDate    time.Time
	EndDate      time.Time
}

func (q *Query) MatchGroupSortProject() (bson.M, bson.M, bson.M, bson.M) {
	group := bson.M{}
	sort := bson.M{}
	match := bson.M{}
	project := bson.M{}

	startDate := time.Time{}
	endDate := time.Now()
	now := time.Now()

	if false == q.StartDate.IsZero() {
		startDate = q.StartDate
	}

	if false == q.EndDate.IsZero() {
		endDate = q.EndDate
	}

	// query for one extra period to ensure we dont miss it by a half period
	// the limit action that happens in actual Pipe then corrects it back
	periods := q.Periods + 1

	switch q.PeriodName {
	case "minute":
		if true == startDate.IsZero() {
			startDate = now.Add(-1 * (time.Duration(periods*q.PeriodNumber) * time.Minute))
		}

		group = bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"year":      bson.M{"$year": "$date"},
					"dayOfYear": bson.M{"$dayOfYear": "$date"},
					"hour":      bson.M{"$hour": "$date"},
					"interval": bson.M{
						"$subtract": []bson.M{
							{"$minute": "$date"},
							{"$mod": []interface{}{
								bson.M{"$minute": "$date"},
								q.PeriodNumber,
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

		sort = bson.M{"$sort": bson.M{
			"_id.year":      1,
			"_id.dayOfYear": 1,
			"_id.hour":      1,
			"_id.interval":  1,
		}}

		project = bson.M{"$project": bson.M{
			"_id": false,
			"date": bson.M{
				"$add": []interface{}{
					time.Date(startDate.Year(), time.January, 0, 0, 0, 0, 0, time.UTC),
					bson.M{
						"$multiply": []interface{}{
							bson.M{"$subtract": []interface{}{"$_id.year", startDate.Year()}},
							365 * 24 * 60 * 60 * 1000,
						},
					},
					bson.M{
						"$multiply": []interface{}{"$_id.dayOfYear", 24 * 60 * 60 * 1000},
					},
					bson.M{
						"$multiply": []interface{}{"$_id.hour", 60 * 60 * 1000},
					},
					bson.M{
						"$multiply": []interface{}{"$_id.interval", 60 * 1000},
					},
				},
			},
			"interval": "$_id.interval",
			"open":     "$open",
			"close":    "$close",
			"high":     "$high",
			"low":      "$low",
			"ticker":   "$ticker",
			"volume":   "$volume",
		}}

	case "hour":
		startDate = now.Add(-1 * (time.Duration(periods*q.PeriodNumber) * time.Hour))

		group = bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"year":      bson.M{"$year": "$date"},
					"dayOfYear": bson.M{"$dayOfYear": "$date"},
					"interval": bson.M{
						"$subtract": []bson.M{
							{"$hour": "$date"},
							{"$mod": []interface{}{
								bson.M{"$hour": "$date"},
								q.PeriodNumber,
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

		sort = bson.M{"$sort": bson.M{
			"_id.year":      1,
			"_id.dayOfYear": 1,
			"_id.interval":  1,
		}}

		project = bson.M{"$project": bson.M{
			"_id": false,
			"date": bson.M{
				"$add": []interface{}{
					time.Date(startDate.Year(), time.January, 0, 0, 0, 0, 0, time.UTC),
					bson.M{
						"$multiply": []interface{}{
							bson.M{"$subtract": []interface{}{"$_id.year", startDate.Year()}},
							365 * 24 * 60 * 60 * 1000,
						},
					},
					bson.M{
						"$multiply": []interface{}{"$_id.dayOfYear", 24 * 60 * 60 * 1000},
					},
					bson.M{
						"$multiply": []interface{}{"$_id.interval", 60 * 1000},
					},
				},
			},
			"interval": "$_id.interval",
			"open":     "$open",
			"close":    "$close",
			"high":     "$high",
			"low":      "$low",
			"ticker":   "$ticker",
			"volume":   "$volume",
		}}

	case "day":
		if true == startDate.IsZero() {
			startDate = now.Add(-1 * (time.Duration(periods*q.PeriodNumber*24) * time.Hour))
		}

		group = bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"year":      bson.M{"$year": "$date"},
					"dayOfYear": bson.M{"$dayOfYear": "$date"},
				},
				"high":   bson.M{"$max": "$high"},
				"low":    bson.M{"$min": "$low"},
				"open":   bson.M{"$first": "$open"},
				"close":  bson.M{"$last": "$close"},
				"volume": bson.M{"$sum": "$volume"},
			},
		}

		sort = bson.M{"$sort": bson.M{
			"_id.year":      1,
			"_id.dayOfYear": 1,
		}}

		project = bson.M{"$project": bson.M{
			"_id": false,
			"date": bson.M{
				"$add": []interface{}{
					time.Date(startDate.Year(), time.January, 0, 0, 0, 0, 0, time.UTC),
					bson.M{
						"$multiply": []interface{}{
							bson.M{"$subtract": []interface{}{"$_id.year", startDate.Year()}},
							365 * 24 * 60 * 60 * 1000,
						},
					},
					bson.M{
						"$multiply": []interface{}{"$_id.dayOfYear", 24 * 60 * 60 * 1000},
					},
				},
			},
			"interval": "$_id.interval",
			"open":     "$open",
			"close":    "$close",
			"high":     "$high",
			"low":      "$low",
			"ticker":   "$ticker",
			"volume":   "$volume",
		}}

	case "week":
		if true == startDate.IsZero() {
			startDate = now.Add(-1 * (time.Duration(periods*q.PeriodNumber*24*7) * time.Hour))
		}

		group = bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"year": bson.M{"$year": "$date"},
					"week": bson.M{"$week": "$date"},
				},
				"high":   bson.M{"$max": "$high"},
				"low":    bson.M{"$min": "$low"},
				"open":   bson.M{"$first": "$open"},
				"close":  bson.M{"$last": "$close"},
				"volume": bson.M{"$sum": "$volume"},
			},
		}

		sort = bson.M{"$sort": bson.M{
			"_id.year": 1,
			"_id.week": 1,
		}}

		project = bson.M{"$project": bson.M{
			"_id": false,
			"date": bson.M{
				"$add": []interface{}{
					time.Date(startDate.Year(), time.January, 0, 0, 0, 0, 0, time.UTC),
					bson.M{
						"$multiply": []interface{}{
							bson.M{"$subtract": []interface{}{"$_id.year", startDate.Year()}},
							365 * 24 * 60 * 60 * 1000,
						},
					},
					bson.M{
						"$multiply": []interface{}{"$_id.week", 7 * 24 * 60 * 60 * 1000},
					},
				},
			},
			"interval": "$_id.interval",
			"open":     "$open",
			"close":    "$close",
			"high":     "$high",
			"low":      "$low",
			"ticker":   "$ticker",
			"volume":   "$volume",
		}}

	case "month":
		if true == startDate.IsZero() {
			startDate = now.Add(-1 * (time.Duration(periods*q.PeriodNumber*24*31) * time.Hour))
		}

		group = bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"year":  bson.M{"$year": "$date"},
					"month": bson.M{"$month": "$date"},
				},
				"high":   bson.M{"$max": "$high"},
				"low":    bson.M{"$min": "$low"},
				"open":   bson.M{"$first": "$open"},
				"close":  bson.M{"$last": "$close"},
				"volume": bson.M{"$sum": "$volume"},
			},
		}

		sort = bson.M{"$sort": bson.M{
			"_id.year":  1,
			"_id.month": 1,
		}}

		project = bson.M{"$project": bson.M{
			"_id": false,
			"date": bson.M{
				"$add": []interface{}{
					time.Date(startDate.Year(), time.January, 0, 0, 0, 0, 0, time.UTC),
					bson.M{
						"$multiply": []interface{}{
							bson.M{"$subtract": []interface{}{"$_id.year", startDate.Year()}},
							365 * 24 * 60 * 60 * 1000,
						},
					},
					bson.M{
						"$multiply": []interface{}{"$_id.month", 31 * 24 * 60 * 60 * 1000},
					},
				},
			},
			"interval": "$_id.interval",
			"open":     "$open",
			"close":    "$close",
			"high":     "$high",
			"low":      "$low",
			"ticker":   "$ticker",
			"volume":   "$volume",
		}}

	}

	fmt.Println("startDate", startDate)
	fmt.Println("endDate", endDate)

	match = bson.M{"$match": bson.M{
		"ticker": q.TickerName,
		"date": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}}

	return match, group, sort, project
}

type TickAggregator struct {
	db   *mgo.Database
	hash *PeriodHash
}

func (t *TickAggregator) setDB(db *mgo.Database) {
	t.db = db
}

func (t *TickAggregator) setKV(hash *PeriodHash) {
	t.hash = hash
}

func (t *TickAggregator) currentTicker(ticker string) *Period {
	return t.hash.get(ticker)
}

func (t *TickAggregator) query(q Query) []interface{} {
	match, group, sort, project := q.MatchGroupSortProject()

	c := t.db.C("ticks")
	pipe := c.Pipe([]bson.M{
		match,
		{"$sort": bson.M{"date": 1}},
		group,
		sort,
		project,
	})

	var results []interface{}
	err := pipe.All(&results)
	if err != nil {
		fmt.Println("TODO: fault tolerance needed; ", err)
	}

	// ADD REDIS PERIOD IF IN QUERY
	endDate := time.Now()
	if false == q.EndDate.IsZero() {
		endDate = q.EndDate
	}
	currentTicker := t.currentTicker(q.TickerName)
	if endDate.Unix() > currentTicker.Date.Unix() {
		results = append(results, currentTicker)
	}

	limitStart := 0
	if len(results) > q.Periods {
		limitStart = len(results) - q.Periods
	}
	return results[limitStart:]
}
