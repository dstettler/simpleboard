import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map } from 'rxjs/operators';
import { AuthStateService } from './auth-state.service';

@Injectable({
  providedIn: 'root'
})
export class GameService {
  private http = inject(HttpClient);
  private authState = inject(AuthStateService);
  private baseUrl = 'http://localhost:8080/api';

  private authPayload() {
    const userId = this.authState.userId();

    if (this.authState.isGuest()) {
      return { guest_id: userId };
    }

    const playerId = Number(userId);
    return Number.isNaN(playerId) || playerId <= 0 ? {} : { player_id: playerId };
  }

  createGame() {
    return this.http.post<any>(`${this.baseUrl}/game`, {
      action: 'create',
      starting_side: 'w',
      ...this.authPayload()
    }).pipe(
      map((res: any) => {
        const gameId = res.state?.game_id;

        if (!gameId) {
          throw new Error('Backend did not return a game id');
        }

        return String(gameId);
      })
    );
  }

  joinGame(gameId: string) {
    return this.http.post<any>(`${this.baseUrl}/game`, {
      action: 'join',
      game_id: gameId,
      ...this.authPayload()
    }).pipe(
      map((res: any) => {
        const id = res.state?.game_id;

        if (!id) {
          throw new Error('Join failed: no game id');
        }

        return String(id);
      })
    );
  }

  getGameState(gameId: string) {
    return this.http.post<any>(`${this.baseUrl}/game`, {
      action: 'state',
      game_id: gameId,
      ...this.authPayload()
    });
  }

  makeMove(gameId: string, move: string) {
    return this.http.post<any>(`${this.baseUrl}/game`, {
      action: 'move',
      game_id: gameId,
      move,
      ...this.authPayload()
    });
  }
}