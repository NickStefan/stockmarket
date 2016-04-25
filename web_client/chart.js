
function Chart(options){
    //var data = options.data;

    var margin = {top: 20, right: 20, bottom: 30, left: 50};
    var width = 960 - margin.left - margin.right;
    var height = 500 - margin.top - margin.bottom;
    this.width = width;
    this.height = height;

    this.count = 20;

    this.x = techan.scale.financetime().range([0, width]);
    this.y = d3.scale.linear().range([height, 0]);
    var ohlc = techan.plot.ohlc().xScale(this.x).yScale(this.y);

    this.yVolume = d3.scale.linear().range([this.y(0), this.y(0.2)]);


    this.candlestick = techan.plot.candlestick().xScale(this.x).yScale(this.y);

    this.sma0 = techan.plot.sma().xScale(this.x).yScale(this.y);

    this.sma0Calculator = techan.indicator.sma().period(10);

    this.sma1 = techan.plot.sma().xScale(this.x).yScale(this.y);

    this.sma1Calculator = techan.indicator.sma().period(20);


    this.accessor = this.candlestick.accessor();

    // Set the accessor to a ohlc accessor so we get highlighted bars
    this.volume = techan.plot.volume().accessor(ohlc.accessor()).xScale(this.x).yScale(this.yVolume);

    this.xAxis = d3.svg.axis().scale(this.x).orient("bottom");

    this.yAxis = d3.svg.axis().scale(this.y).orient("left");

    this.volumeAxis = d3.svg.axis().scale(this.yVolume).orient("right").ticks(3).tickFormat(d3.format(",.3s"));

    var timeAnnotation = techan.plot.axisannotation().axis(this.xAxis).format(d3.time.format('%Y-%m-%d')).width(65).translate([0, height]);

    var ohlcAnnotation = techan.plot.axisannotation().axis(this.yAxis).format(d3.format(',.2fs'));

    var volumeAnnotation = techan.plot.axisannotation().axis(this.volumeAxis).width(35);

    this.crosshair = techan.plot.crosshair().xScale(this.x).yScale(this.y).xAnnotation(timeAnnotation).yAnnotation([ohlcAnnotation, volumeAnnotation]);

    this.svg = d3.select("body").append("svg").attr("width", width + margin.left + margin.right).attr("height", height + margin.top + margin.bottom);

    var defs = this.svg.append("defs");

    defs.append("clipPath")
        .attr("id", "ohlcClip")
        .append("rect")
            .attr("x", 0)
            .attr("y", 0)
            .attr("width", width)
            .attr("height", height);

    this.svg = this.svg.append("g").attr("transform", "translate(" + margin.left + "," + margin.top + ")");

    this.createCandleSelection();

    this.svg.append("g").attr("class", "x axis").attr("transform", "translate(0," + height + ")");

    this.svg.append("g")
        .attr("class", "y axis")
        .append("text")
            .attr("transform", "rotate(-90)")
            .attr("y", 6)
            .attr("dy", ".71em")
            .style("text-anchor", "end")
            .text("Price ($)");

    this.svg.append("g").attr("class", "volume axis");

    this.svg.append('g').attr("class", "crosshair ohlc");

};

function refreshIndicator(selection, indicator, data) {
    var datum = selection.datum();
    // Some trickery to remove old and insert new without changing array reference,
    // so no need to update __data__ in the DOM
    datum.splice.apply(datum, [0, datum.length].concat(data));
    selection.call(indicator);
}

Chart.prototype.createCandleSelection = function(){
    this.candleSelection = this.svg.append("g").attr("class", "ohlc").attr("transform", "translate(0,0)");
    this.candleSelection.append("g").attr("class", "volume").attr("clip-path", "url(#ohlcClip)");
    this.candleSelection.append("g").attr("class", "candlestick").attr("clip-path", "url(#ohlcClip)");
    this.candleSelection.append("g").attr("class", "indicator sma ma-0").attr("clip-path", "url(#ohlcClip)");
    this.candleSelection.append("g").attr("class", "indicator sma ma-1").attr("clip-path", "url(#ohlcClip)");
};

