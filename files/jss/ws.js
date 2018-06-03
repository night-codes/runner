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
		var requestID = 0;
		var readers = {}; // удалить
		var subscriptions = {}; // нужны для того, чтобы переотправлять после реконнекта
		var self = this;
		var cid = "" + (new Date().valueOf()) + url;


		function on(type, fn) {
			document.addEventListener(cid + type, function (e, result) {
				if (result && result.data) {
					fn(result.data);
				} else {
					fn();
				}
			});
		}
		function one(type, fn) {
			var cb = function (e, result) {
				e.target.removeEventListener(cid + type, cb);
				if (result && result.data) {
					fn(result.data);
				} else {
					fn();
				}
			}
			document.addEventListener(cid + type, cb);
		}
		function trigger(type, data) {
			var e = new Event(cid + type)
			if (typeof data !== "undefined") {
				var e = new CustomEvent(cid + type, { "data": data })
			}
			document.dispatchEvent(e);
		}



		(function () {
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
				} else if (readers[result.command]) {
					readers[result.command](result.data, result.command)
				}
				if (result && result.command) {
					subscriptions[result.command] = true;
					trigger('wsSubscibe', result.command);
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
		self.request = function (command, msg, callback) {
			var reqID = ++requestID;
			if (callback) {
				cbcs[reqID] = {
					fn: callback,
					time: (new Date()).getTime() / 1000
				};
			}
			self.send(command, msg, reqID)
		}

		// отправить сообщение в сокет и получить ответ в коллбек
		self.send = function (command, msg, reqID) {
			reqID = reqID || 0;
			var msg = JSON.stringify(msg);
			if (global.dev) {
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
				one('wsConnect', function () {
					sock.send(msg);
				});
			}
		}

		// повесить обработчик на сообщения, санкционированные сервером (без запроса)
		self.subscribe = function (command, callback) {
			readers[command] = callback;
			if (typeof autoReconnect === "undefined" || autoReconnect) {
				on('wsConnect', function () {
					self.request(command, msg, callback);
				});
			}
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