var $ = require("jquery");
var anchorme = require("anchorme");
var mainWS = require("ws").getChannel("/ws");

var logs = {};
var ends = {};
var down = {};

exports.typeInfo = 0;
exports.typeTitle = 1;
exports.typeError = 2;

function printlogs(l, $console) {
    $.each(l, function (_, val) {
        var sett = {
            "html": anchorme(val.Message, {
                "attributes": [
                    {
                        "name": "target",
                        "value": "_blank"
                    }]
            }),
        };
        if (val.Type == exports.typeTitle || val.Type == exports.typeError) {
            sett.class = "header" + (val.Type == exports.typeError ? " error" : "");
        }
        $("<li/>", sett).appendTo($console);
    });
}

function getConsole(id) {
    var sel = "#console>ul[data-id=" + id + "]";
    if ($(sel).length === 0) {
        $("<ul/>", {
            "data-id": id,
        }).hide().appendTo($("#console"));
    }
    return $(sel);
}

exports.load = function (id) {
    $("#console>ul").hide();
    var $console = getConsole(id);
    $console.show();

    if (!logs[id]) {
        $.getJSON("/logs/" + id, function (data) {
            if (data) {
                logs[id] = data.logs;
                ends[id] = data.end;
                down[id] = true;
                printlogs(data.logs, $console);
                exports.scrollDown();
            }
        });
    } else {
        exports.scrollDown();
    }
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

mainWS.read("log", function (data) {
    var $console = getConsole(data.service);
    printlogs([data.message], $console);
    if (require("app").getActive() == data.service && down[data.service]) {
        exports.scrollDown();
    }
});

$(function () {
    $(window).on("resize", function (e) {
        if (down[require("app").getActive()]) {
            exports.scrollDown();
        }
    });
});
