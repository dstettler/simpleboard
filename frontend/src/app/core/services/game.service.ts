import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class GameService {
  private http = inject(HttpClient);
  private baseUrl = 'http://localhost:8080/api';

  createGame() {
    return this.http.post<any>(`${this.baseUrl}/game`, {
      action: 'create'
    }).pipe(
      map((res: any) => {
        console.log('RAW CREATE GAME RESPONSE:', JSON.stringify(res, null, 2));

        const gameId =
          res.game_id ??
          res.id ??
          res.gameId ??
          res.game?.game_id ??
          res.data?.game_id ??
          res.user?.game_id;

        if (!gameId) {
          throw new Error('Backend did not return a game id');
        }

        return String(gameId);
      })
    );
  }

  getGame(gameId: string) {
    return this.http.get<any>(`${this.baseUrl}/game/${gameId}`);
  }
}