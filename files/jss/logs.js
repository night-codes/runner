var $ = require("jquery");

exports.typeInfo = 0;
exports.typeTitle = 1;
exports.typeError = 2;

exports.load = function (id) {
    var console = $("#console");
    console.html("");

    $.getJSON("/logs/" + id + "?json", function (data) {
        $.each(data, function (key, val) {
            $("<li/>", {
                "class": val.Type == exports.typeTitle ? "header" : (val.Type == exports.typeError ? "error header" : ""),
                "html": val.Message,
            }).appendTo(console);
        });
    });
}