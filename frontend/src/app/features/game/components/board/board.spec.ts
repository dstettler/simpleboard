import { ComponentFixture, TestBed } from '@angular/core/testing';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { provideHttpClient } from '@angular/common/http';

import { Board } from './board';
import { API_ENDPOINT } from '../../../../app.constants';

export const mockBoardStateResponse = { "fenString": "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1" };

describe('Board', () => {
  let component: Board;
  let fixture: ComponentFixture<Board>;
  let httpMock: HttpTestingController;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [Board],
      providers: [
        provideHttpClient(),
        provideHttpClientTesting()
      ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(Board);
    component = fixture.componentInstance;

    httpMock = TestBed.inject(HttpTestingController);
    await fixture.whenStable();
  });

  it('should create', () => {
    const ex = httpMock.expectOne(`${API_ENDPOINT}/api/mock-board`);
    ex.flush(mockBoardStateResponse);

    expect(component).toBeTruthy();
  });
});
