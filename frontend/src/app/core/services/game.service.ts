import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { tap } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class GameService {
  private http = inject(HttpClient);

  
  private baseUrl = 'http://localhost:8080/api';

  
  currentGame: any = null;

  
  createGame() {
    return this.http.post(`${this.baseUrl}/game`, {
      action: 'create'
    }).pipe(
      tap((res: any) => {
        console.log('GameService: storing game', res);
        this.currentGame = res;
      })
    );
  }

  
  setGame(game: any) {
    this.currentGame = game;
  }

  getGame() {
    return this.currentGame;
  }

  
  hasGame() {
    return this.currentGame !== null;
  }

  
  clearGame() {
    this.currentGame = null;
  }
}