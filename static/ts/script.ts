import { SolidWebSocket } from '@/websocket.js';

const socket = new SolidWebSocket('/ws');

socket.on("open", function(_: Event) {
    console.log('WebSocket connection opened');

    socket.sendTextMessage('Hello, Server!');
});

socket.on('message', function(event: MessageEvent) {
    console.log(event);
});
