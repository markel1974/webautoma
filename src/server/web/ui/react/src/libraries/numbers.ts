
export function roundOneDigit(n: number) : number {
    return Math.round((n + Number.EPSILON) * 10) / 10
}

export function roundTwoDigits(n: number) : number {
    return Math.round((n + Number.EPSILON) * 100) / 100
}

export function roundThreeDigits(n: number) : number {
    return Math.round((n + Number.EPSILON) * 1000) / 1000
}