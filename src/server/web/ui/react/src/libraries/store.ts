import {defaultEmitter} from "./bus";

export class store {
    private readonly storeContainer = new Map<string, any>();

    constructor() {
    }

    public commit(id: string, data: any) {
        this.storeContainer.set(id, data);
        defaultEmitter().emit('store', '*', 'add', id, data)
    }

    public get(id: string) : any {
        return this.storeContainer.get(id);
    }

    public has(id: string) : boolean {
        return this.storeContainer.has(id);
    }

    public delete(id: string) : boolean {
        const ret = this.storeContainer.delete(id);
        defaultEmitter().emit('store', '*', 'delete', id, null)
        return ret
    }
}

const _defaultStore = new store();

export function defaultStore() : store {
    return _defaultStore
}
