import { Component, effect, inject } from '@angular/core';
import { AsyncPipe } from '@angular/common';

import { BoardLoadService } from '../../services/board-load-service';
import { Observable } from 'rxjs';
import { Piece } from '../piece/piece';
import { ChessPiece } from '../../services/pieces/ChessPiece';
import { CdkDrag } from '@angular/cdk/drag-drop';

@Component({
  selector: 'app-board',
  imports: [AsyncPipe, Piece, CdkDrag],
  templateUrl: './board.html',
  styleUrl: './board.css',
})
export class Board {
  rows = Array.from({ length: 8 });
  cols = Array.from({ length: 8 });

  private loadService = inject(BoardLoadService);
  boardState$: Observable<Map<string, ChessPiece>> = this.loadService.boardLoad();

  constructor() {}

  getItem(itemMap: Map<string, ChessPiece>, row: number, col: number): ChessPiece | null {
    return itemMap.get(`${row},${col}`) ?? null;
  }
}
