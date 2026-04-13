import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class GameService {

  createGame(): Observable<{ gameId: string }> {
    const fakeId = Math.random().toString(36).substring(2, 9);
    return of({ gameId: fakeId });
  }
}