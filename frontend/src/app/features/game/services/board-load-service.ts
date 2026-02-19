import { inject, Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map, Observable } from 'rxjs';

import { Bishop, ChessPiece, King, Knight, Pawn, Queen, Rook } from './pieces/ChessPiece';
import { Position } from './pieces/Position';

interface BoardApiRequest {
  fenString: string
}

interface UpdatedRequest {
  arbitraryValue: boolean
}

@Injectable({
  providedIn: 'root',
})
export class BoardLoadService {
  private http = inject(HttpClient);
  
  positionsArray: ChessPiece[]|null = null;

  /**
   * @returns {Map<string, ChessPiece} Indexed map of pieces on board with key of "[Position.x],[Position.y]".
   */
  boardLoad(): Observable<ChessPiece[]> {
    // Returns an observable after sequentially decoding JSON string and filtering into the map via rxjs pipe.
    return this.http.get<BoardApiRequest>(`http://127.0.0.1:5000/api/mock-board`).pipe(
      map(state => this.fenDecode(state.fenString))
    );
  }

  updatePiecePosition(piece: ChessPiece, newPos: Position): Observable<ChessPiece[]> {
    return this.http.get<UpdatedRequest>(`http://127.0.0.1:5000/api/update-board`).pipe(
      map(_state => this.updatePos(piece, newPos))
    )
  }

  private updatePos(piece: ChessPiece, newPos: Position): ChessPiece[] {
    console.log(`${piece.id} new pos: ${newPos.x}, ${newPos.y}`)
    if (this.positionsArray !== null) {
      this.positionsArray[piece.id].position = newPos;
    } else {
      console.warn("Attempted to update pos before init")
      this.positionsArray = this.mockPositions();
    }

    return this.positionsArray;
  }

  private mockPositions(): ChessPiece[] {
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

  private fenDecode(fenString: string): ChessPiece[] {
    // TODO add proper decoding. https://github.com/dstettler/simpleboard/issues/35
    // Mock return.
    console.log(`fenString: ${fenString}`);
    this.positionsArray = this.mockPositions();
    return this.positionsArray
  }
}
