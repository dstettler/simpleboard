import { inject, Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map, Observable } from 'rxjs';

import { ChessPiece, getPieceFromFenCharacter } from './pieces/ChessPiece';
import { Position } from './pieces/Position';
import { mockPositions } from './BoardState';
import { API_ENDPOINT } from '../../../app.constants';

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
    return this.http.get<BoardApiRequest>(`${API_ENDPOINT}/api/mock-board`).pipe(
      map(state => this.fenDecode(state.fenString))
    );
  }

  updatePiecePosition(piece: ChessPiece, newPos: Position): Observable<ChessPiece[]> {
    return this.http.get<UpdatedRequest>(`${API_ENDPOINT}/api/update-board`).pipe(
      map(_state => this.updatePos(piece, newPos))
    )
  }

  private updatePos(piece: ChessPiece, newPos: Position): ChessPiece[] {
    console.log(`${piece.id} new pos: ${newPos.x}, ${newPos.y}`)
    if (this.positionsArray !== null) {
      this.positionsArray[piece.id].position = newPos;
    } else {
      console.warn("Attempted to update pos before init")
      this.positionsArray = mockPositions();
    }

    return this.positionsArray;
  }

  public fenDecode(fenString: string): ChessPiece[] {
    console.log(`fenString is as follows: ${fenString}`);
    const fenFields = fenString.split(' ');

    const validation = this.validateFenFields(fenFields);
    if (!validation[0]) {
      const errorString = `Invalid FEN string provided by server: ${fenString}. Reason: ${validation[1]}`
      console.error(errorString);
      throw new Error(errorString);
    }

    const placement = fenFields[0]
    const activeColor = fenFields[1];
    const castleable = fenFields[2];
    const enPassant = fenFields[3];
    const halfmoveClock = parseInt(fenFields[4]);
    const fullmoveNumber = parseInt(fenFields[5]);

    const placementRanks = placement.split('/');

    let pieces: ChessPiece[] = [];
    let currentId = 0, currentX = 0, currentY = 0;
    for (const rank of placementRanks) {
      for (const char of rank) {
        if (Number.isNaN(parseInt(char))) {
          pieces.push(getPieceFromFenCharacter(char, currentId, currentX, currentY));
          currentId++;
          currentY++;
        } else {
          const offset = parseInt(char);
          currentY + offset;
        }
      }

      currentY = 0;
      currentX++;
    }

    this.positionsArray = pieces;
    return pieces;
  }

  private validatePlacementField(field: string): [boolean, string] {
    const ranks = field.split('/');
    const validPieces = "pnbrqkPNBRQK";

    if (ranks.length != 8) {
      return [false, `Invalid number of ranks: ${ranks.length}`];
    }
    
    let placementFreqDict: { [key: string]: number } = {}

    for (const rank of ranks) {
      let rankWidth = 0;

      for (const char of rank) {
        const charAsInt = parseInt(char);

        if (validPieces.includes(char.toUpperCase())) {
          if (char in placementFreqDict) {
            placementFreqDict[char] += 1;
          } else {
            placementFreqDict[char] = 1;
          }

          rankWidth++;
        } else if (!Number.isNaN(charAsInt)) {
          rankWidth += charAsInt;
        } else {
          return [false, `Invalid character ${char} in rank: ${rank}`]
        }
      }

      if (rankWidth != 8) {
        return [false, `Invalid rank size ${rankWidth} in rank: ${rank}`];
      }
    }

    const isPawn = (key: string) => {
      return "pP".includes(key);
    }

    for (const [key, val] of Object.entries(placementFreqDict)) {
      if (val > 1 && "qQkK".includes(key)) {
        return [false, `Invalid number of piece: ${key}, ${val}`];
      } else if (val > 8 && isPawn(key)) {
        return [false, `Invalid number of piece: ${key}, ${val}`];
      } else if (val > 2 && !isPawn(key)) {
        return [false, `Invalid number of piece: ${key}, ${val}`];
      }
    }

    return [true, ""];
  }

  public validateFenFields(fields: string[]): [boolean, string|null] {
    if (fields.length != 6) {
      return [false, 'Invalid number of FEN fields'];
    }
    
    // Placement
    const placementValidation = this.validatePlacementField(fields[0]);
    if (!placementValidation[0]) {
      return placementValidation;
    }

    // Active color
    if (fields[1] != 'w' && fields[1] != 'b') {
      return [false, 'Invalid active color field'];
    }

    // Castling
    const castleRegex = /^([KkQq]{1,4}|-)$/gm;
    if (!castleRegex.test(fields[2])) {
      return [false, 'Invalid castle field'];
    }
    
    // En passant
    const enPassantRegex = /^([a-hA-H][1-8])|-$/gm;
    if (!enPassantRegex.test(fields[3])) {
      return [false, 'Invalid en passant field'];
    }

    // Halfmove clock
    if (Number.isNaN(parseInt(fields[4]))) {
      return [false, 'Invalid halfmove clock field'];
    }

    // Fullmove number
    if (Number.isNaN(parseInt(fields[5]))) {
      return [false, 'Invalid fullmove number field'];
    }

    return [true, null];
  }
}
