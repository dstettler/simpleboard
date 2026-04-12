import { ChessPiece, Rook, Knight, Bishop, Queen, King, Pawn } from './pieces/ChessPiece';
import { Position } from './pieces/Position';

const gameStatusMembers = ["InProgress", "Waiting", "Error"] as const;
export type GameStatus = typeof gameStatusMembers[number];

export interface BoardState {
    pieces: ChessPiece[];
    isWhiteMove: boolean;
    castleables: string;
    enPassant: Position|null;
    halfmoveClock: number;
    fullmoveNum: number;
    userColor: string;
    nextMoves: string[];
    gameStatus: GameStatus;
}

export function parseGameStatus(statusStr: string): GameStatus {
  const found = gameStatusMembers.find((matched) => matched === statusStr);
  if (found) {
    return found;
  }

  return "Error";
}

export function emptyState(): BoardState {
    return {
      pieces: [],
      isWhiteMove: false,
      castleables: '',
      enPassant: null,
      halfmoveClock: -1,
      fullmoveNum: -1,
      userColor: '',
      nextMoves: [],
      gameStatus: "Waiting"
    };
}

export function mockPositions(): ChessPiece[] {
    return [
      new Rook(0, false, 0, 0),
      new Knight(1, false, 1, 0),
      new Bishop(2, false, 2, 0),
      new Queen(3, false, 3, 0),
      new King(4, false, 4, 0),
      new Pawn(5, false, 0, 1),
      new Pawn(6, true, 0, 6),
      new Rook(7, true, 0, 7),
      new Knight(8, true, 1, 7),
      new Bishop(9, true, 2, 7),
      new Queen(10, true, 3, 7),
      new King(11, true, 4, 7),
    ]
  }
