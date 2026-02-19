import { Component, Input } from '@angular/core';
import { ChessPiece, Pawn } from '../../services/pieces/ChessPiece';

@Component({
  selector: 'app-piece',
  imports: [],
  templateUrl: './piece.html',
  styleUrl: './piece.css',
})
export class Piece {
  @Input()
  // The default value is set to a Pawn, but this field should be set by a parent.
  public piece: ChessPiece = new Pawn(false, 0, 0);

  constructor() {
  }
}
