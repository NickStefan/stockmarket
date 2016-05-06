
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
    client.send(JSON.stringify(msg), function(err){
        if (err) client.close();
    });
};

Hub.prototype.sendByTicker = function(ticker, msg){
    msg = JSON.stringify(msg);
    if (undefined === this.clientsByTicker[ticker]){
        return;
    }

    Object.keys(this.clientsByTicker[ticker])
    .forEach(function(clientName){
        this.clientsByTicker[ticker][clientName].send(msg, function(err){
            if(err) client.close();
        });
    }.bind(this));
};

Hub.prototype.pingPong = function(client){
    client.pingssent = 0;
    client.on('error', function(){
        client.close();
    });

    var interval = setInterval(function() {
        if (client.pingssent >= 1) {
            client.close();
            clearInterval(interval);
        } else {
            try{
                client.ping();
            } catch(e){
                client.close();
            }
            client.pingssent++;
        }
    }, 10*1000);

    client.on("pong", function() { 
        client.pingssent = 0; 
    });
};

module.exports.Hub = Hub;
