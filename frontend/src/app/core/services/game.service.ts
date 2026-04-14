import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map } from 'rxjs/operators';

import { AuthStateService } from './auth-state.service';

@Injectable({
  providedIn: 'root'
})
export class GameService {
  private http = inject(HttpClient);
  private authService = inject(AuthStateService);
  private baseUrl = 'http://localhost:8080/api';

  createGame() {
    console.log(`using uid ${this.authService.userId()}`);
    const playerId = Number(this.authService.userId());
    let otherId;
    if (playerId == 1) {
      otherId = 2;
    } else {
      otherId = 1;
    }

    return this.http.post<any>(`${this.baseUrl}/game`, {
      action: 'create',
      player_id: playerId,
      other_id: otherId,
      starting_side: 'w'
    }).pipe(
      map((res: any) => {
        console.log('RAW CREATE GAME RESPONSE:', JSON.stringify(res, null, 2));

        const gameId = res.state.game_id;

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
