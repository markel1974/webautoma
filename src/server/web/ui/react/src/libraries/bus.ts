

export interface busMessage {
    from: string
    to: string
    method: string
    type: string
    data: any
}

export class emitter {
    protected readonly id: string
    protected readonly channel: BroadcastChannel;

    constructor(id: string) {
        this.id = id;
        this.channel = new BroadcastChannel(id);
    }

    public getId() {
        return this.id;
    }

    public emit(from: string, to: string, method: string, type: string, data: any) {
        const msg : busMessage = { from: from, to: to, method: method, type: type, data: data };
        this.channel.postMessage(msg);
    }

    public disconnect() {
        this.channel.close()
    }
}

export class bus extends emitter {
    constructor(id: string) {
        super(id)
    }

    public bind(onmessage : (this: BroadcastChannel, ev: MessageEvent) => any) {
        this.channel.onmessage = onmessage
    }
}

function initialize() : bus {
    const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    const length = 16
    let result = ''
    while (result.length < length) {
        result += characters.charAt(Math.floor(Math.random() * characters.length));
    }
    const defaultBusId = 'default_channel_' + result
    return new bus(defaultBusId)
}

const _defaultBus = initialize()

export function defaultEmitter() : emitter {
    return _defaultBus
}