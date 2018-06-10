(function (factory) {
	if (typeof define === 'function' && define.amd) {
		define("ws", ["exports"], factory);
	} else if (typeof exports === 'object') {
		factory(exports);
	} else {
		factory(global || window);
	}
}(function (exports) {
	var channels = {};

	function createWebSocket(path) {
		return new WebSocket((location.protocol === 'https:' ? 'wss://' : 'ws://') + location.host + path);
	}

	function Channel(url) {
		var global = global || window;
		var document = global.document;
		var sock = null;
		var prevID = 0;
		var requestTimeout = 0;
		var readers = {}; // удалить
		var subscriptions = {}; // нужны для того, чтобы переотправлять после реконнекта
		var self = this;
		var cid = "" + (Math.random().toFixed(16).substring(2) + new Date().valueOf()) + url;


		function on(type, fn) {
			document.addEventListener(cid + type, function (e) {
				if (e && e.result) {
					fn(e.result);
				} else {
					fn();
				}
			});
		}
		function one(type, fn) {
			var cb = function (e) {
				e.target.removeEventListener(cid + type, cb);
				if (e && e.result) {
					fn(e.result);
				} else {
					fn();
				}
			}
			document.addEventListener(cid + type, cb);
		}

		function trigger(type, data) {
			var t = cid + type;
			var event;

			if (document.dispatchEvent) {
				if (typeof CustomEvent === "function") {
					event = new CustomEvent(t, { "result": data })
				} else if (typeof Event === "function") {
					event = new Event(t, { "result": data })
				} else if (document.createEvent) {
					event = document.createEvent('HTMLEvents');
					event.initEvent(t, true, true);
				}
				if (event) {
					event.result = data;
					document.dispatchEvent(event);
					return;
				}
			}

			event = document.createEventObject();
			event.result = data;
			document.fireEvent('on' + t, event);
			return
		};

		(function connect() {
			function done(result) {
				try {
					result = JSON.parse(result);
				} catch (err) {
					return
				}

				if (result.requestID) {
					trigger("request:" + result.command + ":" + result.requestID, result.data)
				} else {
					trigger("read:" + result.command, result.data)
				}
			}

			sock = createWebSocket(url);
			sock.onopen = function () {
				trigger('wsConnect');
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
		}());

		// отправить сообщение в сокет и получить ответ в коллбек
		self.send = function (command, msg, requestID) {
			command = command.replace(/\:/g, '_');
			requestID = requestID || 0;
			var msg = JSON.stringify(msg);
			if (global.dev) {
				msg = [requestID, command, msg].join(":");
			} else {
				try {
					msg = new Blob([requestID, ":", command, ":", msg]);
				} catch (e) {
					if (e.name == "InvalidStateError") {
						msg = [requestID, command, msg].join(":");
						var bb = new MSBlobBuilder();
						bb.append(msg);
						msg = bb.getBlob('text/csv;charset=utf-8');
					}
				}
			}
			if (sock && sock.readyState === WebSocket.OPEN) {
				sock.send(msg);
			} else {
				one('wsConnect', function () {

					sock.send(msg);
				});
			}
		}


		// server messages handler
		self.read = function (command, callback) {
			on("read:" + command, callback)
		}

		// send request to server and wait answer to handler
		self.request = function (command, msg, callback, timeout) {
			var requestID = 0;
			command = command.replace(/\:/g, '_');

			if (callback) {
				requestID = ++prevID;
				timeout = timeout || requestTimeout;
				one("request:" + command + ":" + requestID, callback)
				if (timeout > 0) {
					setTimeout(function () {
						trigger("request:" + result.command + ":" + result.requestID)
					}, timeout);
				}
			}
			self.send(command, msg, requestID)
		}

		// set request timeout
		self.setRequestTimeout = function (timeout) {
			requestTimeout = timeout;
		}

		// повесить обработчик на сообщения, санкционированные сервером (без запроса)
		self.subscribe = function (command, callback) {
			self.send("subscribe", command);
			on('wsConnect', function () {
				self.send("subscribe", command);
			});
		}

		self.wait = function (commands, callback) {
			var cmds = {}
			commands.forEach(function (command) {
				if (!subscriptions[command]) {
					cmds[command] = true;
				}
			});

			function wt(e, data) {
				if (data && data.command && cmds[data.command]) {
					delete cmds[data.command];
				}
				if (Object.keys(cmds).length < 1) {
					setTimeout(callback, 1);
				} else {
					one('wsSubscibe', wt);
				}
			}
			wt();
		}

	};

	exports.getChannel = function (url) {
		if (typeof url !== 'string') {
			return null;
		}

		if (!channels[url]) {
			channels[url] = new Channel(url);
		}

		return channels[url];
	}

}));