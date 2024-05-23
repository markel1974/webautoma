import {defaultRest} from "./rest";
import {htmlIdGenerator} from "@elastic/eui";

export function fileReadLines(file: Blob, validator: (d: string, x: number) => number, callback:(preview: string, count: number, header: boolean, valid: number, invalid : number) => void) {
    const previewMax = 1000
    const reader = new FileReader();
    reader.addEventListener('load', (event) => {
        if (event.target && event.target.result) {
            let buf = event.target.result.toString()
            const lines = buf.split("\n")
            let count = 0
            let valid = 0
            let invalid = 0
            let header = false
            lines.forEach((v: string, x: number) => {
                v = v.replace(/(\r)/gm, "");
                if (v.length > 0) {
                    count++
                    switch (validator(v, x)) {
                        case 1: valid++; break
                        case -1: invalid++; break
                        case 2: header = true; break
                    }
                }
            })
            if (buf.length > previewMax) {
                buf = buf.substring(0, previewMax) + " ...."
            }
            callback(buf, count, header, valid, invalid)
        }
    });
    /*
    reader.addEventListener('progress', (event) => {
        if (event.loaded && event.total) {
            const percent = (event.loaded / event.total) * 100;
            //this.setState({ progressValue: percent})
        }
    });
    */
    reader.readAsText(file,'utf8')
}


export function fileUpload(endpoint: string, qs: string, f: File | null, callback: (end: boolean, err: string, progress: number)=>void) {
    if (f == null) {
        callback(true, "file is null", 100)
        return;
    }
    const formData = new FormData();
    formData.append('file', f);
    const req = defaultRest().POST('endpoint/upload/' + endpoint + qs, formData);
    req.setDataType('binary');
    req.onProgress((percentage: any) => {
        const v = typeof percentage === 'number' ? percentage : 0;
        callback(false, "", v)
    });
    req.onFailed(() => {
        callback(true, "error sending file", 100)
    })
    req.onSucceed(200, (resp: unknown) => {
        setTimeout(()=> {
            callback(true, "", 100)
        }, 1000)
    });
    req.do();
}

export function fileDownload(filename: string, body: string) {
    const element = document.createElement('a');
    element.setAttribute('href', 'data:text/plain;charset=utf-8,' + encodeURIComponent(body));
    element.setAttribute('download', filename);
    element.style.display = 'none';
    document.body.appendChild(element);
    element.click();
    document.body.removeChild(element);
}