var expressWsConstructor = require('express-ws');

module.exports = function(app){
    var expressWs = expressWsConstructor(app);
    var hub = new Hub();

    app.hub = hub;      
    
    // receive messages intended for connected clients
    app.post('/msg', function(req, res){
        app.hub.sendByTicker(req.body.ticker, req.body);
        res.sendStatus(200);
    });

    // receive client connections
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

function Hub(){
    this.clientsByUserId = {};
    this.clientsByTicker = {};
}

Hub.prototype.addClient = function(msg, client){
    this.clientsByUserId[ msg.user_id ] = client;

    msg.tickers
    .forEach(function(ticker){
        if (undefined === this.clientsByTicker[ticker]){
            this.clientsByTicker[ticker] = {};
        }
        this.clientsByTicker[ticker][msg.user_id] = client;
    }.bind(this));

    this.pingPong(client);
};

Hub.prototype.removeClient = function(msg){
    delete this.clientsByUserId[msg.user_id];

    Object.keys(this.clientsByTicker)
    .forEach(function(ticker){
        delete this.clientsByTicker[ticker][msg.user_id];
    }.bind(this));
};

Hub.prototype.sendByUser = function(user_id, msg){
    var client = this.clientsByUserId[user_id];
    if (!client){
        return;
    }
    client.send(JSON.stringify(msg));
};

Hub.prototype.sendByTicker = function(ticker, msg){
    msg = JSON.stringify(msg);
    if (undefined === this.clientsByTicker[ticker]){
        return;
    }

    Object.keys(this.clientsByTicker[ticker])
    .forEach(function(clientName){
        this.clientsByTicker[ticker][clientName].send(msg);
    }.bind(this));
};

Hub.prototype.pingPong = function(client){
    client.pingssent = 0;

    var interval = setInterval(function() {
        if (client.pingssent >= 2) {
            client.close();
            clearInterval(interval);
        } else {
            client.ping();
            client.pingssent++;
        }
    }, 75*1000);

    client.on("pong", function() { 
        client.pingssent = 0; 
    });
}

