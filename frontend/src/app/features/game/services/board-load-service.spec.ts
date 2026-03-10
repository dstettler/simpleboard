import { TestBed } from '@angular/core/testing';

import { BoardLoadService } from './board-load-service';
import { mockPositions } from './BoardState';

describe('BoardLoadService', () => {
  let service: BoardLoadService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(BoardLoadService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should decode fen strings', () => {
    const knownPositions = mockPositions();

    const mockPositionString = 'rnbqk3/p7/8/8/8/8/P7/RNBQK3 b KQkq - 0 1';

    const createdPositions = service.fenDecode(mockPositionString);

    expect(createdPositions).toStrictEqual(knownPositions);
  })
});
