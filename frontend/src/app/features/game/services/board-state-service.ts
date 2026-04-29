import { inject, Injectable, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { EMPTY, map, Observable, switchMap, timer } from 'rxjs';

import { ChessPiece, getPieceFromFenCharacter } from './pieces/ChessPiece';
import { algebraicToPosition, Position, positionsEqual, positionToAlgebraic } from './pieces/Position';
import { API_ENDPOINT, BACKEND_PING_RATE_MS } from '../../../app.constants';
import { GameStatus, parseGameStatus } from './BoardState';

import { takeUntilDestroyed, toObservable } from '@angular/core/rxjs-interop';

type GameRequest = {
  action: string;
  game_id: number;
  player_id: number;
  move: string;
}

type ResponseUser = {
  state: GameApiResponse
}

type GameApiResponse = {
  state: string;
  status: string;
  side: string;
  black_player_id: string;
  white_player_id: string;
  next_moves: string[];
  prev_moves: string[];
  white_remaining_ms: number;
  black_remaining_ms: number;
  last_move_at: string;
  server_time: string;
  updated_at: string;
}

type GameApiError = {
  error: string
}

interface Move {
  pieceIdx: number;
  targetPos: Position;
}

export type PlayerColor = 'w' | 'b' | null;

@Injectable({ providedIn: 'root' })
export class BoardStateService {
  private http = inject(HttpClient);

  private _pieces = signal<ChessPiece[]>([]);
  private _isOwnMove = signal<boolean>(false);
  private _castleables = signal<string>('');
  private _enPassant = signal<Position|null>(null);
  private _halfmoveClock = signal<number>(-1);
  private _fullmoveNum = signal<number>(-1);
  private _userColor = signal<PlayerColor>(null);
  private _nextMoves = signal<Move[]>([]);
  private _gameStatus = signal<GameStatus>("Waiting");
  private _playerId = signal<number>(-1);
  private _gameId = signal<number>(-1);
  private _timerRemainingMs = signal<number>(-1);
  private _pollBackend = signal<boolean>(false);
  private _gameTimerRunning = signal<boolean>(false);

  readonly pieces = this._pieces.asReadonly();
  readonly isOwnMove = this._isOwnMove.asReadonly();
  readonly castleables = this._castleables.asReadonly();
  readonly enPassant = this._enPassant.asReadonly();
  readonly halfmoveClock = this._halfmoveClock.asReadonly();
  readonly fullmoveNum = this._fullmoveNum.asReadonly();
  readonly userColor = this._userColor.asReadonly();
  readonly nextMoves = this._nextMoves.asReadonly();
  readonly gameStatus = this._gameStatus.asReadonly();
  readonly playerId = this._playerId.asReadonly();
  readonly gameId = this._gameId.asReadonly();
  readonly timerRemainingMs = this._timerRemainingMs.asReadonly();
  private readonly pollBackend = this._pollBackend.asReadonly();
  private readonly gameTimerRunning = this._gameTimerRunning.asReadonly();

  private poll$ = toObservable(this.pollBackend).pipe(
    switchMap(p => p ? timer(0, BACKEND_PING_RATE_MS) : EMPTY),
    switchMap(() => {
      if (this.playerId() != -1 && this.gameId() != -1) {
        return this.boardLoad(this.gameId(), this.playerId())
      } else {
        return EMPTY;
      }
    }),
    takeUntilDestroyed()
  )

  private gameTimer$ = toObservable(this.gameTimerRunning).pipe(
    // If timer is runing (p), wait one second, update GUI signal
    switchMap(p => p ? timer(0, 1000) : EMPTY),
    switchMap(() => {
      if (this.playerId() != -1 && this.gameId() != -1) {
        this._timerRemainingMs.update(prevTime => prevTime - 1000);
        console.log(`tick ${this.timerRemainingMs()}`);
      }

      return EMPTY;
    }),
    takeUntilDestroyed()
  )

  constructor() {
    this.poll$.subscribe();
    this.gameTimer$.subscribe();
  }

  /**
   * @returns Observable<void>, so the caller may make use of the completion event after request completion and state update.
   */
  boardLoad(gameId: number, playerId: number): Observable<void> {
    // Load initial state
    const req: GameRequest = {
      action:"state",
      game_id: gameId,
      player_id: playerId,
      move: ""
    };

    this._playerId.update(_i => playerId);
    this._gameId.update(_i => gameId);

    return this.gameRequest(req);
  }

  gameRequest(reqBody: GameRequest): Observable<void> {
    // Returns an observable after sequentially decoding JSON string and filtering into the map via rxjs pipe.
    return this.http.post<ResponseUser|GameApiError>(`${API_ENDPOINT}/api/game`, reqBody).pipe(
      map(state => {
        if ("error" in state) {
          const err = state as GameApiError;
          // Illegal operation
          console.error(err.error);
          return;
        } else {
          const resp = state as ResponseUser;

          // User color should only need to be updated once.
          // This must happen before decoding state so pieces can be enabled.
          if (this.userColor() == null) {
            console.log('setting user color')
            switch (this.playerId()) {
              case Number(resp.state.black_player_id):
                this._userColor.update(_ => 'b');
                break;
              case Number(resp.state.white_player_id):
                this._userColor.update(_ => 'w');
                break;
              default:
                console.error("Player id matches neither side.");
                throw new Error("Player id matches neither side.");
            }
          }

          this.fenDecode(resp.state.state);
          this._nextMoves.update(_p => resp.state.next_moves.map(move_str => {
            let start: Position, finish: Position;

            if (move_str.length == 5) {
              // move_str == "a1xb1"
              start = algebraicToPosition(move_str.slice(0, 2));
              finish = algebraicToPosition(move_str.slice(3));
            } else if (move_str.length == 4) {
              // move_str == "a1b1"
              start = algebraicToPosition(move_str.slice(0, 2));
              finish = algebraicToPosition(move_str.slice(2));
            } else {
              console.error(`Invalid next move str: ${move_str}`);
              throw new Error("Invalid next move str");
            }

            for (const [i, piece] of this._pieces().entries()) {
              const pieceIsUserColor = piece.isWhite == (this.userColor() == 'w');
              const pieceEnableable = this.isOwnMove() && pieceIsUserColor;
              if (positionsEqual(piece.position, start)) {
                this._pieces.update(p => { p[i].enabled = pieceEnableable; return p; });
                return {pieceIdx: i, targetPos: finish};
              }
            }

            console.error(`No piece found for target move ${move_str}`);
            throw new Error(`No piece found for target move ${move_str}`);
          }));

         const serverTimeStr = resp.state.server_time;
          const serverTime = new Date(serverTimeStr + "Z");
          console.log(serverTime)

          const timeDelta = Date.now() - serverTime.getTime();

          if (this.userColor() == 'w') {
            const trueRemaining = resp.state.white_remaining_ms;
            this._timerRemainingMs.update(_ => trueRemaining - timeDelta);
          } else if (this.userColor() == 'b') {
            const trueRemaining = resp.state.black_remaining_ms;
            this._timerRemainingMs.update(_ => trueRemaining - timeDelta);
          }

          this._gameStatus.update(_p => parseGameStatus(resp.state.status));
          if (this.isOwnMove()) {
            this._gameTimerRunning.update(_ => true);
            this._pollBackend.update(_ => false);
          } else {
            this._gameTimerRunning.update(_ => false);
            this._pollBackend.update(_ => true);
          }

          return;
        }
      })
    );
  }

  updatePiecePosition(gameId: number, playerId: number, piece: ChessPiece, newPos: Position): Observable<void> {
    let captureChar = '';
    for (const piece of this._pieces()) {
      if (positionsEqual(piece.position, newPos))
        captureChar = 'x';
    }

    const moveStr = `${positionToAlgebraic(piece.position)}${captureChar}${positionToAlgebraic(newPos)}`;
    console.log(`Moving: ${moveStr}`);

    const req: GameRequest = {
      action: "move",
      game_id: gameId,
      player_id: playerId,
      move: moveStr
    };

    return this.gameRequest(req);
  }

  public fenDecode(fenString: string) {
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
          currentX++;
        } else {
          const offset = parseInt(char);
          currentX = currentX + offset;
        }
      }

      currentX = 0;
      currentY++;
    }

    this._pieces.update(_p => pieces);

    this._isOwnMove.update(_m => this.userColor() == activeColor);
    this._castleables.update(_c => castleable);
    if (enPassant != '-') {
      this._enPassant.update(_e => algebraicToPosition(enPassant));
    } else {
      this._enPassant.update(_e => null);
    }

    this._halfmoveClock.update(_c => halfmoveClock);
    this._fullmoveNum.update(_n => fullmoveNumber);
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
