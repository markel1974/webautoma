export class cookies {
    private defaultExpire = 7 * 24 * 60 * 60 * 1000;
    private readonly expire : number;

    constructor(expire?: number) {
        this.expire = expire && expire > 0 ? expire : this.defaultExpire;
    }

    public set(name: string, val: string) {
        const date = new Date();
        const value = val;
        date.setTime(date.getTime() + this.expire);
        document.cookie = name + "=" + value + "; expires=" + date.toUTCString() + "; path=/";
    }

    public get(name: string) : string {
        const value = "; " + document.cookie;
        const parts = value.split("; " + name + "=");
        if (parts.length == 2) {
            const part = parts.pop();
            if (part) {
                const res = part.split(";").shift();
                if (res) {
                    return res
                }
            }
        }
        return ""
    }

    public delete(name: string) {
        const date = new Date();
        date.setTime(date.getTime() + (-1 * 24 * 60 * 60 * 1000));
        document.cookie = name + "=; expires=" + date.toUTCString() + "; path=/";
    }
}

const _defaultCookies = new cookies();

export function defaultCookies() : cookies {
    return _defaultCookies
}