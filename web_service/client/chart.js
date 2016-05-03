function Chart(options){
    options = options || {};
    this._data = options.data;
    this.selector = options.selector;
    this.label = options.label;
    this.periodMs = options.periodMs;
}

Chart.prototype.addData = function(data){
    this._data.push(data);
    this._data.shift();
}

// will need to setInterval,
// where on periodMs, add new whole node to end of this._data stack
// in between that interval, call this method to update in place ticks
Chart.prototype.addPartialData = function(data){
    var last = this._data[ this._data.length - 1];

    if (last.volume === 0){
        last.high = data.high;
        last.low = data.low;
        last.open = data.open;
        last.close = data.close;
    }

    if (last.low > data.low){
        last.low = data.low;
    }

    if (last.high < data.high){
        last.high = data.high;
    }

    last.close = data.close;
    last.volume = last.volume + data.volume;
}

Chart.prototype.draw = function(){

    // compute the bollinger bands
    var bollingerAlgorithm = fc.indicator.algorithm.bollingerBands();
    bollingerAlgorithm(this._data);

    // Offset the range to include the full bar for the latest value
    var xExtent = fc.util.extent()
        .fields(["date"])
        .padUnit("domain")
        .pad([this.periodMs * -bollingerAlgorithm.windowSize()(this._data), this.periodMs]);

    // ensure y extent includes the bollinger bands
    var yExtent = fc.util.extent().fields([
        function(d) { return d.bollingerBands.upper; },
        function(d) { return 0; }// return d.bollingerBands.lower; }
    ]);

    // create a chart
    var chart = fc.chart.cartesian(
            fc.scale.dateTime(),
            d3.scale.linear()
        )
        .xDomain(xExtent(this._data))
        .yDomain(yExtent(this._data))
        .yNice()
        .chartLabel(this.label)
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

    d3.select(this.selector)
        .datum(this._data)
        .call(chart);
};
