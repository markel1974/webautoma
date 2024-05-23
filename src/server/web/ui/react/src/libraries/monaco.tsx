import * as monacoEditor from "monaco-editor";
import {editor, languages, Position} from "monaco-editor";



function templateString(line: string, pos: number) : string | null {
    const target = 2
    let counter = 0
    if (pos < 2) {
        return null
    }
    const end = pos - 2
    for (let x = end; x >= 0; x--) {
        const curr = line[x]
        switch (curr) {
            case "{":
                counter ++;
                if (counter == target) {
                    const start = x + counter
                    return line.substring(start, end)
                }
                break
            case "}":
                return null
            default:
                counter = 0
        }
    }
    return null
}

export function ConfigureTemplate(monaco: typeof monacoEditor, docs: any) {
    if (!monaco) {
        return
    }
    const helper = new Map<string, any>()
    const helperLower = new Map<string, any>()

    docs.forEach((v: any) => {
        helper.set(v.name, v)
        helperLower.set(v.name.toLowerCase(), v)
    })

    const createDependencyProposals = (model: editor.ITextModel, position: Position) => {
        const proposals: Array<languages.CompletionItem> = [];
        const lc = templateString(model.getLineContent(position.lineNumber), position.column)
        if (lc === null) {
            return proposals
        }

        const wup = model.getWordUntilPosition(position)
        const range = {startLineNumber: position.lineNumber, endLineNumber: position.lineNumber, startColumn: wup.startColumn, endColumn: wup.endColumn};
        const w = wup.word.toLowerCase()
        helperLower.forEach((v: any, k) => {
            if (k.includes(w)) {
                proposals.push({
                    label: v.name,
                    documentation: v.documentation,
                    detail: v.documentation,
                    kind: languages.CompletionItemKind.Function,
                    insertText: v.name,
                    range: range});
            }
        });
        return proposals;
    };

    monaco.languages.registerCompletionItemProvider('yaml', {
        provideCompletionItems: (model: editor.ITextModel, position: Position) => {
            return { suggestions: createDependencyProposals(model, position) };
        }
    });

    monaco.languages.registerSignatureHelpProvider('yaml', {
        signatureHelpTriggerCharacters:[' '],
        provideSignatureHelp: (model: editor.ITextModel, position: Position, a, shc) => {
            const lc = templateString(model.getLineContent(position.lineNumber), position.column)
            if (lc === null) {
                return null
            }
            let signature : any
            let target = ""
            let start = lc.length - 1

            for (; start >= 0; start--) {
                if (lc[start] === ' ') {
                    if (target.length > 0) {
                        signature = helper.get(target)
                        if (signature) {
                            start++
                            break
                        }
                    }
                    target = ""
                } else {
                    target = lc[start] + target
                }
            }

            if (!signature && target.length > 0) {
                signature = helper.get(target)
            }

            if (!signature) { return null }

            if (start < 0) { start = 0 }

            let quoter = false
            let sentence = lc.slice(start)//model.getLineContent(position.lineNumber).slice(start, position.column)
            sentence = sentence.replace(/\s+/g, ' ')
            if (sentence.includes('"')) {
                sentence = sentence.replace(/"/g, "'")
                quoter = true
            } else if (sentence.includes("'")) {
                quoter = true
            }
            if (quoter) {
                sentence += "'"
                sentence = sentence.replace(/'.+?'/g, 'QUOTER')
                sentence = sentence.slice(0, -1)
            }

            const targets = sentence.split(' ')
            const pos = targets.length >= 2 ? targets.length -2 : 0

            return {
                value: {
                    signatures: [ signature ],
                    activeSignature: 0,
                    activeParameter: pos
                },
                dispose: () => {
                }
            }
        }
    });
}