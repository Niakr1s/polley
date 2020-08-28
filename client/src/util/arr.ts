export function count<T>(arr: T[], etalon: T): number {
    return arr.reduce((acc, v) => {
        if (v === etalon) acc++;
        return acc
    }, 0)
}