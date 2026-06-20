import { SolidWebSocket as e } from "/lib/websocket.js";
//#endregion
//#region static/ts/script.ts
var t = /* @__PURE__ */ ((e, t) => () => (t || (e((t = { exports: {} }).exports, t), e = null), t.exports))((() => {
	var t = new e("/ws");
	t.on("open", function(e) {
		console.log("WebSocket connection opened"), t.sendTextMessage("Hello, Server!");
	}), t.on("message", function(e) {
		console.log(e);
	});
}));
//#endregion
export default t();
