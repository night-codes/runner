var $ = require("jquery");
var actions = require("actions");

exports.setActive = function (id) {
    var list = $("#list");
    var active = list.children("[data-id=" + id + "]");
    if (active) {
        list.children().removeClass("active");
        active.addClass("active");
        actions.setStatus(active.data("status"));
    }
}
exports.current = function () {
    var list = $("#list");
    var active = list.children(".active");
    if (active) {
        return active.data("id");
    }
    return 0;
}

$(function () {
    var list = $("#list");
    list.find("li>a").on("click", function (e) {
        e.preventDefault();
        exports.setActive($(this).parent("li").data("id"));
    });
    exports.setActive(exports.current());
})