import { Component, EventEmitter, Input, Output } from '@angular/core';
import { CdkDrag, CdkDragEnd } from '@angular/cdk/drag-drop';

import { ChessPiece } from '../../services/pieces/ChessPiece';
import { newPosition, Position } from '../../services/pieces/Position';

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

  dragPosition = { x: 0, y: 0 };

  onDragEnded(event: CdkDragEnd) {
    const boardSize = this.boardRef.getBoundingClientRect();
    const squareSize = boardSize.width / 8;

    const newRow = Math.round(this.piece.position.x + event.distance.y / squareSize);
    const newCol = Math.round(this.piece.position.y + event.distance.x / squareSize);

    const clampedRow = Math.max(0, Math.min(7, newRow));
    const clampedCol = Math.max(0, Math.min(7, newCol));

    this.moved.emit(newPosition(clampedRow, clampedCol));

    // This is necessary to reset cdkDrag's internal offset each time state is updated, 
    // since the position is determined via style in the board component.
    event.source._dragRef.reset();
  }
}
