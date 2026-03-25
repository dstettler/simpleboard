import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './dashboard.html',
  styleUrl: './dashboard.css'
})
export class DashboardComponent {
  stats = {
    games: 42,
    wins: 25,
    losses: 17,
    winRate: 59,
    timePlayed: '12h 35m'
  };

  get winRank(): { icon: string; title: string } {
    const rate = this.stats.winRate;

    if (rate < 40) return { icon: '♙', title: 'Pawn' };
    if (rate < 55) return { icon: '♘', title: 'Knight' };
    if (rate < 65) return { icon: '♗', title: 'Bishop' };
    if (rate < 75) return { icon: '♖', title: 'Rook' };
    if (rate < 90) return { icon: '♕', title: 'Queen' };

    return { icon: '♔', title: 'King' };
  }
}