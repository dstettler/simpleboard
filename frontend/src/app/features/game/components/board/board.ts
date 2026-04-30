import { Component, inject, Input, signal, Signal } from '@angular/core';
import { AsyncPipe } from '@angular/common';
import { toObservable } from '@angular/core/rxjs-interop';

import { BoardStateService, PlayerColor } from '../../services/board-state-service';
import { Piece } from '../piece/piece';
import { ChessPiece } from '../../services/pieces/ChessPiece';
import { Position, positionsEqual, positionToAlgebraic } from '../../services/pieces/Position';

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

  // Readonly signals
  private boardPieces: Signal<ChessPiece[]> = this.stateService.pieces;
  private side: Signal<PlayerColor> = this.stateService.userColor;
  private boardTimer: Signal<number> = this.stateService.timerRemainingMs;

  // Writeable signals
  private targetable = signal<Position[]|null>(null);

  // Derived observables
  readonly boardPiecesObservable$ = toObservable(this.boardPieces);
  readonly sideObservable$ = toObservable(this.side);
  readonly boardTimerObservable$ = toObservable(this.boardTimer);
  readonly targetableObservable$ = toObservable(this.targetable);

  ngOnInit() {
    const userId = Number(this.authService.userId());
    this.stateService.boardLoad(this.gameId, userId).subscribe();
  }

  onPieceSelected(piece: ChessPiece) {
    const targetables = this.stateService.getTargetables(piece.id);
    this.targetable.update(_t => targetables);
  }

  onPieceMoved(piece: ChessPiece, target: Position) {
    const userId = Number(this.authService.userId());
    this.stateService.updatePiecePosition(this.gameId, userId, piece, target).subscribe();
  }

  isDarkSquare(i: number): boolean {
    const row = Math.floor(i / 8);
    const col = i % 8;
    return (row + col) % 2 !== 0;
  }

  isTargetable(i: number, targetables: Position[]): boolean {
    const row = Math.floor(i / 8);
    const col = i % 8;
    const squarePos: Position = {x: col, y: row};
    for (const targetablePos of targetables) {
      if (positionsEqual(targetablePos, squarePos)) {
        return true;
      }
    }

    return false;
  }

  getTimerText(i: number): string {
    const wholeNum = Math.round(i / 1000);
    return `Your Time: ${wholeNum}s`
  }
}
