export interface Position {
    x: number,
    y: number
}

export function newPosition(x: number, y: number): Position {
    return {
        x: x,
        y: y
    }
}