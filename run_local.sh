cd ticker_service && go install && cd ..
cd ledger_service && go install && cd ..
cd orderbook_service && go install && cd ..

$GOPATH/bin/ticker_service &
$GOPATH/bin/ledger_service &
$GOPATH/bin/orderbook_service