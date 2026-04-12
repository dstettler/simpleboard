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

export function algebraicToPosition(algebraicStr: string): Position {
  const x = algebraicStr.toLowerCase().charCodeAt(0) - 97;
  const y = Number(algebraicStr[1]);
  console.log(`${algebraicStr}: ${x}, ${y}`);
  return {x: x, y: y};
}

export function positionsEqual(pos1: Position, pos2: Position): boolean {
  const ret = pos1.x == pos2.x && pos1.y == pos2.y;
  console.log(`(${pos1.x}, ${pos2.x}), (${pos1.y}, ${pos2.y}), ${ret}`);
  return ret;
}
