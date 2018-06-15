var $ = require("jquery");
var anchorme = require("anchorme");

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
                "html": anchorme(val.Message, {
                    "attributes": [
                        {
                            "name": "target",
                            "value": "_blank"
                        }]
                }),
            }).appendTo(console);
        });
        exports.scrollDown();
    });
}

exports.scrollDown = function () {
    var console = $("#console");
    console.scrollTop(console[0].scrollHeight);
}

exports.scrollDownAnima = function () {
    var console = $("#console");
    console.animate({
        scrollTop: console[0].scrollHeight,
    }, 1200);
}