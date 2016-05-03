process.title = 'web_service';

var path = require('path');
var express = require('express');
var bodyParser = require('body-parser');

var messageLayer = require('./message-layer');

var app = express();

app.use('/assets', express.static(path.join(__dirname, '../client')));

app.use(bodyParser.json());

app.set('view engine', 'html');
app.set('views', path.join(__dirname, '../client'));
app.engine('html', require('hbs').__express);

app.get('/', function (req, res) {
    res.render('index', {});
});

messageLayer(app);

app.listen(8004, function () {});
