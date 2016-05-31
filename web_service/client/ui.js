var orderbookAPI = (
    window.location.protocol +
    "//" +
    window.location.hostname +
    window.location.port +
    "/orderbook"
);

var lastTradeSharesDOM = document.querySelector(".last-shares-value");
var lastTradePriceDOM = document.querySelector(".last-price-value");

var bidButton = document.querySelector(".bid");
var bidSharesDOM = document.querySelector(".bid-shares-value");
var bidPriceDOM = document.querySelector(".bid-price-value");

var askButton = document.querySelector(".ask");
var askSharesDOM = document.querySelector(".ask-shares-value");
var askPriceDOM= document.querySelector(".ask-price-value");

var buyButton = document.querySelector(".buy");
var buySharesInputDOM = document.querySelector(".buy-shares-value > input");
var buyPriceInputDOM = document.querySelector(".buy-price-value > input");

var sellButton = document.querySelector(".sell");
var sellSharesInputDOM = document.querySelector(".sell-shares-value > input");
var sellPriceInputDOM = document.querySelector(".sell-price-value > input");

function createOrder(intent){
    var order = {
        actor: "BOB",
        ticker: "STOCK",
        kind: "LIMIT",
        state: "OPEN",
        timecreated: (new Date()).getTime()
    };

    if ("buy" === intent.toLowerCase()){
        order.intent = "BUY"; 
        order.shares = parseInt(buySharesInputDOM.value);
        order.bid = parseFloat(buyPriceInputDOM.value);
    } else {
        order.intent = "SELL"; 
        order.shares = parseInt(sellSharesInputDOM.value);
        order.ask = parseFloat(sellPriceInputDOM.value);
    }

    return {
        uuid: Math.floor(Math.random() * 1000 * 1000),
        ticker: "STOCK",
        orders: [order]
    };
}

$(buyButton).on("click", function(){
    $.ajax({
        url: orderbookAPI,
        method: "POST",
        contentType: 'application/json; charset=utf-8',
        data: JSON.stringify(createOrder("buy")),
        success: function(data){
            //console.log(data);
        },
        error: function(){
            console.log(arguments);
        }
    });
});

$(sellButton).on("click", function(){
    $.ajax({
        url: orderbookAPI,
        method: "POST",
        contentType: 'application/json; charset=utf-8',
        data: JSON.stringify(createOrder("sell")),
        success: function(data){
            //console.log(data);
        },
        error: function(){
            console.log(arguments);
        }
    });

});

function setLastBidAsk(payload){
    bidSharesDOM.innerHTML = payload[0].shares;
    bidPriceDOM.innerHTML = payload[0].bid;
    askSharesDOM.innerHTML = payload[1].shares;
    askPriceDOM.innerHTML = payload[1].ask;
}

function setLastTrade(payload){
    lastTradeSharesDOM.innerHTML = payload.shares;
    lastTradePriceDOM.innerHTML = payload.price;
}
