window._data = [
    {
        "high" : 11,
        "low" : 2,
        "open" : 5,
        "close" : 10,
        "volume" : 1300,
        "ticker" : "STOCK",
        "time" : 1460610409
    },
    {
        "high" : 19,
        "low" : 11,
        "open" : 11,
        "close" : 19,
        "volume" : 2000,
        "ticker" : "STOCK",
        "time" : 1460610469
    },
    {
        "high" : 27,
        "low" : 15,
        "open" : 19,
        "close" : 23,
        "volume" : 2000,
        "ticker" : "STOCK",
        "time" : 1460610531
    },
    {
        "high" : 22,
        "low" : 18,
        "open" : 22,
        "close" : 22,
        "volume" : 1200,
        "ticker" : "STOCK",
        "time" : 1460610592
    },
    {
        "high" : 22,
        "low" : 15,
        "open" : 22,
        "close" : 18,
        "volume" : 1200,
        "ticker" : "STOCK",
        "time" : 1460610652
    },
    {
        "high" : 19,
        "low" : 15,
        "open" : 18,
        "close" : 16,
        "volume" : 1200,
        "ticker" : "STOCK",
        "time" : 1460610712
    },
    {
        "high" : 18,
        "low" : 15,
        "open" : 16,
        "close" : 18,
        "volume" : 1200,
        "ticker" : "STOCK",
        "time" : 1460610872
    },
    {
        "high" : 22,
        "low" : 18,
        "open" : 19,
        "close" : 21,
        "volume" : 1200,
        "ticker" : "STOCK",
        "time" : 1460610932
    },
    {
        "high" : 24,
        "low" : 21,
        "open" : 21,
        "close" : 24,
        "volume" : 1200,
        "ticker" : "STOCK",
        "time" : 1460611092
    },
    {
        "high" : 27,
        "low" : 23,
        "open" : 24,
        "close" : 27,
        "volume" : 1200,
        "ticker" : "STOCK",
        "time" : 1460611152
    },
    {
        "high" : 27,
        "low" : 15,
        "open" : 26,
        "close" : 18,
        "volume" : 1200,
        "ticker" : "STOCK",
        "time" : 1460611213
    },
    {
        "high" : 27,
        "low" : 15,
        "open" : 19,
        "close" : 23,
        "volume" : 2000,
        "ticker" : "STOCK",
        "time" : 1460611273
    },
    {
        "high" : 22,
        "low" : 18,
        "open" : 22,
        "close" : 22,
        "volume" : 1200,
        "ticker" : "STOCK",
        "time" : 1460611333
    },
    {
        "high" : 22,
        "low" : 15,
        "open" : 22,
        "close" : 18,
        "volume" : 1200,
        "ticker" : "STOCK",
        "time" : 1460611399
    },
    {
        "high" : 19,
        "low" : 15,
        "open" : 18,
        "close" : 16,
        "volume" : 1200,
        "ticker" : "STOCK",
        "time" : 1460611459
    },
    {
        "high" : 18,
        "low" : 15,
        "open" : 16,
        "close" : 18,
        "volume" : 1200,
        "ticker" : "STOCK",
        "time" : 1460611519
    },
    {
        "high" : 22,
        "low" : 18,
        "open" : 19,
        "close" : 21,
        "volume" : 1200,
        "ticker" : "STOCK",
        "time" : 1460611579
    },
    {
        "high" : 24,
        "low" : 21,
        "open" : 21,
        "close" : 24,
        "volume" : 1200,
        "ticker" : "STOCK",
        "time" : 1460611629
    },
    {
        "high" : 27,
        "low" : 23,
        "open" : 24,
        "close" : 27,
        "volume" : 1200,
        "ticker" : "STOCK",
        "time" : 1460611689
    },
    {
        "high" : 27,
        "low" : 15,
        "open" : 26,
        "close" : 18,
        "volume" : 1200,
        "ticker" : "STOCK",
        "time" : 1460611759
    },
];

window._data = fc.data.random.financial().stream().take(200);

var dte = new Date().getTime();
for (var i = 0; i < window._data.length; i++){
    window._data[i].date = new Date(dte + (i * 1000 * 60))
}

var _stream = window._data.splice(40, 160);
var stream = new Stream(_stream);

function Stream(_stream){
    this.stream = _stream; //.map(mapData);
    this.length = this.stream.length;
}

function mapData (d){
    d.date = (new Date(d.time));
    delete d.time;
    return d;
}

Stream.prototype.next = function(){
    this.length--;
    return this.stream.shift();
}
