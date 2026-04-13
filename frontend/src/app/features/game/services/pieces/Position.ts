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
  return `${String.fromCharCode(97 + pos.x)}${8 - pos.y}`;
}

export function algebraicToPosition(algebraicStr: string): Position {
  const x = algebraicStr.toLowerCase().charCodeAt(0) - 97;
  const y = 8 - Number(algebraicStr[1]);
  return {x: x, y: y};
}

export function positionsEqual(pos1: Position, pos2: Position): boolean {
  return pos1.x == pos2.x && pos1.y == pos2.y;
}
