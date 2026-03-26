import { inject, Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map, Observable, switchMap } from 'rxjs';

import { ChessPiece, getPieceFromFenCharacter } from './pieces/ChessPiece';
import { Position, positionToAlgebraic } from './pieces/Position';
import { API_ENDPOINT } from '../../../app.constants';

type GameRequest = {
  action: string;
  game_id: number;
  player_id: number;
  move?: string;
}

type GameApiResponse = {
  state: string;
  status: string;
  side: string;
  next_moves: string[]
  prev_moves: string[]
}

type GameApiError = {
  error: string
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
  boardLoad(gameId: number, playerId: number): Observable<ChessPiece[]> {
    // Load initial state
    const req: GameRequest = {
      action: "state",
      game_id: gameId,
      player_id: playerId
    };

    return this.gameRequest(req);
  }

  gameRequest(reqBody: GameRequest): Observable<ChessPiece[]> {
    // Returns an observable after sequentially decoding JSON string and filtering into the map via rxjs pipe.
    return this.http.post<GameApiResponse|GameApiError>(`${API_ENDPOINT}/api/game`, reqBody).pipe(
      map(state => {
        if ("error" in state) {
          const err = state as GameApiError;
          // Illegal operation
          console.error(err.error);
          if (this.positionsArray == null) {
            return [];
          } else {
            return this.positionsArray;
          }
        } else {
          const resp = state as GameApiResponse;
          const ret = this.fenDecode(resp.state);
          this.positionsArray = ret;
          return this.fenDecode(resp.state);
        }
      })
    );
  }

  updatePiecePosition(gameId: number, playerId: number, piece: ChessPiece, newPos: Position): Observable<ChessPiece[]> {
    let captureChar = '';
    if (this.positionsArray) {
      console.log(newPos);
      for (const piece of this.positionsArray) {
        const isSamePos = piece.position.x == newPos.x && piece.position.y == newPos.y;
        if (isSamePos)
          captureChar = 'x';
      }
    }

    const moveStr = `${positionToAlgebraic(piece.position)}${captureChar}${positionToAlgebraic(newPos)}`;
    console.log(`Moving: ${moveStr}`);

    const req: GameRequest = {
      action: "move",
      game_id: gameId,
      player_id: playerId,
      move: moveStr
    };

    const stateReq: GameRequest = {
      action: "state",
      game_id: gameId,
      player_id: playerId,
    };


    return this.gameRequest(req).pipe(switchMap(() => this.gameRequest(stateReq)));
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
          currentY = currentY + offset;
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
