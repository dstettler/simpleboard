import { ComponentFixture, TestBed } from '@angular/core/testing';

import { Game } from './game';
import { provideHttpClient } from '@angular/common/http';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';

import { mockBoardStateResponse } from './components/board/board.spec';

import { API_ENDPOINT } from '../../app.constants'

describe('Game', () => {
  let component: Game;
  let fixture: ComponentFixture<Game>;
  let httpMock: HttpTestingController;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [Game],
      providers: [
        provideHttpClient(),
        provideHttpClientTesting()
      ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(Game);
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
