import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map, switchMap } from 'rxjs/operators';
import { Observable, of } from 'rxjs';

import { API_ENDPOINT } from '../../app.constants';

import { AuthStateService } from './auth-state.service';

@Injectable({
  providedIn: 'root'
})
export class GameService {
  private http = inject(HttpClient);
  private authService = inject(AuthStateService);

  private baseUrl = 'http://localhost:8080/api';

  private guestId: string | null = null;
  private guestToken: string | null = null;

  private ensureGuestIfNeeded(): Observable<{ id: string; isGuest: boolean }> {
  const userId = this.authService.userId?.();

  if (userId) {
    return of({ id: String(userId), isGuest: false });
  }

  if (this.guestId && this.guestToken) {
    return of({ id: this.guestId, isGuest: true });
  }

  return this.http.get<any>(`${this.baseUrl}/guest`).pipe(
    map((res: any) => {
      const guestId = res.user?.guest_id;
      const token = res.token;

      if (!guestId || !token) {
        throw new Error('Guest endpoint did not return guest_id or token');
      }

      this.guestId = guestId;
      this.guestToken = token;

      localStorage.setItem('token', token);

      return { id: guestId, isGuest: true };
    })
  );
}

  createGame() {
    console.log('creating game')
    return this.http.get<any>(`${API_ENDPOINT}/api/health`).pipe(switchMap(() => {
      console.log('health pinged')
        const payload = this.authService.isGuest()
          ? {
              action: 'create',
              guest_id: this.authService.userId(),
              starting_side: 'w'
            }
          : {
              action: 'create',
              player_id: Number(this.authService.userId()),
              starting_side: 'w'
            };

        return this.http.post<any>(`${this.baseUrl}/game`, payload);
    }),
                     map((res: any) => {
      const gameId = res.state?.game_id;

      if (!gameId) {
          throw new Error('Backend did not return a game id');
        }

        return String(gameId);
      })
   )
  }
  // return this.ensureGuestIfNeeded().pipe(
    //   switchMap((identity) => {
    //     const payload = identity.isGuest
    //       ? {
    //           action: 'create',
    //           guest_id: identity.id,
    //           starting_side: 'w'
    //         }
    //       : {
    //           action: 'create',
    //           player_id: Number(identity.id),
    //           starting_side: 'w'
    //         };
    //
    //     return this.http.post<any>(`${this.baseUrl}/game`, payload);
    //   }),
    //   map((res: any) => {
    //     const gameId = res.state?.game_id;
    //
    //     if (!gameId) {
    //       throw new Error('Backend did not return a game id');
    //     }
    //
    //     return String(gameId);
    //   })
    // );
  //}

  joinGame(gameId: string) {
    return this.http.get<any>(`${API_ENDPOINT}/api/health`).pipe(
      switchMap(() => {
        const payload = this.authService.isGuest()
          ? {
              action: 'join',
              game_id: gameId,
              guest_id: this.authService.userId()
            }
          : {
              action: 'join',
              game_id: gameId,
              player_id: Number(this.authService.userId())
            };

        return this.http.post<any>(`${this.baseUrl}/game`, payload);
      }),
      map((res: any) => {
        const joinedGameId = res.state?.game_id;

        if (!joinedGameId) {
          throw new Error('Backend did not return joined game id');
        }

        return String(joinedGameId);
      })
    );
  }

  getGameState(gameId: string) {
    return this.http.get<any>(`${API_ENDPOINT}/api/health`).pipe(
      switchMap(() => {
        const payload = this.authService.isGuest()
          ? {
              action: 'state',
              game_id: gameId,
              guest_id: this.authService.userId()
            }
          : {
              action: 'state',
              game_id: gameId,
              player_id: Number(this.authService.userId())
            };

        return this.http.post<any>(`${this.baseUrl}/game`, payload);
      })
    );
  }

  makeMove(gameId: string, move: string) {
    return this.http.get<any>(`${API_ENDPOINT}/api/health`).pipe(
      switchMap(() => {
        const payload = this.authService.isGuest()
          ? {
              action: 'move',
              game_id: gameId,
              guest_id: this.authService.userId(),
              move
            }
          : {
              action: 'move',
              game_id: gameId,
              player_id: Number(this.authService.userId()),
              move
            };

        return this.http.post<any>(`${this.baseUrl}/game`, payload);
      })
    );
  }
}
