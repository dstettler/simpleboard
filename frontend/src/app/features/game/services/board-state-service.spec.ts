import { TestBed } from '@angular/core/testing';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { provideHttpClient } from '@angular/common/http';

import { BoardStateService } from './board-state-service';
import { ChessPiece, Rook, Knight, Bishop, Queen, King, Pawn } from './pieces/ChessPiece';
import { positionToAlgebraic, algebraicToPosition, positionsEqual } from './pieces/Position';
import { API_ENDPOINT } from '../../../app.constants';

export const mockBoardStateResponse = { state: {
  "state": "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
    "next_moves": ["a2a3"],
    "side": "w",
    "white_player_id": 0,
    "black_player_id": 1 } };

function mockPositions(): ChessPiece[] {
    return [
      new Rook(0, false, 0, 0),
      new Knight(1, false, 1, 0),
      new Bishop(2, false, 2, 0),
      new Queen(3, false, 3, 0),
      new King(4, false, 4, 0),
      new Pawn(5, false, 0, 1),
      new Pawn(6, true, 0, 6),
      new Rook(7, true, 0, 7),
      new Knight(8, true, 1, 7),
      new Bishop(9, true, 2, 7),
      new Queen(10, true, 3, 7),
      new King(11, true, 4, 7),
    ]
  }

describe('BoardStateService', () => {
  let service: BoardStateService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        provideHttpClient(),
        provideHttpClientTesting()
      ]
    });
    service = TestBed.inject(BoardStateService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should decode fen strings', () => {
    const knownPositions = mockPositions();

    const mockPositionString = 'rnbqk3/p7/8/8/8/8/P7/RNBQK3 b KQkq - 0 1';

    service.fenDecode(mockPositionString);

    expect(service.pieces()).toStrictEqual(knownPositions);
  });

  it('should load state from request', () => {
    service.boardLoad(0, 1).subscribe(_ => {
      expect(service.userColor()).toBe('b');
      expect(service.fullmoveNum()).toBe(1);

      // Only one next move is provided, so we can directly index this
      expect(service.nextMoves()[0].pieceIdx).toBe(16);
      expect(service.nextMoves()[0].targetPos).toStrictEqual({x: 0, y: 5})
    });

    const ex = httpMock.expectOne(`${API_ENDPOINT}/api/game`);
    ex.flush(mockBoardStateResponse);
  });

  it('should decode positions properly', () => {
    expect(algebraicToPosition('a8')).toStrictEqual({x: 0, y: 0});
    expect(algebraicToPosition('a1')).toStrictEqual({x: 0, y: 7});
    expect(algebraicToPosition('h8')).toStrictEqual({x: 7, y: 0});
    expect(algebraicToPosition('h1')).toStrictEqual({x: 7, y: 7});
  });

  it('should encode positions properly', () => {
    expect(positionToAlgebraic({x: 0, y: 0})).toBe('a8');
    expect(positionToAlgebraic({x: 0, y: 7})).toBe('a1');
    expect(positionToAlgebraic({x: 7, y: 0})).toBe('h8');
    expect(positionToAlgebraic({x: 7, y: 7})).toBe('h1');
  });

  it('should determine position equality properly', () => {
    expect(positionsEqual({x: 0, y: 1}, {x: 0, y: 1})).toBe(true);
    expect(positionsEqual({x: 0, y: 0}, {x: 0, y: 1})).toBe(false);
    expect(positionsEqual({x: 1, y: 1}, {x: 0, y: 1})).toBe(false);
  });
});
