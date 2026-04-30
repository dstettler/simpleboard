import { Component, inject } from '@angular/core';
import { Router } from '@angular/router';
import { GameService } from '../../core/services/game.service';
import { GameEntryModalComponent } from '../../shared/game-entry-modal/game-entry-modal';
@Component({
  selector: 'app-home',
  standalone: true,
  imports: [GameEntryModalComponent],
  templateUrl: './home.html',
  styleUrl: './home.css'
})
export class HomeComponent {
  Math = Math;

  private router = inject(Router);
  private gameService = inject(GameService);

  isCreating = false;

  onStartGame() {
    if (this.isCreating) return;

    this.isCreating = true;

    this.gameService.createGame().subscribe({
      next: (gameId: string) => {
        console.log('game created id', gameId);
        this.router.navigate(['/game', gameId]);
        this.isCreating = false;
      },
      error: (err: any) => {
        console.error('game creation failed', err);
        this.isCreating = false;
      }
    });
  }
  onCreateGame(): void {
  if (this.isCreating) return;

  this.isCreating = true;

  this.gameService.createGame().subscribe({
    next: (gameId: string) => {
      this.router.navigate(['/game', gameId]);
      this.isCreating = false;
    },
    error: (err: any) => {
      console.error('home create game failed', err);
      this.isCreating = false;
    }
  });
}
}