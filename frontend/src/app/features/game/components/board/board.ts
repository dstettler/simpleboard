import { Component, effect, inject } from '@angular/core';
import { AsyncPipe } from '@angular/common';
import { Observable } from 'rxjs';

import { BoardLoadService } from '../../services/board-load-service';
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

  private loadService = inject(BoardLoadService);
  boardState$: Observable<ChessPiece[]> = this.loadService.boardLoad();

  constructor() {}

  onPieceMoved(piece: ChessPiece, target: Position) {
    this.boardState$ = this.loadService.updatePiecePosition(piece, target)
  }

  isDarkSquare(i: number): boolean {
    const row = Math.floor(i / 8)
    const col = i % 8;
    return (row + col) % 2 !== 0
  }
}
