import { Component, inject, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { GameService } from '../../core/services/game.service';
import { Board } from './components/board/board';

@Component({
  selector: 'app-game',
  standalone: true,
  imports: [Board],
  templateUrl: './game.html',
  styleUrl: './game.css'
})
export class Game implements OnInit {
  private route = inject(ActivatedRoute);
  private router = inject(Router);
  private gameService = inject(GameService);

  gameId = '';
  game: any;

  ngOnInit(): void {
    const paramId = this.route.snapshot.paramMap.get('id');

    if (!paramId) {
      this.router.navigate(['/']);
      return;
    }

    this.gameId = paramId;

    this.gameService.getGameState(this.gameId).subscribe({
      next: (res: any) => {
        console.log('Game state loaded:', res);
        this.game = res.state;
      },
      error: (err: any) => {
        console.error('Failed to load game state:', err);
      }
    });
  }
}