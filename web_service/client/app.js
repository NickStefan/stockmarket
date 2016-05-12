
// query localhost:8003/query
// TickerName:   "STOCK",
// Periods:      2,
// PeriodNumber: 1,
// PeriodName:   "minute",

// we dont want to closure in the chart
// we want to be able to add and remove tickers
// we want to render charts decoupled from socket

// maybe for now, we just async auto block the websockets until the chart
// loads, then chart will work for now, can come back to the front end later

async.auto({
    _data: function(done){
        $.ajax({
            url:"http://localhost:8003/query",
            method: "POST",
            data: {},
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
            socket = new WebSocket("ws://localhost:8004/ws");

            socket.onopen = function(e){
                var msg = JSON.stringify({
                    user_id: 'userNick',
                    tickers: ['STOCK']
                });
                socket.send(msg);
            };

            socket.onmessage = function(e){
                console.log(e.data);
                var msg = JSON.parse(e.data);

                switch (msg.api){
                    case 'ticker':
                        if (!msg.payload.volume){
                        return;
                    }
                    results.chart.addPartialData(msg.payload);
                    results.chart.draw();
                    break;
                default:
                    break;
                }
            };

            socket.onclose = function(e){
                console.log("connection closed");
            };
        }
        done(null);
    }]
}, function(err, results){
    
});
