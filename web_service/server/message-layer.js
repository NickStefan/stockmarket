var expressWsConstructor = require('express-ws');
var Hub = require('./hub').Hub;

module.exports = function(app){
    var expressWs = expressWsConstructor(app);
    var hub = new Hub();

    app.hub = hub;      
    
    // receive messages intended for connected clients
    app.post('/msg/ticker/:ticker', function(req, res){
        app.hub.sendByTicker(req.params.ticker, req.body);
        res.sendStatus(200);
    });

    app.post('/msg/user/:user_id', function(req, res){
        app.hub.sendByUser(req.params.user_id, req.body);
        res.sendStatus(200);
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


