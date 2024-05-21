export function positive(x: number): number {
    return x >= 0 ? x : 0
}

export function asEuro(x: number): string {
    return x.toFixed(2) + ' €'
}

export function roundUpTo(x: number, divisor: number): number {
    return Math.ceil(x / divisor) * divisor
}

export function int(x: number): string {
    return x.toFixed(0)
}

export function round2(x: number): string {
    let y = x.toFixed(2)
    if (y.endsWith('.00')) {
        return y.substring(0, y.length - 3)
    }
    if (y.endsWith('0')) {
        return y.substring(0, y.length - 1)
    }
    return y
}

export function makePlural(x: number, singular: string, plural: string): string {
    return x === 1 ? singular : plural
}