var expressWsConstructor = require('express-ws');
var Hub = require('./hub').Hub;

module.exports = function(app){
    var expressWs = expressWsConstructor(app);
    var hub = new Hub();

    app.hub = hub;      
    
    var start;
    var end;
    
    // receive messages intended for connected clients
    app.post('/msg/ticker/:ticker', function(req, res){
        if (req.body.payload.price === 2){
            start = new Date().getTime();
            console.log(start);
        }
        if (req.body.payload.price === 70){
            end = new Date().getTime();
            console.log(end);
        }
        app.hub.sendByTicker(req.params.ticker, req.body);
        res.sendStatus(200);
    });

    app.post('/msg/user/:user_id', function(req, res){
        app.hub.sendByUser(req.params.user_id, req.body);
        res.sendStatus(200);
    });

    // DEBUG AWS
    app.get('/msg/info', function(req, res){
        res.json({
            start: start,
            end: end
        });
    });

    // receive client connections
    // for now client only sends an initial message with identifying info
    app.ws('/ws', function(client, req) {
        client.on('message', function(msg){
            msg = JSON.parse(msg);

            client.on('close', function(){
                hub.removeClient(msg);
            });

            hub.addClient(msg, client);
        });
    });
};


