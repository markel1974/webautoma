
export class request {
    private readonly _url : string;
    private readonly _method : string;
    private readonly _data : string;
    private readonly _specialHandlers: any;
    private readonly _headers: any;

    private _dataType: string;
    private _successCode: number;
    private _onSucceed: any;
    private _onFailed: any;
    private _onException: any;
    private _onEnd: any;
    private _onProgress: any;
    private _body: string;

    constructor(url: string, method: string, data: any, specialHandlers: any) {
        this._url = url;
        this._method = method;
        this._data = data;
        this._specialHandlers = specialHandlers;
        this._headers = {};
        this._dataType = 'json';
        this._successCode = -1;
        this._body = "";
    }

    setHeader(key: string, val: any) {
        this._headers[key] = val
    }

    setDataType(type: string) {
        this._dataType = type
    }

    do(){
        const self = this;
        const xhr = new XMLHttpRequest();
        xhr.open(this._method, this._url, true);

        if (self._headers) {
            for (const id in self._headers) {
                if (self._headers.hasOwnProperty(id)) {
                    xhr.setRequestHeader(id, self._headers[id])
                }
            }
        }

        if (typeof self._onException === 'function') {
            const warpExceptionHandler = (msg: any)=>{
                self._onException(msg);
                typeof self._onEnd === 'function' && self._onEnd(xhr);
            };
            xhr.onabort=()=>{ warpExceptionHandler('request aborted.') };
            xhr.onerror=()=>{ warpExceptionHandler('request error.') };
            xhr.ontimeout=()=>{ warpExceptionHandler('request timeout.') };
        }

        xhr.upload.onprogress = (event) => {
            const progress = Math.round((100 * event.loaded) / event.total);
            //console.log("loaded: " + event.loaded);
            //console.log("total: " + event.total);
            //console.log("progress: " + progress);
            typeof self._onProgress === 'function' && self._onProgress(progress);
        };

        xhr.onreadystatechange = function(){
            if (xhr.readyState !== XMLHttpRequest.DONE) {
                return;
            }
            let data;
            if (typeof xhr.response != 'object') {
                try {
                    data = JSON.parse(xhr.response)
                } catch(e) {
                    data = xhr.response;
                }
            } else {
                data = xhr.response;
            }
            if (xhr.status !== self._successCode) {
                typeof self._onFailed === 'function' && self._onFailed(data);
            } else if (xhr.status === self._successCode && typeof self._onSucceed === 'function') {
                self._onSucceed(data, xhr);
            } else if (self._specialHandlers && typeof self._specialHandlers[xhr.status] === 'function') {
                self._specialHandlers[xhr.status](data, xhr);
            }
            typeof self._onEnd === 'function' && self._onEnd(data);
        };
        if (self._dataType === 'json' && typeof self._data === 'object' && self._data) {
            self._body = JSON.stringify(self._data);
        } else {
            self._body = self._data;
        }
        xhr.send(self._body);
    }

    public exec(successCode: number) {
        this._successCode = successCode;
        return new Promise((resolve, reject) => {
            this._onSucceed = (resp: any) => {
                resolve(resp)
            };
            this._onFailed = (resp: any) => {
                reject(resp)
            };
            this._onException = (resp: any) => {
                reject(resp)
            };
            this.do();
        })
    }

    public onSucceed(successCode: number, f: any) {
        this._successCode = successCode;
        this._onSucceed = f;
        return this;
    }

    public onFailed(f: any){
        this._onFailed = f;
        return this;
    }

    public onException(f: any){
        this._onException = f;
        return this;
    }

    public onEnd(f: any){
        this._onEnd = f;
        return this;
    }

    public onProgress(f: any){
        this._onProgress = f;
        return this;
    }
}

