import { inject, Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map, Observable } from 'rxjs';
import { Bishop, ChessPiece, King, Knight, Pawn, Queen, Rook } from './pieces/ChessPiece';

interface BoardApiRequest {
  fenString: string
}

@Injectable({
  providedIn: 'root',
})
export class BoardLoadService {
  private http = inject(HttpClient);

  /**
   * @returns {Map<string, ChessPiece} Indexed map of pieces on board with key of "[Position.x],[Position.y]".
   */
  boardLoad(): Observable<Map<string, ChessPiece>> {
    // Returns an observable after sequentially decoding JSON string and filtering into the map via rxjs pipe.
    return this.http.get<BoardApiRequest>(`http://127.0.0.1:5000/api/mock-board`).pipe(
      map(state => this.fenDecode(state.fenString)), 
      map(items => {
          const itemMap = new Map<string, ChessPiece>();
          items.forEach(item => itemMap.set(`${item.position.x},${item.position.y}`, item))
          return itemMap;
        }
      )
    );
  }

  private fenDecode(fenString: string): ChessPiece[] {
    // TODO add proper decoding. https://github.com/dstettler/simpleboard/issues/35
    // Mock return.
    console.log(fenString);
    return [
      new Rook(false, 0, 0),
      new Knight(false, 0, 1),
      new Bishop(false, 0, 2),
      new Queen(false, 0, 3),
      new King(false, 0, 4),
      new Pawn(false, 1, 0),
      new Pawn(true, 6, 0),
      new Rook(true, 7, 0),
      new Knight(true, 7, 1),
      new Bishop(true, 7, 2),
      new Queen(true, 7, 3),
      new King(true, 7, 4),
    ]
  }
}
