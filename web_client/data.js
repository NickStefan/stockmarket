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
