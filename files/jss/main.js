var $ = require("jquery");

$(function () {
    var actions = require("actions");
    var list = require("list");

    actions.setStatusById(list.getActiveId());
})
