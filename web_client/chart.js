window._data = fc.data.random.financial().stream().take(200);

var dte = new Date().getTime();
for (var i = 0; i < window._data.length; i++){
    window._data[i].date = new Date(dte + (i * 1000 * 60))
}

var _stream = window._data.splice(40, 160);
var stream = new Stream(_stream);
//window._data = stream.take(110);

//window._data = window._data.map(mapData);

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

function renderChart(data) {
    if (data){
    window._data.push(data);
    window._data.shift();
    }
    var data = window._data;

    // compute the bollinger bands
    var bollingerAlgorithm = fc.indicator.algorithm.bollingerBands();
    bollingerAlgorithm(data);

    // Offset the range to include the full bar for the latest value
    var DAY_MS = 1000 * 60// * 60 * 24;
    var xExtent = fc.util.extent()
        .fields(["date"])
        .padUnit("domain")
        .pad([DAY_MS * -bollingerAlgorithm.windowSize()(data), DAY_MS]);

    // ensure y extent includes the bollinger bands
    var yExtent = fc.util.extent().fields([
        function(d) { return d.bollingerBands.upper; },
        function(d) { return d.bollingerBands.lower; }
    ]);

    // create a chart
    var chart = fc.chart.cartesian(
            fc.scale.dateTime(),
            d3.scale.linear()
        )
        .xDomain(xExtent(data))
        .yDomain(yExtent(data))
        .yNice()
        .chartLabel("Streaming Candlestick")
        .margin({left: 30, right: 30, bottom: 20, top: 30});

    // obtain ticks from the underlying scales
    var xTicks = chart.xScaleTicks(10);
    var yTicks = chart.yScaleTicks(10);

    // render a reduced number of ticks on each axis
    //chart
        //.xTickValues(xTicks.filter(function(d) { return d.getDate() % 2 === 0; }))
        //.yTickValues(yTicks.filter(function(d, i) { return i % 2 === 0; }));

    // Create the gridlines and series
    var gridlines = fc.annotation.gridline()
        .xTickValues(xTicks)
        .yTickValues(yTicks);
    var candlestick = fc.series.candlestick();
    var bollingerBands = fc.indicator.renderer.bollingerBands();

    // add them to the chart via a multi-series
    var multi = fc.series.multi()
        .series([gridlines, bollingerBands, candlestick]);

    chart.plotArea(multi);

    d3.select("#streaming-chart")
        .datum(data)
        .call(chart);
}

renderChart(stream.next());

setInterval(function(){
    if (stream.length !== 0){
        renderChart(stream.next());
    }
}, 1000);
