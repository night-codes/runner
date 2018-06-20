var $ = require("jquery");

exports.confirm = function (title, msg, callback) {
    var $dialog = $("#dialog");
    $dialog.find('.header h3').text(title);
    $dialog.find('.dialog-msg').text(msg);
    $dialog.css({ "display": "flex" });
    $dialog.find('.doAction').on("click", function afn() {
        $dialog.hide();
        $(document).off("click", afn);
        setTimeout(callback, 10);
    });
};

$(function () {
    var $dialog = $("#dialog");

    $dialog.find('.cancelAction, .fa-close').on("click", function () {
        $dialog.hide();
    });

    $dialog.find('.dialog').on("click", function (e) {
        e.preventDefault();
        e.stopPropagation();
    });

    $dialog.on("click", function (e) {
        $dialog.hide();
    });
});