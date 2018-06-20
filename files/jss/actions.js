var $ = require("jquery");
var mainWS = require("ws").getChannel("/ws");

$(function () {
    var actions = $("#actions");
    var list = require("list");
    actions.find(".stop").on("click", function (e) {
        mainWS.send("stop", list.getActiveId());
    });
    actions.find(".restart").on("click", function (e) {
        mainWS.send("restart", list.getActiveId());
    });
    actions.find(".start").on("click", function (e) {
        mainWS.send("start", list.getActiveId());
    });
});

exports.setStatus = function (status) {
    var app = require("app");
    var actions = $("#actions");

    actions.hide();
    if (typeof status !== "undefined" && status !== app.statusIgnored) {
        actions.show();
        actions.children().show();
    }

    switch (status) {
        case app.statusStopped:
            actions.find(".stop").hide();
            break;
        case app.statusWaiting:
        case app.statusRunned:
            actions.find(".start").hide();
            break;
    }
};

exports.setStatusById = function (id) {
    var item = require("list").getById(id);
    if (item) {
        exports.setStatus(item.data("status"));
    }
};