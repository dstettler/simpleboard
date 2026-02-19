import { TestBed } from '@angular/core/testing';

import { BoardLoadService } from './board-load-service';

describe('BoardLoadService', () => {
  let service: BoardLoadService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(BoardLoadService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
