
/*
export interface iResp {
    __id: string
    __index: string
    __score: string
    __type: string
    nm_cleaned_address: string
}
*/

export interface IData {
    key: string,
    val: string
}

export interface IResult {
    label: string
    title: string
    id: string
    data: Array<IData>
}

function parseResultEntry(id: string, label: string, container: Array<IResult>, obj: any) {
    if (typeof obj === 'object') {
        const data: Array<IData> = [];
        let targetId = '';
        let targetLabel = '';
        for (const [k, v] of Object.entries(obj)) {
            const val = '' + v + '';
            data.push({key: k, val: val});
            if (k === label) {
                targetLabel = val
            } else if (k === id) {
                targetId = val
            }
        }
        container.push({
            label: targetLabel,
            id: targetId,
            title: targetLabel,
            data: data
        });
    }
}

export function parseResult(id: string, label: string, r: any) {
    let results: Array<IResult> = [];
    if (typeof r === 'string') {
        const resp = r ? r.split('\n') : [];
        resp.forEach((d: string) => {
            try {
                if (d) {
                    parseResultEntry(id, label, results, JSON.parse(d))
                }
            } catch (e) {
            }
        });
    } else if (typeof r === 'object') {
        parseResultEntry(id, label, results, r)
    }
    return results;
}