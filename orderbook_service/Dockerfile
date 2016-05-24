FROM golang:1.6.2

COPY . /go/src/github.com/nickstefan/market/orderbook_service

RUN go get github.com/garyburd/redigo/redis
RUN go get gopkg.in/redsync.v1

RUN go install github.com/nickstefan/market/orderbook_service

EXPOSE 8001
CMD ["/go/bin/orderbook_service"]
