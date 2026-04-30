import { Component, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';

import { HttpClient } from '@angular/common/http';

import { API_ENDPOINT } from '../app.constants';
import { firstValueFrom } from 'rxjs';
import { toObservable } from '@angular/core/rxjs-interop';

type Stats = {
 games: number;
    wins: number;
    losses:number;
    winRate:number;
    timePlayed:string;
}

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './dashboard.html',
  styleUrl: './dashboard.css'
})
export class DashboardComponent {
  private http = inject(HttpClient);

  _stats = signal<Stats>({
    games: 42,
    wins: 25,
    losses: 17,
    winRate: 59,
    timePlayed: '12h 35m'
  });


  statsObservable$ = toObservable(this._stats);

  ngOnInit() {
    firstValueFrom(this.http.get(`${API_ENDPOINT}/api/dashboard`)).then((resp: any) => {
      const updated = {
        games: resp.total_games,
        wins: resp.wins,
        losses: resp.losses,
        timePlayed: '12h 35m',
        winRate: resp.win_rate,
      };
      this._stats.update(_ => updated);

    })
  }

  get winRank(): { icon: string; title: string } {
    const rate = this._stats().winRate;

    if (rate < 40) return { icon: '♙', title: 'Pawn' };
    if (rate < 55) return { icon: '♘', title: 'Knight' };
    if (rate < 65) return { icon: '♗', title: 'Bishop' };
    if (rate < 75) return { icon: '♖', title: 'Rook' };
    if (rate < 90) return { icon: '♕', title: 'Queen' };

    return { icon: '♔', title: 'King' };
  }

  round(n: number): number {
    return Math.round(n);
  }
}
