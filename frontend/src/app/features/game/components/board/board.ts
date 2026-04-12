import { Component, inject, Signal } from '@angular/core';
import { AsyncPipe } from '@angular/common';
import { toObservable } from '@angular/core/rxjs-interop';

import { BoardStateService } from '../../services/board-state-service';
import { Piece } from '../piece/piece';
import { ChessPiece } from '../../services/pieces/ChessPiece';
import { Position } from '../../services/pieces/Position';

@Component({
  selector: 'app-board',
  imports: [AsyncPipe, Piece],
  templateUrl: './board.html',
  styleUrl: './board.css',
})
export class Board {
  grid = Array.from({ length: 64 });

  private stateService = inject(BoardStateService);
  boardPieces: Signal<ChessPiece[]> = this.stateService.pieces;
  boardPiecesObservable$ = toObservable(this.boardPieces);

  constructor() {
    // TODO this needs to be replaced. Discuss @ #74 (https://github.com/dstettler/simpleboard/issues/74)
    this.stateService.boardLoad(1, 0, 'w').subscribe();
  }

  onPieceMoved(piece: ChessPiece, target: Position) {
    // TODO this needs to be replaced. Discuss @ #74 (https://github.com/dstettler/simpleboard/issues/74)
    this.stateService.updatePiecePosition(1, 0, piece, target).subscribe();
  }

  isDarkSquare(i: number): boolean {
    const row = Math.floor(i / 8)
    const col = i % 8;
    return (row + col) % 2 !== 0
  }
}
