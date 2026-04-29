import { Component, inject, Input, Signal } from '@angular/core';
import { AsyncPipe } from '@angular/common';
import { toObservable } from '@angular/core/rxjs-interop';

import { BoardStateService, PlayerColor } from '../../services/board-state-service';
import { Piece } from '../piece/piece';
import { ChessPiece } from '../../services/pieces/ChessPiece';
import { Position } from '../../services/pieces/Position';

import { AuthStateService } from '../../../../core/services/auth-state.service';
@Component({
  selector: 'app-board',
  imports: [AsyncPipe, Piece],
  templateUrl: './board.html',
  styleUrl: './board.css',
})
export class Board {
  grid = Array.from({ length: 64 });

  @Input() gameId!: string;

  private stateService = inject(BoardStateService);
  private authService = inject(AuthStateService);

  boardPieces: Signal<ChessPiece[]> = this.stateService.pieces;
  side: Signal<PlayerColor> = this.stateService.userColor;
  boardPiecesObservable$ = toObservable(this.boardPieces);
  sideObservable$ = toObservable(this.side);

  ngOnInit() {
    const userId = Number(this.authService.userId());
    this.stateService.boardLoad(this.gameId, userId).subscribe();
  }

  onPieceMoved(piece: ChessPiece, target: Position) {
    const userId = Number(this.authService.userId());
    this.stateService.updatePiecePosition(this.gameId, userId, piece, target).subscribe();
  }

  isDarkSquare(i: number): boolean {
    const row = Math.floor(i / 8)
    const col = i % 8;
    return (row + col) % 2 !== 0
  }
}
