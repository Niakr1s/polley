export const isExpired = (date: string): boolean => {
    return strToDate(date).valueOf() < Date.now()
}

export const secondsRemained = (date: string): number => {
    return Math.round((strToDate(date).valueOf() - Date.now()) / 1000)
}

const strToDate = (dateStr: string): Date => {
    return new Date(dateStr)
}