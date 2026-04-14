import { Component, inject } from '@angular/core';
import { Router } from '@angular/router';
import { GameService } from '../../core/services/game.service';

@Component({
  selector: 'app-home',
  standalone: true,
  templateUrl: './home.html'
})
export class HomeComponent {
  Math= Math;
  private router = inject(Router);
  private gameService = inject(GameService);

  onStartGame() {
    this.gameService.createGame().subscribe({
      next: (res: any) => {
        console.log('game created', res);

        this.gameService.setGame(res);

        this.router.navigate(['/game']);
      },
      error: (err) => {
        console.error('game error', err);
      }
    });
  }
}