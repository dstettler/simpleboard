import { Component, inject } from '@angular/core';
import { Router, RouterLink, RouterLinkActive } from '@angular/router';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { GameService } from '../../core/services/game.service';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [CommonModule, RouterLink, RouterLinkActive, FormsModule],
  templateUrl: './navbar.html',
  styleUrl: './navbar.css'
})
export class NavbarComponent {
  private router = inject(Router);
  private gameService = inject(GameService);

  isCreating = false;
  showGameModal = false;
  joinGameId = '';

  openGameModal(): void {
    this.showGameModal = true;
  }

  closeGameModal(): void {
    this.showGameModal = false;
    this.joinGameId = '';
  }

  createShareGame(): void {
    if (this.isCreating) return;

    this.isCreating = true;

    this.gameService.createGame().subscribe({
      next: (gameId: string) => {
        console.log('game created id', gameId);
        this.closeGameModal();
        this.router.navigate(['/game', gameId]);
        this.isCreating = false;
      },
      error: (err: any) => {
        console.error('game creation failed', err);
        this.isCreating = false;
      }
    });
  }

  joinGame(): void {
    const id = this.extractGameId(this.joinGameId);

    if (!id) return;

    this.closeGameModal();
    this.router.navigate(['/game', id]);
  }

  private extractGameId(value: string): string {
    const trimmed = value.trim();

    if (!trimmed) return '';

    if (trimmed.includes('/game/')) {
      return trimmed.split('/game/')[1].split(/[?#]/)[0];
    }

    return trimmed;
  }
}