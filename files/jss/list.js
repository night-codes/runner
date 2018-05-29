var $ = require("jquery");


// set active list item
exports.setActive = function (id) {
    var list = $("#list");
    var item = list.children("[data-id=" + id + "]");
    if (item) {
        list.children().removeClass("active");
        item.addClass("active");
        window.history.pushState("", "", '/service/' + id);
    }
}

// get id of active list item
exports.getActiveId = function () {
    var list = $("#list");
    var active = list.children(".active").one();
    if (active) {
        return active.data("id");
    }
    return 0;
}

// get list item by id
exports.getById = function (id) {
    return $("#list").children("[data-id=" + id + "]");
}

$(function () {
    var app = require("app");
    var list = $("#list");
    list.find("li>a").on("click", function (e) {
        e.preventDefault();
        app.setActive($(this).parent("li").data("id"));
    });
})