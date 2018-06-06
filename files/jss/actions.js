var $ = require("jquery");


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
}

exports.setStatusById = function (id) {
    var item = require("list").getById(id);
    if (item) {
        exports.setStatus(item.data("status"))
    }
}