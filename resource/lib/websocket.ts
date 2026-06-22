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
    url: string
    socket: WebSocket
    private readyState: 0 | 1 | 2 | 3

    constructor(url: string) {
        this.url = url;
        this.socket = new WebSocket(window.location.origin + this.url);
        this.readyState = this.socket.readyState;
    }

    private send(data: object) {
        if (this.socket.readyState !== WebSocket.OPEN) {
            this.datas.push(JSON.stringify(data));
        } else {
            this.socket.send(JSON.stringify(data));
        }
    }

    private push() {
        for (const data of this.datas) {
            this.socket.send(data);
        }

        this.datas = []
    }

    private async reconnect() {
        let n = 0;
        while (this.socket.readyState !== WebSocket.OPEN && n < 5) {
            console.warn(`WebSocket connection lost. Attempting to reconnect... (Attempt ${n + 1})`);
            this.socket = new WebSocket(window.location.origin + this.url);
            await new Promise((resolve) => setTimeout(resolve, 1000 * ++n));
        }
    }

    public sendTextMessage(data: string) {
        const message = {
            type: "text",
            data: data
        };

        this.send(message);
    }

    public sendBinaryMessage(data: any) {
        const message = {
            type: "binary",
            data: data
        };

        this.send(message);
    }

    public sendJSONMessage(data: string) {
        const message = {
            type: "json",
            data: data
        };

        this.send(message);
    }

    public close() {
        this.socket.close();
    }

    public on(event: string, callback: Function) {
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

    public getReadyState(): 0 | 1 | 2 | 3 {
        return this.readyState
    }
}

export { SolidWebSocket };