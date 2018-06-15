global.dev = true;

var $ = require("jquery");

$(function () {
    var app = require("app");
    var list = require("list");
    app.setActive(list.getActiveId(), true);
})
