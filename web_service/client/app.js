
// TickerName:   "STOCK",
// Periods:      2,
// PeriodNumber: 1,
// PeriodName:   "minute",

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
    data: ['_data', function(done, results){
        var data = results._data.map(function(p){
            p.date = new Date(p.date);
            return p;
        });
        done(null, data);
    }],
    chart: ['data', function(done, results){
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
    sockets: ['chart', function(done, results){
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
                    if (msg && msg.payload && msg.payload.volume){
                        console.log(msg.payload.volume, msg.payload.date);
                    }
                    switch (msg.api){
                        case 'ticker':
                            if (!msg.payload.volume){
                            return;
                        }
                        chart.addPartialData(msg.payload);
                        chart.draw();
                        break;
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
