import { newPosition, Position } from "./Position";

export abstract class ChessPiece {
    id: number;
    isWhite: boolean;
    position: Position;

    constructor(id: number, isWhite: boolean, posX: number, posY: number) {
        this.id = id;
        this.isWhite = isWhite;
        this.position = newPosition(posX, posY)
    }

    abstract getClass(): string;
    abstract getPossibleMoves(): Position[];

    getImageUrl(): string {
        let colorPrefix = this.isWhite ? "w" : "b";
        return `assets/${colorPrefix}_${this.getClass()}.png`
    }
}

export class Rook extends ChessPiece {
    override getClass(): string {
        return "Rook"
    }

    // TODO: Implement https://github.com/dstettler/simpleboard/issues/11
    override getPossibleMoves(): Position[] {
        return [];
    }
}

export class Knight extends ChessPiece {
    override getClass(): string {
        return "Knight"
    }

    // TODO: Implement https://github.com/dstettler/simpleboard/issues/11
    override getPossibleMoves(): Position[] {
        return [];
    }
}

export class Bishop extends ChessPiece {
    override getClass(): string {
        return "Bishop"
    }

    // TODO: Implement https://github.com/dstettler/simpleboard/issues/11
    override getPossibleMoves(): Position[] {
        return [];
    }
}

export class Pawn extends ChessPiece {
    override getClass(): string {
        return "Pawn"
    }

    // TODO: Implement https://github.com/dstettler/simpleboard/issues/11
    override getPossibleMoves(): Position[] {
        return [];
    }
}

export class King extends ChessPiece {
    override getClass(): string {
        return "King"
    }

    // TODO: Implement https://github.com/dstettler/simpleboard/issues/11
    override getPossibleMoves(): Position[] {
        return [];
    }
}

export class Queen extends ChessPiece {
    override getClass(): string {
        return "Queen"
    }

    // TODO: Implement https://github.com/dstettler/simpleboard/issues/11
    override getPossibleMoves(): Position[] {
        return [];
    }
}

export function getPieceFromFenCharacter(char: string, id: number, x: number, y: number): ChessPiece {
    const pieceConstructors = {
        P: Pawn,
        N: Knight,
        B: Bishop,
        R: Rook,
        Q: Queen,
        K: King
    }

    const pieceCtor = pieceConstructors[char.toUpperCase() as keyof typeof pieceConstructors];

    // This case should never be hit, since we validate the FEN string before calling this function.
    // Just in case something goes wrong return *something*.
    if (pieceCtor == undefined) {
        return new Pawn(-1,true,0,0);
    }

    console.log(`${char} X: ${x}, Y:${y}`);

    return new pieceCtor(id, char == char.toUpperCase(), x, y);
}
