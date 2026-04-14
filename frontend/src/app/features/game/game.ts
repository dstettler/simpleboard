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

  game: any;

  ngOnInit(): void {
    const gameId = this.route.snapshot.paramMap.get('id');

    if (!gameId) {
      console.warn('No game id found');
      this.router.navigate(['/']);
      return;
    }

    console.log('Game ID:', gameId);

    this.gameService.getGame(gameId).subscribe({
      next: (res: any) => {
        console.log('Loaded game:', res);
        this.game = res;
      },
      error: (err: any) => {
        console.error('Failed to load game:', err);
      }
    });
  }
}