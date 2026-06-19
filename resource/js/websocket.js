function handleEvent(callback) {
    return (event) => {
        const data = JSON.parse(event.data ?? "{}");

        callback({
            target: event.target,
            type: data.type,
            requestId: data.requestId,
            data: data.data,
        });
    } 
}

class SolidWebSocket {
    datas = []

    constructor(url) {
        this.url = url;
        this.socket = new WebSocket(window.location.origin + this.url);
        this.readyState = this.socket.readyState;
    }

    send(data) {
        if (this.socket.readyState !== WebSocket.OPEN) {
            this.datas.push(JSON.stringify(data));
        } else {
            this.socket.send(JSON.stringify(data));
        }
    }

    sendTextMessage(data) {
        data = {
            type: "text",
            data: data
        };

        this.send(data);
    }

    sendBinaryMessage(data) {
        data = {
            type: "binary",
            data: data
        };

        this.send(data);
    }

    sendJSONMessage(data) {
        data = {
            type: "json",
            data: data
        };

        this.send(data);
    }

    close() {
        this.socket.close();
    }

    push() {
        for (const data of this.datas) {
            this.socket.send(data);
        }

        this.datas = []
    }

    async reconnect() {
        let n = 0;
        while (this.socket.readyState !== WebSocket.OPEN && n < 5) {
            console.warn(`WebSocket connection lost. Attempting to reconnect... (Attempt ${n + 1})`);
            this.socket = new WebSocket(window.location.origin + this.url);
            await new Promise((resolve) => setTimeout(resolve, 1000 * ++n));
        }
    }

    on(event, callback) {
        switch (event) {
            case 'open':
                this.socket.onopen = (event) => {
                    this.readyState = this.socket.readyState;

                    this.push();

                    handleEvent(callback)(event);
                };
                break;

            case 'message':
                this.socket.onmessage = handleEvent(callback);
                break;

            case 'close':
                this.socket.onclose = (event) => {
                    this.readyState = this.socket.readyState;

                    this.reconnect();

                    handleEvent(callback)(event);
                };
                break;

            case 'error':
                this.socket.onerror = (event) => {
                    this.readyState = this.socket.readyState;
                    handleEvent(callback)(event);
                };
                break;

            default:
                console.warn(`Unsupported event type: ${event}`);
        }
    }
}

export { SolidWebSocket };