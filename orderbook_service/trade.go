package main

type Trade struct {
	Actor  string  `json:"actor"`
	Shares int     `json:"shares"`
	Ticker string  `json:"ticker"`
	Price  float64 `json:"price"`
	Intent string  `json:"intent"`
	Kind   string  `json:"kind"`
	State  string  `json:"state"`
	Time   int64   `json:"time"`
}

type AnonymizedTrade struct {
	Shares int     `json:"shares"`
	Ticker string  `json:"ticker"`
	Price  float64 `json:"price"`
	Kind   string  `json:"kind"`
	Time   int64   `json:"time"`
}

func AnonymizeTrade(trade Trade) AnonymizedTrade {
	return AnonymizedTrade{
		Shares: trade.Shares,
		Ticker: trade.Ticker,
		Price:  trade.Price,
		Kind:   trade.Kind,
		Time:   trade.Time,
	}
}
