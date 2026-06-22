//#region resource/lib/websocket.ts
function e(e) {
	return (t) => {
		let n = JSON.parse(t.data ?? "{}");
		e({
			target: t.target,
			type: n.type,
			requestId: n.requestId,
			data: n.data
		});
	};
}
function t(e) {
	return (t) => {
		e({ target: t.target });
	};
}
var n = class {
	datas = [];
	url;
	socket;
	readyState;
	constructor(e) {
		this.url = e, this.socket = new WebSocket(window.location.origin + this.url), this.readyState = this.socket.readyState;
	}
	send(e) {
		this.socket.readyState === WebSocket.OPEN ? this.socket.send(JSON.stringify(e)) : this.datas.push(JSON.stringify(e));
	}
	push() {
		for (let e of this.datas) this.socket.send(e);
		this.datas = [];
	}
	async reconnect() {
		let e = 0;
		for (; this.socket.readyState !== WebSocket.OPEN && e < 5;) console.warn(`WebSocket connection lost. Attempting to reconnect... (Attempt ${e + 1})`), this.socket = new WebSocket(window.location.origin + this.url), await new Promise((t) => setTimeout(t, 1e3 * ++e));
	}
	sendTextMessage(e) {
		let t = {
			type: "text",
			data: e
		};
		this.send(t);
	}
	sendBinaryMessage(e) {
		let t = {
			type: "binary",
			data: e
		};
		this.send(t);
	}
	sendJSONMessage(e) {
		let t = {
			type: "json",
			data: e
		};
		this.send(t);
	}
	close() {
		this.socket.close();
	}
	on(n, r) {
		switch (n) {
			case "open":
				this.socket.onopen = (e) => {
					this.readyState = this.socket.readyState, this.push(), t(r)(e);
				};
				break;
			case "message":
				this.socket.onmessage = e(r);
				break;
			case "close":
				this.socket.onclose = (e) => {
					this.readyState = this.socket.readyState, this.reconnect(), t(r)(e);
				};
				break;
			case "error":
				this.socket.onerror = (e) => {
					this.readyState = this.socket.readyState, t(r)(e);
				};
				break;
			default: console.warn(`Unsupported event type: ${n}`);
		}
	}
	getReadyState() {
		return this.readyState;
	}
};
//#endregion
export { n as SolidWebSocket };
