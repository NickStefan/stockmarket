function Chart(options){
    options = options || {};
    this._data = options.data;
    this.selector = options.selector;
    this.label = options.label;
    this.periodMs = options.periodMs;
    this.periods = options.periods;

    //this.addPeriod(); // should definitely check how recent last period is!!!

    // but this interval part is mostly right
    this.interval = setInterval(function(){
        this.addPeriod();
        this.draw();
    }.bind(this), this.periodMs);
}

Chart.prototype.cleanup = function(){   
    clearInterval(this.interval);
}

Chart.prototype.addData = function(data){
    this._data.push(data);
    this._data.shift();
}

Chart.prototype.addPeriod = function(){
    var lastPeriod = this._data[ this._data.length - 1];

    var newPeriod = {
        date: new Date(),
        high: lastPeriod.close,
        low: lastPeriod.close,
        open: lastPeriod.close,
        close: lastPeriod.close,
        volume: 0
    };

    this._data.push(newPeriod);

    if (this._data > this.periods){
        this._data.shift();
    }
};

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

Chart.prototype.getVisibleRange = function(){
    var beginning  = this._data[0].date; // should default to new Date() ?
    var end = new Date( beginning.getTime() + this.periods * this.periodMs);
    return [beginning, end];
};

// TODO make bollinger bands optional
// good way to flush out making the chart extendible 
Chart.prototype.draw = function(){

    // compute the bollinger bands
    //var bollingerAlgorithm = fc.indicator.algorithm.bollingerBands();
    //bollingerAlgorithm(this._data);

    // Offset the range to include the full bar for the latest value
    var xExtent = fc.util.extent()
        .fields(["date"])
        .padUnit("domain")
        .pad([this.periodMs, this.periodMs])
        //.pad([this.periodMs * -bollingerAlgorithm.windowSize()(this._data), this.periodMs])
        .include(this.getVisibleRange(this._data));

    var yExtent = fc.util.extent().fields(
        ['low', 'high']
    );

    // ensure y extent includes the bollinger bands
    //var yExtent = fc.util.extent().fields([
        //function(d) { return d.bollingerBands.upper; },
        //function(d) { return 0; }// return d.bollingerBands.lower; }
    //]);

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
    //var bollingerBands = fc.indicator.renderer.bollingerBands();

    // add them to the chart via a multi-series
    var multi = fc.series.multi()
        .series([gridlines, candlestick]);

    chart.plotArea(multi);

    d3.select(this.selector)
        .datum(this._data)
        .call(chart);
};
