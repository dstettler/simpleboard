import { Component, inject } from '@angular/core';
import { Router } from '@angular/router';
import { GameService } from '../../core/services/game.service';

@Component({
  selector: 'app-home',
  standalone: true,
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
      next: (res) => {
        this.router.navigate(['/game', res.gameId]);
      },
      error: () => {
        this.isCreating = false;
      }
    });
  }
}