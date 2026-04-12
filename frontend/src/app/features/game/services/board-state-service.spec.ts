import { TestBed } from '@angular/core/testing';

import { BoardStateService } from './board-state-service';
import { mockPositions } from './BoardState';

describe('BoardLoadService', () => {
  let service: BoardStateService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(BoardStateService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should decode fen strings', () => {
    const knownPositions = mockPositions();

    const mockPositionString = 'rnbqk3/p7/8/8/8/8/P7/RNBQK3 b KQkq - 0 1';

    service.fenDecode(mockPositionString);

    expect(service.pieces()).toStrictEqual(knownPositions);
  })
});
