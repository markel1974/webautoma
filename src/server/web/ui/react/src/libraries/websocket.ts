import {busMessage, defaultEmitter} from "./bus";

export class websocket {
    private ws: WebSocket | null = null;
    private readonly url: string;
    private readonly subProtocol: string | undefined;
    private readonly uesAutoReconnect : boolean;
    private autoReconnect : boolean;
    private count : number;
    private open: boolean;

    constructor(url: string, subProtocol: string | undefined, reconnect: boolean) {
        //this.connectionId = "";
        this.url = url;
        this.subProtocol = subProtocol;
        this.uesAutoReconnect = reconnect
        this.autoReconnect = this.uesAutoReconnect;
        this.open = false;
        this.count = 0;
    }

    public connect() {
        this.autoReconnect = this.uesAutoReconnect;
        this.count = 0;
        console.log('Websocket Connection Status: connecting to ' + this.url + ' [auto reconnect: ' + this.autoReconnect + ']')
        if (this.ws === null) {
            this.ws = new WebSocket(this.url, this.subProtocol);
            this.ws.onopen = (e) => { this.onOpen(e) };
            this.ws.onmessage = (e) => { this.onMessage(e) };
            this.ws.onerror = (e) => { this.onError(e) };
            this.ws.onclose = (e) => { this.onClose(e) };
        } else {
            console.log('Websocket Connection Status: already connected');
        }
    }

    private onOpen(event: Event): void {
        this.open = true;
        console.log("Websocket Connection Status: opened");
        console.log(event);
    }

    private onMessage(event: MessageEvent): void {
        this.count++;
        const data = (typeof event.data == "string") ? JSON.parse(event.data) : event.data
        defaultEmitter().emit("websocket", "*", "websocket","websocket", data)
    }

    private onError(event: Event): void {
        console.log('Websocket Connection Status: Error')
        console.log(event);
        if (this.ws) {
            this.ws.close();
        }
    }

    private onClose(event: CloseEvent): void {
        this.open = false;
        this.ws = null;
        console.log("Websocket Connection Status: disconnected:");
        console.log(event);
        console.log('total messages: ' + this.count)
        this.reconnect(3000);
    }

    private reconnect(ms: number) {
        if (this.autoReconnect) {
            console.log("Websocket Auto Reconnect: reconnecting.......")
            setTimeout(() => { this.connect() }, ms);
        } else {
            console.log("Websocket Auto Reconnect: off")
        }
    }

    public close() {
        this.autoReconnect = false
        if (this.ws) {
            this.ws.close();
        }
    }

    public send(data: any) {
        if (this.ws) {
            this.ws.send(JSON.stringify(data));
        }
    }
}

