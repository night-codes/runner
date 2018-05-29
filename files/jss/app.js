var $ = require("jquery");
var list = require("list");
var actions = require("actions");

exports.statusStopped = 0;
exports.statusWaiting = 1;
exports.statusRunned = 2;
exports.statusIgnored = 3;

exports.setActive = function (id) {
    list.setActive(id);
    actions.setStatusById(id);
}