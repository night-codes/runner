var $ = require("jquery");
var anchorme = require("anchorme");
var mainWS = require("ws").getChannel("/ws");
var dialog = require("dialog");

var logs = {};
var starts = {};
var ends = {};
var down = {};

exports.typeInfo = 0;
exports.typeTitle = 1;
exports.typeError = 2;

function printlogs(l, $console, pre) {
    $.each(l, function (i, val) {
        if (!pre || i < pre) {
            var sett = {
                "html": anchorme(val.Message, {
                    "attributes": [
                        {
                            "name": "target",
                            "value": "_blank"
                        }]
                }),
            };
            if (val.Type == exports.typeTitle) {
                sett.class = "header";
            } else if (val.Type == exports.typeError) {
                sett.class = "error";
            }
            if (pre) {
                $("<li/>", sett).prependTo($console);
            } else {
                $("<li/>", sett).appendTo($console);
            }
        }
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
                printlogs(data.logs, $console, starts[id]);
                exports.scrollDown();
            }
        });
    } else {
        exports.scrollDown();
    }
};

exports.scrollDown = function () {
    var console = $("#console");
    console.scrollTop(console[0].scrollHeight);
};

exports.scrollDownAnima = function () {
    var console = $("#console");
    console.animate({
        scrollTop: console[0].scrollHeight,
    }, 1200);
};

mainWS.read("log", function (data) {
    var $console = getConsole(data.service);
    printlogs([data.message], $console);
    if (typeof starts[data.service] === "undefined") {
        starts[data.service] = data.end;
    }
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
    $("#console").on('mousedown', 'li', function (e) {
        var ctrlDown = e.ctrlKey || e.metaKey; // Mac support
        if (ctrlDown) {
            e.preventDefault();
            var opened = $(this).hasClass("opened");
            $("#console li.opened").removeClass("opened");
            if (!opened) {
                $(this).addClass("opened");
            }
        }
    });
    $(document).on('keydown', function (e) {
        e.preventDefault();
        var key = event.key || event.keyCode;
        var ctrlDown = e.ctrlKey || e.metaKey; // Mac support
        if (ctrlDown && key == 'c') {
            document.execCommand('copy');
        }
    });

    $("#console").bind("contextmenu", function (e) {
        e.preventDefault();

        $("#cntnr").css("left", e.pageX - 5 + "px");
        $("#cntnr").css("top", e.pageY - 5 + "px");
        $("#cntnr").show();
        $("#cntnr").on("mousedown", function sfn() {
            e.preventDefault();
            e.stopPropagation();
            $("#cntnr").hide();
            $(document).off("mousedown", sfn);
            return false;
        });
        $("#copyMenuItem").on("mousedown", function () {
            document.execCommand('copy');
        });
        $("#refreshMenuItem").on("mousedown", function () {
            location.reload();
        });
        $("#clearMenuItem").on("mousedown", function () {
            dialog.confirm("Are you sure?", "Do you want to clear all service logs? This action can not be undone.", function () {
                var id = require("app").getActive();
                mainWS.send("clear", id);
                getConsole(id).empty();
            });
        });
        $(document).on("mousedown", function sfn() {
            $("#cntnr").hide();
            $(document).off("mousedown", sfn);
        });
    });
});
