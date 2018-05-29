var $ = require("jquery");


exports.setStatus = function (status) {
    var actions = $("#actions");
    var app = require("app");
    actions.children().show();

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