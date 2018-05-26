var $ = require("jquery");

exports.statusStopped = 0;
exports.statusWaiting = 1;
exports.statusRunned = 2;
exports.statusIgnored = 3;

exports.setStatus = function (status) {
    var actions = $("#actions");
    actions.children().show();

    switch (status) {
        case exports.statusStopped:
            actions.find(".stop").hide();
            break;
        case exports.statusWaiting:
        case exports.statusRunned:
            actions.find(".start").hide();
            break;
        case exports.statusIgnored:
            actions.find(".start").hide();
            actions.find(".stop").hide();
            actions.find(".restart").hide();
            break;
    }
}