FROM golang:1.6.2

COPY . /go/src/github.com/nickstefan/market/ticker_service

RUN go get gopkg.in/mgo.v2
RUN go get github.com/garyburd/redigo/redis
RUN go get gopkg.in/redsync.v1

RUN go install github.com/nickstefan/market/ticker_service

EXPOSE 8080
CMD ["/go/bin/ticker_service"]