Chart.prototype.mapData = function(d){
    return {
        date: new Date(d.time * 1000),
        open: +d.open,
        high: +d.high,
        low: +d.low,
        close: +d.close,
        volume: +d.volume
    };
};

Chart.prototype.addData = function (options){
    this.data = this.data || [];
    if (options !== undefined && options.data){
        this.data.push.apply(this.data, options.data.map(this.mapData));
    } else if (options !== undefined && !options.data) {
        // Simulate intra day updates when no feed is left
        var last = this.data[this.data.length-1];
        // Last must be between high and low
        last.close = Math.round(((last.high - last.low)*Math.random())*10)/10+last.low;
    }
};

Chart.prototype.createMinimumCount = function(){
    this.data = this.data || [];
    this.presentationData = this.data.slice();

    var lastDate;
    if (this.data.length){
        lastDate = this.data[ this.data.length - 1].date.getTime();
    } else {
        lastDate = (new Date()).getTime();
    }

    while (this.presentationData.length < this.count){
        lastDate = lastDate + 6000;
        this.presentationData.push({
            date: lastDate,
            open: 0,
            high: 0,
            low: 0,
            close: 0,
            volume: 0
        });
    }
};

Chart.prototype.draw = function(options){

    this.createMinimumCount();

    //this.presentationData = this.presentationData
    //.map(this.mapData)
    //.sort(function(a, b) {
        //return d3.ascending(this.accessor.d(a), this.accessor.d(b));
    //}.bind(this));

    //this.svg.select("g.candlestick").data(this.presentationData);
    //this.svg.select("g.sma.ma-0").datum(this.sma0Calculator(this.presentationData));
    //this.svg.select("g.sma.ma-1").datum(this.sma1Calculator(this.presentationData));
    //this.svg.select("g.volume").datum(this.presentationData);

    this.svg.selectAll("g.ohlc").remove()

    this.x = techan.scale.financetime().range([0, this.width]);
    this.y = d3.scale.linear().range([this.height, 0]);
    this.candlestick = techan.plot.candlestick().xScale(this.x).yScale(this.y);
    this.accessor = this.candlestick.accessor();

    this.x.domain(this.presentationData.map(this.accessor.d));

    // Show only 150 points on the plot
    // this.x.zoomable().domain([this.presentationData.length-130, this.presentationData.length]);
    //this.x.zoomable().domain([0, this.presentationData.length]);

    // Update y scale min max, only on viewable zoomable.domain()
    // this.y.domain(techan.scale.plot.ohlc(this.presentationData.slice(this.presentationData.length-130, this.presentationData.length)).domain());
    this.y.domain(techan.scale.plot.ohlc(this.presentationData.slice(0, this.presentationData.length)).domain());
    // this.yVolume.domain(techan.scale.plot.volume(this.presentationData.slice(this.presentationData.length-130, this.presentationData.length)).domain());
    //this.yVolume.domain(techan.scale.plot.volume(this.presentationData.slice(0, this.presentationData.length)).domain());

    this.svg.select('g.x.axis').call(this.xAxis);
    this.svg.select('g.y.axis').call(this.yAxis);
    this.svg.select("g.volume.axis").call(this.volumeAxis);

    this.createCandleSelection();
    this.svg.select("g.candlestick").data(this.presentationData);
    this.svg.select("g.candlestick").call(this.candlestick);

    // this.svg.select("g.candlestick").call(this.candlestick);

    // Recalculate indicators and update the SAME array and redraw moving average
    // refreshIndicator(this.svg.select("g.sma.ma-0"), this.sma0, this.sma0Calculator(this.presentationData));
    // refreshIndicator(this.svg.select("g.sma.ma-1"), this.sma1, this.sma1Calculator(this.presentationData));

    // this.svg.select("g.volume").call(this.volume);
    //this.svg.selectAll("g.volume").remove();
    //this.svg.select("g.volume").datum(this.presentationData);

    //this.svg.select("g.crosshair.ohlc").call(this.crosshair);
}
