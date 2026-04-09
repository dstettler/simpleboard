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

export function positionToAlgebraic(pos: Position): string {
  return `${String.fromCharCode(97 + pos.y)}${8 - pos.x}`;
}
