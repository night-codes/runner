var list = require("list");
var actions = require("actions");
var logs = require("logs");
var aid = 0;

exports.statusStopped = 0;
exports.statusWaiting = 1;
exports.statusRunned = 2;
exports.statusIgnored = 3;

exports.getActive = function () {
    return aid;
};

exports.setActive = function (id, onpopstate) {
    aid = id;
    list.setActive(id, onpopstate);
    actions.setStatusById(id);
    logs.load(id);
};

window.onpopstate = function (event) {
    var p = window.location.href.split("/");
    if (p.length) {
        exports.setActive(p[p.length - 1], event.state !== null);
    }
};