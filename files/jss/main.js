global.dev = true;

var $ = require("jquery");
var ws = require("ws")

$(function () {
    var app = require("app");
    var list = require("list");

    app.setActive(list.getActiveId(), true);


    var ch = ws.getChannel("/ws/connect1")
    var ch2 = ws.getChannel("/ws/connect2")
    ch.subscribe("news");
    ch2.subscribe("news");
})
