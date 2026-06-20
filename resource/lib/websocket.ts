function handleMessageEvent(callback: Function): (event: MessageEvent) => void {
    return (event: MessageEvent) => {
        const data = JSON.parse(event.data ?? "{}");

        callback({
            target: event.target,
            type: data.type,
            requestId: data.requestId,
            data: data.data,
        });
    } 
}

function handleEvent(callback: Function): Function {
    return (event: Event) => {
        callback({
            target: event.target,
        });
    } 
}

class SolidWebSocket {
    datas: string[] = []
    url
    socket
    readyState

    constructor(url: string) {
        this.url = url;
        this.socket = new WebSocket(window.location.origin + this.url);
        this.readyState = this.socket.readyState;
    }

    send(data: object) {
        if (this.socket.readyState !== WebSocket.OPEN) {
            this.datas.push(JSON.stringify(data));
        } else {
            this.socket.send(JSON.stringify(data));
        }
    }

    sendTextMessage(data: string) {
        const message = {
            type: "text",
            data: data
        };

        this.send(message);
    }

    sendBinaryMessage(data: any) {
        const message = {
            type: "binary",
            data: data
        };

        this.send(message);
    }

    sendJSONMessage(data: string) {
        const message = {
            type: "json",
            data: data
        };

        this.send(message);
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

    on(event: string, callback: Function) {
        switch (event) {
            case 'open':
                this.socket.onopen = (event) => {
                    this.readyState = this.socket.readyState;

                    this.push();

                    handleEvent(callback)(event);
                };
                break;

            case 'message':
                this.socket.onmessage = handleMessageEvent(callback);
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