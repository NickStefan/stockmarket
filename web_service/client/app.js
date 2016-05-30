// we dont want to closure in the chart
// we want to be able to add and remove tickers
// we want to render charts decoupled from socket

// maybe for now, we just async auto block the websockets until the chart
// loads, then chart will work for now, can come back to the front end later
var tickerAPI = (
    window.location.protocol +
    "//" +
    window.location.hostname +
    window.location.port +
    "/ticker/query"
);

var messageAPI = (
    "ws://" +
    window.location.hostname +
    window.location.port + 
    "/ws"
);

var lastTradeDOM = document.querySelector(".last-value");
var bidDOM = document.querySelector(".bid-value");
var askDOM = document.querySelector(".ask-value");

function setBidAsk(payload){
    bidDOM.innerHTML = payload[0].shares + " @ " + payload[0].bid;
    askDOM.innerHTML = payload[1].shares + " @ " + payload[1].ask;
}

function setLastTrade(payload){
    lastTradeDOM.innerHTML = payload.shares + " @ " + payload.price;
}

async.auto({
    _data: function(done){
        $.ajax({
            url: tickerAPI,
            method: "POST",
            contentType: 'application/json; charset=utf-8',
            dataType: 'json',
            data: JSON.stringify({
                tickerName: "STOCK",
                periods: 50,
                periodNumber: 1,
                periodName: "minute"
            }),
            success: function(data){
                done(null, data);
            },
            error: function(){
                console.log(arguments);
            }
        });
    },
    data: ['_data', function(results, done){
        var data = results._data.map(function(p){
            p.date = new Date(p.date);
            return p;
        });
        done(null, data);
    }],
    chart: ['data', function(results, done){
        var chart = new Chart({
            data: results.data,
            selector: "#streaming-chart",
            label: "STOCK chart",
            periodMs: 1000 * 60, // * 60 * 24;
            periods: 50
        });
        chart.draw();    
        done(null, chart);
    }],
    sockets: ['chart', function(results, done){
        if (window.WebSocket){
            function connectChart(chart){
                var socket = new WebSocket(messageAPI);

                socket.onopen = function(e){
                    console.log("connected");
                    var msg = JSON.stringify({
                        user_id: 'userNick',
                        tickers: ['STOCK']
                    });
                    socket.send(msg);
                };

                socket.onmessage = function(e){
                    var msg = JSON.parse(e.data);

                    switch (msg.api){
                        case 'ticker':
                            console.log(msg.payload.shares, msg.payload.time);
                            setLastTrade(msg.payload);
                            chart.addPartialData(msg.payload);
                            chart.draw();
                            break;
                        case 'bid-ask':
                            setBidAsk(msg.payload);
                        default:
                            break;
                    }
                };

                socket.onclose = function(e){
                    console.log("connection closed");
                    setTimeout(function(){
                        console.log("retrying...");
                        connectChart(chart);
                    }, 10*1000);
                };
            }
            connectChart(results.chart);
        }
        done(null);
    }]
}, function(err, results){

});
