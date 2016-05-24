FROM golang:1.6.2

COPY . /go/src/github.com/nickstefan/market/ledger_service

RUN go install github.com/nickstefan/market/ledger_service

EXPOSE 8002
CMD ["/go/bin/ledger_service"]
