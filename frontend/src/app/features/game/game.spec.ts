import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideHttpClient } from '@angular/common/http';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { provideRouter } from '@angular/router';

import { Game } from './game';
import { API_ENDPOINT } from '../../app.constants'
import { mockBoardStateResponse } from './services/board-state-service.spec';

describe('Game', () => {
  let component: Game;
  let fixture: ComponentFixture<Game>;
  let httpMock: HttpTestingController;

  beforeEach(async () => {
    const localStorage = {
      getItem: (_key: string) => {
        return "1";
      }
    }

    Object.defineProperty(window, 'localStorage', { value:  localStorage });

    await TestBed.configureTestingModule({
      imports: [Game],
      providers: [
        provideHttpClient(),
        provideHttpClientTesting(),
        provideRouter([{path: 'game/:id', component: Game}])
      ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(Game);
    component = fixture.componentInstance;

    httpMock = TestBed.inject(HttpTestingController);
    await fixture.whenStable();
  });

  it('should create', () => {
    const ex = httpMock.expectOne(`${API_ENDPOINT}/api/game`);
    ex.flush(mockBoardStateResponse);

    expect(component).toBeTruthy();
  });
});
