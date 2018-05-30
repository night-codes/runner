var $ = require("jquery");


exports.setStatus = function (status) {
    var app = require("app");
    var actions = $("#actions");

    actions.children().hide();
    if (typeof status !== "undefined") {
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
        case app.statusIgnored:
            actions.find(".start").hide();
            actions.find(".stop").hide();
            actions.find(".restart").hide();
            break;
    }
}

exports.setStatusById = function (id) {
    var item = require("list").getById(id);
    if (item) {
        exports.setStatus(item.data("status"))
    }
}