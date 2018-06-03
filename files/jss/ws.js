(function (factory) {
	if (typeof define === 'function' && define.amd) {
		define(['jquery'], factory);
	} else if (typeof exports === 'object') {
		factory(require('jquery'));
	} else {
		factory($);
	}
}(function ($) {
	$.ws = function () { };
	var url = "/ws/connect";
	var sock = null;
	var requestID = 0;
	var cbcs = {};
	var cmd_cbcs = {};
	var subscriptions = {};

	function createWebSocket(path) {
		return new WebSocket((location.protocol === 'https:' ? 'wss://' : 'ws://') + location.host + path);
	}

	function connect() {
		// данные получены
		function done(result) {
			try {
				result = JSON.parse(result);
			} catch (err) {
				return
			}

			if (result.requestID && cbcs[result.requestID]) {
				cbcs[result.requestID].fn(result.data, result.command)
				delete cbcs[result.requestID]
			} else if (cmd_cbcs[result.command]) {
				cmd_cbcs[result.command](result.data, result.command)
			}
			if (result && result.command) {
				subscriptions[result.command] = true;
				$(document).trigger('wsSubscibe', [result.command]);
			}

		}

		sock = createWebSocket(url);
		sock.onopen = function () {
			$(document).trigger('wsConnect');
		};
		sock.onclose = function (e) {
			setTimeout(connect, 300);
		};
		sock.onmessage = function (e) {
			if (e && typeof e.data === 'string' || e.data instanceof Blob) {
				if (e.data instanceof Blob) { // извлекаем бинарные данные
					var reader = new FileReader();
					reader.onload = function () {
						done(reader.result)
					};
					reader.readAsText(e.data);
				} else { // строковые данные (json)
					done(e.data)
				}
			}
		};
	}

	// отправить сообщение в сокет и получить ответ в коллбек
	$.ws.request = function (command, msg, callback) {
		requestID++;
		var reqID = requestID;
		if (callback) {
			cbcs[reqID] = {
				fn: callback,
				time: (new Date()).getTime() / 1000
			};
		}

		var msg = JSON.stringify(msg);
		if (window.dev) {
			msg = [reqID, command, msg].join(":");
		} else {
			try {
				msg = new Blob([reqID, ":", command, ":", msg]);
			} catch (e) {
				if (e.name == "InvalidStateError") {
					msg = [reqID, command, msg].join(":");
					var bb = new MSBlobBuilder();
					bb.append(msg);
					msg = bb.getBlob('text/csv;charset=utf-8');
				}
			}
		}

		if (sock && sock.readyState === WebSocket.OPEN) {
			sock.send(msg);
		} else {
			$(document).one('wsConnect', function (event) {
				sock.send(msg);
			});
		}
	}

	// повесить обработчик на сообщения, санкционированные сервером (без запроса)
	$.ws.subscribe = function (command, callback) {
		cmd_cbcs[command] = callback;
	}

	// повесить обработчик на сообщения, санкционированные сервером (без запроса)
	// и сделать запрос, чтобы сервер прислал первоначальные данные в коллбек
	$.ws.connect = function (command, msg, callback, autoReconnect) {
		$.ws.subscribe(command, callback);
		$.ws.request(command, msg, callback);
		if (typeof autoReconnect === "undefined" || autoReconnect) {
			$(document).on('wsConnect', function (event) {
				$.ws.request(command, msg, callback);
			});
		}
	}

	$.ws.wait = function (commands, callback) {
		var cmds = {}
		commands.forEach(function (command) {
			if (!subscriptions[command]) {
				cmds[command] = true;
			}
		});

		function wait(event, command) {
			if (cmds[command]) {
				delete cmds[command];
			}
			if (Object.keys(cmds).length < 1) {
				setTimeout(callback, 1);
			} else {
				$(document).one('wsSubscibe', wait);
			}
		}
		wait();
	}

	module.exports = $.ws;
	connect();
}));