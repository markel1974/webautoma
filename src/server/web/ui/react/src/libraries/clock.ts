export function toClockDigit(value: number) : string {
    const v = value | 0
    const digit = v.toString()
    if (digit.length == 0) {
        return "00"
    } else if (digit.length == 1) {
        return "0" + digit
    } else if (digit.length == 2) {
        return digit
    } else {
        return digit.substring(0, 2)
    }
}

export function fromMillisToClockDigit(v: number) : string {
    const seconds = (v / 1000) % 60
    const minutes = (v / (1000 * 60)) % 60
    const hours = (v / (1000 *60 *60)) % 24
    return toClockDigit(hours) + ":" + toClockDigit(minutes) + ":" + toClockDigit(seconds)
}

export function fromSecondsToClockDigit(v: number) : string{
    const seconds = v % 60
    const minutes = (v / 60) % 60
    const hours = (v / (60 * 60)) % 24
    return toClockDigit(hours) + ":" + toClockDigit(minutes) + ":" + toClockDigit(seconds)
}