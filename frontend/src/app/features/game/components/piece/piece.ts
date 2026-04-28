import { Component, EventEmitter, Input, Output } from '@angular/core';
import { CdkDrag, CdkDragEnd, CdkDragStart } from '@angular/cdk/drag-drop';

import { ChessPiece } from '../../services/pieces/ChessPiece';
import { newPosition, Position, positionsEqual } from '../../services/pieces/Position';

@Component({
  selector: 'app-piece',
  imports: [CdkDrag],
  templateUrl: './piece.html',
  styleUrl: './piece.css',
})
export class Piece {
  @Input() piece!: ChessPiece;
  @Input() boardRef!: HTMLElement;

  @Output()
  moved = new EventEmitter<Position>();
  pieceSelected = new EventEmitter<boolean>();

  dragPosition = { x: 0, y: 0 };

  onDragStart(_event: CdkDragStart) {
    console.log('started drag');
    this.pieceSelected.emit(true);
  }

  onDragEnded(event: CdkDragEnd) {
    const boardSize = this.boardRef.getBoundingClientRect();
    const squareSize = boardSize.width / 8;

    const newRow = Math.round(this.piece.position.x + event.distance.x / squareSize);
    const newCol = Math.round(this.piece.position.y + event.distance.y / squareSize);

    const clampedRow = Math.max(0, Math.min(7, newRow));
    const clampedCol = Math.max(0, Math.min(7, newCol));
    const newPos = newPosition(clampedRow, clampedCol);

    if (!positionsEqual(newPos, this.piece.position)) {
      this.moved.emit(newPos);
    }

    // This is necessary to reset cdkDrag's internal offset each time state is updated,
    // since the position is determined via style in the board component.
    event.source._dragRef.reset();
  }
}
