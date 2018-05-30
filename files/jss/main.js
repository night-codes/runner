var $ = require("jquery");

$(function () {
    var actions = require("actions");
    var list = require("list");

    // app.setActive(list.getActiveId());
    actions.setStatusById(list.getActiveId());
})
