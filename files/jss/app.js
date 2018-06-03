var $ = require("jquery");
var list = require("list");
var actions = require("actions");
var logs = require("logs");
var ws = require("ws");

exports.statusStopped = 0;
exports.statusWaiting = 1;
exports.statusRunned = 2;
exports.statusIgnored = 3;

exports.setActive = function (id, onpopstate) {
    list.setActive(id, onpopstate);
    actions.setStatusById(id);
    logs.load(id);
}

window.onpopstate = function (event) {
    var p = window.location.href.split("/");
    if (p.length) {
        exports.setActive(p[p.length - 1], event.state !== null);
    }
}