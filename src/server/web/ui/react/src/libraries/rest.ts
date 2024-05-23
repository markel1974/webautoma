import { request } from './request'
import { defaultEmitter } from "./bus";
import {Toaster} from "../components/Toaster";

export class rest {
    private readonly prefix: string;
    private readonly defaultFailedHandler: any;
    private readonly defaultExceptionHandler: any;
    private readonly specialStatusHandles: any;
    private readonly mh: any;

    constructor(prefix?: string, defaultFailedHandler?: any, defaultExceptionHandler?: any, specialStatusHandles?: any){
        this.prefix = prefix ? prefix : '/v1/';
        this.defaultFailedHandler = defaultFailedHandler;
        this.defaultExceptionHandler = defaultExceptionHandler;
        this.specialStatusHandles = specialStatusHandles;
        this.mh = {}
    }

    public handleSpecialStatus(code: number, h: any) {
        this.mh[code] = h
    }

    public GET(url: string){
        return new request(this.prefix + url, 'GET', null, this.specialStatusHandles)
            .onFailed(this.defaultFailedHandler)
            .onException(this.defaultExceptionHandler);
    }

    public POST(url: string, data: any) {
        return new request(this.prefix + url, 'POST', data, this.specialStatusHandles)
            .onFailed(this.defaultFailedHandler)
            .onException(this.defaultExceptionHandler);
    }

    public PUT(url: string, data: any) {
        return new request(this.prefix + url, 'PUT', data, this.specialStatusHandles)
            .onFailed(this.defaultFailedHandler)
            .onException(this.defaultExceptionHandler);
    };

    public DELETE(url: string) {
        return new request(this.prefix + url, 'DELETE',null, this.specialStatusHandles)
            .onFailed(this.defaultFailedHandler)
            .onException(this.defaultExceptionHandler);
    }
}

const _defaultRest = new rest('/v1/', (msg: any) => {
    const error = msg ? msg.toString() : 'failed';
    Toaster.add('danger', 'globe', 'Error', error)
}, (msg: any) => {
    const error = msg ? msg.toString() : 'exception';
    Toaster.add('danger', 'globe', 'Error', error)
}, {
    401: (data: any, xhr: any) => {
        defaultEmitter().emit('authentication','*','authentication','error', 'unauthorized');
    }
});

export function defaultRest() : rest {
    return _defaultRest
}