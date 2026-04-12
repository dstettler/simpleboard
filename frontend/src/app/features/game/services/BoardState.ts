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
      new Knight(1, false, 0, 1),
      new Bishop(2, false, 0, 2),
      new Queen(3, false, 0, 3),
      new King(4, false, 0, 4),
      new Pawn(5, false, 1, 0),
      new Pawn(6, true, 6, 0),
      new Rook(7, true, 7, 0),
      new Knight(8, true, 7, 1),
      new Bishop(9, true, 7, 2),
      new Queen(10, true, 7, 3),
      new King(11, true, 7, 4),
    ]
  }
